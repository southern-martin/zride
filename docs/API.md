# ZRide API Documentation

## Overview
ZRide provides a comprehensive RESTful API for managing ride-sharing operations through a microservices architecture.

## Base URL
```
Development: http://localhost:8080
Production: https://api.zride.vn
```

## Authentication
All API requests require authentication using JWT tokens obtained from the auth service.

### Headers
```
Authorization: Bearer <jwt_token>
Content-Type: application/json
```

## Services Overview

### 1. Auth Service (Port: 8001)
Handles user authentication and authorization.

#### Endpoints:
- `POST /auth/login` - User login with Zalo OAuth
- `POST /auth/refresh` - Refresh JWT token
- `POST /auth/logout` - User logout
- `GET /auth/me` - Get current user info

### 2. User Service (Port: 8002)
Manages user profiles and driver information.

#### Endpoints:
- `GET /users/profile` - Get user profile
- `PUT /users/profile` - Update user profile
- `POST /users/driver` - Register as driver
- `GET /users/driver/{id}` - Get driver details
- `PUT /users/driver/{id}/rating` - Rate a driver

### 3. Trip Service (Port: 8003)
Manages trip creation, updates, and bookings.

#### Endpoints:
- `GET /trips` - List available trips
- `POST /trips` - Create new trip
- `GET /trips/{id}` - Get trip details
- `PUT /trips/{id}` - Update trip
- `DELETE /trips/{id}` - Cancel trip
- `POST /trips/{id}/book` - Book a trip
- `PUT /trips/{id}/status` - Update trip status

### 4. Matching Service (Port: 8004)
AI-powered trip and passenger matching.

#### Endpoints:
- `POST /matching/suggest-trips` - Get trip suggestions for passenger
- `POST /matching/suggest-passengers` - Get passenger suggestions for driver
- `GET /matching/routes` - Get optimized routes

### 5. Payment Service (Port: 8005)
Handles payments and transactions.

#### Endpoints:
- `POST /payments/create` - Create payment
- `GET /payments/{id}` - Get payment status
- `POST /payments/webhook/zalopay` - ZaloPay webhook
- `GET /payments/history` - Payment history

## Data Models

### User
```json
{
  "id": "uuid",
  "zalo_id": "string",
  "name": "string",
  "phone": "string",
  "email": "string",
  "avatar": "string",
  "is_driver": "boolean",
  "rating": "float",
  "created_at": "timestamp",
  "updated_at": "timestamp"
}
```

### Driver
```json
{
  "id": "uuid",
  "user_id": "uuid",
  "license_number": "string",
  "vehicle_type": "string",
  "vehicle_plate": "string",
  "vehicle_capacity": "integer",
  "is_verified": "boolean",
  "rating": "float",
  "total_trips": "integer"
}
```

### Trip
```json
{
  "id": "uuid",
  "driver_id": "uuid",
  "origin": {
    "latitude": "float",
    "longitude": "float",
    "address": "string"
  },
  "destination": {
    "latitude": "float",
    "longitude": "float",
    "address": "string"
  },
  "departure_time": "timestamp",
  "available_seats": "integer",
  "price_per_seat": "integer",
  "status": "string", // pending, active, completed, cancelled
  "created_at": "timestamp",
  "updated_at": "timestamp"
}
```

### Booking
```json
{
  "id": "uuid",
  "trip_id": "uuid",
  "passenger_id": "uuid",
  "seats_booked": "integer",
  "total_amount": "integer",
  "status": "string", // pending, confirmed, completed, cancelled
  "payment_id": "uuid",
  "created_at": "timestamp"
}
```

## Error Handling

### Standard Error Response
```json
{
  "error": {
    "code": "ERROR_CODE",
    "message": "Human readable error message",
    "details": {}
  }
}
```

### HTTP Status Codes
- `200` - Success
- `201` - Created
- `400` - Bad Request
- `401` - Unauthorized
- `403` - Forbidden
- `404` - Not Found
- `409` - Conflict
- `422` - Unprocessable Entity
- `500` - Internal Server Error

## Rate Limiting
- 100 requests per minute per IP
- Burst limit of 10 requests

## WebSocket Events (Future)
Real-time updates for trip status and location tracking.

### Events:
- `trip.status_updated`
- `booking.confirmed`
- `driver.location_updated`
- `message.received`

## SDK and Examples

### JavaScript/Node.js
```javascript
const ZRideAPI = require('@zride/sdk');

const client = new ZRideAPI({
  baseURL: 'http://localhost:8080',
  apiKey: 'your_api_key'
});

// Get available trips
const trips = await client.trips.list({
  origin: { lat: 10.762622, lng: 106.660172 },
  destination: { lat: 10.8231, lng: 106.6297 }
});
```

### Python
```python
from zride_sdk import ZRideClient

client = ZRideClient(
    base_url='http://localhost:8080',
    api_key='your_api_key'
)

# Create a new trip
trip = client.trips.create({
    'origin': {'lat': 10.762622, 'lng': 106.660172},
    'destination': {'lat': 10.8231, 'lng': 106.6297},
    'departure_time': '2025-10-16T08:00:00Z',
    'available_seats': 3,
    'price_per_seat': 50000
})
```

## Testing
Use the provided Postman collection for API testing:
- Import `docs/postman/ZRide_API.postman_collection.json`
- Set environment variables for base URL and auth token

## Support
For API support, contact: api-support@zride.vn