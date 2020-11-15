package mgmrel_test

import (
	mgmrel "github.com/kamva/mgm-relation"
	"github.com/kamva/mgm/v3"
	f "github.com/kamva/mgm/v3/field"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"testing"
)

func insertHasOneRelation(t *testing.T) (*Doc, *DocAuthor) {
	d := NewDoc("Ali", 12)
	require.NoError(t, mgm.Coll(d).Create(d))
	author := NewDocAuthor("Reza", d.ID)
	require.NoError(t, mgmrel.HasOne(d, author).Sync(author))
	return d, author
}

func TestHasOneRelation_Get_Empty(t *testing.T) {
	setupDefConnection()
	resetCollection()
	d := NewDoc("Ali", 12)
	require.NoError(t, mgm.Coll(d).Create(d))

	foundAuthor := &DocAuthor{}
	require.Equal(t, mongo.ErrNoDocuments, mgmrel.HasOne(d, &DocAuthor{}).Get(foundAuthor))
}

func TestHasOneRelation_Get(t *testing.T) {
	setupDefConnection()
	resetCollection()
	d, author := insertHasOneRelation(t)
	foundAuthor := &DocAuthor{}
	require.NoError(t, mgmrel.HasOne(d, &DocAuthor{}).Get(foundAuthor))
	require.Equal(t, author.ID, foundAuthor.ID)
	require.Equal(t, author.DocID, foundAuthor.DocID)
	require.Equal(t, author.Name, foundAuthor.Name)
}

func TestHasOneRelation_Sync_Insert(t *testing.T) {
	setupDefConnection()
	resetCollection()
	d, author := insertHasOneRelation(t)

	c, err := mgm.Coll(d).CountDocuments(nil, bson.M{})
	require.NoError(t, err)
	require.Equal(t, int64(1), c)

	ca, err := mgm.Coll(author).CountDocuments(nil, bson.M{})
	require.NoError(t, err)
	require.Equal(t, int64(1), ca)
}

func TestHasOneRelation_Sync_DeleteEmpty(t *testing.T) {
	setupDefConnection()
	resetCollection()

	d := NewDoc("Ali", 12)
	require.NoError(t, mgm.Coll(d).Create(d))
	require.NoError(t, mgmrel.HasOne(d, &DocAuthor{}).Sync(nil))

	c, err := mgm.Coll(d).CountDocuments(nil, bson.M{})
	require.NoError(t, err)
	require.Equal(t, int64(1), c)

	ca, err := mgm.Coll(&DocAuthor{}).CountDocuments(nil, bson.M{})
	require.NoError(t, err)
	require.Equal(t, int64(0), ca)
}

func TestHasOneRelation_Sync_DeleteEmpty_Does_Not_Affect_Other(t *testing.T) {
	setupDefConnection()
	resetCollection()

	unrelatedAuthor := NewDocAuthor("unrelated", primitive.NewObjectID())
	require.NoError(t, mgm.Coll(unrelatedAuthor).Create(unrelatedAuthor))

	d := NewDoc("Ali", 12)
	require.NoError(t, mgm.Coll(d).Create(d))
	require.NoError(t, mgmrel.HasOne(d, &DocAuthor{}).Sync(nil))

	results := make([]*DocAuthor, 0)
	require.NoError(t, mgm.Coll(&DocAuthor{}).SimpleFind(&results, bson.M{}))
	require.Equal(t, 1, len(results))
	require.Equal(t, unrelatedAuthor.ID, results[0].ID)
	require.Equal(t, unrelatedAuthor.DocID, results[0].DocID)
	require.Equal(t, unrelatedAuthor.Name, results[0].Name)
}

func TestHasOneRelation_Sync_Delete(t *testing.T) {
	setupDefConnection()
	resetCollection()
	d, author := insertHasOneRelation(t)

	require.NoError(t, mgmrel.HasOne(d, author).Sync(nil))

	c, err := mgm.Coll(d).CountDocuments(nil, bson.M{})
	require.NoError(t, err)
	require.Equal(t, int64(1), c)

	ca, err := mgm.Coll(author).CountDocuments(nil, bson.M{})
	require.NoError(t, err)
	require.Equal(t, int64(0), ca)

	dd := &Doc{}
	err = mgm.Coll(d).First(bson.M{}, dd)
	require.NoError(t, err)
	require.Equal(t, d.ID, dd.ID)
	require.Equal(t, d.Name, dd.Name)
	require.Equal(t, d.Age, dd.Age)
}

