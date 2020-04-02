package mgmrel_test

import (
	"fmt"
	mgmrel "github.com/Kamva/mgm-relation"
	"github.com/Kamva/mgm/v3"
	f "github.com/Kamva/mgm/v3/field"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
)

func insertHasManyRelation(t *testing.T) (*Doc, []*DocAuthor) {
	d := NewDoc("A", 12)
	require.NoError(t, mgm.Coll(d).Create(d))

	authors := []*DocAuthor{
		NewDocAuthor("B1", d.ID),
		NewDocAuthor("B2", d.ID),
	}

	require.NoError(t, mgmrel.HasMany(d, &DocAuthor{}).Sync(authors))
	return d, authors
}

func insertHasManyRelationWithoutRemove(t *testing.T) (*Doc, []*DocAuthor) {
	d := NewDoc("A", 12)
	require.NoError(t, mgm.Coll(d).Create(d))

	authors := []*DocAuthor{
		NewDocAuthor("B1", d.ID),
		NewDocAuthor("B2", d.ID),
	}

	require.NoError(t, mgmrel.HasMany(d, &DocAuthor{}).SyncWithoutRemove(authors))
	return d, authors
}

func TestHasManyRelation_SimpleGet_Empty(t *testing.T) {
	setupDefConnection()
	resetCollection()
	d := NewDoc("Ali", 12)
	require.NoError(t, mgm.Coll(d).Create(d))

	foundAuthors := make([]*DocAuthor, 0)
	require.NoError(t, mgmrel.HasMany(d, &DocAuthor{}).SimpleGet(&foundAuthors, 10))
	assert.Equal(t, 0, len(foundAuthors))
}

func TestHasManyRelation_SimpleGetWithLimit(t *testing.T) {
	setupDefConnection()
	resetCollection()
	d, _ := insertHasManyRelation(t)
	foundAuthors := make([]*DocAuthor, 0)
	require.NoError(t, mgmrel.HasMany(d, &DocAuthor{}).SimpleGet(&foundAuthors, 1))

	assert.Equal(t, 1, len(foundAuthors))
}

func TestHasManyRelation_Get(t *testing.T) {
	setupDefConnection()
	resetCollection()
	d, authors := insertHasManyRelation(t)
	foundAuthors := make([]*DocAuthor, 0)
	require.NoError(t, mgmrel.HasMany(d, &DocAuthor{}).Get(&foundAuthors, "_id", 0, 2))

	assert.Equal(t, len(authors), len(foundAuthors))
	for i, author := range authors {
		foundAuthor := foundAuthors[i]
		assert.Equal(t, author.ID, foundAuthor.ID, fmt.Sprintf("author wiht index: %v", i))
		assert.Equal(t, author.DocID, foundAuthor.DocID, fmt.Sprintf("author wiht index: %v", i))
		assert.Equal(t, author.Name, foundAuthor.Name, fmt.Sprintf("author wiht index: %v", i))
	}
}

func TestHasManyRelation_GetWithLimit(t *testing.T) {
	setupDefConnection()
	resetCollection()
	d, _ := insertHasManyRelation(t)
	foundAuthors := make([]*DocAuthor, 0)
	require.NoError(t, mgmrel.HasMany(d, &DocAuthor{}).Get(&foundAuthors, "_id", 0, 1))

	assert.Equal(t, 1, len(foundAuthors))
}

func TestHasManyRelation_GetWithSort(t *testing.T) {
	setupDefConnection()
	resetCollection()
	d, authors := insertHasManyRelation(t)
	foundAuthors := make([]*DocAuthor, 0)
	require.NoError(t, mgmrel.HasMany(d, &DocAuthor{}).Get(&foundAuthors, "_id", 0, 1))

	assert.Equal(t, 1, len(foundAuthors))
	author := authors[0]
	foundAuthor := foundAuthors[0]
	assert.Equal(t, author.ID, foundAuthor.ID)
	assert.Equal(t, author.DocID, foundAuthor.DocID)
	assert.Equal(t, author.Name, foundAuthor.Name)
}

func TestHasManyRelation_SyncWithoutRemove_Insert(t *testing.T) {
	setupDefConnection()
	resetCollection()
	d, authors := insertHasManyRelationWithoutRemove(t)

	c, err := mgm.Coll(d).CountDocuments(nil, bson.M{})
	require.NoError(t, err)
	require.Equal(t, int64(1), c)

	ca, err := mgm.Coll(&DocAuthor{}).CountDocuments(nil, bson.M{})
	require.NoError(t, err)
	require.Equal(t, int64(len(authors)), ca)
}

