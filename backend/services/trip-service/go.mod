module github.com/southern-martin/zride/backend/services/trip-service

go 1.21

require (
	github.com/google/uuid v1.6.0
	github.com/lib/pq v1.10.9
	github.com/southern-martin/zride/backend/shared v0.0.0-00010101000000-000000000000
)

replace github.com/southern-martin/zride/backend/shared => ../../shared
