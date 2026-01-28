# Chirpy API Documentation

Welcome to the Chirpy API documentation! Chirpy is a Twitter-like microblogging platform built with Go and PostgreSQL.

## Table of Contents

- [Overview](#overview)
- [Getting Started](#getting-started)
- [Authentication](#authentication)
- [API Endpoints](#api-endpoints)
  - [Users](#users)
  - [Chirps](#chirps)
  - [Authentication](#authentication-1)
  - [Admin](#admin)
  - [Webhooks](#webhooks)
- [Error Handling](#error-handling)
- [Rate Limiting](#rate-limiting)
- [Examples](#examples)

## Overview

The Chirpy API provides RESTful endpoints for:
- User registration and authentication
- Creating and managing chirps (posts)
- User profile management
- Admin metrics and operations
- Webhook integration for payments

### Base URL
```
http://localhost:8080
```

### Content Type
All API requests should use `Content-Type: application/json` unless otherwise specified.

## Getting Started

### Prerequisites
- Go 1.25.1+
- PostgreSQL database
- Environment variables configured

### Environment Variables
Create a `.env` file with:
```env
DB_URL=postgresql://user:password@localhost/dbname
PLATFORM=dev
SERVER_SECRET=your-secret-key
POLKA_KEY=your-polka-api-key
```

### Running the Server
```bash
go run .
```

The server will start on port 8080.

### Health Check
```bash
curl http://localhost:8080/api/healthz
```
Returns `200 OK` if the service is running.

## Authentication

Chirpy uses JWT (JSON Web Tokens) for authentication. There are two types of tokens:

### Access Token
- Short-lived (1 hour by default)
- Used for authenticating API requests
- Sent in the `Authorization` header as `Bearer <token>`

### Refresh Token
- Long-lived
- Used to obtain new access tokens
- Stored securely by the client

### Authentication Flow
1. Register or login to receive tokens
2. Include access token in API requests
3. Use refresh token to get new access tokens when needed

## API Endpoints

### Users
- `POST /api/users` - Register new user
- `PUT /api/users` - Update user credentials

### Chirps
- `POST /api/chirps` - Create new chirp
- `GET /api/chirps` - List all chirps
- `GET /api/chirps/{chirpId}` - Get specific chirp
- `DELETE /api/chirps/{chirpId}` - Delete chirp

### Authentication
- `POST /api/login` - User login
- `POST /api/refresh` - Refresh access token
- `POST /api/revoke` - Revoke refresh token

### Admin
- `GET /admin/metrics` - Get system metrics
- `POST /admin/reset` - Reset system (dev only)

### Webhooks
- `POST /api/polka/webhooks` - Handle payment webhooks

## Error Handling

The API returns consistent error responses:

```json
{
  "error": "Error message description"
}
```

Common HTTP status codes:
- `200` - Success
- `201` - Created
- `400` - Bad Request
- `401` - Unauthorized
- `403` - Forbidden
- `404` - Not Found
- `500` - Internal Server Error

## Rate Limiting

Currently, there are no rate limits implemented, but clients should be reasonable with their request frequency.

## Examples

See the [Examples](./examples.md) page for complete usage examples and common workflows.