func TestHasManyRelation_Sync_Insert(t *testing.T) {
	setupDefConnection()
	resetCollection()
	d, authors := insertHasManyRelation(t)

	c, err := mgm.Coll(d).CountDocuments(nil, bson.M{})
	require.NoError(t, err)
	require.Equal(t, int64(1), c)

	ca, err := mgm.Coll(&DocAuthor{}).CountDocuments(nil, bson.M{})
	require.NoError(t, err)
	require.Equal(t, int64(len(authors)), ca)
}

func TestHasManyRelation_Sync_DeleteNil(t *testing.T) {
	setupDefConnection()
	resetCollection()

	d := NewDoc("A", 12)
	require.NoError(t, mgm.Coll(d).Create(d))
	require.NoError(t, mgmrel.HasMany(d, &DocAuthor{}).Sync(nil))

	c, err := mgm.Coll(d).CountDocuments(nil, bson.M{})
	require.NoError(t, err)
	require.Equal(t, int64(1), c)

	ca, err := mgm.Coll(&DocAuthor{}).CountDocuments(nil, bson.M{})
	require.NoError(t, err)
	require.Equal(t, int64(0), ca)
}

func TestHasManyRelation_SyncWithoutRemove_DeleteNil(t *testing.T) {
	setupDefConnection()
	resetCollection()

	d := NewDoc("A", 12)
	require.NoError(t, mgm.Coll(d).Create(d))
	require.NoError(t, mgmrel.HasMany(d, &DocAuthor{}).SyncWithoutRemove(nil))

	c, err := mgm.Coll(d).CountDocuments(nil, bson.M{})
	require.NoError(t, err)
	require.Equal(t, int64(1), c)

	ca, err := mgm.Coll(&DocAuthor{}).CountDocuments(nil, bson.M{})
	require.NoError(t, err)
	require.Equal(t, int64(0), ca)
}

func TestHasManyRelation_Sync_DeleteEmpty(t *testing.T) {
	setupDefConnection()
	resetCollection()

	d := NewDoc("A", 12)
	require.NoError(t, mgm.Coll(d).Create(d))
	require.NoError(t, mgmrel.HasMany(d, &DocAuthor{}).Sync(make([]*DocAuthor, 0)))

	c, err := mgm.Coll(d).CountDocuments(nil, bson.M{})
	require.NoError(t, err)
	require.Equal(t, int64(1), c)

	ca, err := mgm.Coll(&DocAuthor{}).CountDocuments(nil, bson.M{})
	require.NoError(t, err)
	require.Equal(t, int64(0), ca)
}

func TestHasManyRelation_SyncWithoutRemove_DeleteEmpty(t *testing.T) {
	setupDefConnection()
	resetCollection()

	d := NewDoc("A", 12)
	require.NoError(t, mgm.Coll(d).Create(d))
	require.NoError(t, mgmrel.HasMany(d, &DocAuthor{}).SyncWithoutRemove(make([]*DocAuthor, 0)))

	c, err := mgm.Coll(d).CountDocuments(nil, bson.M{})
	require.NoError(t, err)
	require.Equal(t, int64(1), c)

	ca, err := mgm.Coll(&DocAuthor{}).CountDocuments(nil, bson.M{})
	require.NoError(t, err)
	require.Equal(t, int64(0), ca)
}

func TestHasManyRelation_Sync_DeleteEmpty_Does_Not_Affect_Other(t *testing.T) {
	setupDefConnection()
	resetCollection()

	unrelatedAuthor := NewDocAuthor("unrelated", primitive.NewObjectID())
	require.NoError(t, mgm.Coll(unrelatedAuthor).Create(unrelatedAuthor))

	d := NewDoc("Ali", 12)
	require.NoError(t, mgm.Coll(d).Create(d))
	require.NoError(t, mgmrel.HasMany(d, &DocAuthor{}).Sync(nil))

	results := make([]*DocAuthor, 0)
	require.NoError(t, mgm.Coll(&DocAuthor{}).SimpleFind(&results, bson.M{}))
	require.Equal(t, 1, len(results))
	require.Equal(t, unrelatedAuthor.ID, results[0].ID)
	require.Equal(t, unrelatedAuthor.DocID, results[0].DocID)
	require.Equal(t, unrelatedAuthor.Name, results[0].Name)
}

