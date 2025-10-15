-- Create extension for UUID generation
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "postgis";

-- Create match_requests table
CREATE TABLE IF NOT EXISTS match_requests (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    passenger_id UUID NOT NULL,
    pickup_location JSONB NOT NULL,
    dropoff_location JSONB NOT NULL,
    request_time TIMESTAMP WITH TIME ZONE NOT NULL,
    max_wait_time BIGINT NOT NULL DEFAULT 600000000000, -- 10 minutes in nanoseconds
    preferred_car_type VARCHAR(50),
    max_distance DECIMAL(10,2) NOT NULL DEFAULT 15.0,
    price_range JSONB,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Create indexes
CREATE INDEX idx_match_requests_passenger_id ON match_requests(passenger_id);
CREATE INDEX idx_match_requests_status ON match_requests(status);
CREATE INDEX idx_match_requests_created_at ON match_requests(created_at);
CREATE INDEX idx_match_requests_pickup_location ON match_requests USING GIN(pickup_location);

-- Create drivers table
CREATE TABLE IF NOT EXISTS drivers (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL UNIQUE,
    current_location JSONB NOT NULL,
    is_available BOOLEAN NOT NULL DEFAULT false,
    car_type VARCHAR(50) NOT NULL,
    rating DECIMAL(3,2) DEFAULT 5.0,
    completed_trips INTEGER DEFAULT 0,
    last_active_time TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    max_distance DECIMAL(10,2) NOT NULL DEFAULT 15.0,
    preferred_areas JSONB DEFAULT '[]'::jsonb
);

-- Create indexes for drivers
CREATE INDEX idx_drivers_user_id ON drivers(user_id);
CREATE INDEX idx_drivers_is_available ON drivers(is_available);
CREATE INDEX idx_drivers_current_location ON drivers USING GIN(current_location);
CREATE INDEX idx_drivers_last_active_time ON drivers(last_active_time);

-- Create spatial index for driver locations
CREATE INDEX idx_drivers_location_spatial ON drivers 
USING GIST(ST_SetSRID(ST_MakePoint(
    (current_location->>'longitude')::float, 
    (current_location->>'latitude')::float
), 4326));

-- Create match_results table
CREATE TABLE IF NOT EXISTS match_results (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    match_request_id UUID NOT NULL REFERENCES match_requests(id) ON DELETE CASCADE,
    driver_id UUID NOT NULL REFERENCES drivers(id) ON DELETE CASCADE,
    score DECIMAL(5,4) NOT NULL,
    estimated_distance DECIMAL(10,2) NOT NULL,
    estimated_time BIGINT NOT NULL, -- Duration in nanoseconds
    estimated_price DECIMAL(12,2) NOT NULL,
    match_time TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Create indexes for match_results
CREATE INDEX idx_match_results_match_request_id ON match_results(match_request_id);
CREATE INDEX idx_match_results_driver_id ON match_results(driver_id);
CREATE INDEX idx_match_results_status ON match_results(status);
CREATE INDEX idx_match_results_score ON match_results(score DESC);
CREATE INDEX idx_match_results_created_at ON match_results(created_at);

-- Add constraints
ALTER TABLE match_requests ADD CONSTRAINT chk_status_values 
    CHECK (status IN ('pending', 'matched', 'expired', 'cancelled'));

ALTER TABLE match_results ADD CONSTRAINT chk_result_status_values 
    CHECK (status IN ('pending', 'accepted', 'rejected', 'expired'));

ALTER TABLE drivers ADD CONSTRAINT chk_rating_range 
    CHECK (rating >= 0 AND rating <= 5);

ALTER TABLE drivers ADD CONSTRAINT chk_car_type_values 
    CHECK (car_type IN ('motorbike', 'car_4_seat', 'car_7_seat'));

-- Insert sample drivers for testing
INSERT INTO drivers (id, user_id, current_location, is_available, car_type, rating, completed_trips, last_active_time, max_distance) VALUES
    (uuid_generate_v4(), uuid_generate_v4(), '{"latitude": 10.7769, "longitude": 106.7009, "address": "District 1, Ho Chi Minh City"}', true, 'car_4_seat', 4.8, 150, NOW(), 20.0),
    (uuid_generate_v4(), uuid_generate_v4(), '{"latitude": 10.7829, "longitude": 106.6926, "address": "District 3, Ho Chi Minh City"}', true, 'motorbike', 4.6, 85, NOW(), 15.0),
    (uuid_generate_v4(), uuid_generate_v4(), '{"latitude": 10.7626, "longitude": 106.6822, "address": "District 5, Ho Chi Minh City"}', true, 'car_7_seat', 4.9, 220, NOW(), 25.0),
    (uuid_generate_v4(), uuid_generate_v4(), '{"latitude": 10.8003, "longitude": 106.6593, "address": "Tan Binh District, Ho Chi Minh City"}', true, 'car_4_seat', 4.7, 95, NOW(), 18.0),
    (uuid_generate_v4(), uuid_generate_v4(), '{"latitude": 10.7320, "longitude": 106.7019, "address": "District 7, Ho Chi Minh City"}', true, 'motorbike', 4.5, 60, NOW(), 12.0)
ON CONFLICT (user_id) DO NOTHING;