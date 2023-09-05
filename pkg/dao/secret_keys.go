package dao

import (
	"cloud.google.com/go/storage"
	"context"
	"crypto/ed25519"
	goerrors "errors"
	"fmt"
	"github.com/a-novel/go-framework/errors"
	"google.golang.org/api/iterator"
	"io"
	"os"
	"path"
	"sort"
	"strings"
	"time"
)

type SecretKeysRepository interface {
	// Write creates a new entry.
	Write(ctx context.Context, key ed25519.PrivateKey, name string) (*SecretKeyModel, error)
	// Read a key from the specified name.
	Read(ctx context.Context, name string) (*SecretKeyModel, error)
	// List all entries.
	List(ctx context.Context) ([]*SecretKeyModel, error)
	// Delete the specified entry.
	Delete(ctx context.Context, name string) error
}

type SecretKeyModel struct {
	// Key returns the decoded key for the current entry.
	Key ed25519.PrivateKey
	// Date returns the date when the key was created.
	Date time.Time
	// Name of the record (file) that stores the entry.
	Name string
}
type fileSystemRepositoryImpl struct {
	basePath string
	prefix   string
}

func NewFileSystemSecretKeysRepository(basePath, prefix string) SecretKeysRepository {
	return &fileSystemRepositoryImpl{basePath: basePath, prefix: prefix}
}

func (repository *fileSystemRepositoryImpl) getPath(name string) string {
	if strings.HasPrefix(name, repository.prefix+"-") {
		return path.Join(repository.basePath, name)
	}

	return path.Join(repository.basePath, fmt.Sprintf("%s-%s", repository.prefix, name))
}

func (repository *fileSystemRepositoryImpl) removePrefix(name string) string {
	return strings.TrimPrefix(name, repository.prefix+"-")
}

func (repository *fileSystemRepositoryImpl) Write(_ context.Context, key ed25519.PrivateKey, name string) (*SecretKeyModel, error) {
	fileWriter, err := os.Create(repository.getPath(name))
	if err != nil {
		return nil, fmt.Errorf("failed to create file %q: %w", name, err)
	}

	defer fileWriter.Close()

	if err = writeKeyToOutput(fileWriter, key); err != nil {
		return nil, fmt.Errorf("failed to write key to file %q: %w", name, err)
	}

	stat, err := fileWriter.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to read stats of file %q: %w", name, err)
	}

	return &SecretKeyModel{
		Key:  key,
		Date: stat.ModTime(),
		Name: repository.removePrefix(stat.Name()),
	}, nil
}

func (repository *fileSystemRepositoryImpl) Read(_ context.Context, name string) (*SecretKeyModel, error) {
	fileReader, err := os.Open(repository.getPath(name))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, errors.ErrNotFound
		}

		return nil, fmt.Errorf("failed to open file %q: %w", repository.getPath(name), err)
	}

	fileData, err := io.ReadAll(fileReader)
	if err != nil {
		return nil, fmt.Errorf("failed to read content of file %q: %w", repository.getPath(name), err)
	}

	key, err := unmarshalPrivateKey(fileData)
	if err != nil {
		return nil, fmt.Errorf("failed to decode file %q: %w", repository.getPath(name), err)
	}

	stat, err := fileReader.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to read stats of file %q: %w", name, err)
	}

	return &SecretKeyModel{
		Key:  key,
		Date: stat.ModTime(),
		Name: repository.removePrefix(stat.Name()),
	}, nil
}

func (repository *fileSystemRepositoryImpl) List(ctx context.Context) ([]*SecretKeyModel, error) {
	entries, err := os.ReadDir(repository.basePath)
	if err != nil {
		return nil, err
	}

	var records []*SecretKeyModel

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		if strings.HasPrefix(entry.Name(), repository.prefix) {
			record, err := repository.Read(ctx, entry.Name())
			if err != nil {
				return nil, err
			}
			if record != nil {
				records = append(records, record)
			}
		}
	}

	sort.SliceStable(records, func(i, j int) bool {
		return records[i].Date.After(records[j].Date)
	})

	return records, nil
}

func (repository *fileSystemRepositoryImpl) Delete(_ context.Context, name string) error {
	if err := os.Remove(repository.getPath(name)); err != nil {
		if goerrors.Is(err, os.ErrNotExist) {
			return errors.ErrNotFound
		}

		return fmt.Errorf("failed to delete file %q: %w", name, err)
	}

	return nil
}

type googleDatastoreRepositoryImpl struct {
	bucket *storage.BucketHandle
}

func NewGoogleDatastoreSecretKeysRepository(bucket *storage.BucketHandle) SecretKeysRepository {
	return &googleDatastoreRepositoryImpl{bucket: bucket}
}

func (repository *googleDatastoreRepositoryImpl) Write(ctx context.Context, key ed25519.PrivateKey, name string) (*SecretKeyModel, error) {
	fileWriter := repository.bucket.Object(name).NewWriter(ctx)
	defer fileWriter.Close()

	if err := writeKeyToOutput(fileWriter, key); err != nil {
		return nil, fmt.Errorf("failed to write key to file %q: %w", name, err)
	}

	return &SecretKeyModel{
		Key:  key,
		Date: fileWriter.Updated,
		Name: fileWriter.Name,
	}, nil
}

func (repository *googleDatastoreRepositoryImpl) Read(ctx context.Context, name string) (*SecretKeyModel, error) {
	fileReader, err := repository.bucket.Object(name).NewReader(ctx)
	if err != nil {
		if goerrors.Is(err, storage.ErrObjectNotExist) {
			return nil, errors.ErrNotFound
		}

		return nil, fmt.Errorf("failed to acquire reader for file %q: %w", name, err)
	}
	defer fileReader.Close()

	data, err := io.ReadAll(fileReader)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %q: %w", name, err)
	}

	key, err := unmarshalPrivateKey(data)
	if err != nil {
		return nil, fmt.Errorf("failed to decode file %q: %w", name, err)
	}
	return &SecretKeyModel{
		Key:  key,
		Date: fileReader.Attrs.LastModified,
		Name: name,
	}, nil
}

func (repository *googleDatastoreRepositoryImpl) List(ctx context.Context) ([]*SecretKeyModel, error) {
	entries := repository.bucket.Objects(ctx, nil)

	var records []*SecretKeyModel
	for {
		entry, err := entries.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to read next entry: %w", err)
		}

		record, err := repository.Read(ctx, entry.Name)
		if err != nil {
			return nil, err
		}

		records = append(records, record)
	}

	sort.SliceStable(records, func(i, j int) bool {
		return records[i].Date.After(records[j].Date)
	})

	return records, nil
}

func (repository *googleDatastoreRepositoryImpl) Delete(ctx context.Context, name string) error {
	if err := repository.bucket.Object(name).Delete(ctx); err != nil {
		if goerrors.Is(err, storage.ErrObjectNotExist) {
			return errors.ErrNotFound
		}

		return fmt.Errorf("failed to acquire reader for file %q: %w", name, err)
	}

	return nil
}