func TestHasManyRelation_Sync_Delete(t *testing.T) {
	setupDefConnection()
	resetCollection()
	d, _ := insertHasManyRelation(t)
	require.NoError(t, mgmrel.HasMany(d, &DocAuthor{}).Sync(nil))

	c, err := mgm.Coll(d).CountDocuments(nil, bson.M{})
	require.NoError(t, err)
	require.Equal(t, int64(1), c)

	ca, err := mgm.Coll(&DocAuthor{}).CountDocuments(nil, bson.M{})
	require.NoError(t, err)
	require.Equal(t, int64(0), ca)
}

func TestHasManyRelation_SyncWithoutRemove_Delete(t *testing.T) {
	setupDefConnection()
	resetCollection()
	d, authors := insertHasManyRelationWithoutRemove(t)
	require.NoError(t, mgmrel.HasMany(d, &DocAuthor{}).SyncWithoutRemove(nil))

	c, err := mgm.Coll(d).CountDocuments(nil, bson.M{})
	require.NoError(t, err)
	require.Equal(t, int64(1), c)

	ca, err := mgm.Coll(&DocAuthor{}).CountDocuments(nil, bson.M{})
	require.NoError(t, err)
	require.Equal(t, int64(len(authors)), ca)
}

func TestHasManyRelation_Sync_Delete_Does_Not_Affect_Other(t *testing.T) {
	setupDefConnection()
	resetCollection()

	unrelatedAuthor := NewDocAuthor("unrelated", primitive.NewObjectID())
	require.NoError(t, mgm.Coll(unrelatedAuthor).Create(unrelatedAuthor))

	d, _ := insertHasManyRelation(t)
	require.NoError(t, mgmrel.HasMany(d, &DocAuthor{}).Sync(nil))

	results := make([]*DocAuthor, 0)
	require.NoError(t, mgm.Coll(&DocAuthor{}).SimpleFind(&results, bson.M{}))
	require.Equal(t, 1, len(results))
	require.Equal(t, unrelatedAuthor.ID, results[0].ID)
	require.Equal(t, unrelatedAuthor.DocID, results[0].DocID)
	require.Equal(t, unrelatedAuthor.Name, results[0].Name)
}

func TestHasManyRelation_Sync_MultiAdd(t *testing.T) {
	setupDefConnection()
	resetCollection()
	d, authors := insertHasManyRelation(t)

	extraAuthor := NewDocAuthor("extra", d.ID)
	require.NoError(t, mgm.Coll(&DocAuthor{}).Create(extraAuthor))

	newAuthor := NewDocAuthor("Omid", d.ID)
	authors = append(authors, newAuthor)
	require.NoError(t, mgmrel.HasMany(d, &DocAuthor{}).Sync(authors))
	c, err := mgm.Coll(d).CountDocuments(nil, bson.M{})
	require.NoError(t, err)
	require.Equal(t, int64(1), c)

	foundDAuthors := make([]*DocAuthor, 0)
	require.NoError(t, mgm.Coll(&DocAuthor{}).SimpleFind(&foundDAuthors, bson.M{}))
	assert.Equal(t, len(authors), len(foundDAuthors))
	for i, author := range authors {
		foundAuthor := foundDAuthors[i]
		assert.Equal(t, author.ID, foundAuthor.ID)
		assert.Equal(t, author.DocID, foundAuthor.DocID)
		assert.Equal(t, author.Name, foundAuthor.Name)
	}
}

func TestHasManyRelation_SyncWithoutRemove_MultiAdd(t *testing.T) {
	setupDefConnection()
	resetCollection()
	d, authors := insertHasManyRelationWithoutRemove(t)

	extraAuthor := NewDocAuthor("extra", d.ID)
	authors = append(authors, extraAuthor)
	require.NoError(t, mgm.Coll(&DocAuthor{}).Create(extraAuthor))

	newAuthor := NewDocAuthor("Omid", d.ID)
	authors = append(authors, newAuthor)
	require.NoError(t, mgmrel.HasMany(d, &DocAuthor{}).SyncWithoutRemove(authors))
	c, err := mgm.Coll(d).CountDocuments(nil, bson.M{})
	require.NoError(t, err)
	require.Equal(t, int64(1), c)

	foundDAuthors := make([]*DocAuthor, 0)
	require.NoError(t, mgm.Coll(&DocAuthor{}).SimpleFind(&foundDAuthors, bson.M{}))
	assert.Equal(t, len(authors), len(foundDAuthors))
	for i, author := range authors {
		foundAuthor := foundDAuthors[i]
		assert.Equal(t, author.ID, foundAuthor.ID)
		assert.Equal(t, author.DocID, foundAuthor.DocID)
		assert.Equal(t, author.Name, foundAuthor.Name)
	}
}

