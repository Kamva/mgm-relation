package mgmrel

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// IDField struct contain model's ID field. it also implements the SyncingHook
// to set id before sync the model. you can use this IDField instead of the mgm
// IDField.
type IDField struct {
	ID primitive.ObjectID `json:"id" bson:"_id,omitempty"`
}

// PrepareID method prepare id value to using it as id in filtering,...
// e.g convert hex-string id value to bson.ObjectId
func (f *IDField) PrepareID(id interface{}) (interface{}, error) {
	if idStr, ok := id.(string); ok {
		return primitive.ObjectIDFromHex(idStr)
	}

	// Otherwise id must be ObjectId
	return id, nil
}

// GetID method return model's id
func (f *IDField) GetID() interface{} {
	return f.ID
}

// SetID set id value of model's id field.
func (f *IDField) SetID(id interface{}) {
	f.ID = id.(primitive.ObjectID)
}

// Syncing set the ID if it's zero(empty ID).
func (f *IDField) Syncing() error {
	if f.ID.IsZero() {
		f.ID = primitive.NewObjectID()
	}
	return nil
}

var _ SyncingHook = &IDField{}
