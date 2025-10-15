# User Service

A comprehensive user management microservice for the Zride ride-sharing platform. This service handles user profiles, driver vehicle management, and rating systems with full CRUD operations.

## Features

- **User Profile Management**: Complete profile CRUD operations with avatar support
- **Vehicle Management**: Driver vehicle registration and management
- **Rating System**: User and driver rating management
- **JWT Authentication**: Secure API endpoints with JWT validation
- **File Upload**: Avatar and document upload with size validation
- **Soft Deletes**: Data retention with soft delete functionality
- **Clean Architecture**: Domain-driven design with clear separation of concerns

## Architecture

```
user-service/
├── internal/
│   ├── domain/          # Business entities and rules
│   ├── application/     # Business logic and DTOs
│   ├── infrastructure/  # Data access and external services
│   └── interfaces/      # HTTP handlers and routing
├── shared/             # Shared utilities and error handling
├── uploads/            # File upload storage
├── main.go            # Application entry point
└── config/            # Configuration management
```

## API Endpoints

### User Profiles

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| POST | `/api/v1/users/profile` | Create user profile | Yes |
| GET | `/api/v1/users/profile` | Get own profile | Yes |
| PUT | `/api/v1/users/profile` | Update own profile | Yes |
| DELETE | `/api/v1/users/profile` | Delete own profile | Yes |
| POST | `/api/v1/users/profile/avatar` | Upload avatar | Yes |

### Vehicles (Driver Only)

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| POST | `/api/v1/users/vehicles` | Add vehicle | Yes |
| GET | `/api/v1/users/vehicles` | Get own vehicles | Yes |
| GET | `/api/v1/users/vehicles/:id` | Get vehicle by ID | Yes |
| PUT | `/api/v1/users/vehicles/:id` | Update vehicle | Yes |
| DELETE | `/api/v1/users/vehicles/:id` | Delete vehicle | Yes |

### Ratings

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| POST | `/api/v1/users/ratings` | Create rating | Yes |
| GET | `/api/v1/users/ratings` | Get user ratings | Yes |
| GET | `/api/v1/users/ratings/:id` | Get rating by ID | Yes |
| PUT | `/api/v1/users/ratings/:id` | Update rating | Yes |
| DELETE | `/api/v1/users/ratings/:id` | Delete rating | Yes |

### Health Check

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| GET | `/health` | Health check | No |

## Data Models

### UserProfile

```json
{
  "id": "uuid",
  "user_id": "uuid",
  "first_name": "string",
  "last_name": "string",
  "phone": "string",
  "avatar_url": "string",
  "date_of_birth": "date",
  "gender": "male|female|other",
  "address": "string",
  "emergency_contact": {
    "name": "string",
    "phone": "string",
    "relationship": "string"
  },
  "preferences": {
    "language": "string",
    "currency": "string",
    "notifications": "object"
  },
  "verification_status": "unverified|pending|verified|rejected",
  "is_driver": "boolean",
  "driver_license": "string",
  "created_at": "datetime",
  "updated_at": "datetime"
}
```

### Vehicle

```json
{
  "id": "uuid",
  "user_id": "uuid",
  "make": "string",
  "model": "string",
  "year": "integer",
  "license_plate": "string",
  "color": "string",
  "vehicle_type": "car|motorcycle|bicycle",
  "seats": "integer",
  "features": ["string"],
  "documents": {
    "registration": "string",
    "insurance": "string"
  },
  "status": "active|inactive|maintenance",
  "created_at": "datetime",
  "updated_at": "datetime"
}
```

### Rating

```json
{
  "id": "uuid",
  "rated_user_id": "uuid",
  "rater_user_id": "uuid",
  "trip_id": "uuid",
  "rating": "integer (1-5)",
  "comment": "string",
  "rating_type": "passenger|driver",
  "created_at": "datetime",
  "updated_at": "datetime"
}
```

## Configuration

### Environment Variables

| Variable | Description | Required | Default |
|----------|-------------|----------|---------|
| `PORT` | Server port | No | 8082 |
| `DB_HOST` | Database host | Yes | - |
| `DB_PORT` | Database port | Yes | - |
| `DB_USER` | Database user | Yes | - |
| `DB_PASSWORD` | Database password | Yes | - |
| `DB_NAME` | Database name | Yes | - |
| `JWT_SECRET` | JWT secret key | Yes | - |
| `AUTH_SERVICE_URL` | Auth service URL | Yes | - |
| `UPLOAD_MAX_SIZE` | Max upload size in MB | No | 10 |
| `UPLOAD_PATH` | Upload directory path | No | ./uploads |

