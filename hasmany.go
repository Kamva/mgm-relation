package mgmrel

import (
	"github.com/Kamva/gutil"
	"github.com/Kamva/mgm/v3"
	f "github.com/Kamva/mgm/v3/field"
	o "github.com/Kamva/mgm/v3/operator"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type HasManyRelation struct {
	m       mgm.Model
	related mgm.Model
	// foreignKey uses in filters.
	foreignKey string
}

// Get method get the list of related models with provided filter,limit,...
// if not found, returns the Mongo Go driver not found error.
func (r *HasManyRelation) GetWithOptions(results interface{}, options ...*options.FindOptions) error {
	return mgm.Coll(r.related).SimpleFind(results, r.filterByRelation(nil), options...)
}

// Get method get the list of related models with provided filter,limit,...
// if not found, returns the Mongo Go driver not found error.
func (r *HasManyRelation) Get(results interface{}, sort string, skip, limit int64) error {
	return r.GetWithOptions(results, &options.FindOptions{
		Limit: &limit,
		Skip:  &skip,
		Sort:  r.sortFieldToBsonD(sort),
	})
}

// SimpleGet method get the list of related models
// if not found, returns the Mongo Go driver not found error.
// sort is the sort field. you can sort descending by adding a `-` to the sort field. e.g `-created_at`
func (r *HasManyRelation) SimpleGet(results interface{}, limit int64) error {
	return r.Get(results, "-_id", 0, limit)
}

// SyncWithoutRemove method sync the relations without
// removing items that are not in the provided list.
func (r *HasManyRelation) SyncWithoutRemove(docs interface{}) error {
	if gutil.IsNil(docs) {
		return nil
	}
	models := r.toModel(gutil.InterfaceToSlice(docs))
	if len(models) == 0 {
		return nil
	}
	upsert := true
	for _, m := range models {
		if err := callToBeforeSyncHooks(m); err != nil {
			return err
		}
		_, err := mgm.Coll(r.related).UpdateOne(mgm.Ctx(), bson.M{f.ID: m.GetID()}, bson.M{o.Set: m}, &options.UpdateOptions{
			Upsert: &upsert,
		})
		if err != nil {
			return err
		}
		if err := callToAfterSyncHooks(m); err != nil {
			return err
		}
	}
	return nil
}

// Sync method sync the relations:
// If provided models is nil(or length is zero): it remove the related models in the DB.
// If provided models is not nil and length is not zero: udpate new items, and remove
// items that are not in the provided list.
// Use sync just when your 1-m model contains just few m mdoel. otherwise use SyncWithoutRemove
func (r *HasManyRelation) Sync(docs interface{}) error {
	if gutil.IsNil(docs) {
		_, err := r.delete(nil)
		return err
	}
	models := r.toModel(gutil.InterfaceToSlice(docs))
	if len(models) == 0 {
		_, err := r.delete(nil)
		return err
	}
	upsert := true
	for _, m := range models {
		if err := callToBeforeSyncHooks(m); err != nil {
			return err
		}
		_, err := mgm.Coll(r.related).UpdateOne(mgm.Ctx(), bson.M{f.ID: m.GetID()}, bson.M{o.Set: m}, &options.UpdateOptions{
			Upsert: &upsert,
		})
		if err != nil {
			return err
		}
		if err := callToAfterSyncHooks(m); err != nil {
			return err
		}
	}
	// Delete All other models that are not in provided models.
	_, err := r.delete(r.extractIDs(models))
	return err
}

func (r *HasManyRelation) delete(exceptIDs []interface{}) (*mongo.DeleteResult, error) {
	return mgm.Coll(r.related).DeleteMany(mgm.Ctx(), r.filterByRelation(exceptIDs))
}

func (r *HasManyRelation) filterByRelation(exceptIDs []interface{}) bson.M {
	filter := bson.M{r.foreignKey: r.m.GetID()}
	if len(exceptIDs) != 0 {
		filter[f.ID] = bson.M{o.Nin: exceptIDs}
	}

	return filter
}

func (r *HasManyRelation) toModel(docs []interface{}) []mgm.Model {
	models := make([]mgm.Model, len(docs))
	for i, doc := range docs {
		models[i] = doc.(mgm.Model)
	}

	return models
}

func (r *HasManyRelation) extractIDs(models []mgm.Model) []interface{} {
	ids := make([]interface{}, len(models))
	for i, m := range models {
		ids[i] = m.GetID()
	}
	return ids
}

// sortFieldToBsonD converts the string sort field to bson D.
func (r *HasManyRelation) sortFieldToBsonD(field string) bson.D {
	// Ascending order
	order := 1
	if field[0] == '-' {
		order = -1
		field = field[1:]
	}

	return bson.D{
		{field, order},
	}
}

// HasMany returns new instance of the "has many" relation ship.
func HasMany(model mgm.Model, related mgm.Model) *HasManyRelation {
	return HasManyByOptions(model, related, foreignKeyName(model))
}

// HasManyByOptions gets HasManyRelation options and returns new instance of it.
func HasManyByOptions(model mgm.Model, related mgm.Model, foreignKey string) *HasManyRelation {
	return &HasManyRelation{
		m:          model,
		related:    related,
		foreignKey: foreignKey,
	}
}
