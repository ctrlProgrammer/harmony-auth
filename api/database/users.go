package database

import (
	"auth/api/types"
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	COLLECTION_USERS = "COLLECTION_USERS"
)

func AddUser(database *mongo.Database, user types.User) (bool, error) {
	col := database.Collection(COLLECTION_USERS)

	role, err := GetUserByEmail(database, user.Email)

	if err != nil && err != mongo.ErrNoDocuments {
		return false, err
	}

	if role != nil {
		return false, errors.New("you already have one user with the same email")
	}

	insert, err := col.InsertOne(context.Background(), user)

	if err != nil {
		return false, err
	}

	return insert.InsertedID != nil, nil
}

func RemoveUser(database *mongo.Database, id string) (bool, error) {
	col := database.Collection(COLLECTION_USERS)

	deleted, err := col.DeleteOne(context.Background(), bson.M{"_id": id})

	if err != nil {
		return false, err
	}

	return deleted.DeletedCount == 1, nil
}

func GetUserByEmail(database *mongo.Database, email string) (*types.User, error) {
	col := database.Collection(COLLECTION_USERS)

	var result types.User

	err := col.FindOne(context.Background(), bson.M{"email": email}).Decode(&result)

	if err != nil {
		return nil, err
	}

	return &result, err
}

func GetUserRole(database *mongo.Database, id string) (*types.User, error) {
	col := database.Collection(COLLECTION_USERS)

	var result types.User

	err := col.FindOne(context.Background(), bson.M{"_id": id}).Decode(&result)

	if err != nil {
		return nil, err
	}

	return &result, err
}

func GetUserRoleByEmail(database *mongo.Database, email string) (string, error) {
	col := database.Collection(COLLECTION_USERS)

	var result types.User

	err := col.FindOne(context.Background(), bson.M{"email": email}).Decode(&result)

	if err != nil {
		return "", err
	}

	return result.Role, err
}

func GetUsers(database *mongo.Database) ([]types.User, error) {
	col := database.Collection(COLLECTION_USERS)

	cursor, err := col.Find(context.Background(), bson.M{}, options.Find().SetProjection(bson.M{"password": 0}))

	if err != nil {
		return nil, err
	}

	var results []types.User

	if err = cursor.All(context.TODO(), &results); err != nil {
		return nil, err
	}

	return results, nil
}

func UpdateUserRoleByEmail(database *mongo.Database, email string, role string) (bool, error) {
	col := database.Collection(COLLECTION_USERS)

	updated, err := col.UpdateOne(context.Background(), bson.M{"email": email}, bson.M{"$set": bson.M{"role": role}})

	if err != nil {
		return false, err
	}

	return updated.ModifiedCount == 1, nil
}
