package mongo

import (
	"context"
	"github.com/YspCoder/simple/mongo/field"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func create(ctx context.Context, c *Collection, model Model, opts ...*options.InsertOneOptions) (interface{}, error) {
	// Call to saving hook
	if err := beforeCreateHooks(ctx, model); err != nil {
		return nil, err
	}

	res, err := c.InsertOne(ctx, model, opts...)

	if err != nil {
		return nil, err
	}

	// Set new id
	model.SetID(res.InsertedID.(string))

	err = afterCreateHooks(ctx, model)
	if err != nil {
		return nil, err
	}
	return res.InsertedID, nil
}

func first(ctx context.Context, c *Collection, filter interface{}, model Model, opts ...*options.FindOneOptions) error {
	return c.FindOne(ctx, filter, opts...).Decode(model)
}

func update(ctx context.Context, c *Collection, model Model, opts ...*options.UpdateOptions) error {
	// Call to saving hook
	if err := beforeUpdateHooks(ctx, model); err != nil {
		return err
	}

	res, err := c.UpdateOne(ctx, bson.M{field.ID: model.GetID()}, bson.M{"$set": model}, opts...)

	if err != nil {
		return err
	}

	return afterUpdateHooks(ctx, res, model)
}

func patch(ctx context.Context, c *Collection, model Model, fields map[string]interface{}, opts ...*options.UpdateOptions) error {
	// Call to saving hook
	if err := beforeUpdateHooks(ctx, model); err != nil {
		return err
	}

	res, err := c.UpdateOne(ctx, bson.M{field.ID: model.GetID()}, bson.M{"$set": fields}, opts...)

	if err != nil {
		return err
	}

	return afterUpdateHooks(ctx, res, model)
}

func deleteByID(ctx context.Context, c *Collection, model Model) error {
	if err := beforeDeleteHooks(ctx, model); err != nil {
		return err
	}
	res, err := c.DeleteOne(ctx, bson.M{field.ID: model.GetID()})
	if err != nil {
		return err
	}

	return afterDeleteHooks(ctx, res, model)
}
