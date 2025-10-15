// Package domain contains rating-related entities
package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Rating represents a user rating
type Rating struct {
	ID        uuid.UUID `json:"id" db:"id"`
	RaterID   uuid.UUID `json:"rater_id" db:"rater_id"`     // Who gave the rating
	RatedID   uuid.UUID `json:"rated_id" db:"rated_id"`     // Who received the rating
	TripID    uuid.UUID `json:"trip_id" db:"trip_id"`       // Associated trip
	Score     float64   `json:"score" db:"score"`           // 1.0 - 5.0
	Comment   string    `json:"comment" db:"comment"`       // Optional comment
	Tags      []string  `json:"tags" db:"tags"`             // Rating tags (friendly, punctual, clean, etc.)
	IsVisible bool      `json:"is_visible" db:"is_visible"` // Whether the rating is public
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	Version   int       `json:"version" db:"version"`
}

// NewRating creates a new rating
func NewRating(raterID, ratedID, tripID uuid.UUID, score float64, comment string) (*Rating, error) {
	if raterID == uuid.Nil {
		return nil, errors.New("rater ID is required")
	}
	if ratedID == uuid.Nil {
		return nil, errors.New("rated user ID is required")
	}
	if tripID == uuid.Nil {
		return nil, errors.New("trip ID is required")
	}
	if raterID == ratedID {
		return nil, errors.New("cannot rate yourself")
	}
	if score < 1.0 || score > 5.0 {
		return nil, errors.New("score must be between 1.0 and 5.0")
	}

	now := time.Now()
	return &Rating{
		ID:        uuid.New(),
		RaterID:   raterID,
		RatedID:   ratedID,
		TripID:    tripID,
		Score:     score,
		Comment:   comment,
		Tags:      make([]string, 0),
		IsVisible: true, // Default to visible
		CreatedAt: now,
		UpdatedAt: now,
		Version:   1,
	}, nil
}

// UpdateRating updates the rating score and comment
func (r *Rating) UpdateRating(score float64, comment string) error {
	if score < 1.0 || score > 5.0 {
		return errors.New("score must be between 1.0 and 5.0")
	}

	r.Score = score
	r.Comment = comment
	r.UpdatedAt = time.Now()
	return nil
}

// AddTag adds a tag to the rating
func (r *Rating) AddTag(tag string) error {
	if tag == "" {
		return errors.New("tag cannot be empty")
	}
	
	// Check if tag already exists
	for _, t := range r.Tags {
		if t == tag {
			return errors.New("tag already exists")
		}
	}
	
	r.Tags = append(r.Tags, tag)
	r.UpdatedAt = time.Now()
	return nil
}

// RemoveTag removes a tag from the rating
func (r *Rating) RemoveTag(tag string) {
	for i, t := range r.Tags {
		if t == tag {
			r.Tags = append(r.Tags[:i], r.Tags[i+1:]...)
			r.UpdatedAt = time.Now()
			break
		}
	}
}

// SetTags sets all tags for the rating
func (r *Rating) SetTags(tags []string) {
	r.Tags = tags
	r.UpdatedAt = time.Now()
}

// Hide hides the rating from public view
func (r *Rating) Hide() {
	r.IsVisible = false
	r.UpdatedAt = time.Now()
}

// Show shows the rating in public view
func (r *Rating) Show() {
	r.IsVisible = true
	r.UpdatedAt = time.Now()
}

// Validate validates the rating data
func (r *Rating) Validate() error {
	if r.RaterID == uuid.Nil {
		return errors.New("rater ID is required")
	}
	if r.RatedID == uuid.Nil {
		return errors.New("rated user ID is required")
	}
	if r.TripID == uuid.Nil {
		return errors.New("trip ID is required")
	}
	if r.RaterID == r.RatedID {
		return errors.New("cannot rate yourself")
	}
	if r.Score < 1.0 || r.Score > 5.0 {
		return errors.New("score must be between 1.0 and 5.0")
	}
	return nil
}

// GetID returns the rating ID as string
func (r *Rating) GetID() string {
	return r.ID.String()
}