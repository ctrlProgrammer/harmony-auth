package database

import (
	"auth/api/types"
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Roles

const (
	COLLECTION_ROLES = "COLLECTION_ROLES"
)

func AddRole(database *mongo.Database, name string) (bool, error) {
	// Add new documento to the database, it will create a basic role without any kind of access
	// To add access you will need to configure the role with an admin
	col := database.Collection(COLLECTION_ROLES)

	role, err := GetRoleByName(database, name)

	if err != nil && err != mongo.ErrNoDocuments {
		return false, err
	}

	if role != nil {
		return false, errors.New("you already have one role with the same name")
	}

	config := make(map[string]bool)

	insert, err := col.InsertOne(context.Background(), types.Role{
		Name:   name,
		Config: config,
	})

	if err != nil {
		return false, err
	}

	return insert.InsertedID != nil, nil
}

func RemoveRole(database *mongo.Database, id string) (bool, error) {
	col := database.Collection(COLLECTION_ROLES)

	deleted, err := col.DeleteOne(context.Background(), bson.M{"_id": id})

	if err != nil {
		return false, err
	}

	return deleted.DeletedCount == 1, nil
}

func GetRoleByName(database *mongo.Database, name string) (*types.Role, error) {
	col := database.Collection(COLLECTION_ROLES)

	var result types.Role

	err := col.FindOne(context.Background(), bson.M{"name": name}).Decode(&result)

	if err != nil {
		return nil, err
	}

	return &result, err
}

func GetRoles(database *mongo.Database) ([]types.Role, error) {
	col := database.Collection(COLLECTION_ROLES)

	cursor, err := col.Find(context.Background(), bson.M{})

	if err != nil {
		return nil, err
	}

	var results []types.Role

	if err = cursor.All(context.TODO(), &results); err != nil {
		return nil, err
	}

	return results, nil
}

func UpdateRoleConfig(database *mongo.Database, id string, config map[string]bool) (bool, error) {
	// Update all configuration on only one role
	col := database.Collection(COLLECTION_ROLES)

	updated, err := col.UpdateOne(context.Background(), bson.M{"_id": id}, bson.M{"$set": bson.M{"config": config}})

	if err != nil {
		return false, err
	}

	return updated.ModifiedCount == 1, nil
}
