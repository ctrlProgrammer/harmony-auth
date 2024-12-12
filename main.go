package main

import (
	"auth/api"
	"context"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

func initializeLogger() (*zap.SugaredLogger, error) {
	logger, err := zap.NewProduction()

	if err != nil {
		return nil, err
	}

	sugar := logger.Sugar()

	return sugar, nil
}

func loadEnv() error {
	err := godotenv.Load("local.env")

	if err != nil {
		return err
	}

	return nil
}

func connectToDatabase(logger *zap.SugaredLogger) *mongo.Database {
	logger.Info("Trying to connect to mongodb at " + os.Getenv("MONGODB"))

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(os.Getenv("MONGODB")))

	if err != nil {
		logger.Error(err.Error())
		return nil
	}

	err = client.Ping(context.Background(), nil)

	if err != nil {
		logger.Error(err.Error())
		return nil
	}

	return client.Database(os.Getenv("MONGODB_DATABASE"))
}

func main() {
	logger, err := initializeLogger()

	if err != nil {
		return
	}

	err = loadEnv()

	if err != nil {
		logger.Error("Error loading env variables (" + err.Error() + ")")
		return
	}

	database := connectToDatabase(logger)

	if database == nil {
		logger.Error("Error loading the database connection")
		return
	}

	apiEngine := api.API{}
	apiEngine.Initialize(logger, database)
}
