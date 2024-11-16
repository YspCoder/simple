package mdu

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"os"
	"reflect"
	"time"
)

var (
	db *mongo.Database
)

// Init initializes the client and database using the specified configuration values, or default.
func Init(conf *ConfigEntity, models ...interface{}) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var clientOptions options.ClientOptions

	if conf.Account != "" && conf.Password != "" {
		clientOptions.SetAuth(options.Credential{
			Username: conf.Account,
			Password: conf.Password,
		})
	}
	if conf.Tls.CaCert != "" && conf.Tls.ClientCert != "" && conf.Tls.ClientCertKey != "" {
		certPool := x509.NewCertPool()
		CAFile, err := os.ReadFile(conf.Tls.CaCert)
		if err != nil {
			return err
		}
		certPool.AppendCertsFromPEM(CAFile)

		clientCert, err := tls.LoadX509KeyPair(conf.Tls.ClientCert, conf.Tls.ClientCertKey)
		if err != nil {
			return err
		}

		tlsConfig := tls.Config{
			Certificates: []tls.Certificate{clientCert},
			RootCAs:      certPool,
		}
		clientOptions.SetTLSConfig(&tlsConfig)
	}

	uri := fmt.Sprintf("mongodb://%s", conf.Address)
	clientOptions.ApplyURI(uri)

	clientOptions.SetBSONOptions(&options.BSONOptions{
		UseLocalTimeZone: true,
	})

	clientOptions.SetMaxConnecting(uint64(conf.MaxOpenConnects))
	clientOptions.SetMaxPoolSize(uint64(conf.MaxIdleConnects))
	clientOptions.SetMaxConnIdleTime(time.Second * time.Duration(conf.ConnMaxLifeTime))

	loger := log.New()

	var statement string
	clientOptions.Monitor = &event.CommandMonitor{
		Started: func(ctx context.Context, event *event.CommandStartedEvent) {
			statement = event.Command.String()
		},
		Succeeded: func(ctx context.Context, event *event.CommandSucceededEvent) {
			loger.Trace(ctx, event.RequestID, event.Duration, statement, "")
		},
		Failed: func(ctx context.Context, event *event.CommandFailedEvent) {
			loger.Trace(ctx, event.RequestID, event.Duration, statement, event.Failure)
		},
	}

	client, err := mongo.Connect(ctx, &clientOptions)
	if err != nil {
		return err
	}

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return err
	}

	db = client.Database(conf.Database)

	for _, model := range models {
		collectionName := getCollectionName(model)
		exists, err := isCollectionExists(db, collectionName)
		if err != nil {
			log.Fatalf("isCollectionExists err : %v", err)
		}

		if !exists {
			err = db.CreateCollection(context.TODO(), collectionName)
			if err != nil {
				log.Fatalf("create %s err: %v", collectionName, err)
			}
		}

		// 创建索引
		collectionRef := db.Collection(collectionName)
		indexes := extractIndexes(model)
		if len(indexes) > 0 {
			if _, err := collectionRef.Indexes().CreateMany(context.TODO(), indexes); err != nil {
				log.Fatalf("create %s index err: %v", collectionName, err)
			}
		}
	}

	return nil
}

// 获取集合名称
func getCollectionName(model interface{}) string {
	// 可以通过反射获取结构体名称作为集合名，或自定义处理
	return reflect.TypeOf(model).Name()
}

// 提取结构体中的索引信息
func extractIndexes(model interface{}) []mongo.IndexModel {
	var indexes []mongo.IndexModel

	modelType := reflect.TypeOf(model)
	for i := 0; i < modelType.NumField(); i++ {
		field := modelType.Field(i)
		indexTag := field.Tag.Get("index")
		if indexTag != "" {
			keys := bson.D{{Key: field.Tag.Get("bson"), Value: 1}}
			opts := options.Index()
			if indexTag == "unique" {
				opts.SetUnique(true)
			}
			indexes = append(indexes, mongo.IndexModel{
				Keys:    keys,
				Options: opts,
			})
		}
	}
	return indexes
}

// 检查集合是否存在
func isCollectionExists(database *mongo.Database, collectionName string) (bool, error) {
	collections, err := database.ListCollectionNames(context.TODO(), bson.D{{Key: "name", Value: collectionName}})
	if err != nil {
		return false, err
	}
	return len(collections) > 0, nil
}

func DB() *mongo.Database {
	return db
}

func CollectionByName(name string, opts ...*options.CollectionOptions) *Collection {
	coll := db.Collection(name, opts...)
	return &Collection{Collection: coll}
}