func TestHasOneRelation_Sync_Delete_Does_Not_Affect_Other(t *testing.T) {
	setupDefConnection()
	resetCollection()

	unrelatedAuthor := NewDocAuthor("unrelated", primitive.NewObjectID())
	require.NoError(t, mgm.Coll(unrelatedAuthor).Create(unrelatedAuthor))

	d, author := insertHasOneRelation(t)
	require.NoError(t, mgmrel.HasOne(d, author).Sync(nil))

	results := make([]*DocAuthor, 0)
	require.NoError(t, mgm.Coll(author).SimpleFind(&results, bson.M{}))
	require.Equal(t, 1, len(results))
	require.Equal(t, unrelatedAuthor.ID, results[0].ID)
	require.Equal(t, unrelatedAuthor.DocID, results[0].DocID)
	require.Equal(t, unrelatedAuthor.Name, results[0].Name)
}

func TestHasOneRelation_Sync_MultiAdd(t *testing.T) {
	setupDefConnection()
	resetCollection()
	d, author := insertHasOneRelation(t)

	extraAuthor := NewDocAuthor("extra", d.ID)
	require.NoError(t, mgm.Coll(author).Create(extraAuthor))

	newAuthor := NewDocAuthor("Omid", d.ID)
	require.NoError(t, mgmrel.HasOne(d, author).Sync(newAuthor))
	c, err := mgm.Coll(d).CountDocuments(nil, bson.M{})
	require.NoError(t, err)
	require.Equal(t, int64(1), c)

	ca, err := mgm.Coll(author).CountDocuments(nil, bson.M{})
	require.NoError(t, err)
	require.Equal(t, int64(1), ca)

	foundDAuthor := &DocAuthor{}
	err = mgm.Coll(author).First(bson.M{}, foundDAuthor)
	require.NoError(t, err)
	require.Equal(t, newAuthor.ID, foundDAuthor.ID)
	require.Equal(t, newAuthor.DocID, foundDAuthor.DocID)
	require.Equal(t, newAuthor.Name, foundDAuthor.Name)
}

func TestHasOneRelation_Sync_Update(t *testing.T) {
	setupDefConnection()
	resetCollection()
	d, author := insertHasOneRelation(t)
	author.Name = "Haamed"
	require.NoError(t, mgmrel.HasOne(d, author).Sync(author))
	c, err := mgm.Coll(d).CountDocuments(nil, bson.M{})
	require.NoError(t, err)
	require.Equal(t, int64(1), c)

	authorCount, err := mgm.Coll(author).CountDocuments(nil, bson.M{})
	require.NoError(t, err)
	require.Equal(t, int64(1), authorCount)

	foundDAuthor := &DocAuthor{}
	err = mgm.Coll(author).First(bson.M{}, foundDAuthor)
	require.NoError(t, err)
	require.Equal(t, author.ID, foundDAuthor.ID)
	require.Equal(t, author.DocID, foundDAuthor.DocID)
	require.Equal(t, author.Name, foundDAuthor.Name)
}

func TestHasOneRelation_Sync_Update_Does_Not_Affect_Other(t *testing.T) {
	setupDefConnection()
	resetCollection()

	unrelatedAuthor := NewDocAuthor("unrelated", primitive.NewObjectID())
	require.NoError(t, mgm.Coll(unrelatedAuthor).Create(unrelatedAuthor))

	d, author := insertHasOneRelation(t)
	author.Name = "Haamed"
	require.NoError(t, mgmrel.HasOne(d, author).Sync(author))

	foundAuthor := &DocAuthor{}
	require.NoError(t, mgm.Coll(&DocAuthor{}).First(bson.M{f.ID: unrelatedAuthor.ID}, foundAuthor))
	assert.Equal(t, unrelatedAuthor.ID, foundAuthor.ID)
	assert.Equal(t, unrelatedAuthor.DocID, foundAuthor.DocID)
	assert.Equal(t, unrelatedAuthor.Name, foundAuthor.Name)
}
