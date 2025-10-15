-- Initialize Zride Database Schema
-- This script creates the necessary databases and tables for all microservices

-- Create databases for different services
CREATE DATABASE zride_auth;
CREATE DATABASE zride_users;
CREATE DATABASE zride_trips;

-- Switch to auth database
\c zride_auth;

-- Auth service tables
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    zalo_id VARCHAR(100) UNIQUE,
    email VARCHAR(255) UNIQUE,
    phone VARCHAR(20) UNIQUE,
    password_hash VARCHAR(255),
    is_verified BOOLEAN DEFAULT false,
    verification_token VARCHAR(255),
    reset_token VARCHAR(255),
    reset_token_expires TIMESTAMP WITH TIME ZONE,
    last_login TIMESTAMP WITH TIME ZONE,
    login_attempts INTEGER DEFAULT 0,
    locked_until TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE refresh_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id),
    token_hash VARCHAR(255) NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create indexes for auth service
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_phone ON users(phone);
CREATE INDEX idx_users_zalo_id ON users(zalo_id);
CREATE INDEX idx_refresh_tokens_user_id ON refresh_tokens(user_id);
CREATE INDEX idx_refresh_tokens_expires_at ON refresh_tokens(expires_at);

-- Switch to users database
\c zride_users;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- User service tables
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL UNIQUE, -- References auth service user ID
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    phone VARCHAR(20),
    avatar_url TEXT,
    date_of_birth DATE,
    gender VARCHAR(10) CHECK (gender IN ('male', 'female', 'other')),
    address TEXT,
    emergency_contact JSONB DEFAULT '{}',
    preferences JSONB DEFAULT '{}',
    verification_status VARCHAR(20) DEFAULT 'unverified' CHECK (verification_status IN ('unverified', 'pending', 'verified', 'rejected')),
    is_driver BOOLEAN DEFAULT false,
    driver_license VARCHAR(50),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE vehicles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(user_id),
    make VARCHAR(50) NOT NULL,
    model VARCHAR(50) NOT NULL,
    year INTEGER NOT NULL CHECK (year >= 1900 AND year <= EXTRACT(YEAR FROM NOW()) + 1),
    license_plate VARCHAR(20) NOT NULL UNIQUE,
    color VARCHAR(30) NOT NULL,
    vehicle_type VARCHAR(20) NOT NULL CHECK (vehicle_type IN ('car', 'motorcycle', 'bicycle')),
    seats INTEGER DEFAULT 4 CHECK (seats >= 1 AND seats <= 50),
    features JSONB DEFAULT '[]',
    documents JSONB DEFAULT '{}',
    status VARCHAR(20) DEFAULT 'active' CHECK (status IN ('active', 'inactive', 'maintenance')),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE ratings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    rated_user_id UUID NOT NULL, -- User being rated
    rater_user_id UUID NOT NULL, -- User giving the rating
    trip_id UUID NOT NULL, -- Trip this rating is for
    rating INTEGER NOT NULL CHECK (rating >= 1 AND rating <= 5),
    comment TEXT,
    rating_type VARCHAR(20) NOT NULL CHECK (rating_type IN ('passenger', 'driver')),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE,
    
    -- Prevent duplicate ratings for the same trip and type
    UNIQUE(trip_id, rater_user_id, rating_type)
);

-- Create indexes for user service
CREATE INDEX idx_users_user_id ON users(user_id);
CREATE INDEX idx_users_phone ON users(phone);
CREATE INDEX idx_users_is_driver ON users(is_driver);
CREATE INDEX idx_users_verification_status ON users(verification_status);

CREATE INDEX idx_vehicles_user_id ON vehicles(user_id);
CREATE INDEX idx_vehicles_license_plate ON vehicles(license_plate);
CREATE INDEX idx_vehicles_status ON vehicles(status);
CREATE INDEX idx_vehicles_type ON vehicles(vehicle_type);

