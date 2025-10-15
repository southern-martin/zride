// Package domain contains the core business entities and value objectspackage domain

// This layer is independent of any external concerns and contains pure business logic
package domain

import (
	"time"
	"github.com/google/uuid"
)

// Entity represents a domain entity with identity
type Entity struct {
	ID        uuid.UUID `json:"id" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// NewEntity creates a new entity with generated ID and timestamps
func NewEntity() Entity {
	now := time.Now()
	return Entity{
		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// UpdateTimestamp updates the entity's UpdatedAt field
func (e *Entity) UpdateTimestamp() {
	e.UpdatedAt = time.Now()
}

// ValueObject represents a domain value object
type ValueObject interface {
	Validate() error
	Equals(other ValueObject) bool
}

// AggregateRoot represents a domain aggregate root
type AggregateRoot interface {
	GetID() uuid.UUID
	GetVersion() int
	MarkAsModified()
}

// DomainEvent represents a domain event that occurred
type DomainEvent interface {
	GetEventID() uuid.UUID
	GetEventType() string
	GetAggregateID() uuid.UUID
	GetEventData() interface{}
	GetOccurredAt() time.Time
}

// BaseDomainEvent provides a base implementation for domain events
type BaseDomainEvent struct {
	EventID     uuid.UUID   `json:"event_id"`
	EventType   string      `json:"event_type"`
	AggregateID uuid.UUID   `json:"aggregate_id"`
	EventData   interface{} `json:"event_data"`
	OccurredAt  time.Time   `json:"occurred_at"`
}

func NewDomainEvent(eventType string, aggregateID uuid.UUID, eventData interface{}) *BaseDomainEvent {
	return &BaseDomainEvent{
		EventID:     uuid.New(),
		EventType:   eventType,
		AggregateID: aggregateID,
		EventData:   eventData,
		OccurredAt:  time.Now(),
	}
}

func (e *BaseDomainEvent) GetEventID() uuid.UUID     { return e.EventID }
func (e *BaseDomainEvent) GetEventType() string      { return e.EventType }
func (e *BaseDomainEvent) GetAggregateID() uuid.UUID { return e.AggregateID }
func (e *BaseDomainEvent) GetEventData() interface{} { return e.EventData }
func (e *BaseDomainEvent) GetOccurredAt() time.Time  { return e.OccurredAt }