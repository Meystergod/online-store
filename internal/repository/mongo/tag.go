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

type tagRepository struct {
	collection *mongo.Collection
}

func NewTagRepository(storage *mongo.Database, collection string) repository.TagRepository {
	return &tagRepository{
		collection: storage.Collection(collection),
	}
}

func (tagRepository *tagRepository) GetTag(ctx context.Context, uuid string) (*model.Tag, error) {
	var tag *model.Tag

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)

	defer cancel()

	oid, err := primitive.ObjectIDFromHex(uuid)
	if err != nil {
		return tag, errors.Wrap(utils.ErrorConvert, "error convert hex to oid")
	}

	filter := bson.M{"_id": oid}

	result := tagRepository.collection.FindOne(ctx, filter)
	if result.Err() != nil {
		if errors.Is(result.Err(), mongo.ErrNoDocuments) {
			return tag, errors.Wrap(utils.ErrorExecuteQuery, err.Error())
		}

		return tag, errors.Wrap(utils.ErrorExecuteQuery, err.Error())
	}

	if err = result.Decode(&tag); err != nil {
		return tag, errors.Wrap(utils.ErrorDecode, err.Error())
	}

	return tag, nil
}

func (tagRepository *tagRepository) GetAllTags(ctx context.Context) (*[]model.Tag, error) {
	var tags []model.Tag

	filter := bson.M{}

	cursor, err := tagRepository.collection.Find(ctx, filter)
	if err != nil {
		return &tags, errors.Wrap(utils.ErrorExecuteQuery, err.Error())
	}

	if err = cursor.All(ctx, &tags); err != nil {
		return &tags, errors.Wrap(utils.ErrorDecode, err.Error())
	}

	return &tags, nil
}

func (tagRepository *tagRepository) CreateTag(ctx context.Context, tag *model.Tag) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)

	defer cancel()

	result, err := tagRepository.collection.InsertOne(ctx, tag)
	if err != nil {
		return utils.EmptyString, errors.Wrap(utils.ErrorExecuteQuery, err.Error())
	}

	oid, ok := result.InsertedID.(primitive.ObjectID)
	if ok {
		return oid.Hex(), nil
	}

	return utils.EmptyString, errors.Wrap(utils.ErrorConvert, "error convert oid to hex")
}

func (tagRepository *tagRepository) UpdateTag(ctx context.Context, tag *model.Tag) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)

	defer cancel()

	oid, err := primitive.ObjectIDFromHex(tag.ID)
	if err != nil {
		return errors.Wrap(utils.ErrorConvert, "error convert hex to oid")
	}

	filter := bson.M{"_id": oid}

	tagByte, err := bson.Marshal(tag)
	if err != nil {
		return errors.Wrap(utils.ErrorMarshal, err.Error())
	}

	var object bson.M

	err = bson.Unmarshal(tagByte, &object)
	if err != nil {
		return errors.Wrap(utils.ErrorUnmarshal, err.Error())
	}

	delete(object, "_id")

	update := bson.M{
		"$set": object,
	}

	result, err := tagRepository.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return errors.Wrap(utils.ErrorExecuteQuery, err.Error())
	}

	if result.MatchedCount == 0 {
		return errors.Wrap(utils.ErrorExecuteQuery, "not found")
	}

	return nil
}

func (tagRepository *tagRepository) DeleteTag(ctx context.Context, uuid string) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)

	defer cancel()

	oid, err := primitive.ObjectIDFromHex(uuid)
	if err != nil {
		return errors.Wrap(utils.ErrorConvert, "error convert hex to oid")
	}

	filter := bson.M{"_id": oid}

	result, err := tagRepository.collection.DeleteOne(ctx, filter)
	if err != nil {
		return errors.Wrap(utils.ErrorExecuteQuery, err.Error())
	}

	if result.DeletedCount == 0 {
		return errors.Wrap(utils.ErrorExecuteQuery, "not found")
	}

	return nil
}
