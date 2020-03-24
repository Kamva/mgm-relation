package mgmrel_test

import (
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
)

func TestIDField_PrepareID(t *testing.T) {
	id := primitive.NewObjectID()
	d := &Doc{}

	preparedId, err := d.PrepareID(id)
	assert.NoError(t, err)
	assert.Equal(t, id, preparedId)

	preparedId, err = d.PrepareID(id.Hex())
	assert.NoError(t, err)
	assert.Equal(t, id, preparedId)
}
