# Admin and Webhook Endpoints

This document covers administrative endpoints and webhook integrations in the Chirpy API.

## Table of Contents

- [Admin Endpoints](#admin-endpoints)
  - [Get Metrics](#get-metrics)
  - [Reset System](#reset-system)
- [Webhook Endpoints](#webhook-endpoints)
  - [Polka Payment Webhook](#polka-payment-webhook)

## Admin Endpoints

### Get Metrics

Retrieve system metrics and statistics for monitoring purposes.

**Endpoint:** `GET /admin/metrics`

**Authentication:** Not required (but should be protected in production)

**Response (200 OK):**
```html
<!DOCTYPE html>
<html>
    <body>
        <h1>Welcome, Chirpy Admin</h1>
        <p>Chirpy has been viewed 1234 times!</p>
    </body>
</html>
```

**Error Responses:**
- `500 Internal Server Error` - Server error

**Notes:**
- Returns HTML content (not JSON)
- Tracks file server hits for the `/app/` route only
- Simple counter that can be reset
- Consider adding authentication in production environments

---

### Reset System

Reset all data in the system (development only).

**Endpoint:** `POST /admin/reset`

**Authentication:** None (environment-protected)

**Response (200 OK):**
```
Reset hits to 0 and deleted all users
```

**Error Responses:**
- `403 Forbidden` - Not in development environment
- `500 Internal Server Error` - Database error

**Environment Requirements:**
- Only works when `PLATFORM=dev` environment variable is set
- Blocked in production environments for safety
- Deletes all users, chirps, and refresh tokens
- Resets metrics counter to zero

**⚠️ Safety Warning:**
This endpoint is extremely destructive and should only be used in development environments. Never expose this endpoint in production.

---

## Webhook Endpoints

### Polka Payment Webhook

Handle payment webhooks from Polka for premium user upgrades.

**Endpoint:** `POST /api/polka/webhooks`

**Authentication:** Required (API Key)

**Headers:**
```
Authorization: ApiKey <polka_api_key>
Content-Type: application/json
```

**Request Body:**
```json
{
  "event": "user.upgraded",
  "data": {
    "user_id": "123e4567-e89b-12d3-a456-426614174000"
  }
}
```

**Response (204 No Content):**
```
(no body)
```

**Error Responses:**
- `401 Unauthorized` - Invalid or missing API key
- `404 Not Found` - User does not exist
- `500 Internal Server Error` - Database error

**Supported Events:**

#### user.upgraded
- **Description**: User has been upgraded to premium (Chirpy Red)
- **Action**: Sets `is_chirpy_red` flag to `true` for the specified user
- **User ID**: Must be a valid UUID that exists in the database

#### Other Events
- **Description**: Any other event types are ignored
- **Response**: Returns 204 No Content without making changes

---

## Webhook Security

### API Key Authentication
- Uses `ApiKey` prefix in Authorization header (not `Bearer`)
- API key must match `POLKA_KEY` environment variable
- Rejects requests with invalid or missing keys

### Event Validation
- Validates JSON structure before processing
- Checks that user_id is a valid UUID
- Verifies user exists in database before updating

### Idempotency
- Webhook processing is idempotent
- Multiple identical requests have same effect as single request
- Setting `is_chirpy_red` to `true` multiple times is safe

---

## Usage Examples

### Get System Metrics
```bash
curl http://localhost:8080/admin/metrics
```

### Reset Development Data
```bash
# Only works in development environment
curl -X POST http://localhost:8080/admin/reset
```

### Process Payment Webhook
```bash
curl -X POST http://localhost:8080/api/polka/webhooks \
  -H "Content-Type: application/json" \
  -H "Authorization: ApiKey your-polka-api-key" \
  -d '{
    "event": "user.upgraded",
    "data": {
      "user_id": "123e4567-e89b-12d3-a456-426614174000"
    }
  }'
```

### Webhook with Invalid User
```bash
# This will return 404 if user doesn't exist
curl -X POST http://localhost:8080/api/polka/webhooks \
  -H "Content-Type: application/json" \
  -H "Authorization: ApiKey your-polka-api-key" \
  -d '{
    "event": "user.upgraded",
    "data": {
      "user_id": "non-existent-user-id"
    }
  }'
```

---

## Implementation Notes

### Metrics Tracking
- Uses atomic counter for thread-safe incrementing
- Only counts requests to `/app/` route (static files)
- Stored in memory (resets on server restart)
- Can be reset via admin endpoint

### Webhook Processing
- Asynchronous processing not required (fast operations)
- Minimal validation for performance
- No retry mechanism implemented
- Consider adding request logging for debugging

### Environment Safety
- Reset endpoint protected by environment variable
- Metrics endpoint should be protected in production
- API key validation prevents unauthorized webhook calls

---

## Production Considerations

### Security Enhancements
1. **Metrics Endpoint**: Add authentication for production
2. **Reset Endpoint**: Remove or add strong authentication
3. **Webhooks**: Add request logging and monitoring
4. **API Keys**: Use key rotation and secure storage

### Monitoring
1. **Metrics**: Implement proper metrics collection
2. **Webhooks**: Track success/failure rates
3. **Performance**: Monitor response times
4. **Errors**: Alert on high error rates

### Reliability
1. **Idempotency**: Ensure webhook reprocessing is safe
2. **Retries**: Consider webhook retry mechanisms
3. **Fallbacks**: Handle payment provider downtime
4. **Data Integrity**: Audit premium status changes