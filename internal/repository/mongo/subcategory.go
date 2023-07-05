package mongo

import (
	"context"
	"time"

	"github.com/Meystergod/online-store/internal/domain/model"
	"github.com/Meystergod/online-store/internal/repository"
	"github.com/Meystergod/online-store/internal/utils"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type subcategoryRepository struct {
	collection *mongo.Collection
}

func NewSubcategoryRepository(storage *mongo.Database, collection string) repository.SubcategoryRepository {
	return &subcategoryRepository{
		collection: storage.Collection(collection),
	}
}

func (subcategoryRepository *subcategoryRepository) GetSubcategory(ctx context.Context, uuid string) (*model.Subcategory, error) {
	var subcategory *model.Subcategory

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)

	defer cancel()

	oid, err := primitive.ObjectIDFromHex(uuid)
	if err != nil {
		return subcategory, errors.Wrap(err, utils.ErrorConvert.Error())
	}

	filter := bson.M{"_id": oid}

	result := subcategoryRepository.collection.FindOne(ctx, filter)
	if result.Err() != nil {
		return subcategory, errors.Wrap(result.Err(), utils.ErrorExecuteQuery.Error())
	}

	if err = result.Decode(&subcategory); err != nil {
		return subcategory, errors.Wrap(err, utils.ErrorDecode.Error())
	}

	return subcategory, nil
}

func (subcategoryRepository *subcategoryRepository) GetSubcategoryByTitle(ctx context.Context, title string) (*model.Subcategory, error) {
	var subcategory *model.Subcategory

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)

	defer cancel()

	filter := bson.M{"title": title}

	result := subcategoryRepository.collection.FindOne(ctx, filter)
	if result.Err() != nil {
		return subcategory, errors.Wrap(result.Err(), utils.ErrorExecuteQuery.Error())
	}

	if err := result.Decode(&subcategory); err != nil {
		return subcategory, errors.Wrap(err, utils.ErrorDecode.Error())
	}

	return subcategory, nil
}

func (subcategoryRepository *subcategoryRepository) GetAllSubcategories(ctx context.Context) (*[]model.Subcategory, error) {
	var subcategories []model.Subcategory

	filter := bson.M{}

	cursor, err := subcategoryRepository.collection.Find(ctx, filter)
	if err != nil {
		return &subcategories, errors.Wrap(err, utils.ErrorExecuteQuery.Error())
	}

	if err = cursor.All(ctx, &subcategories); err != nil {
		return &subcategories, errors.Wrap(err, utils.ErrorDecode.Error())
	}

	return &subcategories, nil
}

func (subcategoryRepository *subcategoryRepository) CreateSubcategory(ctx context.Context, subcategory *model.Subcategory) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)

	defer cancel()

	result, err := subcategoryRepository.collection.InsertOne(ctx, subcategory)
	if err != nil {
		return utils.EmptyString, errors.Wrap(err, utils.ErrorExecuteQuery.Error())
	}

	oid, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return utils.EmptyString, errors.Wrap(errors.New("error convert hex to oid"), utils.ErrorConvert.Error())
	}

	return oid.Hex(), nil
}

func (subcategoryRepository *subcategoryRepository) UpdateSubcategory(ctx context.Context, category *model.Subcategory) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)

	defer cancel()

	oid, err := primitive.ObjectIDFromHex(category.ID)
	if err != nil {
		return errors.Wrap(err, utils.ErrorConvert.Error())
	}

	filter := bson.M{"_id": oid}

	categoryByte, err := bson.Marshal(category)
	if err != nil {
		return errors.Wrap(err, utils.ErrorMarshal.Error())
	}

	var object bson.M

	err = bson.Unmarshal(categoryByte, &object)
	if err != nil {
		return errors.Wrap(err, utils.ErrorUnmarshal.Error())
	}

	delete(object, "_id")

	update := bson.M{
		"$set": object,
	}

	result, err := subcategoryRepository.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return errors.Wrap(err, utils.ErrorExecuteQuery.Error())
	}

	if result.MatchedCount == 0 {
		return errors.Wrap(errors.New("not found"), utils.ErrorExecuteQuery.Error())
	}

	return nil
}

func (subcategoryRepository *subcategoryRepository) DeleteSubcategory(ctx context.Context, uuid string) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)

	defer cancel()

	oid, err := primitive.ObjectIDFromHex(uuid)
	if err != nil {
		return errors.Wrap(errors.New("error convert hex to oid"), utils.ErrorConvert.Error())
	}

	filter := bson.M{"_id": oid}

	result, err := subcategoryRepository.collection.DeleteOne(ctx, filter)
	if err != nil {
		return errors.Wrap(err, utils.ErrorExecuteQuery.Error())
	}

	if result.DeletedCount == 0 {
		return errors.Wrap(errors.New("not found"), utils.ErrorExecuteQuery.Error())
	}

	return nil
}
