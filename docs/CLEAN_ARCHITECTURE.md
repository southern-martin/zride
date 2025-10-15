# Clean Architecture Implementation Guide

## Overview
This document outlines the clean architecture implementation for the ZRide project, following Domain-Driven Design (DDD) principles and SOLID design patterns.

## Architecture Layers

### 1. Domain Layer (`/domain`)
**Purpose**: Contains pure business logic, entities, and domain rules

**Components**:
- **Entities**: Core business objects with identity (`User`, `Trip`, `Booking`)
- **Value Objects**: Immutable objects that describe aspects of domain
- **Aggregates**: Cluster of entities treated as single unit
- **Domain Events**: Events that represent business occurrences
- **Repository Interfaces**: Contracts for data persistence (no implementation)
- **Domain Services**: Business logic that doesn't belong to entities

**Key Principles**:
- ✅ No dependencies on external layers
- ✅ Contains business rules and validations
- ✅ Framework-agnostic
- ✅ Testable without infrastructure

**Example Structure**:
```
domain/
├── entities/
│   ├── user.go
│   ├── trip.go
│   └── booking.go
├── value_objects/
│   ├── location.go
│   └── money.go
├── events/
│   └── domain_events.go
├── repositories/
│   └── interfaces.go
└── services/
    └── domain_services.go
```

### 2. Application Layer (`/application`)
**Purpose**: Orchestrates business workflows and use cases

**Components**:
- **Use Cases**: Application-specific business rules
- **Commands & Queries**: CQRS pattern implementation
- **Command/Query Handlers**: Process commands and queries
- **DTOs**: Data transfer objects for API boundaries
- **Application Services**: Coordinate between domain and infrastructure

**Key Principles**:
- ✅ Depends only on Domain layer
- ✅ Contains application-specific business rules
- ✅ Orchestrates domain objects
- ✅ Defines interfaces for infrastructure dependencies

**Example Structure**:
```
application/
├── usecases/
│   ├── login_usecase.go
│   ├── create_trip_usecase.go
│   └── book_trip_usecase.go
├── commands/
│   └── command_definitions.go
├── queries/
│   └── query_definitions.go
├── handlers/
│   ├── command_handlers.go
│   └── query_handlers.go
└── dto/
    ├── request_dtos.go
    └── response_dtos.go
```

### 3. Infrastructure Layer (`/infrastructure`)
**Purpose**: Implements external concerns and technical details

**Components**:
- **Repository Implementations**: Database access implementations
- **External Service Adapters**: Third-party service integrations
- **Database Configuration**: Connection and migration management
- **Messaging**: Event publishing and consumption
- **Caching**: Redis and other caching implementations

**Key Principles**:
- ✅ Implements interfaces defined in inner layers
- ✅ Contains framework-specific code
- ✅ Handles external system integration
- ✅ Can depend on all other layers

**Example Structure**:
```
infrastructure/
├── repositories/
│   ├── postgresql/
│   ├── redis/
│   └── mongodb/
├── external_services/
│   ├── zalo_service.go
│   ├── zalopay_service.go
│   └── google_maps_service.go
├── database/
│   ├── connection.go
│   └── migrations.go
└── messaging/
    ├── event_publisher.go
    └── event_consumer.go
```

### 4. Interface Layer (`/interfaces`)
**Purpose**: Handles external communication (HTTP, gRPC, etc.)

**Components**:
- **HTTP Handlers**: REST API endpoints
- **Middleware**: Authentication, logging, CORS
- **Request/Response Models**: API-specific data structures
- **Route Configuration**: URL routing setup
- **Input Validation**: Request validation logic

**Key Principles**:
- ✅ Translates external requests to application layer calls
- ✅ Handles protocol-specific concerns
- ✅ Validates input data
- ✅ Formats responses

**Example Structure**:
```
interfaces/
├── http/
│   ├── handlers/
│   ├── middleware/
│   ├── routes/
│   └── models/
├── grpc/
│   ├── handlers/
│   └── proto/
└── cli/
    └── commands/
```

## Implementation Patterns

