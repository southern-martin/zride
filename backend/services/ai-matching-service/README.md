# AI Matching Service

The AI Matching Service is a sophisticated microservice that implements intelligent driver-passenger matching using various algorithms including machine learning approaches for optimal ride assignments in the Zride ride-sharing platform.

## Features

### Core Functionality
- **Multi-Algorithm Matching**: Supports multiple matching algorithms
  - Nearest Driver Algorithm
  - Weighted Scoring Algorithm  
  - Machine Learning Algorithm
  - Hybrid Algorithm (combines multiple approaches)

### AI/ML Capabilities
- **Geospatial Analysis**: Uses PostGIS for efficient location-based queries
- **Real-time Scoring**: Dynamic scoring based on multiple factors
- **Feature Engineering**: Extracts meaningful features for ML algorithms
- **Predictive Matching**: Uses historical data patterns for better matches

### Matching Criteria
- **Distance Optimization**: Haversine distance calculations
- **Driver Rating**: Weighted by driver experience and ratings
- **Time Estimation**: Dynamic ETA calculations
- **Price Optimization**: Competitive pricing within user ranges
- **Preference Matching**: Car type and area preferences

## Architecture

```
┌─────────────────────┐
│   HTTP Layer        │
│  (Gin Framework)    │
├─────────────────────┤
│  Application Layer  │
│   (Use Cases)       │
├─────────────────────┤
│   Domain Layer      │
│  (Business Logic)   │
├─────────────────────┤
│ Infrastructure      │
│ PostgreSQL + Redis  │
└─────────────────────┘
```

### Components
- **Domain Models**: Match requests, drivers, match results
- **Matching Service**: Core AI algorithms implementation
- **Use Cases**: Business logic orchestration
- **Repositories**: Data persistence abstractions
- **HTTP Handlers**: REST API endpoints

## API Endpoints

### Match Requests
```
POST   /api/v1/matching/requests                    # Create match request
GET    /api/v1/matching/passengers/{id}/requests    # Get passenger requests
```

### Driver Operations
```
GET    /api/v1/matching/drivers/{id}/matches        # Get available matches
POST   /api/v1/matching/drivers/{id}/matches/{mid}/accept  # Accept match
```

### Health Check
```
GET    /api/v1/health                               # Service health
```

## Algorithms

### 1. Nearest Driver Algorithm
- Simple distance-based matching
- Fastest response time
- Best for high-demand scenarios

### 2. Weighted Scoring Algorithm
```go
score = (distance_score * 0.4) + 
        (rating_score * 0.2) + 
        (time_score * 0.2) + 
        (price_score * 0.1) + 
        (experience_score * 0.1)
```

### 3. Machine Learning Algorithm
- Feature extraction and engineering
- Non-linear scoring with interactions
- Sigmoid normalization
- Historical pattern learning

### 4. Hybrid Algorithm
- Combines multiple algorithm results
- Weighted ensemble approach
- 30% nearest + 40% weighted + 30% ML

## Configuration

### Environment Variables
```bash
DATABASE_URL=postgres://user:pass@host:port/db
REDIS_URL=redis:6379
PORT=8084
LOG_LEVEL=info
MAX_DRIVERS=10
MAX_SEARCH_RADIUS=15.0
MIN_DRIVER_RATING=3.0
DEFAULT_ALGORITHM=weighted
```

### Matching Configuration
```go
type MatchingConfig struct {
    Algorithm       MatchingAlgorithm
    MaxDrivers      int     // Default: 10
    MaxSearchRadius float64 // Default: 15km  
    MinDriverRating float64 // Default: 3.0
    Criteria        MatchingCriteria
    RealTimeEnabled bool
}
```

## Database Schema

### Tables
- **match_requests**: Passenger ride requests
- **drivers**: Available driver pool with locations
- **match_results**: Matching algorithm results

### Key Features
- PostGIS spatial indexing
- JSON fields for flexible location data
- Optimized indexes for geospatial queries
- Foreign key constraints for data integrity

## Development

### Prerequisites
- Go 1.21+
- PostgreSQL 15+ with PostGIS
- Redis 6+
- Docker (for containerized deployment)

### Local Development
```bash
# Clone repository
git clone <repository-url>
cd zride/backend/services/ai-matching-service

# Install dependencies
go mod tidy

# Copy environment file
cp .env.example .env

# Run database migrations
psql -h localhost -U zride_user -d zride_ai_matching -f migrations/001_init.sql

# Start service
go run cmd/main.go
```

### Docker Deployment
```bash
# Build image
docker build -t zride-ai-matching .

# Run with docker-compose
cd ../../
docker-compose up ai-matching-service
```

### Testing
```bash
# Run unit tests
go test ./...

# Run integration tests
go test -tags=integration ./...

# Test with sample data
curl -X POST http://localhost:8084/api/v1/matching/requests \
  -H "Content-Type: application/json" \
  -d '{
    "passenger_id": "uuid-here",
    "pickup_location": {"latitude": 10.7769, "longitude": 106.7009},
    "dropoff_location": {"latitude": 10.7626, "longitude": 106.6822},
    "max_distance": 15.0
  }'
```

## Performance Characteristics

### Optimizations
- **Spatial Indexing**: PostGIS GIST indexes for O(log n) location queries
- **Redis Caching**: Driver locations cached for 5 minutes
- **Connection Pooling**: PostgreSQL connection pool (25 max connections)
- **Async Processing**: Background match processing to reduce response time

### Scalability
- Horizontal scaling through multiple service instances
- Database read replicas for query distribution
- Redis cluster for cache distribution
- Load balancing through API gateway

### Performance Metrics
- **Match Response Time**: < 2 seconds for 95th percentile
- **Database Query Time**: < 100ms for spatial queries
- **Algorithm Processing**: < 500ms for 50 drivers
- **Memory Usage**: < 100MB per instance

## Monitoring and Observability

### Logging
- Structured JSON logging
- Request/response logging
- Algorithm performance metrics
- Error tracking and alerting

### Health Checks
- Database connectivity
- Redis connectivity  
- Algorithm performance validation
- Service dependency checks

## Future Enhancements

### Planned Features
- **Deep Learning**: Neural network-based matching
- **Real-time Learning**: Online learning from user feedback
- **Multi-objective Optimization**: Pareto optimal solutions
- **Predictive Analytics**: Demand forecasting
- **A/B Testing**: Algorithm performance comparison

### Integration Points
- **Trip Service**: Real-time trip updates
- **User Service**: User preferences and history
- **Payment Service**: Dynamic pricing integration
- **Notification Service**: Real-time match notifications

## Vietnamese Market Adaptations

### Local Considerations
- **Traffic Patterns**: Ho Chi Minh City traffic optimization
- **Weather Integration**: Monsoon season adjustments
- **Cultural Preferences**: Vietnamese user behavior patterns
- **Pricing Strategy**: VND currency optimizations
- **Peak Hours**: Local rush hour patterns (7-9 AM, 5-7 PM)

## License

This project is part of the Zride ride-sharing platform.