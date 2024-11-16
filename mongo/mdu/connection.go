package mdu

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"os"
	"time"
)

type MongoDB struct {
	Db  *mongo.Database
	Ctx context.Context
}

// Init initializes the client and database using the specified configuration values, or default.
func Init(conf *ConfigEntity) (*MongoDB, error) {
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
		CAFile, CAErr := os.ReadFile(conf.Tls.CaCert)
		if CAErr != nil {
			return nil, CAErr
		}
		certPool.AppendCertsFromPEM(CAFile)

		clientCert, clientCertErr := tls.LoadX509KeyPair(conf.Tls.ClientCert, conf.Tls.ClientCertKey)
		if clientCertErr != nil {
			return nil, clientCertErr
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

	client, cErr := mongo.Connect(ctx, &clientOptions)
	if cErr != nil {
		return nil, cErr
	}

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, err
	}

	db := client.Database(conf.Database)
	return &MongoDB{Db: db, Ctx: ctx}, nil
}

func (c *MongoDB) CollectionByName(name string, opts ...*options.CollectionOptions) *Collection {
	coll := c.Db.Collection(name, opts...)
	return &Collection{Collection: coll, ctx: c.Ctx}
}
