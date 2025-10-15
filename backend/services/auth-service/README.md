# Authentication Service

A microservice for handling user authentication using Zalo OAuth integration with JWT token management.

## Features

- **Zalo OAuth Integration**: Secure login using Zalo Mini App authentication
- **JWT Token Management**: Access and refresh token generation with proper expiration
- **User Management**: User profile creation and management
- **Clean Architecture**: Domain-driven design with clear separation of concerns
- **PostgreSQL Integration**: Persistent user data storage
- **Docker Ready**: Containerized for easy deployment
- **Health Checks**: Service monitoring and health endpoints

## API Endpoints

### Public Endpoints

- `GET /health` - Health check
- `POST /api/v1/auth/login` - User login with Zalo code
- `POST /api/v1/auth/refresh` - Refresh access token
- `POST /api/v1/auth/logout` - User logout
- `POST /api/v1/auth/validate` - Token validation

### Authentication Flow

1. User authenticates with Zalo Mini App
2. Frontend receives authorization code from Zalo
3. Send code to `/api/v1/auth/login` endpoint
4. Service exchanges code for Zalo access token
5. Retrieves user profile from Zalo API
6. Creates or updates user in database
7. Returns JWT access and refresh tokens

## Environment Variables

```env
# Server Configuration
PORT=8081

# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=zride_user
DB_PASSWORD=zride_password
DB_NAME=zride
DB_SSL_MODE=disable

# JWT Configuration
JWT_SECRET=your-super-secret-jwt-key
JWT_ACCESS_EXPIRY=15m
JWT_REFRESH_EXPIRY=7d

# Zalo OAuth Configuration
ZALO_APP_ID=your_zalo_app_id
ZALO_APP_SECRET=your_zalo_app_secret
```

## Running the Service

### Local Development

1. Copy environment variables:
   ```bash
   cp .env.example .env
   ```

2. Update the `.env` file with your actual Zalo app credentials

3. Ensure PostgreSQL is running with the zride database

4. Run the service:
   ```bash
   go run main.go
   ```

### Docker

1. Build the image:
   ```bash
   docker build -t zride-auth-service .
   ```

2. Run the container:
   ```bash
   docker run -p 8081:8081 --env-file .env zride-auth-service
   ```

## Database Schema

The service requires the following PostgreSQL table:

```sql
-- Users table (see database/schema/001_init.sql)
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    zalo_id VARCHAR(255) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    phone VARCHAR(20),
    email VARCHAR(255),
    picture TEXT,
    user_type VARCHAR(20) NOT NULL DEFAULT 'passenger',
    is_active BOOLEAN NOT NULL DEFAULT true,
    last_login_at TIMESTAMP WITH TIME ZONE,
    refresh_token TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE,
    version INTEGER NOT NULL DEFAULT 1
);
```

## Project Structure

```
.
├── main.go                     # Application entry point
├── internal/
│   ├── application/           # Business logic layer
│   │   └── auth_service.go   # Auth service implementation
│   ├── domain/               # Domain entities and interfaces
│   │   ├── user.go          # User aggregate
│   │   └── repository.go    # Repository interfaces
│   ├── infrastructure/      # External services and data access
│   │   ├── config.go       # Configuration management
│   │   ├── jwt_service.go  # JWT token service
│   │   ├── user_repository.go # PostgreSQL user repository
│   │   └── zalo_oauth_service.go # Zalo API integration
│   └── interfaces/         # HTTP handlers and routes
│       ├── auth_handler.go # Authentication endpoints
│       ├── middleware.go   # Authentication middleware
│       └── router.go       # Route setup
├── cmd/
│   └── test_client.go     # Simple test client
├── .env.example           # Environment variables template
└── Dockerfile             # Container configuration
```

## Integration with Other Services

This service is designed to be used by:

- **API Gateway**: For request authentication
- **User Service**: For user profile management
- **Trip Service**: For user identification in bookings
- **Payment Service**: For transaction user context

## Security Features

- **JWT Secret**: Configurable secret key for token signing
- **Token Expiration**: Configurable access and refresh token lifetimes
- **HTTPS Ready**: Supports SSL/TLS in production
- **CORS Middleware**: Cross-origin request handling
- **Input Validation**: Request payload validation
- **Error Handling**: Structured error responses

## Next Steps

1. Add comprehensive unit and integration tests
2. Implement rate limiting
3. Add request logging and monitoring
4. Implement token blacklisting for secure logout
5. Add user role-based access control
6. Implement account deactivation/suspension features