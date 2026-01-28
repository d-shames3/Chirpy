# Chirps Endpoints

This document covers all chirp-related endpoints in the Chirpy API. Chirps are the core content type - similar to tweets or posts.

## Table of Contents

- [Create Chirp](#create-chirp)
- [List All Chirps](#list-all-chirps)
- [Get Specific Chirp](#get-specific-chirp)
- [Delete Chirp](#delete-chirp)

## Create Chirp

Create a new chirp (post) with text content.

**Endpoint:** `POST /api/chirps`

**Authentication:** Required (Bearer token)

**Request Body:**
```json
{
  "body": "This is my first chirp! #excited"
}
```

**Response (201 Created):**
```json
{
  "id": "123e4567-e89b-12d3-a456-426614174000",
  "created_at": "2023-01-01T12:00:00Z",
  "updated_at": "2023-01-01T12:00:00Z",
  "body": "This is my first chirp! #excited",
  "user_id": "456e7890-e89b-12d3-a456-426614174111"
}
```

**Error Responses:**
- `400 Bad Request` - Missing body, exceeds length limit, or invalid format
- `401 Unauthorized` - Missing or invalid authentication token
- `500 Internal Server Error` - Database error

**Validation Rules:**
- **Max Length:** 140 characters
- **Content Filtering:** Profanity is replaced with `****`
- **Required Fields:** `body` must not be empty

**Content Filtering:**
The system automatically filters inappropriate language by replacing profanity with `****`. For example:
- Input: "This is damn awesome"
- Output: "This is **** awesome"

---

## List All Chirps

Retrieve all chirps from the system, ordered by creation date (newest first).

**Endpoint:** `GET /api/chirps`

**Authentication:** Not required

**Query Parameters:**
- `author_id` (optional) - Filter by specific user ID

**Response (200 OK):**
```json
[
  {
    "id": "123e4567-e89b-12d3-a456-426614174000",
    "created_at": "2023-01-01T12:00:00Z",
    "updated_at": "2023-01-01T12:00:00Z",
    "body": "Latest chirp content",
    "user_id": "456e7890-e89b-12d3-a456-426614174111"
  },
  {
    "id": "789e0123-e89b-12d3-a456-426614174222",
    "created_at": "2023-01-01T11:30:00Z",
    "updated_at": "2023-01-01T11:30:00Z",
    "body": "Earlier chirp content",
    "user_id": "456e7890-e89b-12d3-a456-426614174111"
  }
]
```

**Error Responses:**
- `500 Internal Server Error` - Database error

**Notes:**
- Returns empty array `[]` if no chirps exist
- Results are sorted by creation date (newest first)
- Use `author_id` parameter to filter by specific user

**Example with author filter:**
```
GET /api/chirps?author_id=456e7890-e89b-12d3-a456-426614174111
```

---

## Get Specific Chirp

Retrieve a single chirp by its ID.

**Endpoint:** `GET /api/chirps/{chirpId}`

**Authentication:** Not required

**Path Parameters:**
- `chirpId` - UUID of the chirp to retrieve

**Response (200 OK):**
```json
{
  "id": "123e4567-e89b-12d3-a456-426614174000",
  "created_at": "2023-01-01T12:00:00Z",
  "updated_at": "2023-01-01T12:00:00Z",
  "body": "This is a specific chirp",
  "user_id": "456e7890-e89b-12d3-a456-426614174111"
}
```

**Error Responses:**
- `400 Bad Request` - Invalid chirp ID format
- `404 Not Found` - Chirp does not exist
- `500 Internal Server Error` - Database error

**Notes:**
- Uses path-based routing (Gorilla Mux style)
- Chirp ID must be a valid UUID format

---

## Delete Chirp

Delete a chirp (only by the original author).

**Endpoint:** `DELETE /api/chirps/{chirpId}`

**Authentication:** Required (Bearer token)

**Path Parameters:**
- `chirpId` - UUID of the chirp to delete

**Response (204 No Content):**
```
(no body)
```

**Error Responses:**
- `400 Bad Request` - Invalid chirp ID format
- `401 Unauthorized` - Missing or invalid authentication token
- `403 Forbidden` - User is not the author of the chirp
- `404 Not Found` - Chirp does not exist
- `500 Internal Server Error` - Database error

**Authorization Rules:**
- Only the original author can delete their chirp
- Admin users cannot delete other users' chirps
- Deleted chirps are permanently removed from the database

---

## Data Schema

### Chirp Object
```json
{
  "id": "uuid",
  "created_at": "datetime",
  "updated_at": "datetime",
  "body": "string",
  "user_id": "uuid"
}
```

### Field Descriptions
- `id` - Unique identifier (UUID v4)
- `created_at` - Timestamp when chirp was created (ISO 8601)
- `updated_at` - Timestamp when chirp was last modified (ISO 8601)
- `body` - The text content (max 140 chars, filtered)
- `user_id` - ID of the user who created the chirp

---

## Usage Examples

### Creating and Managing Chirps

```bash
# Create a chirp
curl -X POST http://localhost:8080/api/chirps \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <access_token>" \
  -d '{"body":"Hello world! #firstpost"}'

# Get all chirps
curl http://localhost:8080/api/chirps

# Get specific chirp
curl http://localhost:8080/api/chirps/123e4567-e89b-12d3-a456-426614174000

# Get chirps by specific author
curl "http://localhost:8080/api/chirps?author_id=456e7890-e89b-12d3-a456-426614174111"

# Delete your own chirp
curl -X DELETE http://localhost:8080/api/chirps/123e4567-e89b-12d3-a456-426614174000 \
  -H "Authorization: Bearer <access_token>"
```

### Content Filtering Example

```bash
# Request with profanity
curl -X POST http://localhost:8080/api/chirps \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <access_token>" \
  -d '{"body":"This is damn good!"}'

# Response (filtered)
{
  "id": "uuid",
  "created_at": "2023-01-01T12:00:00Z",
  "updated_at": "2023-01-01T12:00:00Z",
  "body": "This is **** good!",
  "user_id": "uuid"
}
```

---

## Best Practices

1. **Character Limits**: Keep chirps under 140 characters
2. **Content**: Be mindful of content filtering
3. **Authentication**: Always use valid tokens for create/delete operations
4. **Error Handling**: Check response codes and handle appropriately
5. **Rate**: Be reasonable with API request frequency