### Example Configuration

```bash
PORT=8082
DB_HOST=localhost
DB_PORT=5432
DB_USER=zride_user
DB_PASSWORD=your_password
DB_NAME=zride_users
JWT_SECRET=your_jwt_secret_key_here
AUTH_SERVICE_URL=http://localhost:8081
UPLOAD_MAX_SIZE=10
UPLOAD_PATH=./uploads
```

## Database Schema

### Users Table
```sql
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL UNIQUE,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    phone VARCHAR(20),
    avatar_url TEXT,
    date_of_birth DATE,
    gender VARCHAR(10),
    address TEXT,
    emergency_contact JSONB,
    preferences JSONB DEFAULT '{}',
    verification_status VARCHAR(20) DEFAULT 'unverified',
    is_driver BOOLEAN DEFAULT false,
    driver_license VARCHAR(50),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);
```

### Vehicles Table
```sql
CREATE TABLE vehicles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(user_id),
    make VARCHAR(50) NOT NULL,
    model VARCHAR(50) NOT NULL,
    year INTEGER NOT NULL,
    license_plate VARCHAR(20) NOT NULL UNIQUE,
    color VARCHAR(30) NOT NULL,
    vehicle_type VARCHAR(20) NOT NULL,
    seats INTEGER DEFAULT 4,
    features JSONB DEFAULT '[]',
    documents JSONB DEFAULT '{}',
    status VARCHAR(20) DEFAULT 'active',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);
```

### Ratings Table
```sql
CREATE TABLE ratings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    rated_user_id UUID NOT NULL,
    rater_user_id UUID NOT NULL,
    trip_id UUID NOT NULL,
    rating INTEGER NOT NULL CHECK (rating >= 1 AND rating <= 5),
    comment TEXT,
    rating_type VARCHAR(20) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);
```

## Development

### Prerequisites

- Go 1.21+
- PostgreSQL 13+
- Docker (optional)

### Setup

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd zride/backend/services/user-service
   ```

2. **Install dependencies**
   ```bash
   make deps
   ```

3. **Setup environment**
   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

4. **Setup development tools**
   ```bash
   make dev-setup
   ```

### Running the Service

#### Development Mode (with live reload)
```bash
make dev
```

#### Standard Mode
```bash
make run
```

#### Docker Mode
```bash
make docker-build
make docker-run
```

### Testing

#### Run Tests
```bash
make test
```

#### Run Tests with Coverage
```bash
make test-coverage
```

### Building

#### Local Build
```bash
make build
```

#### Docker Build
```bash
make docker-build
```

## Integration

### Authentication

This service integrates with the Auth Service for JWT token validation. Include the JWT token in the Authorization header:

```
Authorization: Bearer <your-jwt-token>
```

### Service Discovery

The service registers itself with the service discovery mechanism and can be accessed at:
- **Development**: `http://localhost:8082`
- **Production**: Via API Gateway

### Database Integration

Requires PostgreSQL connection with the following features:
- JSON/JSONB support for complex data structures
- UUID support for primary keys
- Timezone-aware timestamps

## Error Handling

The service uses structured error responses:

```json
{
  "error": {
    "code": "USER_NOT_FOUND",
    "message": "User not found",
    "details": "No user found with the provided ID"
  }
}
```

### Error Codes

- `USER_NOT_FOUND` - User does not exist
- `VEHICLE_NOT_FOUND` - Vehicle does not exist
- `RATING_NOT_FOUND` - Rating does not exist
- `UNAUTHORIZED` - Authentication required
- `FORBIDDEN` - Insufficient permissions
- `VALIDATION_ERROR` - Invalid input data
- `FILE_TOO_LARGE` - Upload file exceeds size limit
- `INVALID_FILE_TYPE` - Unsupported file type

## Monitoring

### Health Check

The service exposes a health check endpoint at `/health` that returns:

```json
{
  "status": "healthy",
  "service": "user-service",
  "version": "1.0.0",
  "timestamp": "2024-01-01T00:00:00Z",
  "database": "connected"
}
```

### Metrics

- Request latency
- Error rates
- Database connection status
- Upload success/failure rates

## Security

- JWT token validation for all protected endpoints
- Input validation and sanitization
- File upload size and type restrictions
- SQL injection prevention
- Soft delete for data retention

## Contributing

1. Follow Go conventions and best practices
2. Maintain clean architecture principles
3. Add comprehensive tests for new features
4. Update documentation for API changes
5. Use conventional commit messages

## License

This project is part of the Zride platform. All rights reserved.