### 1. Dependency Injection Pattern
```go
// Interface definition in application layer
type UserRepository interface {
    Save(ctx context.Context, user *domain.User) error
    FindByID(ctx context.Context, id string) (*domain.User, error)
}

// Implementation in infrastructure layer
type PostgreSQLUserRepository struct {
    db *sql.DB
}

func (r *PostgreSQLUserRepository) Save(ctx context.Context, user *domain.User) error {
    // Implementation
}

// Injection in main.go
func main() {
    // Infrastructure
    db := infrastructure.NewDatabase(config)
    userRepo := infrastructure.NewPostgreSQLUserRepository(db)
    
    // Application
    loginUseCase := application.NewLoginUseCase(userRepo, ...)
    
    // Interface
    authHandler := interfaces.NewAuthHandler(loginUseCase)
}
```

### 2. Repository Pattern
```go
// Domain interface (contracts)
type Repository[T AggregateRoot] interface {
    Save(ctx context.Context, entity T) error
    FindByID(ctx context.Context, id string) (T, error)
    Delete(ctx context.Context, id string) error
}

// Specific repository extends base
type UserRepository interface {
    Repository[*User]
    FindByZaloID(ctx context.Context, zaloID string) (*User, error)
    FindByEmail(ctx context.Context, email string) (*User, error)
}
```

### 3. CQRS Pattern
```go
// Command (Write operations)
type CreateTripCommand struct {
    DriverID      string
    Origin        LocationDTO
    Destination   LocationDTO
    DepartureTime time.Time
    AvailableSeats int
    PricePerSeat  decimal.Decimal
}

// Query (Read operations)
type FindTripsQuery struct {
    Origin      LocationDTO
    Destination LocationDTO
    DateRange   DateRange
    Pagination  PaginationParams
}

// Handlers
type CreateTripHandler struct {
    tripRepo domain.TripRepository
}

func (h *CreateTripHandler) Handle(ctx context.Context, cmd CreateTripCommand) error {
    // Implementation
}
```

### 4. Domain Events Pattern
```go
// Domain event
type TripCreated struct {
    *domain.BaseDomainEvent
    TripID   string
    DriverID string
    Route    Route
}

// Event publisher (infrastructure)
type EventPublisher interface {
    Publish(ctx context.Context, event domain.DomainEvent) error
}

// Usage in use case
func (uc *CreateTripUseCase) Execute(ctx context.Context, cmd CreateTripCommand) error {
    trip, err := domain.NewTrip(...)
    if err != nil {
        return err
    }
    
    if err := uc.tripRepo.Save(ctx, trip); err != nil {
        return err
    }
    
    event := &TripCreated{
        BaseDomainEvent: domain.NewDomainEvent("trip.created", trip.ID, trip),
        TripID: trip.ID,
        DriverID: trip.DriverID,
        Route: trip.Route,
    }
    
    return uc.eventPublisher.Publish(ctx, event)
}
```

## Service Structure

### Each microservice follows this pattern:

```
service-name/
├── cmd/
│   └── main.go              # Entry point and dependency injection
├── internal/
│   ├── domain/              # Business entities and rules
│   │   ├── entities.go
│   │   ├── repositories.go
│   │   └── services.go
│   ├── application/         # Use cases and application logic
│   │   ├── usecases/
│   │   ├── commands/
│   │   ├── queries/
│   │   └── dto/
│   ├── infrastructure/      # External system implementations
│   │   ├── repositories/
│   │   ├── external_services/
│   │   └── database/
│   └── interfaces/          # HTTP/gRPC handlers
│       ├── http/
│       └── grpc/
├── configs/
├── migrations/
└── tests/
```

## Benefits of This Architecture

### 1. **Testability**
- Each layer can be tested independently
- Domain logic is pure and framework-agnostic
- Easy to mock dependencies

### 2. **Maintainability**
- Clear separation of concerns
- Changes in one layer don't affect others
- Easy to understand and modify

### 3. **Flexibility**
- Can swap implementations (PostgreSQL → MongoDB)
- Can change frameworks without affecting business logic
- Easy to add new features

### 4. **Scalability**
- Each service can be scaled independently
- Clear boundaries enable team autonomy
- Supports microservice architecture

