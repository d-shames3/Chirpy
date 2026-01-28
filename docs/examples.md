# API Examples and Usage Patterns

This document provides practical examples and common usage patterns for the Chirpy API.

## Table of Contents

- [Complete User Workflow](#complete-user-workflow)
- [Common API Patterns](#common-api-patterns)
- [Error Handling Examples](#error-handling-examples)
- [Advanced Scenarios](#advanced-scenarios)
- [Testing with curl](#testing-with-curl)

## Complete User Workflow

### 1. User Registration
```bash
# Register a new user
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{
    "email": "alice@example.com",
    "password": "securePassword123"
  }'

# Response
{
  "id": "123e4567-e89b-12d3-a456-426614174000",
  "created_at": "2023-01-01T12:00:00Z",
  "updated_at": "2023-01-01T12:00:00Z",
  "email": "alice@example.com",
  "is_chirpy_red": false,
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "dGhpcy1pcy1hLXJlZnJlc2gtdG9rZW4="
}
```

### 2. Create First Chirp
```bash
# Use the access token from registration
curl -X POST http://localhost:8080/api/chirps \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -d '{"body":"Hello Chirpy world! #firstpost"}'

# Response
{
  "id": "456e7890-e89b-12d3-a456-426614174111",
  "created_at": "2023-01-01T12:05:00Z",
  "updated_at": "2023-01-01T12:05:00Z",
  "body": "Hello Chirpy world! #firstpost",
  "user_id": "123e4567-e89b-12d3-a456-426614174000"
}
```

### 3. View Timeline
```bash
# Get all chirps (timeline)
curl http://localhost:8080/api/chirps

# Response
[
  {
    "id": "456e7890-e89b-12d3-a456-426614174111",
    "created_at": "2023-01-01T12:05:00Z",
    "updated_at": "2023-01-01T12:05:00Z",
    "body": "Hello Chirpy world! #firstpost",
    "user_id": "123e4567-e89b-12d3-a456-426614174000"
  }
]
```

### 4. Token Refresh (After 1 Hour)
```bash
# When access token expires, use refresh token
curl -X POST http://localhost:8080/api/refresh \
  -H "Authorization: Bearer dGhpcy1pcy1hLXJlZnJlc2gtdG9rZW4="

# Response
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...new_token"
}
```

### 5. Logout
```bash
# Revoke the refresh token
curl -X POST http://localhost:8080/api/revoke \
  -H "Authorization: Bearer dGhpcy1pcy1hLXJlZnJlc2gtdG9rZW4="
```

---

## Common API Patterns

### Authentication Pattern
```bash
# Store tokens after login/register
ACCESS_TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
REFRESH_TOKEN="dGhpcy1pcy1hLXJlZnJlc2gtdG9rZW4="

# Use access token for authenticated requests
curl -X POST http://localhost:8080/api/chirps \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -d '{"body":"Authenticated chirp"}'

# Refresh when needed
NEW_TOKEN=$(curl -s -X POST http://localhost:8080/api/refresh \
  -H "Authorization: Bearer $REFRESH_TOKEN" | \
  jq -r '.token')

# Update access token
ACCESS_TOKEN="$NEW_TOKEN"
```

### Error Handling Pattern
```bash
# Function to handle API responses
api_call() {
  local response=$(curl -s -w "%{http_code}" "$@")
  local http_code="${response: -3}"
  local body="${response%???}"
  
  case $http_code in
    200|201|204)
      echo "$body"
      ;;
    400)
      echo "Bad Request: $(echo "$body" | jq -r '.error')" >&2
      return 1
      ;;
    401)
      echo "Unauthorized: Invalid or missing token" >&2
      return 1
      ;;
    404)
      echo "Not Found: Resource does not exist" >&2
      return 1
      ;;
    500)
      echo "Server Error: $(echo "$body" | jq -r '.error')" >&2
      return 1
      ;;
    *)
      echo "Unexpected error: HTTP $http_code" >&2
      return 1
      ;;
  esac
}
```

### Pagination Pattern (Future Enhancement)
```bash
# Currently not implemented, but here's the pattern:
get_chirps() {
  local limit=20
  local offset=0
  
  while true; do
    local chirps=$(curl -s "http://localhost:8080/api/chirps?limit=$limit&offset=$offset")
    
    if [ "$(echo "$chirps" | jq 'length')" -eq 0 ]; then
      break
    fi
    
    echo "$chirps" | jq -c '.[]'
    offset=$((offset + limit))
  done
}
```

---

## Error Handling Examples

### Invalid Authentication
```bash
# Request without token
curl -X POST http://localhost:8080/api/chirps \
  -H "Content-Type: application/json" \
  -d '{"body":"This will fail"}'

# Response
{
  "error": "No Authorization header"
}

# Request with invalid token
curl -X POST http://localhost:8080/api/chirps \
  -H "Authorization: Bearer invalid_token" \
  -d '{"body":"This will also fail"}'

# Response
{
  "error": "Invalid token"
}
```

### Validation Errors
```bash
# Chirp too long
curl -X POST http://localhost:8080/api/chirps \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -d '{"body":"This chirp is way too long and exceeds the 140 character limit that is enforced by the API"}'

# Response
{
  "error": "Chirp is too long"
}

# Empty chirp
curl -X POST http://localhost:8080/api/chirps \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -d '{"body":""}'

# Response
{
  "error": "Chirp body cannot be empty"
}
```

### Content Filtering
```bash
# Chirp with profanity
curl -X POST http://localhost:8080/api/chirps \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
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

## Advanced Scenarios

### Multi-User Session Management
```bash
# Simulate multiple users
USER1_TOKEN="user1_access_token"
USER1_REFRESH="user1_refresh_token"
USER2_TOKEN="user2_access_token"
USER2_REFRESH="user2_refresh_token"

# User 1 creates chirp
curl -X POST http://localhost:8080/api/chirps \
  -H "Authorization: Bearer $USER1_TOKEN" \
  -d '{"body":"User 1 chirp"}'

# User 2 creates chirp
curl -X POST http://localhost:8080/api/chirps \
  -H "Authorization: Bearer $USER2_TOKEN" \
  -d '{"body":"User 2 chirp"}'

# Get timeline (shows both chirps)
curl http://localhost:8080/api/chirps
```

### User Profile Management
```bash
# Update email
curl -X PUT http://localhost:8080/api/users \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -d '{"email":"newemail@example.com"}'

# Update password
curl -X PUT http://localhost:8080/api/users \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -d '{"password":"newSecurePassword456"}'

# Update both
curl -X PUT http://localhost:8080/api/users \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -d '{
    "email":"newemail@example.com",
    "password":"newSecurePassword456"
  }'
```

### Chirp Management
```bash
# Create chirp
CHIRP_ID=$(curl -s -X POST http://localhost:8080/api/chirps \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -d '{"body":"My test chirp"}' | jq -r '.id')

# Get specific chirp
curl "http://localhost:8080/api/chirps/$CHIRP_ID"

# Delete chirp
curl -X DELETE "http://localhost:8080/api/chirps/$CHIRP_ID" \
  -H "Authorization: Bearer $ACCESS_TOKEN"
```

---

## Testing with curl

### Quick Test Script
```bash
#!/bin/bash

# Configuration
BASE_URL="http://localhost:8080"
EMAIL="test$(date +%s)@example.com"
PASSWORD="testPassword123"

echo "=== Chirpy API Test ==="

# 1. Health Check
echo "1. Health Check..."
curl -s "$BASE_URL/api/healthz" && echo " ✓"

# 2. Register User
echo "2. Registering user..."
REGISTER_RESPONSE=$(curl -s -X POST "$BASE_URL/api/users" \
  -H "Content-Type: application/json" \
  -d "{\"email\":\"$EMAIL\",\"password\":\"$PASSWORD\"}")

echo "$REGISTER_RESPONSE" | jq .

# Extract tokens
ACCESS_TOKEN=$(echo "$REGISTER_RESPONSE" | jq -r '.token')
REFRESH_TOKEN=$(echo "$REGISTER_RESPONSE" | jq -r '.refresh_token')
USER_ID=$(echo "$REGISTER_RESPONSE" | jq -r '.id')

# 3. Create Chirp
echo "3. Creating chirp..."
CHIRP_RESPONSE=$(curl -s -X POST "$BASE_URL/api/chirps" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -d '{"body":"Hello from test script!"}')

echo "$CHIRP_RESPONSE" | jq .

# 4. Get Timeline
echo "4. Getting timeline..."
curl -s "$BASE_URL/api/chirps" | jq .

# 5. Refresh Token
echo "5. Refreshing token..."
REFRESH_RESPONSE=$(curl -s -X POST "$BASE_URL/api/refresh" \
  -H "Authorization: Bearer $REFRESH_TOKEN")

echo "$REFRESH_RESPONSE" | jq .

# 6. Revoke Token
echo "6. Revoking token..."
curl -s -X POST "$BASE_URL/api/revoke" \
  -H "Authorization: Bearer $REFRESH_TOKEN" && echo " ✓"

echo "=== Test Complete ==="
```

### Load Testing (Simple)
```bash
#!/bin/bash

# Simple load test - create multiple chirps
BASE_URL="http://localhost:8080"
ACCESS_TOKEN="your_access_token"

echo "Starting load test..."

for i in {1..10}; do
  curl -s -X POST "$BASE_URL/api/chirps" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $ACCESS_TOKEN" \
    -d "{\"body\":\"Load test chirp #$i\"}" > /dev/null
  
  echo "Created chirp #$i"
  sleep 0.1
done

echo "Load test complete!"
```

---

## Integration Examples

### JavaScript/Fetch
```javascript
class ChirpyAPI {
  constructor(baseURL = 'http://localhost:8080') {
    this.baseURL = baseURL;
    this.accessToken = null;
    this.refreshToken = null;
  }

  async login(email, password) {
    const response = await fetch(`${this.baseURL}/api/login`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ email, password })
    });
    
    if (response.ok) {
      const data = await response.json();
      this.accessToken = data.token;
      this.refreshToken = data.refresh_token;
      return data;
    }
    throw new Error('Login failed');
  }

  async createChirp(body) {
    const response = await fetch(`${this.baseURL}/api/chirps`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${this.accessToken}`
      },
      body: JSON.stringify({ body })
    });
    
    if (response.ok) {
      return await response.json();
    }
    throw new Error('Failed to create chirp');
  }

  async getChirps() {
    const response = await fetch(`${this.baseURL}/api/chirps`);
    return await response.json();
  }
}

// Usage
const api = new ChirpyAPI();
await api.login('user@example.com', 'password123');
const chirp = await api.createChirp('Hello from JavaScript!');
console.log(chirp);
```

### Python/Requests
```python
import requests
import json

class ChirpyAPI:
    def __init__(self, base_url="http://localhost:8080"):
        self.base_url = base_url
        self.access_token = None
        self.refresh_token = None
    
    def login(self, email, password):
        response = requests.post(
            f"{self.base_url}/api/login",
            json={"email": email, "password": password}
        )
        
        if response.status_code == 200:
            data = response.json()
            self.access_token = data["token"]
            self.refresh_token = data["refresh_token"]
            return data
        else:
            raise Exception("Login failed")
    
    def create_chirp(self, body):
        response = requests.post(
            f"{self.base_url}/api/chirps",
            json={"body": body},
            headers={"Authorization": f"Bearer {self.access_token}"}
        )
        
        if response.status_code == 201:
            return response.json()
        else:
            raise Exception("Failed to create chirp")
    
    def get_chirps(self):
        response = requests.get(f"{self.base_url}/api/chirps")
        return response.json()

# Usage
api = ChirpyAPI()
api.login("user@example.com", "password123")
chirp = api.create_chirp("Hello from Python!")
print(chirp)
```

---

## Best Practices

1. **Token Management**: Store tokens securely, refresh before expiration
2. **Error Handling**: Always check HTTP status codes and handle errors
3. **Rate Limiting**: Be reasonable with request frequency
4. **Content Validation**: Validate input on client side before sending
5. **Security**: Use HTTPS in production, never expose tokens
6. **Testing**: Test error scenarios, not just success cases