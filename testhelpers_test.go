package mgmrel_test

import (
	"github.com/kamva/gutil"
	mgmrel "github.com/kamva/mgm-relation"
	"github.com/kamva/mgm/v3"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"testing"
)

func setupDefConnection() {
	gutil.PanicErr(mgm.SetDefaultConfig(nil, "mgm_relations_test", options.Client().ApplyURI("mongodb://root:12345@localhost:27017")), )
}

func resetCollection() {
	_, err := mgm.Coll(&Doc{}).DeleteMany(mgm.Ctx(), bson.M{})
	_, err2 := mgm.Coll(&DocAuthor{}).DeleteMany(mgm.Ctx(), bson.M{})

	gutil.PanicErr(err)
	gutil.PanicErr(err2)
}

func seed() {
	docs := []interface{}{
		NewDoc("Ali", 24),
		NewDoc("Mehran", 24),
		NewDoc("Reza", 26),
		NewDoc("Omid", 27),
	}
	_, err := mgm.Coll(&Doc{}).InsertMany(mgm.Ctx(), docs)

	gutil.PanicErr(err)
}

func findDoc(t *testing.T) *Doc {
	found := &Doc{}
	err := mgm.Coll(found).FindOne(mgm.Ctx(), bson.M{}).Decode(found)
	assert.NoError(t, err)
	return found
}

type Doc struct {
	mgmrel.IDField `bson:",inline"`

	Name string `bson:"name"`
	Age  int    `bson:"age"`
}

type DocAuthor struct {
	mgmrel.IDField `bson:",inline"`

	Name  string             `bson:"name"`
	DocID primitive.ObjectID `json:"doc_id" bson:"doc_id"` // The foreign key
}

func NewDoc(name string, age int) *Doc {
	return &Doc{
		Name: name,
		Age:  age,
	}
}

func NewDocAuthor(name string, docID primitive.ObjectID) *DocAuthor {
	return &DocAuthor{Name: name, DocID: docID}
}
