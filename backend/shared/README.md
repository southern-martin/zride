# Shared utilities and models for ZRide backend services
This directory contains common packages used across all microservices:

## Packages

### `/models`
- Common data structures and models
- Database entity definitions
- API request/response models

### `/utils`
- Common utility functions
- Helper functions for date, string manipulation
- Validation utilities

### `/middleware`
- Authentication middleware
- Logging middleware
- CORS middleware
- Rate limiting middleware

### `/database`
- Database connection utilities
- Migration helpers
- Common database operations

### `/config`
- Configuration management
- Environment variable handling
- Service discovery

### `/logger`
- Structured logging setup
- Log formatting and levels

### `/errors`
- Custom error types
- Error handling utilities
- HTTP error responses

## Usage
Import shared packages in individual services:

```go
import (
    "github.com/zride/backend/shared/models"
    "github.com/zride/backend/shared/utils"
    "github.com/zride/backend/shared/middleware"
)
```