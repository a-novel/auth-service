DROP INDEX IF EXISTS credentials_email;

DROP INDEX IF EXISTS profiles_slug;

DROP INDEX IF EXISTS user_profile_username_search;
DROP INDEX IF EXISTS user_profile_slug_search;
DROP INDEX IF EXISTS user_identity_name_search;

--bun:split

DROP VIEW IF EXISTS users_view;

DROP FUNCTION IF EXISTS search_field;
DROP FUNCTION IF EXISTS format_search;

--bun:split

DROP TABLE IF EXISTS credentials;
DROP TABLE IF EXISTS identities;
DROP TABLE IF EXISTS profiles;

--bun:split

DROP EXTENSION IF EXISTS unaccent;
DROP EXTENSION IF EXISTS pg_trgm;

DROP TYPE sex;
