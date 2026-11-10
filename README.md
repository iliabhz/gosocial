# Go Social 
A simple social netowrking app (mostly the backend for now) which is Built with Go + PostgresQL + Redis + Mailgun + JWT (For Auth)

It's more like reddit where users can signup and then need to activate their account by clicking the link they get from a email sent to them. After that users can post somethings with a title, content, and tags which you can use to search posts from tags. Also users can comment on such posts and follow eachother to see their followers posts.

Each user can delete and modify their own comments and posts and moderators can only modify other user comments and posts and the admins can both delete and modify other user comments and posts.

It uses a Repostiory pattern for handling storge which is in internal/store and has a postgres impl in internal/store/postgres. Also has a cache system which is implemented in internal/store/cache which has a both Redis impl and a InMemory impl.

File structure:
```
├── bin
├── cmd
│   ├── api
│   └── seed
├── db
│   └── migrations
├── docs
├── internal
│   ├── env
│   ├── mailer
│   │   └── templates
│   ├── ratelimiter
│   └── store
│       ├── cache
│       └── postgres
├── scripts
└── web
```

# GoSocial API


API for GoSocial which is a social network built in go

## `POST` /authentication/token

Creates a token

Creates a token for a user

**Parameters:**

| Name | In | Type | Required | Description |
|------|-----|------|----------|-------------|
| `payload` | body |  | ✓ | User credentials |

**Responses:**

- `200` — Token
- `400` — Bad Request
- `401` — Unauthorized
- `500` — Internal Server Error

## `POST` /authentication/user

Registers a user

**Parameters:**

| Name | In | Type | Required | Description |
|------|-----|------|----------|-------------|
| `payload` | body |  | ✓ | User credentials |

**Responses:**

- `201` — User registered
- `400` — Bad Request
- `500` — Internal Server Error

## `GET` /health

Healthcheck

Healthcheck endpoint

**Responses:**

- `200` — ok

## `POST` /posts

Creates a post

**Parameters:**

| Name | In | Type | Required | Description |
|------|-----|------|----------|-------------|
| `payload` | body |  | ✓ | Post payload |

**Responses:**

- `201` — Created
- `400` — Bad Request
- `401` — Unauthorized
- `500` — Internal Server Error

## `GET` /posts/{id}

Fetches a post

Fetches a post by ID

**Parameters:**

| Name | In | Type | Required | Description |
|------|-----|------|----------|-------------|
| `id` | path | integer | ✓ | Post ID |

**Responses:**

- `200` — OK
- `404` — Not Found
- `500` — Internal Server Error

## `DELETE` /posts/{id}

Deletes a post

Delete a post by ID

**Parameters:**

| Name | In | Type | Required | Description |
|------|-----|------|----------|-------------|
| `id` | path | integer | ✓ | Post ID |

**Responses:**

- `204` — No Content
- `404` — Not Found
- `500` — Internal Server Error

## `PATCH` /posts/{id}

Updates a post

Updates a post by ID

**Parameters:**

| Name | In | Type | Required | Description |
|------|-----|------|----------|-------------|
| `id` | path | integer | ✓ | Post ID |
| `payload` | body |  | ✓ | Post payload |

**Responses:**

- `200` — OK
- `400` — Bad Request
- `401` — Unauthorized
- `404` — Not Found
- `500` — Internal Server Error

## `PUT` /users/activate

Activates/Register a user

Activates/Register a user by invitation token

**Parameters:**

| Name | In | Type | Required | Description |
|------|-----|------|----------|-------------|
| `token` | body |  | ✓ | Invitation token |

**Responses:**

- `204` — User activated
- `404` — Not Found
- `500` — Internal Server Error

## `GET` /users/feed

Fetches the user feed

**Parameters:**

| Name | In | Type | Required | Description |
|------|-----|------|----------|-------------|
| `since` | query | string |  | Since |
| `until` | query | string |  | Until |
| `limit` | query | integer |  | Limit |
| `offset` | query | integer |  | Offset |
| `sort` | query | string |  | Sort |
| `tags` | query | string |  | Tags |
| `search` | query | string |  | Search |

**Responses:**

- `200` — OK
- `400` — Bad Request
- `500` — Internal Server Error

## `GET` /users/{id}

Fetches a user profile

Fetches a user profile by ID

**Parameters:**

| Name | In | Type | Required | Description |
|------|-----|------|----------|-------------|
| `id` | path | integer | ✓ | User ID |

**Responses:**

- `200` — OK
- `400` — Bad Request
- `404` — Not Found
- `500` — Internal Server Error

## `PUT` /users/{userID}/follow

Follows a user

Follows a user by ID

**Parameters:**

| Name | In | Type | Required | Description |
|------|-----|------|----------|-------------|
| `userID` | path | integer | ✓ | User ID |

**Responses:**

- `204` — User followed
- `400` — User payload missing
- `404` — User not found
- `442` — Couldn't process the user

## `PUT` /users/{userID}/unfollow

Unfollow a user

Unfollow a user by ID

**Parameters:**

| Name | In | Type | Required | Description |
|------|-----|------|----------|-------------|
| `userID` | path | integer | ✓ | User ID |

**Responses:**

- `204` — User unfollowed
- `400` — User payload missing
- `404` — User not found
