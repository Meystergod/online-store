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

type categoryRepository struct {
	collection *mongo.Collection
}

func NewCategoryRepository(storage *mongo.Database, collection string) repository.CategoryRepository {
	return &categoryRepository{
		collection: storage.Collection(collection),
	}
}

func (categoryRepository *categoryRepository) GetCategory(ctx context.Context, uuid string) (*model.Category, error) {
	var category *model.Category

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)

	defer cancel()

	oid, err := primitive.ObjectIDFromHex(uuid)
	if err != nil {
		return category, errors.Wrap(utils.ErrorConvert, "error convert hex to oid")
	}

	filter := bson.M{"_id": oid}

	result := categoryRepository.collection.FindOne(ctx, filter)
	if result.Err() != nil {
		if errors.Is(result.Err(), mongo.ErrNoDocuments) {
			return category, errors.Wrap(utils.ErrorExecuteQuery, err.Error())
		}

		return category, errors.Wrap(utils.ErrorExecuteQuery, err.Error())
	}

	if err = result.Decode(&category); err != nil {
		return category, errors.Wrap(utils.ErrorDecode, err.Error())
	}

	return category, nil
}

func (categoryRepository *categoryRepository) GetAllCategories(ctx context.Context) (*[]model.Category, error) {
	var categories []model.Category

	filter := bson.M{}

	cursor, err := categoryRepository.collection.Find(ctx, filter)
	if err != nil {
		return &categories, errors.Wrap(utils.ErrorExecuteQuery, err.Error())
	}

	if err = cursor.All(ctx, &categories); err != nil {
		return &categories, errors.Wrap(utils.ErrorDecode, err.Error())
	}

	return &categories, nil
}

func (categoryRepository *categoryRepository) CreateCategory(ctx context.Context, category *model.Category) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)

	defer cancel()

	result, err := categoryRepository.collection.InsertOne(ctx, category)
	if err != nil {
		return utils.EmptyString, errors.Wrap(utils.ErrorExecuteQuery, err.Error())
	}

	oid, ok := result.InsertedID.(primitive.ObjectID)
	if ok {
		return oid.Hex(), nil
	}

	return utils.EmptyString, errors.Wrap(utils.ErrorConvert, "error convert oid to hex")
}

func (categoryRepository *categoryRepository) UpdateCategory(ctx context.Context, category *model.Category) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)

	defer cancel()

	oid, err := primitive.ObjectIDFromHex(category.ID)
	if err != nil {
		return errors.Wrap(utils.ErrorConvert, "error convert hex to oid")
	}

	filter := bson.M{"_id": oid}

	categoryByte, err := bson.Marshal(category)
	if err != nil {
		return errors.Wrap(utils.ErrorMarshal, err.Error())
	}

	var object bson.M

	err = bson.Unmarshal(categoryByte, &object)
	if err != nil {
		return errors.Wrap(utils.ErrorUnmarshal, err.Error())
	}

	delete(object, "_id")

	update := bson.M{
		"$set": object,
	}

	result, err := categoryRepository.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return errors.Wrap(utils.ErrorExecuteQuery, err.Error())
	}

	if result.MatchedCount == 0 {
		return errors.Wrap(utils.ErrorExecuteQuery, "not found")
	}

	return nil
}

func (categoryRepository *categoryRepository) DeleteCategory(ctx context.Context, uuid string) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)

	defer cancel()

	oid, err := primitive.ObjectIDFromHex(uuid)
	if err != nil {
		return errors.Wrap(utils.ErrorConvert, "error convert hex to oid")
	}

	filter := bson.M{"_id": oid}

	result, err := categoryRepository.collection.DeleteOne(ctx, filter)
	if err != nil {
		return errors.Wrap(utils.ErrorExecuteQuery, err.Error())
	}

	if result.DeletedCount == 0 {
		return errors.Wrap(utils.ErrorExecuteQuery, "not found")
	}

	return nil
}
