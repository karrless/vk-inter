package mongo

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoConfig struct {
	Host     string `env:"MONGO_HOST" env-default:"127.0.0.1"`
	Port     int    `env:"MONGO_PORT" env-default:"27017"`
	User     string `env:"MONGO_USER" env-default:"admin"`
	Password string `env:"MONGO_PASS" env-default:"admin"`
	DBName   string `env:"MONGO_DB" env-default:"vk-inter"`
}

// MongoDB structure for working with MongoDB
type MongoDB struct {
	Client   *mongo.Client
	Database *mongo.Database
}

// Get new MongoDB instanse
func New(ctx context.Context, cfg MongoConfig) (*MongoDB, error) {
	var uri string

	if cfg.User != "" && cfg.Password != "" {
		uri = fmt.Sprintf("mongodb://%s:%s@%s:%d", cfg.User, cfg.Password, cfg.Host, cfg.Port)
	} else {
		uri = fmt.Sprintf("mongodb://%s:%d", cfg.Host, cfg.Port)
	}

	clientOpts := options.Client().ApplyURI(uri).SetConnectTimeout(10 * time.Second)

	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		return nil, fmt.Errorf("mongo connect error: %w", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("mongo ping error: %w", err)
	}

	db := client.Database(cfg.DBName)
	return &MongoDB{Client: client, Database: db}, nil
}

// Close the MongoDB connection
func (m *MongoDB) Disconnect(ctx context.Context) error {
	if err := m.Client.Disconnect(ctx); err != nil {
		return fmt.Errorf("mongo disconnect error: %w", err)
	}
	return nil
}

func (m *MongoDB) CreateUniqueIndex(ctx context.Context, collectionName, field string, unique bool) error {
	collection := m.Database.Collection(collectionName)

	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: field, Value: 1}},
		Options: options.Index().SetUnique(unique),
	}

	_, err := collection.Indexes().CreateOne(ctx, indexModel)

	return err
}

func (m *MongoDB) Collection(name string) *mongo.Collection {
	return m.Database.Collection(name)
}

func (m *MongoDB) IsDuplicateKeyError(err error) bool {
	if we, ok := err.(mongo.WriteException); ok {
		for _, e := range we.WriteErrors {
			if e.Code == 11000 {
				return true
			}
		}
	}
	return false
}

func (m *MongoDB) SetupValidation(ctx context.Context, schema bson.M, collectionName string) error {
	cmd := bson.D{
		{Key: "collMod", Value: collectionName},
		{Key: "validator", Value: bson.D{
			{Key: "$jsonSchema", Value: schema},
		}},
		{Key: "validationLevel", Value: "strict"},
	}

	var result bson.M
	err := m.Database.RunCommand(ctx, cmd).Decode(&result)
	if err != nil {
		return fmt.Errorf("failed to set collection validator: %w", err)
	}

	// Проверяем успешность выполнения команды
	if ok, _ := result["ok"].(float64); ok != 1 {
		return fmt.Errorf("command failed: %v", result)
	}

	return nil
}
