package mdu

import (
	"context"
)

type TLSEntity struct {
	CaCert        string `json:"ca_cert" bson:"ca_cert" yaml:"ca_cert" mapstructure:"ca_cert"`
	ClientCert    string `json:"client_cert" bson:"client_cert" yaml:"client_cert" mapstructure:"client_cert"`
	ClientCertKey string `json:"client_cert_key" bson:"client_cert_key" yaml:"client_cert_key" mapstructure:"client_cert_key"`
}

type ConfigEntity struct {
	Tls TLSEntity `json:"tls" bson:"tls" yaml:"tls" mapstructure:"tls"`

	Account  string `json:"account" bson:"account" yaml:"account" mapstructure:"account"`
	Password string `json:"password" bson:"password" yaml:"password" mapstructure:"password"`

	Address  string `json:"address" yaml:"address" mapstructure:"address"`
	Database string `json:"database" yaml:"database" mapstructure:"database"`

	Mode bool `json:"mode" yaml:"mode" mapstructure:"mode"` // Mode is true cluster

	MaxOpenConnects int `json:"max_open_connects" bson:"max_open_connects" yaml:"max_open_connects" mapstructure:"max_open_connects"`
	MaxIdleConnects int `json:"max_idle_connects" bson:"max_idle_connects" yaml:"max_idle_connects" mapstructure:"max_idle_connects"`
	ConnMaxLifeTime int `json:"conn_max_life_time" bson:"conn_max_life_time" yaml:"conn_max_life_time" mapstructure:"conn_max_life_time"`
}

// CollectionGetter interface contains a method to return
// a model's custom collection.
type CollectionGetter interface {
	// Collection method return collection
	Collection() *Collection
}

// CollectionNameGetter interface contains a method to return
// the collection name of a model.
type CollectionNameGetter interface {
	// CollectionName method return model collection's name.
	CollectionName() string
}

// Model interface contains base methods that must be implemented by
// each model. If you're using the `DefaultModel` struct in your model,
// you don't need to implement any of these methods.
type Model interface {
	// PrepareID converts the id value if needed, then
	// returns it (e.g.convert string to objectId).
	PrepareID(id string) (string, error)

	GetID() string
	SetID(id string)
}

// DefaultModel struct contains a model's default fields.
type DefaultModel struct {
	IDField    `bson:",inline"`
	DateFields `bson:",inline"`
}

// DefaultTenantModel struct contains a model's default fields. This is useful for multi tenant systems.
type DefaultTenantModel struct {
	IDField       `bson:",inline"`
	DateFields    `bson:",inline"`
	TenantIdField `bson:",inline"`
}

// Creating function calls the inner fields' defined hooks
func (model *DefaultModel) Creating(ctx context.Context) error {
	return model.DateFields.Creating(ctx)
}

// Saving function calls the inner fields' defined hooks
func (model *DefaultModel) Saving(ctx context.Context) error {
	return model.DateFields.Saving(ctx)
}
