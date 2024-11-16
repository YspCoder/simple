package mdu

import (
	"context"
	"github.com/google/uuid"
	"time"
)

type IDField struct {
	ID string `json:"id" bson:"_id,omitempty"`
}

// DateFields struct contains the `created_at` and `updated_at`
// fields that autofill when inserting or updating a model.
type DateFields struct {
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

type TenantIdField struct {
	TenantId string `json:"tenantId" bson:"tenantId,omitempty"`
}

// PrepareID method prepares the ID value to be used for filtering
// generates uuid if not given id is empty
func (f *IDField) PrepareID(id string) (string, error) {
	if id != "" {
		return id, nil
	}
	return uuid.NewString(), nil
}

// GetID method returns a model's ID
func (f *IDField) GetID() string {
	return f.ID
}

// SetID sets the value of a model's ID field.
func (f *IDField) SetID(id string) {
	f.ID = id
}

// GetTenantId method returns a model's TenantId
func (f *TenantIdField) GetTenantId() string {
	return f.TenantId
}

// SetTenantId sets the value of a model's TenantId field.
func (f *TenantIdField) SetTenantId(tenantId string) {
	f.TenantId = tenantId
}

//--------------------------------
// DateField methods
//--------------------------------

// Creating hook is used here to set the `created_at` field
// value when inserting a new model into the database.
func (f *DateFields) Creating(ctx context.Context) error {
	f.CreatedAt = time.Now().UTC()
	return nil
}

// Saving hook is used here to set the `updated_at` field
// value when creating or updating a model.
// TODO: get context as param the next version(4).
func (f *DateFields) Saving(ctx context.Context) error {
	f.UpdatedAt = time.Now().UTC()
	return nil
}
