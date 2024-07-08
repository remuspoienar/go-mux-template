# HTTP Mux Starter

### Simple MUX template project

## Provides starting point for a REST API

- context
- JWT Auth + middleware
- bcrypt helper
- webhooks integration
- godotenv (env sample file provided)
- static file serving
- fs storage (JSON format)

## Auth endpoints

### [POST /api/login](./api/authRoutes.go)

```http request
Content-Type: application/json
```

```json
{
  "email": "mike@example.com",
  "password": "hunter200"
}
```

Returns <span style="color:green">200</span> OR <span style="color:red">401</span>

```json
{
  "id": 1,
  "email": "mike@example.com",
  "is_chirpy_red": false,
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c",
  "refresh_token": "56aa826d22baab4b5ec2cea41a59ecbba03e542aedbb31d9b80326ac8ffcfa2a"
}
```

### [POST /api/refresh](./api/authRoutes.go)

Returns a JWT with a refresh token passed in headers

```http request
Authorization: Bearer b056cec97994e74e6695eba938cee97513eaeb720829559812066272f5529c17
```

Returns <span style="color:green">200</span> OR <span style="color:red">401</span>

```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
}
```

### [POST /api/revoke](./api/authRoutes.go)

Revokes a refresh token

```http request
Authorization: Bearer b056cec97994e74e6695eba938cee97513eaeb720829559812066272f5529c17
```

Returns <span style="color:green">204</span> OR <span style="color:red">401</span>

## User endpoints

### [POST /api/users](./api/userRoutes.go)

Creates a user

```http request
Content-Type: application/json
```

```json
{
  "password": "hunter200",
  "email": "mike@example.com"
}
```

Returns <span style="color:yellow">201</span>

```json
{
  "email": "mike@example.com",
  "id": 6,
  "is_chirpy_red": false
}
```

### [PUT /api/users/{id}](./api/userRoutes.go)

Updates a user

```http request
Content-Type: application/json
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJjaGlycHkiLCJzdWIiOiIxIiwiZXhwIjoxNzIwMjE1NTIxLCJpYXQiOjE3MjAyMTE5MjF9.dpLB975ewERQ5tXrYF7oLYicQyzXl27Rt3HadrgkJzg
```

```json
{
  "email": "new_mike@example.com"
}
```

Returns <span style="color:green">200</span> OR <span style="color:red">401</span>

```json
{
  "email": "new_mike@example.com",
  "id": 6,
  "is_chirpy_red": false
}
```

## Chirps

### [GET /api/chirps](./api/chirpRoutes.go)

Gets all chirps

Optional query parameters

```http request
GET /api/chirps?sort=asc
```

```http request
GET /api/chirps?sort=desc
```

```http request
GET /api/chirps?author_id=3
```

Returns <span style="color:green">200</span> OR <span style="color:red">400</span> (if query params are invalid)

```json
[
  {
    "body": "Darn that fly, I just wanna cook",
    "id": 4,
    "author_id": 1
  },
  {
    "body": "Cmon Pinkman",
    "id": 3,
    "author_id": 2
  },
  {
    "body": "Gale!",
    "id": 2,
    "author_id": 1
  },
  {
    "body": "I'm the one who knocks!",
    "id": 1,
    "author_id": 4
  }
]
```

### [POST /api/chirps](./api/chirpRoutes.go)

Creates a chirp

```http request
Content-Type: application/json
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJjaGlycHkiLCJzdWIiOiIxIiwiZXhwIjoxNzIwMjE1NTIxLCJpYXQiOjE3MjAyMTE5MjF9.dpLB975ewERQ5tXrYF7oLYicQyzXl27Rt3HadrgkJzg
```

```json

{
  "body": "Last time i went to Rustbucks ! #til"
}
```

Returns <span style="color:green">200</span> OR <span style="color:red">400</span>(if length is > 140)

```json
{
  "id": 22,
  "body": "Last time i went to Rustbucks ! #til",
  "author_id": 2
}
```

### [GET /api/chirps/{id}](./api/chirpRoutes.go)

Gets a chirp by id

Returns <span style="color:green">200</span> OR <span style="color:red">404</span>

```json
{
  "id": 4,
  "body": "Darn that fly, I just wanna cook",
  "author_id": 1
}
```

### [DELETE /api/chirps/{id}](./api/chirpRoutes.go)

Delete a chirp by id - Only the author can perform this and will result in a 403 otherwise

Returns <span style="color:green">204</span>, <span style="color:red">403</span> OR <span style="color:red">404</span>

## Payment event webhook
Uses api key auth

### [POST /api/polka/webhooks](./api/webhook.go)

```http request
Authorization: ApiKey <string>
```
```json
{
  "event": "user.upgraded",
  "data": {
    "user_id": 3
  }
}
```

Returns <span style="color:green">204</span> OR <span style="color:red">404</span>
