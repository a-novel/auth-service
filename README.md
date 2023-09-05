# Auth service

Manage platform users and authentication.

- [Endpoints](#endpoints)
  - [`GET /ping`](#get-ping)
  - [`GET /auth`](#get-authintrospect-token)
  - [`POST /auth`](#post-authlogin)
  - [`PUT /auth`](#put-authregister)
  - [`DELETE /email`](#delete-email)
  - [`DELETE /password`](#delete-password)
  - [`PATCH /email`](#patch-email)
  - [`PATCH /identity`](#patch-identity)
  - [`PATCH /profile`](#patch-profile)
  - [`PATCH /password`](#patch-password)
  - [`PATCH /email/validation`](#patch-emailvalidation)
  - [`PATCH /email/pending/validation`](#patch-emailpendingvalidation)
  - [`GET /email/validation`](#get-emailvalidation)
  - [`GET /email/pending/validation`](#get-emailpendingvalidation)
  - [`GET /email/exists`](#get-emailexists)
  - [`GET /slug/exists`](#get-slugexists)
  - [`GET /users`](#get-users)
  - [`GET /users/search`](#get-userssearch)
  - [`GET /user`](#get-user)
  - [`GET /user/me`](#get-userme)
- [Installation](#installation)
- [Commands](#commands)

## Endpoints

### `GET /ping`

Check if the service is up and running.

- **Response**:

  | Status | Reason        |
  |:-------|:--------------|
  | 200    | ✅︎ Success.   |

  | Response key |  Type  | Description               |
  |:-------------|:------:|:--------------------------|
  | `message`    | string | Always contains `"pong"`. |


### `GET /auth/introspect-token`

Return information about the state of the current user token. Automatically refreshes the token when needed.

- **Headers**:

  | Key              |  Type   | Required | Description   |
  |:-----------------|:-------:|:--------:|:--------------|
  | `Authorization`  | string  |          | Bearer token. |

- **Response**:

  | Status | Reason                |
  |:-------|:----------------------|
  | 200    | ✅︎ Success.           |
 
  | Response key       |   Type    | Description                                                                                                          |
  |:-------------------|:---------:|:---------------------------------------------------------------------------------------------------------------------|
  | `ok`               |  boolean  | True when the token is valid, False otherwise.                                                                       |
  | `expired`          |  boolean  | Indicates if the token has expired.                                                                                  |
  | `notIssued`        |  boolean  | Indicates if the token is not available yet.                                                                         |
  | `malformed`        |  boolean  | Indicates if the data encoded in the token is correctly formatted.                                                   |
  | `tokenRaw`         |  string   | The raw token string, as given in the request headers. If the token was refreshed, it contains the up-to-date value. |
  | `token`            |  object   | The decoded content of the token.                                                                                    |
  | `token.header`     |  object   | The decoded header of the token.                                                                                     |
  | `token.header.iat` | timestamp | The unix timestamp, in milliseconds, when the token becomes available.                                               |
  | `token.header.exp` | timestamp | The unix timestamp, in milliseconds, when the token becomes expired.                                                 |
  | `token.header.id`  |   uuid    | A unique identifier for the token. The identifier is kept when the token is refreshed.                               |
  | `token.payload`    |  object   | The content payload carried by the token.                                                                            |
  | `token.payload.id` |   uuid    | The unique identifier of the user who the token belongs to.                                                          |

### `POST /auth/login`

Logs the user in the platform.

- **Body**:

  | Key        |  Type   | Required | Description           | Condition                                                                                              |
  |:-----------|:-------:|:--------:|:----------------------|:-------------------------------------------------------------------------------------------------------|
  | `email`    | string  |    ✔     | Email of the user.    | - `[3-128]` characters long.<br/>- ``/^[a-zA-Z\d!#$%&'*+/=?^_`{\|}.~-]+@[a-z\d]{2,}(.[a-z\d]{2,})+$/`` |
  | `password` | string  |    ✔     | Password of the user. | - `[2-256]` characters long.                                                                           |

- **Response**:

  | Status | Reason                                 |
  |:-------|:---------------------------------------|
  | 200    | ✅︎ Success.                            |
  | 400    | The body is not valid.                 |
  | 403    | The password is not correct.           |
  | 404    | No user matches the given email.       |
  | 422    | Body data does not match requirements. |

  | Response key |  Type  | Description                                                     |
  |:-------------|:------:|:----------------------------------------------------------------|
  | `token`      | string | A brand new authorization token to connect to the private apis. |

### `PUT /auth/register`

Registers the user in the platform. An email is sent to the provided address, with a validation code.

- **Body**:

  | Key         |   Type    | Required | Description                                          | Condition                                                                                               |
  |:------------|:---------:|:--------:|:-----------------------------------------------------|:--------------------------------------------------------------------------------------------------------|
  | `email`     |  string   |    ✔     | Email of the user.                                   | - `[3-128]` characters long.<br/>- ``/^[a-zA-Z\d!#$%&'*+/=?^_`{\|}.~-]+@[a-z\d]{2,}(.[a-z\d]{2,})+$/``  |
  | `password`  |  string   |    ✔     | Password of the user.                                | - `[2-256]` characters long.                                                                            |
  | `firstName` |  string   |    ✔     | First name of the user.                              | - `[1-32]` characters long.<br/>-``/^\p{L}+([- ']\p{L}+)*$/``                                           |
  | `lastName`  |  string   |    ✔     | Last name of the user.                               | - `[1-32]` characters long.<br/>-``/^\p{L}+([- ']\p{L}+)*$/``                                           |
  | `slug`      |  string   |    ✔     | Unique URL identifier for the user profile.          | - `[1-64]` characters long.<br/>-``/^[a-z\d]+(-[a-z\d]+)*$/``                                           |
  | `username`  |  string   |          | Optional display name, to hide first and last names. | - `[1-64]` characters long.<br/>-``/^[\p{L}\p{N}\p{S}\p{P}\p{M}]+( ([\p{L}\p{N}\p{S}\p{P}\p{M}]+))*$/`` |
  | `sex`       |  string   |    ✔     | Biological sex of the user.                          | - `"male"` or `"female"`                                                                                |
  | `birthday`  | timestamp |    ✔     | Birthdate of the user.                               | - The resulting age must be between `[16-150]` years old.                                               |

- **Response**:

  | Status | Reason                                              |
  |:-------|:----------------------------------------------------|
  | 201    | ✅︎ Success.                                         |
  | 400    | The body is not valid.                              |
  | 409    | Either the email or slug are taken by another user. |
  | 422    | Body data does not match requirements.              |

  | Response key |  Type  | Description                                                     |
  |:-------------|:------:|:----------------------------------------------------------------|
  | `token`      | string | A brand new authorization token to connect to the private apis. |

### `DELETE /email`

Cancel the pending email update, if any.

- **Headers**:

  | Key              |  Type   | Required | Description   |
  |:-----------------|:-------:|:--------:|:--------------|
  | `Authorization`  | string  |    ✔     | Bearer token. |

- **Response**:

  | Status | Reason                |
  |:-------|:----------------------|
  | 204    | ✅︎ Success.           |
  | 403    | The token is invalid. |

### `DELETE /password`

Sends a link to the user email to update its password, without requiring the current one. It does not prevent login
with the current password while the link has not been opened.

- **Query**:

  | Key     |  Type   | Required | Description                  |
  |:--------|:-------:|:--------:|:-----------------------------|
  | `email` | string  |    ✔     | Email of the target account. |

- **Response**:

  | Status | Reason                                                           |
  |:-------|:-----------------------------------------------------------------|
  | 202    | ✅︎ The update link has been created. It is scheduled to be sent. |
  | 400    | The email value is not valid.                                    |
  | 404    | There is no user associated with the provided email.             |

### `PATCH /email`

Create a new email update pending request. Request remains in pending state until the user clicks the validation link.

- **Headers**:

  | Key              |  Type   | Required | Description   |
  |:-----------------|:-------:|:--------:|:--------------|
  | `Authorization`  | string  |    ✔     | Bearer token. |

- **Body**:

  | Key        |  Type  | Required | Description                                                                          | 
  |:-----------|:------:|:--------:|:-------------------------------------------------------------------------------------|
  | `newEmail` | string |    ✔     | The new value of the email. Current email remains used, until this one is validated. |

- **Response**:

  | Status | Reason                                                                                                                                                                                      |
  |:-------|:--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
  | 202    | ✅︎ Success. The validation link for the new email is scheduled to be sent.                                                                                                                  |
  | 403    | The token is invalid.                                                                                                                                                                       |
  | 409    | The new email is already used by another user, as their main email. The same email can be pending validation for multiple accounts, but validating one will block validation of the others. |
  | 422    | Body data does not match requirements.                                                                                                                                                      |

### `PATCH /identity`

Update the user identity information.

- **Headers**:

  | Key              |  Type   | Required | Description   |
  |:-----------------|:-------:|:--------:|:--------------|
  | `Authorization`  | string  |    ✔     | Bearer token. |

- **Body**:

  | Key         |   Type    | Required | Description                 | Condition                                                     | 
  |:------------|:---------:|:--------:|:----------------------------|:--------------------------------------------------------------|
  | `firstName` |  string   |    ✔     | First name of the user.     | - `[1-32]` characters long.<br/>-``/^\p{L}+([- ']\p{L}+)*$/`` |
  | `lastName`  |  string   |    ✔     | Last name of the user.      | - `[1-32]` characters long.<br/>-``/^\p{L}+([- ']\p{L}+)*$/`` |
  | `sex`       |  string   |    ✔     | Biological sex of the user. | - `"male"` or `"female"`                                      |
  | `birthday`  | timestamp |    ✔     | Birthdate of the user.      | - The resulting age must be between `[16-150]` years old.     |

- **Response**:

  | Status | Reason                                 |
  |:-------|:---------------------------------------|
  | 201    | ✅︎ Success.                            |
  | 403    | The token is invalid.                  |
  | 422    | Body data does not match requirements. |

### `PATCH /profile`

Update the user profile information.

- **Headers**:

  | Key              |  Type   | Required | Description   |
  |:-----------------|:-------:|:--------:|:--------------|
  | `Authorization`  | string  |    ✔     | Bearer token. |

- **Body**:

  | Key        |   Type    | Required | Description                                          | Condition                                                                                               |
  |:-----------|:---------:|:--------:|:-----------------------------------------------------|:--------------------------------------------------------------------------------------------------------|
  | `slug`     |  string   |    ✔     | Unique URL identifier for the user profile.          | - `[1-64]` characters long.<br/>-``/^[a-z\d]+(-[a-z\d]+)*$/``                                           |
  | `username` |  string   |          | Optional display name, to hide first and last names. | - `[1-64]` characters long.<br/>-``/^[\p{L}\p{N}\p{S}\p{P}\p{M}]+( ([\p{L}\p{N}\p{S}\p{P}\p{M}]+))*$/`` |

- **Response**:

  | Status | Reason                                     |
  |:-------|:-------------------------------------------|
  | 201    | ✅︎ Success.                                |
  | 403    | The token is invalid.                      |
  | 409    | The slug is already taken by another user. |
  | 422    | Body data does not match requirements.     |

### `PATCH /password`

Update the user password. 

Password update does not require authentication because it can be done from an email link.
Requires either the current password, or an active and valid code sent by email. Both values are considered secure
(password is securely stored as a hash, and the email link is one-use only).

User ID can be retrieved either from the private user API (when logged in), or from the validation link sent via
email.

- **Body**:

  | Key           |  Type  | Required | Description                                                                                                             | Condition                    |
  |:--------------|:------:|:--------:|:------------------------------------------------------------------------------------------------------------------------|:-----------------------------|
  | `id`          |  uuid  |    ✔     | Unique URL identifier for the user.                                                                                     |                              |
  | `code`        | string |          | If using a link to update the password, this value is available from the link query parameters. It is usable only once. |                              |
  | `oldPassword` | string |          | The current password of the user. Required, unless a valid code is provided.                                            |                              |
  | `newPassword` | string |    ✔     | The new password to use.                                                                                                | - `[2-256]` characters long. |

- **Response**:

  | Status | Reason                                                                                                         |
  |:-------|:---------------------------------------------------------------------------------------------------------------|
  | 201    | ✅︎ Success.                                                                                                    |
  | 403    | The validation is not valid: either the current password is wrong, or the provided code is misspelled/expired. |
  | 422    | Body data does not match requirements.                                                                         |

### `PATCH /email/validation`

Resend a validation link to the main email, if available.

- **Headers**:

  | Key              |  Type   | Required | Description   |
  |:-----------------|:-------:|:--------:|:--------------|
  | `Authorization`  | string  |    ✔     | Bearer token. |

- **Response**:

  | Status | Reason                                        |
  |:-------|:----------------------------------------------|
  | 202    | ✅︎ Success. New link is scheduled to be sent. |
  | 403    | The token is invalid.                         |
  | 404    | The email has already been validated.         |

### `PATCH /email/pending/validation`

Resend a validation link to the pending email, if available.

- **Headers**:

  | Key              |  Type   | Required | Description   |
  |:-----------------|:-------:|:--------:|:--------------|
  | `Authorization`  | string  |    ✔     | Bearer token. |

- **Response**:

  | Status | Reason                                        |
  |:-------|:----------------------------------------------|
  | 202    | ✅︎ Success. New link is scheduled to be sent. |
  | 403    | The token is invalid.                         |
  | 404    | The email has already been validated.         |

### `GET /email/validation`

Validate the main email.

- **Query**:

  | Key    |  Type  | Required | Description                                 |
  |:-------|:------:|:--------:|:--------------------------------------------|
  | `id`   |  uuid  |    ✔     | The ID of the targeted user.                |
  | `code` | string |    ✔     | The validation code sent to the user email. |

- **Response**:

  | Status | Reason                                                                                           | 
  |:-------|:-------------------------------------------------------------------------------------------------|
  | 204    | ✅︎ Success.                                                                                      |
  | 403    | The code is invalid, or the target user does not exist, or the email has already been validated. |

### `GET /email/pending/validation`

Validate the pending email.

- **Query**:

  | Key    |  Type  | Required | Description                                 |
  |:-------|:------:|:--------:|:--------------------------------------------|
  | `id`   |  uuid  |    ✔     | The ID of the targeted user.                |
  | `code` | string |    ✔     | The validation code sent to the user email. |

- **Response**:

  | Status | Reason                                                                                           | 
  |:-------|:-------------------------------------------------------------------------------------------------|
  | 204    | ✅︎ Success.                                                                                      |
  | 403    | The code is invalid, or the target user does not exist, or the email has already been validated. |
  | 409    | The email address has been taken by someone else.                                                |

### `GET /email/exists`

Check whether the email value is taken by a user.

- **Query**:

  | Key     |  Type   | Required | Description                   |
  |:--------|:-------:|:--------:|:------------------------------|
  | `email` | string  |    ✔     | Value of the email to verify. |

- **Response**:

  | Status | Reason                 |
  |:-------|:-----------------------|
  | 204    | ✅︎ The email is taken. |
  | 404    | ✅︎ The email is free.  |

### `GET /slug/exists`

Check whether the slug value is taken by a user.

- **Query**:

  | Key    |  Type   | Required | Description                   | 
  |:-------|:-------:|:--------:|:------------------------------|
  | `slug` | string  |    ✔     | Value of the slug to verify.  |

- **Response**:

  | Status | Reason                |
  |:-------|:----------------------|
  | 204    | ✅︎ The slug is taken. |
  | 404    | ✅︎ The slug is free.  |

### `GET /users`

Return basic information for every user ID given in the query. Supports empty queries.

- **Query**:

  | Key   |   Type   | Required | Description          |
  |:------|:--------:|:--------:|:---------------------|
  | `ids` | []string |    ✔     | List of users uuids. |

- **Response**:

  | Status | Reason                                 |
  |:-------|:---------------------------------------|
  | 200    | ✅︎ Success.                            |
  | 400    | The value of any user id is not valid. |

  | Response key        |   Type    | Description                                                                      |
  |:--------------------|:---------:|:---------------------------------------------------------------------------------|
  | `users`             | []object  | List of users.                                                                   |
  | `users[].firstName` |  string?  | The first name of the user. This value is empty if a `username` is available.    |
  | `users[].lastName`  |  string?  | The last name of the user. This value is empty if a `username` is available.     |
  | `users[].username`  |  string?  | The optional username of the user. Hides `firstName` and `lastName` if provided. |
  | `users[].slug`      |  string   | The public URL identifier of the user.                                           |
  | `users[].createdAt` | timestamp | The date at which the user joined the platform.                                  |

### `GET /users/search`

Return basic information for every user matching the search. Results are sorted by relevancy.

- **Query**:

  | Key      |  Type  | Required | Description                                                                                                |
  |:---------|:------:|:--------:|:-----------------------------------------------------------------------------------------------------------|
  | `query`  | string |          | A search value. Matches both slug and display name (username if set, first/last name otherwise).           |
  | `limit`  | number |    ✔     | The maximum number of search results to return. This value is required, and must be in the range [1-1000]. |
  | `offset` | number |          | Offset of the returned results, to skip already fetched ones.                                              |

- **Response**:

  | Status | Reason                                                                           |
  |:-------|:---------------------------------------------------------------------------------|
  | 200    | ✅︎ Success.                                                                      |
  | 400    | The query is not correctly formatted, or the `limit` value is incorrect/missing. |

  | Response key      |   Type    | Description                                                                      |
  |:------------------|:---------:|:---------------------------------------------------------------------------------|
  | `res`             | []object  | List of users.                                                                   |
  | `res[].firstName` |  string?  | The first name of the user. This value is empty if a `username` is available.    |
  | `res[].lastName`  |  string?  | The last name of the user. This value is empty if a `username` is available.     |
  | `res[].username`  |  string?  | The optional username of the user. Hides `firstName` and `lastName` if provided. |
  | `res[].slug`      |  string   | The public URL identifier of the user.                                           |
  | `res[].createdAt` | timestamp | The date at which the user joined the platform.                                  |
  | `total`           |  number   | The total number of results matching the given query.                            |

### `GET /user`

Return basic information about a given user.

- **Query**:

  | Key    |  Type   | Required | Description       |
  |:-------|:-------:|:--------:|:------------------|
  | `slug` | string  |    ✔     | Slug of the user. |

- **Response**:

  | Status | Reason                            |
  |:-------|:----------------------------------|
  | 200    | ✅︎ Success.                       |
  | 404    | No user found for the given slug. |

  | Response key |   Type    | Description                                                                      |
  |:-------------|:---------:|:---------------------------------------------------------------------------------|
  | `firstName`  |  string?  | The first name of the user. This value is empty if a `username` is available.    |
  | `lastName`   |  string?  | The last name of the user. This value is empty if a `username` is available.     |
  | `username`   |  string?  | The optional username of the user. Hides `firstName` and `lastName` if provided. |
  | `slug`       |  string   | The public URL identifier of the user.                                           |
  | `createdAt`  | timestamp | The date at which the user joined the platform.                                  |

### `GET /user/me`

Return basic information about the current user.

- **Headers**:

  | Key              |  Type   | Required | Description   |
  |:-----------------|:-------:|:--------:|:--------------|
  | `Authorization`  | string  |    ✔     | Bearer token. |

- **Response**:

  | Status | Reason                |
  |:-------|:----------------------|
  | 200    | ✅︎ Success.           |
  | 403    | The token is invalid. |

  | Response key |   Type    | Description                                                                      |
  |:-------------|:---------:|:---------------------------------------------------------------------------------|
  | `id`         |   uuid    | The unique identifier for the current user.                                      |
  | `email`      |  string   | The main email of the current user.                                              |
  | `newEmail`   |  string?  | The value of the new email pending validation, if any.                           |
  | `validated`  |  boolean  | Indicates if the main email of the current user has been validated.              |
  | `firstName`  |  string?  | The first name of the user. This value is empty if a `username` is available.    |
  | `lastName`   |  string?  | The last name of the user. This value is empty if a `username` is available.     |
  | `username`   |  string?  | The optional username of the user. Hides `firstName` and `lastName` if provided. |
  | `slug`       |  string   | The public URL identifier of the user.                                           |
  | `createdAt`  | timestamp | The date at which the user joined the platform.                                  |

## Installation

Set the database up.

```bash
make db-setup
```

## Commands

Run the API:

```bash
make run
```

Run the internal API (used by google cloud internal services):

```bash
make run-internal
```

Trigger keys rotation (required on first run, api must be up):

```bash
make rotate-keys
```

Run tests:

```bash
make test
```

Connect to the database:

```bash
make db
```

Connect to the test database:

```bash
make db-test
```
