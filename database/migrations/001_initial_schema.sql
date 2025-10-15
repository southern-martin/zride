-- ZRide Database Schema
-- Initial migration for all tables

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Users table (Auth Service)
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    zalo_id VARCHAR(255) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    phone VARCHAR(20),
    email VARCHAR(255),
    avatar TEXT,
    is_active BOOLEAN DEFAULT TRUE,
    last_login_at TIMESTAMP WITH TIME ZONE,
    refresh_token TEXT,
    version INTEGER DEFAULT 1,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Auth sessions table
CREATE TABLE auth_sessions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    access_token TEXT NOT NULL,
    refresh_token TEXT NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    device_info TEXT,
    ip_address INET,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Driver profiles table (User Service)
CREATE TABLE driver_profiles (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    license_number VARCHAR(50) NOT NULL,
    license_expiry DATE,
    vehicle_type VARCHAR(50) NOT NULL,
    vehicle_brand VARCHAR(100),
    vehicle_model VARCHAR(100),
    vehicle_year INTEGER,
    vehicle_plate VARCHAR(20) NOT NULL,
    vehicle_color VARCHAR(50),
    vehicle_capacity INTEGER NOT NULL DEFAULT 4,
    is_verified BOOLEAN DEFAULT FALSE,
    verification_documents JSONB,
    rating DECIMAL(3,2) DEFAULT 0.0,
    total_trips INTEGER DEFAULT 0,
    total_earnings DECIMAL(12,2) DEFAULT 0.0,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Trips table (Trip Service)
CREATE TABLE trips (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    driver_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    origin_latitude DECIMAL(10,8) NOT NULL,
    origin_longitude DECIMAL(11,8) NOT NULL,
    origin_address TEXT NOT NULL,
    destination_latitude DECIMAL(10,8) NOT NULL,
    destination_longitude DECIMAL(11,8) NOT NULL,
    destination_address TEXT NOT NULL,
    departure_time TIMESTAMP WITH TIME ZONE NOT NULL,
    arrival_time TIMESTAMP WITH TIME ZONE,
    available_seats INTEGER NOT NULL DEFAULT 1,
    occupied_seats INTEGER DEFAULT 0,
    price_per_seat DECIMAL(10,2) NOT NULL,
    total_distance DECIMAL(8,2), -- in kilometers
    estimated_duration INTEGER, -- in minutes
    trip_type VARCHAR(50) DEFAULT 'return_trip', -- return_trip, scheduled, regular
    status VARCHAR(50) DEFAULT 'pending', -- pending, active, completed, cancelled
    cancellation_reason TEXT,
    notes TEXT,
    vehicle_info JSONB,
    route_waypoints JSONB, -- Array of lat/lng waypoints
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Bookings table (Trip Service)
CREATE TABLE bookings (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    trip_id UUID NOT NULL REFERENCES trips(id) ON DELETE CASCADE,
    passenger_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    seats_booked INTEGER NOT NULL DEFAULT 1,
    pickup_latitude DECIMAL(10,8),
    pickup_longitude DECIMAL(11,8),
    pickup_address TEXT,
    dropoff_latitude DECIMAL(10,8),
    dropoff_longitude DECIMAL(11,8),
    dropoff_address TEXT,
    total_amount DECIMAL(10,2) NOT NULL,
    booking_fee DECIMAL(10,2) DEFAULT 0.0,
    status VARCHAR(50) DEFAULT 'pending', -- pending, confirmed, in_progress, completed, cancelled
    payment_status VARCHAR(50) DEFAULT 'pending', -- pending, paid, refunded, failed
    payment_id UUID,
    pickup_time TIMESTAMP WITH TIME ZONE,
    dropoff_time TIMESTAMP WITH TIME ZONE,
    cancellation_reason TEXT,
    passenger_notes TEXT,
    driver_notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Ratings table (User Service)
CREATE TABLE ratings (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    trip_id UUID NOT NULL REFERENCES trips(id),
    booking_id UUID REFERENCES bookings(id),
    rater_id UUID NOT NULL REFERENCES users(id),
    rated_id UUID NOT NULL REFERENCES users(id), -- driver or passenger
    rating INTEGER NOT NULL CHECK (rating >= 1 AND rating <= 5),
    comment TEXT,
    rating_type VARCHAR(20) NOT NULL, -- driver_rating, passenger_rating
    is_anonymous BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Payments table (Payment Service)
CREATE TABLE payments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    booking_id UUID NOT NULL REFERENCES bookings(id),
    payer_id UUID NOT NULL REFERENCES users(id),
    recipient_id UUID NOT NULL REFERENCES users(id),
    amount DECIMAL(10,2) NOT NULL,
    currency VARCHAR(3) DEFAULT 'VND',
    payment_method VARCHAR(50) NOT NULL, -- zalopay, momo, bank_transfer, cash
    external_transaction_id VARCHAR(255),
    status VARCHAR(50) DEFAULT 'pending', -- pending, processing, completed, failed, refunded
    gateway_response JSONB,
    fee_amount DECIMAL(10,2) DEFAULT 0.0,
    net_amount DECIMAL(10,2) NOT NULL,
    processed_at TIMESTAMP WITH TIME ZONE,
    refunded_at TIMESTAMP WITH TIME ZONE,
    refund_reason TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Notifications table (shared across services)
CREATE TABLE notifications (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    message TEXT NOT NULL,
    notification_type VARCHAR(50) NOT NULL, -- trip_update, booking_confirmed, payment_success, etc.
    related_entity_id UUID, -- trip_id, booking_id, payment_id, etc.
    related_entity_type VARCHAR(50), -- trip, booking, payment, etc.
    is_read BOOLEAN DEFAULT FALSE,
    is_sent BOOLEAN DEFAULT FALSE,
    sent_at TIMESTAMP WITH TIME ZONE,
    delivery_method VARCHAR(20) DEFAULT 'push', -- push, sms, email, zalo_oa
    metadata JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Messages table (for driver-passenger communication)
CREATE TABLE messages (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    trip_id UUID NOT NULL REFERENCES trips(id),
    booking_id UUID REFERENCES bookings(id),
    sender_id UUID NOT NULL REFERENCES users(id),
    recipient_id UUID NOT NULL REFERENCES users(id),
    message_text TEXT NOT NULL,
    message_type VARCHAR(20) DEFAULT 'text', -- text, location, image
    attachment_url TEXT,
    is_read BOOLEAN DEFAULT FALSE,
    read_at TIMESTAMP WITH TIME ZONE,
    is_system_message BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Trip matching requests (AI Matching Service)
CREATE TABLE trip_matching_requests (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id),
    request_type VARCHAR(20) NOT NULL, -- find_trip, find_passengers
    origin_latitude DECIMAL(10,8) NOT NULL,
    origin_longitude DECIMAL(11,8) NOT NULL,
    origin_address TEXT NOT NULL,
    destination_latitude DECIMAL(10,8) NOT NULL,
    destination_longitude DECIMAL(11,8) NOT NULL,
    destination_address TEXT NOT NULL,
    preferred_departure_time TIMESTAMP WITH TIME ZONE,
    flexible_time_range INTEGER DEFAULT 60, -- minutes
    max_price_per_seat DECIMAL(10,2),
    required_seats INTEGER DEFAULT 1,
    preferences JSONB, -- smoking, music, etc.
    status VARCHAR(20) DEFAULT 'active', -- active, fulfilled, expired, cancelled
    expires_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Matching results cache (AI Matching Service)
CREATE TABLE matching_results (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    request_id UUID NOT NULL REFERENCES trip_matching_requests(id),
    trip_id UUID REFERENCES trips(id),
    match_score DECIMAL(5,4), -- 0.0 to 1.0
    distance_score DECIMAL(5,4),
    time_score DECIMAL(5,4),
    price_score DECIMAL(5,4),
    rating_score DECIMAL(5,4),
    estimated_pickup_time TIMESTAMP WITH TIME ZONE,
    estimated_dropoff_time TIMESTAMP WITH TIME ZONE,
    additional_distance DECIMAL(8,2), -- km
    additional_time INTEGER, -- minutes
    is_presented BOOLEAN DEFAULT FALSE,
    presented_at TIMESTAMP WITH TIME ZONE,
    user_action VARCHAR(20), -- viewed, interested, booked, dismissed
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for better performance
CREATE INDEX idx_users_zalo_id ON users(zalo_id);
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_phone ON users(phone);
CREATE INDEX idx_users_active ON users(is_active);

CREATE INDEX idx_auth_sessions_user_id ON auth_sessions(user_id);
CREATE INDEX idx_auth_sessions_access_token ON auth_sessions(access_token);
CREATE INDEX idx_auth_sessions_refresh_token ON auth_sessions(refresh_token);
CREATE INDEX idx_auth_sessions_active ON auth_sessions(is_active);

CREATE INDEX idx_driver_profiles_user_id ON driver_profiles(user_id);
CREATE INDEX idx_driver_profiles_active ON driver_profiles(is_active);
CREATE INDEX idx_driver_profiles_verified ON driver_profiles(is_verified);

CREATE INDEX idx_trips_driver_id ON trips(driver_id);
CREATE INDEX idx_trips_status ON trips(status);
CREATE INDEX idx_trips_departure_time ON trips(departure_time);
CREATE INDEX idx_trips_origin_coords ON trips(origin_latitude, origin_longitude);
CREATE INDEX idx_trips_destination_coords ON trips(destination_latitude, destination_longitude);

CREATE INDEX idx_bookings_trip_id ON bookings(trip_id);
CREATE INDEX idx_bookings_passenger_id ON bookings(passenger_id);
CREATE INDEX idx_bookings_status ON bookings(status);
CREATE INDEX idx_bookings_payment_status ON bookings(payment_status);

CREATE INDEX idx_ratings_trip_id ON ratings(trip_id);
CREATE INDEX idx_ratings_rater_id ON ratings(rater_id);
CREATE INDEX idx_ratings_rated_id ON ratings(rated_id);

CREATE INDEX idx_payments_booking_id ON payments(booking_id);
CREATE INDEX idx_payments_payer_id ON payments(payer_id);
CREATE INDEX idx_payments_status ON payments(status);

CREATE INDEX idx_notifications_user_id ON notifications(user_id);
CREATE INDEX idx_notifications_is_read ON notifications(is_read);
CREATE INDEX idx_notifications_type ON notifications(notification_type);

CREATE INDEX idx_messages_trip_id ON messages(trip_id);
CREATE INDEX idx_messages_booking_id ON messages(booking_id);
CREATE INDEX idx_messages_sender_id ON messages(sender_id);
CREATE INDEX idx_messages_recipient_id ON messages(recipient_id);

CREATE INDEX idx_trip_matching_requests_user_id ON trip_matching_requests(user_id);
CREATE INDEX idx_trip_matching_requests_status ON trip_matching_requests(status);
CREATE INDEX idx_trip_matching_requests_coords ON trip_matching_requests(origin_latitude, origin_longitude, destination_latitude, destination_longitude);

CREATE INDEX idx_matching_results_request_id ON matching_results(request_id);
CREATE INDEX idx_matching_results_trip_id ON matching_results(trip_id);
CREATE INDEX idx_matching_results_score ON matching_results(match_score);

-- Update triggers for updated_at columns
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_auth_sessions_updated_at BEFORE UPDATE ON auth_sessions FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_driver_profiles_updated_at BEFORE UPDATE ON driver_profiles FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_trips_updated_at BEFORE UPDATE ON trips FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_bookings_updated_at BEFORE UPDATE ON bookings FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_ratings_updated_at BEFORE UPDATE ON ratings FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_payments_updated_at BEFORE UPDATE ON payments FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_notifications_updated_at BEFORE UPDATE ON notifications FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_trip_matching_requests_updated_at BEFORE UPDATE ON trip_matching_requests FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Sample data constraints and additional indexes for geospatial queries (if needed)
-- CREATE INDEX idx_trips_origin_gist ON trips USING GIST (ll_to_earth(origin_latitude, origin_longitude));
-- CREATE INDEX idx_trips_destination_gist ON trips USING GIST (ll_to_earth(destination_latitude, destination_longitude));