func TestHasManyRelation_Sync_Update(t *testing.T) {
	setupDefConnection()
	resetCollection()
	d, authors := insertHasManyRelation(t)
	for i, author := range authors {
		author.Name = fmt.Sprintf("New-%v", i)
	}
	require.NoError(t, mgmrel.HasMany(d, &DocAuthor{}).Sync(authors))
	c, err := mgm.Coll(d).CountDocuments(nil, bson.M{})
	require.NoError(t, err)
	require.Equal(t, int64(1), c)

	foundDAuthors := make([]*DocAuthor, 0)
	require.NoError(t, mgm.Coll(&DocAuthor{}).SimpleFind(&foundDAuthors, bson.M{}))
	assert.Equal(t, len(authors), len(foundDAuthors))
	for i, author := range authors {
		foundAuthor := foundDAuthors[i]
		assert.Equal(t, author.ID, foundAuthor.ID)
		assert.Equal(t, author.DocID, foundAuthor.DocID)
		assert.Equal(t, author.Name, foundAuthor.Name)
	}
}

func TestHasManyRelation_SyncWithoutRemove_Update(t *testing.T) {
	setupDefConnection()
	resetCollection()
	d, authors := insertHasManyRelationWithoutRemove(t)
	for i, author := range authors {
		author.Name = fmt.Sprintf("New-%v", i)
	}
	require.NoError(t, mgmrel.HasMany(d, &DocAuthor{}).SyncWithoutRemove(authors))

	c, err := mgm.Coll(d).CountDocuments(nil, bson.M{})
	require.NoError(t, err)
	require.Equal(t, int64(1), c)

	foundDAuthors := make([]*DocAuthor, 0)
	require.NoError(t, mgm.Coll(&DocAuthor{}).SimpleFind(&foundDAuthors, bson.M{}))
	assert.Equal(t, len(authors), len(foundDAuthors))
	for i, author := range authors {
		foundAuthor := foundDAuthors[i]
		assert.Equal(t, author.ID, foundAuthor.ID)
		assert.Equal(t, author.DocID, foundAuthor.DocID)
		assert.Equal(t, author.Name, foundAuthor.Name)
	}
}

func TestHasManyRelation_Sync_Update_Does_Not_Affect_Other(t *testing.T) {
	setupDefConnection()
	resetCollection()

	unrelatedAuthor := NewDocAuthor("unrelated", primitive.NewObjectID())
	require.NoError(t, mgm.Coll(unrelatedAuthor).Create(unrelatedAuthor))

	d, authors := insertHasManyRelation(t)
	for i, author := range authors {
		author.Name = fmt.Sprintf("New-%v", i)
	}
	require.NoError(t, mgmrel.HasMany(d, &DocAuthor{}).Sync(authors))
	c, err := mgm.Coll(d).CountDocuments(nil, bson.M{})
	require.NoError(t, err)
	require.Equal(t, int64(1), c)

	foundAuthor := &DocAuthor{}
	require.NoError(t, mgm.Coll(&DocAuthor{}).First(bson.M{f.ID: unrelatedAuthor.ID}, foundAuthor))
	assert.Equal(t, unrelatedAuthor.ID, foundAuthor.ID)
	assert.Equal(t, unrelatedAuthor.DocID, foundAuthor.DocID)
	assert.Equal(t, unrelatedAuthor.Name, foundAuthor.Name)
}

func TestHasManyRelation_SyncWithoutRemove_Update_Does_Not_Affect_Other(t *testing.T) {
	setupDefConnection()
	resetCollection()

	unrelatedAuthor := NewDocAuthor("unrelated", primitive.NewObjectID())
	require.NoError(t, mgm.Coll(unrelatedAuthor).Create(unrelatedAuthor))

	d, authors := insertHasManyRelationWithoutRemove(t)
	for i, author := range authors {
		author.Name = fmt.Sprintf("New-%v", i)
	}
	require.NoError(t, mgmrel.HasMany(d, &DocAuthor{}).SyncWithoutRemove(authors))
	c, err := mgm.Coll(d).CountDocuments(nil, bson.M{})
	require.NoError(t, err)
	require.Equal(t, int64(1), c)

	foundAuthor := &DocAuthor{}
	require.NoError(t, mgm.Coll(&DocAuthor{}).First(bson.M{f.ID: unrelatedAuthor.ID}, foundAuthor))
	assert.Equal(t, unrelatedAuthor.ID, foundAuthor.ID)
	assert.Equal(t, unrelatedAuthor.DocID, foundAuthor.DocID)
	assert.Equal(t, unrelatedAuthor.Name, foundAuthor.Name)
}
