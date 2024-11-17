package mongo

import (
	"context"
	"github.com/YspCoder/simple/mongo/builder"
	"github.com/YspCoder/simple/mongo/field"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Collection 类用于对模型和指定的 MongoDB 集合进行操作
// 这里将 mongo.Collection 包装一层，提供更方便的操作方法

// Collection performs operations on models and the given Mongodb collection
type Collection struct {
	*mongo.Collection
}

// FindByIDWithCtx 根据给定的 ID 查找文档并解析为模型
// id 可以是任意类型（e.g. string, bson.ObjectId）
func (c *Collection) FindByIDWithCtx(ctx context.Context, id interface{}, model Model, opts ...*options.FindOneOptions) error {
	return first(ctx, c, bson.M{field.ID: id}, model, opts...)
}

// FindByID 与 FindByIDWithCtx 相同，使用给定的 ID 进行文档查找
func (c *Collection) FindByID(ctx context.Context, id interface{}, model Model, opts ...*options.FindOneOptions) error {
	return first(ctx, c, bson.M{field.ID: id}, model, opts...)
}

// First 方法使用过滤器查找文档，并返回第一个结果
func (c *Collection) First(ctx context.Context, filter interface{}, model Model, opts ...*options.FindOneOptions) error {
	return first(ctx, c, filter, model, opts...)
}

// FirstWithCtx 与 First 相同，使用指定的过滤器查找第一个文档
func (c *Collection) FirstWithCtx(ctx context.Context, filter interface{}, model Model, opts ...*options.FindOneOptions) error {
	return first(ctx, c, filter, model, opts...)
}

// Create 方法将模型插入到数据库
func (c *Collection) Create(ctx context.Context, model Model, opts ...*options.InsertOneOptions) (interface{}, error) {
	return createWithCtx(ctx, c, model, opts...)
}

// CreateWithCtx 与 Create 相同，使用指定的模型创建文档
func (c *Collection) CreateWithCtx(ctx context.Context, model Model, opts ...*options.InsertOneOptions) (interface{}, error) {
	return createWithCtx(ctx, c, model, opts...)
}

// createWithCtx 进行将模型插入到数据库的实际操作
func createWithCtx(ctx context.Context, c *Collection, model Model, opts ...*options.InsertOneOptions) (interface{}, error) {
	id, err := model.PrepareID(model.GetID())

	if err != nil {
		return nil, err
	}
	model.SetID(id)
	return create(ctx, c, model, opts...)
}

// Update 方法使用指定的模型更新数据库中的文档
// 该方法会触发更新前和更新后的钩子
func (c *Collection) Update(ctx context.Context, model Model, opts ...*options.UpdateOptions) error {
	return update(ctx, c, model, opts...)
}

// UpdateWithCtx 与 Update 相同，使用指定的模型更新文档
func (c *Collection) UpdateWithCtx(ctx context.Context, model Model, opts ...*options.UpdateOptions) error {
	return update(ctx, c, model, opts...)
}

// Patch 方法将特定的字段更新到数据库中
// 该方法会触发更新前和更新后的钩子
func (c *Collection) Patch(ctx context.Context, model Model, fields map[string]interface{}, opts ...*options.UpdateOptions) error {
	return patch(ctx, c, model, fields, opts...)
}

// PatchWithCtx 与 Patch 相同，使用特定的字段来更新文档
func (c *Collection) PatchWithCtx(ctx context.Context, model Model, fields map[string]interface{}, opts ...*options.UpdateOptions) error {
	return patch(ctx, c, model, fields, opts...)
}

// Delete 方法使用指定的模型从数据库中删除文档
func (c *Collection) Delete(ctx context.Context, model Model) error {
	return deleteByID(ctx, c, model)
}

// DeleteWithCtx 与 Delete 相同，从数据库中删除文档
func (c *Collection) DeleteWithCtx(ctx context.Context, model Model) error {
	return deleteByID(ctx, c, model)
}

// FindAll 方法使用指定的过滤器找出所有结果，并返回
func (c *Collection) FindAll(ctx context.Context, results interface{}, filter interface{}, opts ...*options.FindOptions) error {
	return findAll(ctx, c, results, filter, opts...)
}

// FindAllWithCtx 与 FindAll 相同，找出所有结果并返回
func (c *Collection) FindAllWithCtx(ctx context.Context, results interface{}, filter interface{}, opts ...*options.FindOptions) error {
	return findAll(ctx, c, results, filter, opts...)
}

// findAll 实际实现了 FindAll 方法，通过过滤器获取所有文档
func findAll(ctx context.Context, c *Collection, results interface{}, filter interface{}, opts ...*options.FindOptions) error {
	cur, err := c.Find(ctx, filter, opts...)

	if err != nil {
		return err
	}

	return cur.All(ctx, results)
}

//--------------------------------
// 聚合操作方法
//--------------------------------

// SimpleAggregateFirst 进行一个简单的聚合，并返回第一个聚合结果
// `stages` 可以是 Operator | bson.M
// 请注意：该方法不能用于事务中，请使用正规聚合方法
func (c *Collection) SimpleAggregateFirst(ctx context.Context, result interface{}, stages ...interface{}) (bool, error) {
	return simpleAggregateFirst(ctx, c, result, stages...)
}

// SimpleAggregateFirstWithCtx 与 SimpleAggregateFirst 相同，返回聚合第一个结果
func (c *Collection) SimpleAggregateFirstWithCtx(ctx context.Context, result interface{}, stages ...interface{}) (bool, error) {
	return simpleAggregateFirst(ctx, c, result, stages...)
}

// simpleAggregateFirst 实现 SimpleAggregateFirst 的聚合操作
func simpleAggregateFirst(ctx context.Context, c *Collection, result interface{}, stages ...interface{}) (bool, error) {
	cur, err := c.SimpleAggregateCursorWithCtx(ctx, stages...)
	if err != nil {
		return false, err
	}
	if cur.Next(ctx) {
		return true, cur.Decode(result)
	}
	return false, nil
}

// SimpleAggregate 进行一个简单的聚合，并返回聚合结果的列表
// `stages` 可以是 Operator | bson.M
// 请注意：该方法不能用于事务中，请使用正规聚合方法
func (c *Collection) SimpleAggregate(ctx context.Context, results interface{}, stages ...interface{}) error {
	return simpleAggregate(ctx, c, results, stages...)
}

// SimpleAggregateWithCtx 与 SimpleAggregate 相同，并返回聚合结果的列表
func (c *Collection) SimpleAggregateWithCtx(ctx context.Context, results interface{}, stages ...interface{}) error {
	return simpleAggregate(ctx, c, results, stages...)
}

// simpleAggregate 实现 SimpleAggregate 的聚合操作
func simpleAggregate(ctx context.Context, c *Collection, results interface{}, stages ...interface{}) error {
	cur, err := c.SimpleAggregateCursorWithCtx(ctx, stages...)
	if err != nil {
		return err
	}

	return cur.All(ctx, results)
}

// SimpleAggregateCursor 进行一个简单的聚合，并返回文档的遍历器
// 请注意：该方法不能用于事务中，请使用正规聚合方法
func (c *Collection) SimpleAggregateCursor(ctx context.Context, stages ...interface{}) (*mongo.Cursor, error) {
	return simpleAggregateCursor(ctx, c, stages...)
}

// SimpleAggregateCursorWithCtx 与 SimpleAggregateCursor 相同，返回聚合的遍历器
func (c *Collection) SimpleAggregateCursorWithCtx(ctx context.Context, stages ...interface{}) (*mongo.Cursor, error) {
	return simpleAggregateCursor(ctx, c, stages...)
}

// simpleAggregateCursor 实现聚合操作，返回文档的遍历器
func simpleAggregateCursor(ctx context.Context, c *Collection, stages ...interface{}) (*mongo.Cursor, error) {
	pipeline := bson.A{}

	for _, stage := range stages {
		if operator, ok := stage.(builder.Operator); ok {
			pipeline = append(pipeline, builder.S(operator))
		} else {
			pipeline = append(pipeline, stage)
		}
	}

	return c.Aggregate(ctx, pipeline, nil)
}
