package mongo

import (
	"context"
	"github.com/YspCoder/simple/mongo/field"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Cnd struct {
	Filter     bson.D   // 过滤条件
	SelectCols []string // 要查询的字段，如果为空，表示查询所有字段
	Sort       bson.D   // 排序
	Paging     *Paging  // 分页
}

type Paging struct {
	Page  int // 页数
	Limit int // 每页的记录数
}

func NewCnd() *Cnd {
	return &Cnd{}
}

func (c *Cnd) Cols(selectCols ...string) *Cnd {
	if len(selectCols) > 0 {
		c.SelectCols = append(c.SelectCols, selectCols...)
	}
	return c
}

func (c *Cnd) Eq(column string, value interface{}) *Cnd {
	c.Filter = append(c.Filter, bson.E{Key: column, Value: value})
	return c
}

func (c *Cnd) NotEq(column string, value interface{}) *Cnd {
	c.Filter = append(c.Filter, bson.E{Key: column, Value: bson.M{field.Ne: value}})
	return c
}

func (c *Cnd) Gt(column string, value interface{}) *Cnd {
	c.Filter = append(c.Filter, bson.E{Key: column, Value: bson.M{field.Gt: value}})
	return c
}

func (c *Cnd) Gte(column string, value interface{}) *Cnd {
	c.Filter = append(c.Filter, bson.E{Key: column, Value: bson.M{field.Gte: value}})
	return c
}

func (c *Cnd) Lt(column string, value interface{}) *Cnd {
	c.Filter = append(c.Filter, bson.E{Key: column, Value: bson.M{field.Lt: value}})
	return c
}

func (c *Cnd) Lte(column string, value interface{}) *Cnd {
	c.Filter = append(c.Filter, bson.E{Key: column, Value: bson.M{field.Lte: value}})
	return c
}

func (c *Cnd) Like(column string, str string) *Cnd {
	c.Filter = append(c.Filter, bson.E{Key: column, Value: bson.M{field.Regex: str, "$options": "i"}})
	return c
}

func (c *Cnd) In(column string, params []interface{}) *Cnd {
	c.Filter = append(c.Filter, bson.E{Key: column, Value: bson.M{field.In: params}})
	return c
}

func (c *Cnd) NotIn(column string, params []interface{}) *Cnd {
	c.Filter = append(c.Filter, bson.E{Key: column, Value: bson.M{field.Nin: params}})
	return c
}

func (c *Cnd) Asc(column string) *Cnd {
	c.Sort = append(c.Sort, bson.E{Key: column, Value: 1})
	return c
}

func (c *Cnd) Desc(column string) *Cnd {
	c.Sort = append(c.Sort, bson.E{Key: column, Value: -1})
	return c
}

func (c *Cnd) Page(page, limit int) *Cnd {
	c.Paging = &Paging{Page: page, Limit: limit}
	return c
}

func (c *Cnd) BuildFindOptions() *options.FindOptions {
	opts := options.Find()

	// Select columns
	if len(c.SelectCols) > 0 {
		projection := bson.D{}
		for _, col := range c.SelectCols {
			projection = append(projection, bson.E{Key: col, Value: 1})
		}
		opts.SetProjection(projection)
	}

	// Sort
	if len(c.Sort) > 0 {
		sort := bson.D{}
		for _, s := range c.Sort {
			sort = append(sort, s)
		}
		opts.SetSort(sort)
	}

	// Paging
	if c.Paging != nil {
		if c.Paging.Limit > 0 {
			opts.SetLimit(int64(c.Paging.Limit))
		}
		if c.Paging.Page > 0 && c.Paging.Limit > 0 {
			offset := int64((c.Paging.Page - 1) * c.Paging.Limit)
			opts.SetSkip(offset)
		}
	}

	return opts
}

func (c *Cnd) Find(ctx context.Context, db *mongo.Collection, results interface{}) error {
	opts := c.BuildFindOptions()
	cur, err := db.Find(ctx, c.Filter, opts)
	if err != nil {
		logrus.Error(err)
		return err
	}
	if err = cur.All(ctx, results); err != nil {
		logrus.Error(err)
		return err
	}
	return nil
}

func (c *Cnd) FindOne(ctx context.Context, db *mongo.Collection, result interface{}) error {
	findOneOpts := options.FindOne()

	// Set select columns
	if len(c.SelectCols) > 0 {
		projection := bson.D{}
		for _, col := range c.SelectCols {
			projection = append(projection, bson.E{Key: col, Value: 1})
		}
		findOneOpts.SetProjection(projection)
	}

	// Set sort options
	if len(c.Sort) > 0 {
		sort := bson.D{}
		for _, s := range c.Sort {
			sort = append(sort, s)
		}
		findOneOpts.SetSort(sort)
	}

	err := db.FindOne(ctx, c.Filter, findOneOpts).Decode(result)
	if err != nil {
		logrus.Error(err)
		return err
	}
	return nil
}

func (c *Cnd) Count(ctx context.Context, db *mongo.Collection) (int64, error) {
	count, err := db.CountDocuments(ctx, c.Filter)
	if err != nil {
		logrus.Error(err)
		return 0, err
	}
	return count, nil
}

// Aggregate performs an aggregation query based on the given pipeline stages.
func (c *Cnd) Aggregate(ctx context.Context, db *mongo.Collection, pipeline []bson.M, results interface{}) error {
	cur, err := db.Aggregate(ctx, pipeline)
	if err != nil {
		logrus.Error(err)
		return err
	}
	if err = cur.All(ctx, results); err != nil {
		logrus.Error(err)
		return err
	}
	return nil
}

// AggregateWithConditions performs an aggregation query with additional filtering, sorting, and paging conditions.
func (c *Cnd) AggregateWithConditions(ctx context.Context, db *mongo.Collection, pipeline []bson.M, results interface{}) error {
	// Add filter stage if there are any conditions
	if len(c.Filter) > 0 {
		pipeline = append([]bson.M{{field.Match: c.Filter}}, pipeline...)
	}

	// Add sort stage if sorting is specified
	if len(c.Sort) > 0 {
		sortStage := bson.M{field.Sort: bson.M{}}
		for _, s := range c.Sort {
			sortStage[field.Sort].(bson.M)[s.Key] = s.Value
		}
		pipeline = append(pipeline, sortStage)
	}

	// Add limit and skip for pagination
	if c.Paging != nil {
		if c.Paging.Limit > 0 {
			pipeline = append(pipeline, bson.M{field.Limit: int64(c.Paging.Limit)})
		}
		if c.Paging.Page > 0 && c.Paging.Limit > 0 {
			offset := int64((c.Paging.Page - 1) * c.Paging.Limit)
			pipeline = append(pipeline, bson.M{field.Skip: offset})
		}
	}

	cur, err := db.Aggregate(ctx, pipeline)
	if err != nil {
		logrus.Error(err)
		return err
	}
	if err = cur.All(ctx, results); err != nil {
		logrus.Error(err)
		return err
	}
	return nil
}
