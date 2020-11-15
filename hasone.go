package mgmrel

import (
	"github.com/kamva/gutil"
	"github.com/kamva/mgm/v3"
	f "github.com/kamva/mgm/v3/field"
	o "github.com/kamva/mgm/v3/operator"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type HasOneRelation struct {
	m       mgm.Model
	related mgm.Model
	// foreignKey uses in filters.
	foreignKey string
}

// SimpleGet method get the single related model.
// if not found, returns the Mongo Go driver not found error.
func (r *HasOneRelation) Get(m mgm.Model) error {
	return mgm.Coll(r.related).First(r.filterByRelation(nil), m)
}

// Sync method sync the relations:
// If provided model is nil: it remove the related model in the DB.
// If provided model is not nil: sync it.
// insert new model, otherwise upsert provided model.
func (r *HasOneRelation) Sync(model mgm.Model) error {
	if gutil.IsNil(model) {
		_, err := r.delete(nil)
		return err
	}
	if err := callToBeforeSyncHooks(model); err != nil {
		return err
	}

	_, err := r.delete(model.GetID())
	upsert := true
	_, err = mgm.Coll(r.related).UpdateOne(mgm.Ctx(), bson.M{f.ID: model.GetID()}, bson.M{o.Set: model}, &options.UpdateOptions{
		Upsert: &upsert,
	})
	if err != nil {
		return err
	}

	return callToAfterSyncHooks(model)
}

func (r *HasOneRelation) delete(exceptID interface{}) (*mongo.DeleteResult, error) {
	return mgm.Coll(r.related).DeleteMany(mgm.Ctx(), r.filterByRelation(exceptID))
}

func (r *HasOneRelation) filterByRelation(exceptionID interface{}) bson.M {
	filter := bson.M{r.foreignKey: r.m.GetID()}
	if !gutil.IsNil(exceptionID) {
		filter[f.ID] = bson.M{o.Ne: exceptionID}
	}
	return filter
}

// HasOne returns new instance of the "has one" relation ship.
func HasOne(model mgm.Model, related mgm.Model) *HasOneRelation {
	return HasOneByOptions(model, related, foreignKeyName(model))
}

// HasOneByOptions gets HasOneRelation options and returns new instance of it.
func HasOneByOptions(model mgm.Model, related mgm.Model, foreignKey string) *HasOneRelation {
	return &HasOneRelation{
		m:          model,
		related:    related,
		foreignKey: foreignKey,
	}
}
