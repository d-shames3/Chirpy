# User Management Endpoints

This document covers user management endpoints in the Chirpy API. Note that user registration is covered in the [Authentication](./authentication.md) documentation.

## Table of Contents

- [Update User Credentials](#update-user-credentials)
- [User Data Schema](#user-data-schema)

## Update User Credentials

Update the authenticated user's email and/or password.

**Endpoint:** `PUT /api/users`

**Authentication:** Required (Bearer token)

**Request Body:**
```json
{
  "email": "newemail@example.com",
  "password": "newSecurePassword456"
}
```

**Response (200 OK):**
```json
{
  "id": "123e4567-e89b-12d3-a456-426614174000",
  "created_at": "2023-01-01T12:00:00Z",
  "updated_at": "2023-01-01T12:30:00Z",
  "email": "newemail@example.com",
  "is_chirpy_red": false
}
```

**Error Responses:**
- `400 Bad Request` - Invalid email format, weak password, or missing fields
- `401 Unauthorized` - Missing or invalid authentication token
- `500 Internal Server Error` - Database error

**Validation Rules:**

### Email Updates
- Must be a valid email format
- Can be the same as current email
- Must be unique across all users

### Password Updates
- Must be at least 8 characters long
- Can be the same as current password
- Hashed using Argon2ID for security

**Partial Updates:**
You can update just email, just password, or both:

```json
// Update email only
{
  "email": "newemail@example.com"
}

// Update password only
{
  "password": "newSecurePassword456"
}

// Update both
{
  "email": "newemail@example.com",
  "password": "newSecurePassword456"
}
```

---

## User Data Schema

### User Object
```json
{
  "id": "uuid",
  "created_at": "datetime",
  "updated_at": "datetime",
  "email": "string",
  "is_chirpy_red": "boolean"
}
```

### Field Descriptions
- `id` - Unique user identifier (UUID v4)
- `created_at` - Timestamp when account was created (ISO 8601)
- `updated_at` - Timestamp when account was last modified (ISO 8601)
- `email` - User's email address (unique)
- `is_chirpy_red` - Premium status flag (true for premium users)

### Authentication-Only Fields
These fields are only included in authentication responses:

```json
{
  "token": "jwt_access_token",        // Only in auth responses
  "refresh_token": "refresh_token"    // Only in auth responses
}
```

---

## Usage Examples

### Update Email Only
```bash
curl -X PUT http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <access_token>" \
  -d '{"email":"newemail@example.com"}'
```

### Update Password Only
```bash
curl -X PUT http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <access_token>" \
  -d '{"password":"newSecurePassword456"}'
```

### Update Both Email and Password
```bash
curl -X PUT http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <access_token>" \
  -d '{
    "email":"newemail@example.com",
    "password":"newSecurePassword456"
  }'
```

---

## Security Considerations

### Password Security
- Passwords are hashed using Argon2ID algorithm
- Passwords are never stored or returned in plain text
- Minimum 8-character length requirement
- No complexity requirements enforced, but recommended

### Authentication
- Users can only update their own profile
- Valid JWT access token required
- Token must belong to the authenticated user

### Email Privacy
- Email addresses are unique across the platform
- Email addresses are returned in user profiles
- Consider privacy implications when sharing user data

---

## User States

### Regular User
```json
{
  "id": "uuid",
  "created_at": "2023-01-01T12:00:00Z",
  "updated_at": "2023-01-01T12:00:00Z",
  "email": "user@example.com",
  "is_chirpy_red": false
}
```

### Premium User (Chirpy Red)
```json
{
  "id": "uuid",
  "created_at": "2023-01-01T12:00:00Z",
  "updated_at": "2023-01-01T12:00:00Z",
  "email": "user@example.com",
  "is_chirpy_red": true
}
```

**Note:** The `is_chirpy_red` status is typically updated through webhook integration with payment providers, not through direct user updates.

---

## Related Endpoints

- **Register User:** `POST /api/users` - [Authentication Documentation](./authentication.md)
- **Login:** `POST /api/login` - [Authentication Documentation](./authentication.md)
- **Get User Chirps:** `GET /api/chirps?author_id=<user_id>` - [Chirps Documentation](./chirps.md)

---

## Best Practices

1. **Validation:** Always validate email format on the client side before sending
2. **Password Strength:** Encourage strong passwords even though not enforced
3. **Token Management:** Ensure valid access tokens are used
4. **Error Handling:** Handle authentication errors gracefully
5. **Security:** Use HTTPS in production to protect credentials