### 5. **Reusability**
- Shared domain models across services
- Common infrastructure components
- Reusable application patterns

## Best Practices

### 1. **Domain Layer**
- ✅ Keep entities rich with behavior
- ✅ Validate invariants in domain objects
- ✅ Use value objects for complex attributes
- ❌ No framework dependencies
- ❌ No infrastructure concerns

### 2. **Application Layer**
- ✅ Keep use cases focused and single-purpose
- ✅ Use dependency injection for repositories
- ✅ Handle application-specific validation
- ❌ No direct database access
- ❌ No HTTP concerns

### 3. **Infrastructure Layer**
- ✅ Implement all repository interfaces
- ✅ Handle connection management
- ✅ Implement proper error handling
- ✅ Use transactions appropriately

### 4. **Interface Layer**
- ✅ Validate input thoroughly
- ✅ Transform domain errors to HTTP errors
- ✅ Use proper HTTP status codes
- ✅ Implement proper middleware

## Error Handling Strategy

```go
// Domain errors
var (
    ErrUserNotFound = domain.NewDomainError("USER_NOT_FOUND", "User not found")
    ErrInvalidInput = domain.NewDomainError("INVALID_INPUT", "Invalid input data")
)

// Application layer error handling
func (uc *LoginUseCase) Execute(ctx context.Context, cmd LoginCommand) (*LoginResponse, error) {
    user, err := uc.userRepo.FindByZaloID(ctx, cmd.ZaloID)
    if err != nil {
        if errors.Is(err, domain.ErrUserNotFound) {
            // Handle user not found
            return nil, application.ErrInvalidCredentials
        }
        return nil, err
    }
    // ... rest of logic
}

// HTTP layer error handling
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
    // ... parse request
    
    result, err := h.loginUseCase.Execute(r.Context(), cmd)
    if err != nil {
        h.handleError(w, err)
        return
    }
    
    h.writeJSON(w, http.StatusOK, result)
}

func (h *AuthHandler) handleError(w http.ResponseWriter, err error) {
    var domainErr *domain.DomainError
    if errors.As(err, &domainErr) {
        h.writeError(w, http.StatusBadRequest, domainErr)
        return
    }
    
    h.writeError(w, http.StatusInternalServerError, domain.ErrInternalError)
}
```

## Testing Strategy

### 1. **Unit Tests** (Domain & Application layers)
```go
func TestUser_UpdateProfile(t *testing.T) {
    user := domain.NewUser("zalo_123", "John Doe", "0901234567", "john@example.com", "")
    
    err := user.UpdateProfile("Jane Doe", "0912345678", "jane@example.com", "avatar.jpg")
    
    assert.NoError(t, err)
    assert.Equal(t, "Jane Doe", user.Name)
    assert.Equal(t, "jane@example.com", user.Email)
}

func TestLoginUseCase_Execute(t *testing.T) {
    mockUserRepo := &mocks.UserRepository{}
    mockZaloService := &mocks.ZaloService{}
    
    useCase := application.NewLoginUseCase(mockUserRepo, mockZaloService)
    
    // Setup mocks and test...
}
```

### 2. **Integration Tests** (Infrastructure layer)
```go
func TestPostgreSQLUserRepository_Save(t *testing.T) {
    db := setupTestDB(t)
    repo := infrastructure.NewPostgreSQLUserRepository(db)
    
    user := domain.NewUser("zalo_123", "Test User", "0901234567", "test@example.com", "")
    
    err := repo.Save(context.Background(), user)
    assert.NoError(t, err)
    
    // Verify data was saved...
}
```

### 3. **End-to-End Tests** (Full application)
```go
func TestAuthAPI_Login(t *testing.T) {
    server := setupTestServer(t)
    
    loginReq := LoginRequest{
        ZaloAccessToken: "valid_token",
    }
    
    resp := server.POST("/auth/login").WithJSON(loginReq).Expect()
    resp.Status(http.StatusOK)
    resp.JSON().Object().ContainsKey("access_token")
}
```

This clean architecture provides a solid foundation for building maintainable, testable, and scalable microservices.