package framework

import (
	"context"
	"github.com/a-novel/auth-service/pkg/models"
	"github.com/a-novel/go-apis"
	"net/http"
	"net/url"
)

type Client interface {
	IntrospectToken(ctx context.Context, token string) (*models.UserTokenStatus, error)
	Ping() error
}

type clientImpl struct {
	url *url.URL
}

func NewClient(url *url.URL) Client {
	return &clientImpl{url: url}
}

func (a *clientImpl) IntrospectToken(ctx context.Context, token string) (*models.UserTokenStatus, error) {
	output := new(models.UserTokenStatus)
	return output, apis.CallHTTP(ctx, apis.CallHTTPConfig{
		Path:            a.url.JoinPath("/auth"),
		Method:          http.MethodGet,
		Headers:         map[string]string{"Authorization": token},
		SuccessStatuses: []int{http.StatusOK},
	}, output)
}

func (a *clientImpl) Ping() error {
	return apis.CallHTTP(context.Background(), apis.CallHTTPConfig{
		Path:            a.url.JoinPath("/ping"),
		Method:          http.MethodGet,
		SuccessStatuses: []int{http.StatusOK},
	}, nil)
}
