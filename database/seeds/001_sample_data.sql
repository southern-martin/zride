-- ZRide Database Seed Data
-- Test data for development and testing

-- Insert test users
INSERT INTO users (id, zalo_id, name, phone, email, avatar, is_active) VALUES
('550e8400-e29b-41d4-a716-446655440001', 'zalo_001', 'Nguyễn Văn Tài Xế', '0901234567', 'taixe@example.com', 'https://example.com/avatar1.jpg', true),
('550e8400-e29b-41d4-a716-446655440002', 'zalo_002', 'Trần Thị Hành Khách', '0912345678', 'hanhkhach@example.com', 'https://example.com/avatar2.jpg', true),
('550e8400-e29b-41d4-a716-446655440003', 'zalo_003', 'Lê Minh Tài Xế 2', '0923456789', 'taixe2@example.com', 'https://example.com/avatar3.jpg', true),
('550e8400-e29b-41d4-a716-446655440004', 'zalo_004', 'Phạm Thị Khách 2', '0934567890', 'khach2@example.com', 'https://example.com/avatar4.jpg', true),
('550e8400-e29b-41d4-a716-446655440005', 'zalo_005', 'Hoàng Văn Driver', '0945678901', 'driver@example.com', 'https://example.com/avatar5.jpg', true);

-- Insert driver profiles
INSERT INTO driver_profiles (
    id, user_id, license_number, license_expiry, vehicle_type, vehicle_brand, 
    vehicle_model, vehicle_year, vehicle_plate, vehicle_color, vehicle_capacity, 
    is_verified, rating, total_trips
) VALUES
(
    '660e8400-e29b-41d4-a716-446655440001', 
    '550e8400-e29b-41d4-a716-446655440001',
    'B1-123456789', 
    '2025-12-31',
    'sedan', 
    'Toyota', 
    'Vios', 
    2020, 
    '51A-12345', 
    'white', 
    4, 
    true, 
    4.5, 
    25
),
(
    '660e8400-e29b-41d4-a716-446655440003', 
    '550e8400-e29b-41d4-a716-446655440003',
    'B1-987654321', 
    '2026-06-30',
    'suv', 
    'Honda', 
    'CR-V', 
    2021, 
    '51B-67890', 
    'black', 
    7, 
    true, 
    4.8, 
    42
),
(
    '660e8400-e29b-41d4-a716-446655440005', 
    '550e8400-e29b-41d4-a716-446655440005',
    'B2-555666777', 
    '2025-08-15',
    'minivan', 
    'Hyundai', 
    'Starex', 
    2019, 
    '52A-11111', 
    'silver', 
    9, 
    true, 
    4.2, 
    18
);

-- Insert sample trips (Saigon to nearby provinces)
INSERT INTO trips (
    id, driver_id, 
    origin_latitude, origin_longitude, origin_address,
    destination_latitude, destination_longitude, destination_address,
    departure_time, available_seats, price_per_seat, 
    total_distance, estimated_duration, status
) VALUES
(
    '770e8400-e29b-41d4-a716-446655440001',
    '550e8400-e29b-41d4-a716-446655440001',
    10.762622, 106.660172, 'Quận 1, TP. Hồ Chí Minh',
    11.312777, 106.477222, 'Thành phố Bình Dương',
    '2025-10-16 07:00:00+07',
    3,
    80000,
    35.5,
    45,
    'pending'
),
(
    '770e8400-e29b-41d4-a716-446655440002',
    '550e8400-e29b-41d4-a716-446655440003',
    10.762622, 106.660172, 'Quận 1, TP. Hồ Chí Minh',
    10.950000, 106.816667, 'Thành phố Thủ Đức',
    '2025-10-16 08:30:00+07',
    5,
    50000,
    25.2,
    35,
    'pending'
),
(
    '770e8400-e29b-41d4-a716-446655440003',
    '550e8400-e29b-41d4-a716-446655440005',
    10.762622, 106.660172, 'Quận 1, TP. Hồ Chí Minh',
    11.533333, 107.183333, 'Tỉnh Bình Phước',
    '2025-10-16 06:00:00+07',
    7,
    150000,
    120.0,
    180,
    'pending'
),
(
    '770e8400-e29b-41d4-a716-446655440004',
    '550e8400-e29b-41d4-a716-446655440001',
    11.312777, 106.477222, 'Thành phố Bình Dương',
    10.762622, 106.660172, 'Quận 1, TP. Hồ Chí Minh',
    '2025-10-16 17:30:00+07',
    2,
    80000,
    35.5,
    45,
    'pending'
);

