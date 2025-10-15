# Zride - Vietnamese Ride-Sharing Platform

A comprehensive ride-sharing platform built with microservices architecture, specifically designed for the Vietnamese market with Zalo integration.

## � Project Overview

Zride is a modern ride-sharing platform that connects passengers with drivers across Vietnam. The platform features:

- **Zalo Integration** - Seamless authentication using Zalo OAuth
- **Microservices Architecture** - Scalable and maintainable service-oriented design
- **Real-time Matching** - AI-powered driver-passenger matching system
- **Vietnamese Market Focus** - Localized features and payment methods
- **Comprehensive Rating System** - Two-way rating for drivers and passengers

## 🏗️ Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                      API Gateway (Nginx)                    │
│                         Port 80/443                         │
└─────────────────┬─────────────────────┬─────────────────────┘
                  │                     │
        ┌─────────▼──────────┐ ┌───────▼─────────┐
        │   Auth Service     │ │   User Service  │
        │     Port 8081      │ │    Port 8082    │
        │                    │ │                 │
        │ • JWT Auth         │ │ • User Profiles │
        │ • Zalo OAuth       │ │ • Vehicle Mgmt  │
        │ • Token Mgmt       │ │ • Rating System │
        └────────────────────┘ └─────────────────┘
                  │                     │
        ┌─────────▼─────────────────────▼─────────┐
        │            PostgreSQL Database           │
        │     Auth DB    │    Users DB    │       │
        └──────────────────────────────────────────┘
```

### Services

1. **Auth Service** (`/backend/services/auth-service/`)
   - JWT token management
   - Zalo OAuth integration
   - User authentication and authorization
   - Password reset and verification

2. **User Service** (`/backend/services/user-service/`)
   - User profile management
   - Driver vehicle registration
   - Rating and review system
   - File upload handling

3. **Trip Service** (Planned)
   - Trip booking and management
   - Real-time tracking
   - Fare calculation
   - Payment processing

4. **Matching Service** (Planned)
   - AI-powered driver matching
   - Route optimization
   - Price estimation

## 🚀 Quick Start

### Prerequisites

- Docker & Docker Compose
- Go 1.21+ (for development)
- PostgreSQL 15+ (for local development)
- Node.js 18+ (for frontend)

### Development Setup

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd zride
   ```

2. **Start all services with Docker Compose**
   ```bash
   cd backend
   docker-compose up -d
   ```

3. **Or run individual services locally**
   
   **Auth Service:**
   ```bash
   cd backend/services/auth-service
   cp .env.example .env
   # Edit .env with your configuration
   make dev
   ```
   
   **User Service:**
   ```bash
   cd backend/services/user-service
   cp .env.example .env
   # Edit .env with your configuration
   make dev
   ```

### Service URLs

- **API Gateway:** http://localhost
- **Auth Service:** http://localhost:8081
- **User Service:** http://localhost:8082
- **Database:** localhost:5432

## 📡 API Documentation

### Authentication Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/auth/register` | Register new user |
| POST | `/api/v1/auth/login` | User login |
| POST | `/api/v1/auth/zalo/login` | Zalo OAuth login |
| POST | `/api/v1/auth/refresh` | Refresh JWT token |
| POST | `/api/v1/auth/logout` | User logout |
| POST | `/api/v1/auth/forgot-password` | Request password reset |
| POST | `/api/v1/auth/reset-password` | Reset password |

### User Management Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/users/profile` | Create user profile |
| GET | `/api/v1/users/profile` | Get user profile |
| PUT | `/api/v1/users/profile` | Update user profile |
| DELETE | `/api/v1/users/profile` | Delete user profile |
| POST | `/api/v1/users/profile/avatar` | Upload avatar |
| GET | `/api/v1/users/vehicles` | Get user vehicles |
| POST | `/api/v1/users/vehicles` | Add vehicle |
| PUT | `/api/v1/users/vehicles/:id` | Update vehicle |
| DELETE | `/api/v1/users/vehicles/:id` | Delete vehicle |
| GET | `/api/v1/users/ratings` | Get user ratings |
| POST | `/api/v1/users/ratings` | Create rating |
- **Rating System**: Driver and passenger rating for quality assurance

