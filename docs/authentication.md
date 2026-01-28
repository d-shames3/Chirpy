# Authentication Endpoints

This document covers all authentication-related endpoints in the Chirpy API.

## Table of Contents

- [Register User](#register-user)
- [Login](#login)
- [Refresh Token](#refresh-token)
- [Revoke Token](#revoke-token)

## Register User

Create a new user account.

**Endpoint:** `POST /api/users`

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "securePassword123"
}
```

**Response (201 Created):**
```json
{
  "id": "123e4567-e89b-12d3-a456-426614174000",
  "created_at": "2023-01-01T12:00:00Z",
  "updated_at": "2023-01-01T12:00:00Z",
  "email": "user@example.com",
  "is_chirpy_red": false,
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "dGhpcy1pcy1hLXJlZnJlc2gtdG9rZW4="
}
```

**Error Responses:**
- `400 Bad Request` - Invalid email or password format
- `500 Internal Server Error` - Database error

**Notes:**
- Password must be at least 8 characters long
- Email must be a valid email format
- Email addresses must be unique
- Returns both access and refresh tokens on successful registration

---

## Login

Authenticate with existing credentials.

**Endpoint:** `POST /api/login`

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "securePassword123"
}
```

**Response (200 OK):**
```json
{
  "id": "123e4567-e89b-12d3-a456-426614174000",
  "created_at": "2023-01-01T12:00:00Z",
  "updated_at": "2023-01-01T12:00:00Z",
  "email": "user@example.com",
  "is_chirpy_red": false,
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "dGhpcy1pcy1hLXJlZnJlc2gtdG9rZW4="
}
```

**Error Responses:**
- `401 Unauthorized` - Invalid email or password
- `500 Internal Server Error` - Database error

**Notes:**
- Returns fresh tokens on successful login
- Any existing refresh tokens for the user are revoked

---

## Refresh Token

Get a new access token using a valid refresh token.

**Endpoint:** `POST /api/refresh`

**Headers:**
```
Authorization: Bearer <refresh_token>
```

**Response (200 OK):**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Error Responses:**
- `401 Unauthorized` - Invalid or expired refresh token
- `500 Internal Server Error` - Server error

**Notes:**
- Only returns a new access token
- Refresh token is not extended
- Use this when access token expires (default: 1 hour)

---

## Revoke Token

Revoke a refresh token to prevent further use.

**Endpoint:** `POST /api/revoke`

**Headers:**
```
Authorization: Bearer <refresh_token>
```

**Response (204 No Content):**
```
(no body)
```

**Error Responses:**
- `401 Unauthorized` - Invalid or missing refresh token
- `500 Internal Server Error` - Server error

**Notes:**
- Revokes the specific refresh token provided
- Use this when logging out or token is no longer needed
- Does not affect existing access tokens (they will expire naturally)

---

## Token Security

### Access Token
- **Expiration:** 1 hour (3600 seconds)
- **Usage:** Include in `Authorization: Bearer <token>` header
- **Scope:** Full API access for the user

### Refresh Token
- **Expiration:** 60 days (configurable)
- **Usage:** Only for token refresh operations
- **Storage:** Store securely (e.g., httpOnly cookies, secure storage)

### Best Practices
1. Store refresh tokens securely on the client
2. Use HTTPS to prevent token interception
3. Implement token refresh before expiration
4. Revoke tokens on logout
5. Handle token expiration gracefully

## Example Flow

```bash
# 1. Register user
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}'

# 2. Use access token (1 hour expiry)
curl -X GET http://localhost:8080/api/chirps \
  -H "Authorization: Bearer <access_token>"

# 3. Refresh when access token expires
curl -X POST http://localhost:8080/api/refresh \
  -H "Authorization: Bearer <refresh_token>"

# 4. Revoke on logout
curl -X POST http://localhost:8080/api/revoke \
  -H "Authorization: Bearer <refresh_token>"
```