CREATE INDEX idx_ratings_rated_user_id ON ratings(rated_user_id);
CREATE INDEX idx_ratings_rater_user_id ON ratings(rater_user_id);
CREATE INDEX idx_ratings_trip_id ON ratings(trip_id);
CREATE INDEX idx_ratings_type ON ratings(rating_type);

-- Switch to trips database (for future trip service)
\c zride_trips;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Trip service tables (placeholder for future implementation)
CREATE TABLE trips (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    passenger_id UUID NOT NULL,
    driver_id UUID,
    pickup_location JSONB NOT NULL,
    dropoff_location JSONB NOT NULL,
    pickup_time TIMESTAMP WITH TIME ZONE,
    dropoff_time TIMESTAMP WITH TIME ZONE,
    distance_km DECIMAL(10,2),
    duration_minutes INTEGER,
    fare_amount DECIMAL(10,2),
    status VARCHAR(20) DEFAULT 'requested' CHECK (status IN ('requested', 'accepted', 'in_progress', 'completed', 'cancelled')),
    vehicle_id UUID,
    payment_method VARCHAR(20),
    payment_status VARCHAR(20) DEFAULT 'pending',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Create indexes for trip service
CREATE INDEX idx_trips_passenger_id ON trips(passenger_id);
CREATE INDEX idx_trips_driver_id ON trips(driver_id);
CREATE INDEX idx_trips_status ON trips(status);
CREATE INDEX idx_trips_created_at ON trips(created_at);
CREATE INDEX idx_trips_pickup_time ON trips(pickup_time);

-- Insert some sample data for development
\c zride_auth;

-- Sample auth user (for testing)
INSERT INTO users (id, email, phone, password_hash, is_verified)
VALUES 
    ('550e8400-e29b-41d4-a716-446655440001', 'test@example.com', '+84123456789', '$2a$12$LQv3c1yqBWVHxkd0LHAkCOYz6TtxMQJqhN8/LewdBdXig2lY6W.Zm', true),
    ('550e8400-e29b-41d4-a716-446655440002', 'driver@example.com', '+84987654321', '$2a$12$LQv3c1yqBWVHxkd0LHAkCOYz6TtxMQJqhN8/LewdBdXig2lY6W.Zm', true);

\c zride_users;

-- Sample user profiles
INSERT INTO users (user_id, first_name, last_name, phone, gender, is_driver, verification_status)
VALUES 
    ('550e8400-e29b-41d4-a716-446655440001', 'John', 'Doe', '+84123456789', 'male', false, 'verified'),
    ('550e8400-e29b-41d4-a716-446655440002', 'Jane', 'Smith', '+84987654321', 'female', true, 'verified');

-- Sample vehicle for the driver
INSERT INTO vehicles (user_id, make, model, year, license_plate, color, vehicle_type, seats, status)
VALUES 
    ('550e8400-e29b-41d4-a716-446655440002', 'Toyota', 'Camry', 2020, '29A-12345', 'White', 'car', 4, 'active');

\c zride_trips;

-- Sample trip
INSERT INTO trips (passenger_id, driver_id, pickup_location, dropoff_location, status, fare_amount)
VALUES 
    ('550e8400-e29b-41d4-a716-446655440001', 
     '550e8400-e29b-41d4-a716-446655440002',
     '{"address": "District 1, Ho Chi Minh City", "lat": 10.7769, "lng": 106.7009}',
     '{"address": "District 3, Ho Chi Minh City", "lat": 10.7879, "lng": 106.6946}',
     'completed', 45000.00);

-- Add corresponding rating
\c zride_users;

INSERT INTO ratings (rated_user_id, rater_user_id, trip_id, rating, comment, rating_type)
VALUES 
    ('550e8400-e29b-41d4-a716-446655440002', 
     '550e8400-e29b-41d4-a716-446655440001',
     (SELECT id FROM zride_trips.trips LIMIT 1),
     5, 'Great driver, very professional!', 'driver');

COMMIT;