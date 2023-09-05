CREATE EXTENSION IF NOT EXISTS unaccent;
CREATE EXTENSION IF NOT EXISTS pg_trgm;

CREATE TYPE sex AS ENUM ('male', 'female');

--bun:split

CREATE TABLE IF NOT EXISTS credentials (
    id uuid PRIMARY KEY NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ,

    email_user VARCHAR(128),
    email_domain VARCHAR(128),
    email_validation_code VARCHAR(256),
    new_email_user VARCHAR(128),
    new_email_domain VARCHAR(128),
    new_email_validation_code VARCHAR(256),
    password_hashed VARCHAR(2048),
    password_validation_code VARCHAR(256),

    UNIQUE(email_user, email_domain),
    CONSTRAINT email_user_filled CHECK ( email_user <> '' AND email_user IS NOT NULL ),
    CONSTRAINT email_domain_filled CHECK ( email_domain <> '' AND email_domain IS NOT NULL ),
    CONSTRAINT require_full_new_email
        CHECK (
            /* Must either be all empty or all filled */
            (
                ( new_email_user = '' OR new_email_user IS NULL ) AND
                ( new_email_domain = '' OR new_email_domain IS NULL ) AND
                ( new_email_validation_code = '' OR new_email_validation_code IS NULL )
            ) OR (
                new_email_user IS NOT NULL AND new_email_user <> '' AND
                new_email_domain IS NOT NULL AND new_email_domain <> '' AND
                new_email_validation_code IS NOT NULL AND new_email_validation_code <> ''
            )
        )
);

CREATE TABLE IF NOT EXISTS identities (
    id uuid PRIMARY KEY NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ,

    first_name VARCHAR(32),
    last_name VARCHAR(32),
    sex sex,
    birthday TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS profiles (
    id uuid PRIMARY KEY NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ,

    username VARCHAR(64),
    slug VARCHAR(64) NOT NULL,

    UNIQUE(slug),
    CONSTRAINT slug_filled CHECK (slug <> '')
);

--bun:split

CREATE FUNCTION format_search(a TEXT) RETURNS TEXT
    LANGUAGE sql IMMUTABLE
    RETURNS NULL ON NULL INPUT
    RETURN unaccent(lower(a));

CREATE FUNCTION search_field(query TEXT, value TEXT) RETURNS REAL
    LANGUAGE sql IMMUTABLE
    RETURN CASE
        WHEN query = '' OR query IS NULL THEN 1.0
        WHEN value = '' OR value IS NULL THEN 0.0
        ELSE similarity(format_search(query), format_search(value))
    END;

CREATE VIEW users_view AS
    SELECT
        credentials.id AS id,
        LEAST(credentials.created_at, identities.created_at, profiles.created_at) AS created_at,
        GREATEST(credentials.updated_at, identities.updated_at, profiles.updated_at) AS updated_at,
        json_build_object(
            'email', json_build_object(
                'user', credentials.email_user,
                'domain', credentials.email_domain
            )
        ) AS credentials,
        json_build_object(
            'firstName', identities.first_name,
            'lastName', identities.last_name,
            'sex', identities.sex,
            'birthday', identities.birthday
        ) AS identity,
        json_build_object(
            'username', profiles.username,
            'slug', profiles.slug
        ) AS profile
    FROM credentials
        INNER JOIN identities ON credentials.id = identities.id
        INNER JOIN profiles ON credentials.id = profiles.id;

--bun:split

CREATE UNIQUE INDEX credentials_email ON credentials (email_user, email_domain);

CREATE UNIQUE INDEX profiles_slug ON profiles (slug);

CREATE INDEX IF NOT EXISTS user_profile_username_search ON profiles USING gin (format_search(username) gin_trgm_ops);
CREATE INDEX IF NOT EXISTS user_profile_slug_search ON profiles USING gin (slug gin_trgm_ops);
CREATE INDEX IF NOT EXISTS user_identity_name_search ON identities USING gin (format_search(first_name || ' ' || last_name) gin_trgm_ops);
