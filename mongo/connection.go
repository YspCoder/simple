package mongo

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/YspCoder/simple/common/utils"
	"github.com/jinzhu/inflection"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"os"
	"reflect"
	"strconv"
	"strings"
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
	name := reflect.TypeOf(model).Elem().Name()

	return inflection.Plural(utils.ToLowerCamelCase(name))
}

func extractIndexes(model interface{}) []mongo.IndexModel {
	var indexes []mongo.IndexModel

	modelValue := reflect.ValueOf(model)
	// If the input is a pointer, get the value it points to
	if modelValue.Kind() == reflect.Ptr {
		modelValue = modelValue.Elem()
	}
	modelType := modelValue.Type()

	// Ensure the model is a struct type
	if modelType.Kind() != reflect.Struct {
		log.Fatalf("The provided model is not a struct type, cannot extract indexes")
	}

	for i := 0; i < modelType.NumField(); i++ {
		field := modelType.Field(i)
		indexTag := field.Tag.Get("index")
		if indexTag != "" {
			// Get the field name from the bson tag
			bsonTag := field.Tag.Get("bson")
			if bsonTag == "" || bsonTag == "-" {
				continue
			}

			var keys bson.D
			opts := options.Index()

			// Parse index type, split by comma to support multiple attributes like "1,unique"
			indexAttributes := strings.Split(indexTag, ",")
			isUnique := false
			for _, attr := range indexAttributes {
				switch strings.TrimSpace(attr) {
				case "1":
					keys = bson.D{{Key: bsonTag, Value: 1}} // Ascending index
				case "-1":
					keys = bson.D{{Key: bsonTag, Value: -1}} // Descending index
				case "unique":
					isUnique = true
				case "sparse":
					opts.SetSparse(true) // Sparse index
				case "text":
					keys = bson.D{{Key: bsonTag, Value: "text"}} // Text index
				case "hashed":
					keys = bson.D{{Key: bsonTag, Value: "hashed"}} // Hashed index
				case "ttl":
					keys = bson.D{{Key: bsonTag, Value: 1}}
					// TTL index requires a ttlSeconds tag
					ttlSecondsTag := field.Tag.Get("ttlSeconds")
					if ttlSecondsTag == "" {
						log.Fatalf("TTL index must have a ttlSeconds tag specified")
					}
					ttlSeconds, err := strconv.Atoi(ttlSecondsTag)
					if err != nil {
						log.Fatalf("The ttlSeconds tag must be an integer: %v", err)
					}
					opts.SetExpireAfterSeconds(int32(ttlSeconds)) // TTL index
				case "2dsphere":
					keys = bson.D{{Key: bsonTag, Value: "2dsphere"}} // 2dsphere geospatial index
				case "2d":
					keys = bson.D{{Key: bsonTag, Value: "2d"}} // 2d geospatial index
				default:
					continue
				}
			}

			if len(keys) > 0 {
				if isUnique {
					opts.SetUnique(true)
				}
				indexes = append(indexes, mongo.IndexModel{
					Keys:    keys,
					Options: opts,
				})
			}
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
