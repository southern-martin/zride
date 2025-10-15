# ZRide Zalo Mini App

## Overview
The frontend application for ZRide ride-sharing platform, built as a Zalo Mini App.

## Features
- User authentication via Zalo OAuth
- Trip browsing and booking
- Real-time chat with drivers
- Payment integration with ZaloPay
- GPS location services
- Push notifications

## Development Setup

### Prerequisites
- Node.js 18+
- Zalo Mini App Developer Account
- Zalo OA (Official Account)

### Installation
```bash
npm install
```

### Configuration
1. Copy `.env.example` to `.env`
2. Update Zalo App credentials
3. Set API endpoints

### Development Server
```bash
npm run dev
```

### Build for Production
```bash
npm run build
```

## Project Structure
```
src/
├── pages/           # Mini app pages
├── components/      # Reusable components  
├── services/        # API services
├── utils/           # Utility functions
├── assets/          # Images, styles
└── config/          # App configuration
```

## Zalo Mini App Configuration
- App ID: Configure in Zalo Developer Console
- Domain Whitelist: Add your API domains
- Permissions: Location, Camera, Payment

## API Integration
All API calls go through the backend services via API Gateway at port 8080.

## Testing
```bash
npm test
```

## Deployment
Deploy to Zalo Mini App platform through developer console.