## 🏗️ Architecture
```
┌─────────────────────┐     ┌──────────────────────┐
│   Zalo Mini App     │────▶│    API Gateway       │
│   (Frontend)        │     │                      │
└─────────────────────┘     └──────────┬───────────┘
                                       │
                            ┌──────────▼───────────┐
                            │   Microservices      │
                            │                      │
                            │ • Auth Service       │
                            │ • User Service       │
                            │ • Trip Service       │
                            │ • Matching Service   │
                            │ • Payment Service    │
                            └──────────┬───────────┘
                                       │
                            ┌──────────▼───────────┐
                            │    Databases         │
                            │                      │
                            │ • PostgreSQL (Main)  │
                            │ • Redis (Cache)      │
                            │ • MongoDB (Logs)     │
                            └──────────────────────┘
```

## 🚀 Technology Stack

### Backend
- **Language**: Go (Golang)
- **Framework**: Gin/Fiber
- **Database**: PostgreSQL, Redis
- **AI Service**: Python with FastAPI
- **Message Queue**: Redis/RabbitMQ

### Frontend
- **Platform**: Zalo Mini App
- **Framework**: React Native/HTML5
- **SDK**: Zalo Mini App SDK

### DevOps
- **Containerization**: Docker
- **Orchestration**: Kubernetes
- **CI/CD**: GitHub Actions
- **Load Balancer**: Traefik

## 📁 Project Structure
```
zride/
├── backend/
│   ├── api-gateway/          # API Gateway service
│   ├── services/
│   │   ├── auth-service/     # Authentication & authorization
│   │   ├── user-service/     # User profile management
│   │   ├── trip-service/     # Trip management
│   │   ├── matching-service/ # AI matching engine
│   │   └── payment-service/  # Payment processing
│   └── shared/               # Shared utilities and models
├── frontend/
│   └── zalo-mini-app/        # Zalo Mini App frontend
├── database/
│   ├── migrations/           # Database migrations
│   └── seeds/                # Initial data
├── deployment/
│   ├── docker/               # Docker configurations
│   └── k8s/                  # Kubernetes manifests
├── docs/                     # Project documentation
└── tests/                    # Integration tests
```

## 🔧 Development Setup

### Prerequisites
- Go 1.21+
- Node.js 18+
- PostgreSQL 14+
- Redis 6+
- Docker & Docker Compose

### Quick Start
```bash
# Clone the repository
git clone <repository-url>
cd zride

# Start development environment
docker-compose up -d

# Run backend services
cd backend && make dev

# Run frontend
cd frontend/zalo-mini-app && npm start
```

## 📋 Development Roadmap

### Phase 1: Foundation (Weeks 1-2)
- [ ] Project setup and architecture design
- [ ] Database schema design
- [ ] Basic authentication service

### Phase 2: Core Services (Weeks 3-6)
- [ ] User and trip management services
- [ ] Basic matching algorithm
- [ ] API Gateway setup

### Phase 3: Frontend Development (Weeks 7-10)
- [ ] Zalo Mini App development
- [ ] User interface and experience
- [ ] Integration with backend services

### Phase 4: Integration & Testing (Weeks 11-12)
- [ ] Payment service integration
- [ ] End-to-end testing
- [ ] Performance optimization

### Phase 5: Deployment (Weeks 13-16)
- [ ] Production deployment setup
- [ ] Monitoring and logging
- [ ] Beta testing in target market

## 👥 Team
- **Project Lead**: Nguyễn Phạm Văn Tân
- **Status**: MVP Planning Phase

## 📄 License
MIT License - see [LICENSE](LICENSE) file for details

## 📞 Contact
For questions and support, please contact: [Contact Information]

---

**Last Updated**: October 15, 2025  
**Version**: 1.0.0-alpha