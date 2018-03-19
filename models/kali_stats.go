package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/satori/go.uuid"
)

// KaliStat contains a date and messagecount
type KaliStat struct {
	ID           uuid.UUID `json:"id" db:"id"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
	MessageCount int64     `json:"message_count" db:"message_count"`
	Date         time.Time `json:"date" db:"date"`
}

// String is not required by pop and may be deleted
func (k KaliStat) String() string {
	jk, _ := json.Marshal(k)
	return string(jk)
}

// KaliStats is not required by pop and may be deleted
type KaliStats []KaliStat

// String is not required by pop and may be deleted
func (k KaliStats) String() string {
	jk, _ := json.Marshal(k)
	return string(jk)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (k *KaliStat) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.TimeIsPresent{Field: k.Date, Name: "Date"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (k *KaliStat) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (k *KaliStat) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