-- Insert sample bookings
INSERT INTO bookings (
    id, trip_id, passenger_id, seats_booked, 
    pickup_latitude, pickup_longitude, pickup_address,
    dropoff_latitude, dropoff_longitude, dropoff_address,
    total_amount, status, payment_status
) VALUES
(
    '880e8400-e29b-41d4-a716-446655440001',
    '770e8400-e29b-41d4-a716-446655440001',
    '550e8400-e29b-41d4-a716-446655440002',
    1,
    10.762622, 106.660172, 'Quận 1, TP. Hồ Chí Minh',
    11.312777, 106.477222, 'Thành phố Bình Dương',
    80000,
    'confirmed',
    'paid'
),
(
    '880e8400-e29b-41d4-a716-446655440002',
    '770e8400-e29b-41d4-a716-446655440002',
    '550e8400-e29b-41d4-a716-446655440004',
    2,
    10.762622, 106.660172, 'Quận 1, TP. Hồ Chí Minh',
    10.950000, 106.816667, 'Thành phố Thủ Đức',
    100000,
    'pending',
    'pending'
);

-- Insert sample payments
INSERT INTO payments (
    id, booking_id, payer_id, recipient_id, 
    amount, payment_method, status, net_amount
) VALUES
(
    '990e8400-e29b-41d4-a716-446655440001',
    '880e8400-e29b-41d4-a716-446655440001',
    '550e8400-e29b-41d4-a716-446655440002',
    '550e8400-e29b-41d4-a716-446655440001',
    80000,
    'zalopay',
    'completed',
    76000
);

-- Insert sample ratings
INSERT INTO ratings (
    id, trip_id, booking_id, rater_id, rated_id, 
    rating, comment, rating_type
) VALUES
(
    'aa0e8400-e29b-41d4-a716-446655440001',
    '770e8400-e29b-41d4-a716-446655440001',
    '880e8400-e29b-41d4-a716-446655440001',
    '550e8400-e29b-41d4-a716-446655440002',
    '550e8400-e29b-41d4-a716-446655440001',
    5,
    'Tài xế thân thiện, lái xe an toàn. Xe sạch sẽ.',
    'driver_rating'
);

-- Insert sample notifications
INSERT INTO notifications (
    id, user_id, title, message, notification_type, 
    related_entity_id, related_entity_type, is_read
) VALUES
(
    'bb0e8400-e29b-41d4-a716-446655440001',
    '550e8400-e29b-41d4-a716-446655440001',
    'Booking Confirmed',
    'Bạn có một booking mới từ Trần Thị Hành Khách',
    'booking_confirmed',
    '880e8400-e29b-41d4-a716-446655440001',
    'booking',
    false
),
(
    'bb0e8400-e29b-41d4-a716-446655440002',
    '550e8400-e29b-41d4-a716-446655440002',
    'Trip Confirmed',
    'Chuyến đi của bạn đã được xác nhận',
    'trip_confirmed',
    '770e8400-e29b-41d4-a716-446655440001',
    'trip',
    true
);

-- Insert sample trip matching requests
INSERT INTO trip_matching_requests (
    id, user_id, request_type, 
    origin_latitude, origin_longitude, origin_address,
    destination_latitude, destination_longitude, destination_address,
    preferred_departure_time, max_price_per_seat, required_seats, status
) VALUES
(
    'cc0e8400-e29b-41d4-a716-446655440001',
    '550e8400-e29b-41d4-a716-446655440004',
    'find_trip',
    10.762622, 106.660172, 'Quận 1, TP. Hồ Chí Minh',
    11.312777, 106.477222, 'Thành phố Bình Dương',
    '2025-10-16 07:30:00+07',
    85000,
    1,
    'active'
);

-- Update occupied seats for trips that have bookings
UPDATE trips SET occupied_seats = (
    SELECT COALESCE(SUM(seats_booked), 0) 
    FROM bookings 
    WHERE bookings.trip_id = trips.id 
    AND bookings.status IN ('confirmed', 'in_progress', 'completed')
);

-- Update available seats
UPDATE trips SET available_seats = vehicle_capacity - occupied_seats
FROM driver_profiles 
WHERE trips.driver_id = driver_profiles.user_id;

-- Update user ratings based on their ratings as drivers
UPDATE users SET 
    last_login_at = CURRENT_TIMESTAMP - INTERVAL '1 day'
WHERE id IN (
    '550e8400-e29b-41d4-a716-446655440001',
    '550e8400-e29b-41d4-a716-446655440002'
);

COMMIT;