# 🚗 ZRide - Zalo Ride Sharing Platform

## 📖 Project Overview
ZRide is a Zalo Mini App that connects drivers with empty return trips to passengers needing rides on the same routes. The platform leverages Zalo's ecosystem for seamless user experience and payment integration.

## 🎯 Key Features
- **Driver-Passenger Matching**: AI-powered route matching for optimal ride sharing
- **Zalo Integration**: Native integration with Zalo OA for messaging and authentication
- **Payment Processing**: ZaloPay integration for secure transactions
- **Real-time Tracking**: Trip status and location updates
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