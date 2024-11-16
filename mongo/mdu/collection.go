package mdu

import (
	"context"
	"github.com/YspCoder/simple/mongo/field"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Collection performs operations on models and the given Mongodb collection
type Collection struct {
	*mongo.Collection
	ctx context.Context
}

func (c *Collection) FindByIDWithCtx(id interface{}, model Model, opts ...*options.FindOneOptions) error {
	return first(c, bson.M{field.ID: id}, model, opts...)
}

// FindByID method finds a doc and decodes it to a model, otherwise returns an error.
// The id field can be any value that if passed to the `PrepareID` method, it returns
// a valid ID (e.g.string, bson.ObjectId).
func (c *Collection) FindByID(id interface{}, model Model, opts ...*options.FindOneOptions) error {
	return first(c, bson.M{field.ID: id}, model, opts...)
}

// First method searches and returns the first document in the search results.
func (c *Collection) First(filter interface{}, model Model, opts ...*options.FindOneOptions) error {
	return first(c, filter, model, opts...)
}

func (c *Collection) FirstWithCtx(filter interface{}, model Model, opts ...*options.FindOneOptions) error {
	return first(c, filter, model, opts...)
}

// Create method inserts a new model into the database.
func (c *Collection) Create(model Model, opts ...*options.InsertOneOptions) (interface{}, error) {
	return createWithCtx(c, model, opts...)
}

func (c *Collection) CreateWithCtx(model Model, opts ...*options.InsertOneOptions) (interface{}, error) {
	return createWithCtx(c, model, opts...)
}

func createWithCtx(c *Collection, model Model, opts ...*options.InsertOneOptions) (interface{}, error) {
	id, err := model.PrepareID(model.GetID())

	if err != nil {
		return nil, err
	}
	model.SetID(id)
	return create(c, model, opts...)
}

// Update function persists the changes made to a model to the database using the specified context.
// Calling this method also invokes the model's mdu updating, updated,
// saving, and saved hooks.
func (c *Collection) Update(model Model, opts ...*options.UpdateOptions) error {
	return update(c, model, opts...)
}

func (c *Collection) UpdateWithCtx(model Model, opts ...*options.UpdateOptions) error {
	return update(c, model, opts...)
}

// Patch function persists the given fields in a model to the database using the specified context.
// Calling this method also invokes the model's mdu updating, updated,
// saving, and saved hooks.
func (c *Collection) Patch(model Model, fields map[string]interface{}, opts ...*options.UpdateOptions) error {
	return patch(c, model, fields, opts...)
}

func (c *Collection) PatchWithCtx(model Model, fields map[string]interface{}, opts ...*options.UpdateOptions) error {
	return patch(c, model, fields, opts...)
}

// Delete method deletes a model (doc) from a collection using the specified context.
// To perform additional operations when deleting a model
// you should use hooks rather than overriding this method.
func (c *Collection) Delete(model Model) error {
	return deleteByID(c, model)
}

func (c *Collection) DeleteWithCtx(model Model) error {
	return deleteByID(c, model)
}

// FindAll finds, decodes and returns the results using the specified context.
func (c *Collection) FindAll(results interface{}, filter interface{}, opts ...*options.FindOptions) error {
	return findAll(c, results, filter, opts...)
}

func (c *Collection) FindAllWithCtx(results interface{}, filter interface{}, opts ...*options.FindOptions) error {
	return findAll(c, results, filter, opts...)
}

func findAll(c *Collection, results interface{}, filter interface{}, opts ...*options.FindOptions) error {
	cur, err := c.Find(c.ctx, filter, opts...)

	if err != nil {
		return err
	}

	return cur.All(c.ctx, results)
}
