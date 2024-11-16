package mdu

import (
	"github.com/YspCoder/simple/mongo/field"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func create(c *Collection, model Model, opts ...*options.InsertOneOptions) (interface{}, error) {
	// Call to saving hook
	if err := beforeCreateHooks(c.ctx, model); err != nil {
		return nil, err
	}

	res, err := c.InsertOne(c.ctx, model, opts...)

	if err != nil {
		return nil, err
	}

	// Set new id
	model.SetID(res.InsertedID.(string))

	err = afterCreateHooks(c.ctx, model)
	if err != nil {
		return nil, err
	}
	return res.InsertedID, nil
}

func first(c *Collection, filter interface{}, model Model, opts ...*options.FindOneOptions) error {
	return c.FindOne(c.ctx, filter, opts...).Decode(model)
}

func update(c *Collection, model Model, opts ...*options.UpdateOptions) error {
	// Call to saving hook
	if err := beforeUpdateHooks(c.ctx, model); err != nil {
		return err
	}

	res, err := c.UpdateOne(c.ctx, bson.M{field.ID: model.GetID()}, bson.M{"$set": model}, opts...)

	if err != nil {
		return err
	}

	return afterUpdateHooks(c.ctx, res, model)
}

func patch(c *Collection, model Model, fields map[string]interface{}, opts ...*options.UpdateOptions) error {
	// Call to saving hook
	if err := beforeUpdateHooks(c.ctx, model); err != nil {
		return err
	}

	res, err := c.UpdateOne(c.ctx, bson.M{field.ID: model.GetID()}, bson.M{"$set": fields}, opts...)

	if err != nil {
		return err
	}

	return afterUpdateHooks(c.ctx, res, model)
}

func deleteByID(c *Collection, model Model) error {
	if err := beforeDeleteHooks(c.ctx, model); err != nil {
		return err
	}
	res, err := c.DeleteOne(c.ctx, bson.M{field.ID: model.GetID()})
	if err != nil {
		return err
	}

	return afterDeleteHooks(c.ctx, res, model)
}
