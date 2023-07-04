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

type discountRepository struct {
	collection *mongo.Collection
}

func NewDiscountRepository(storage *mongo.Database, collection string) repository.DiscountRepository {
	return &discountRepository{
		collection: storage.Collection(collection),
	}
}

func (discountRepository *discountRepository) GetDiscount(ctx context.Context, uuid string) (*model.Discount, error) {
	var discount *model.Discount

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)

	defer cancel()

	oid, err := primitive.ObjectIDFromHex(uuid)
	if err != nil {
		return discount, errors.Wrap(err, utils.ErrorConvert.Error())
	}

	filter := bson.M{"_id": oid}

	result := discountRepository.collection.FindOne(ctx, filter)
	if result.Err() != nil {
		return discount, errors.Wrap(result.Err(), utils.ErrorExecuteQuery.Error())
	}

	if err = result.Decode(&discount); err != nil {
		return discount, errors.Wrap(err, utils.ErrorDecode.Error())
	}

	return discount, nil
}

func (discountRepository *discountRepository) GetDiscountByTitle(ctx context.Context, title string) (*model.Discount, error) {
	var discount *model.Discount

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)

	defer cancel()

	filter := bson.M{"title": title}

	result := discountRepository.collection.FindOne(ctx, filter)
	if result.Err() != nil {
		return discount, errors.Wrap(result.Err(), utils.ErrorExecuteQuery.Error())
	}

	if err := result.Decode(&discount); err != nil {
		return discount, errors.Wrap(err, utils.ErrorDecode.Error())
	}

	return discount, nil
}

func (discountRepository *discountRepository) GetAllDiscounts(ctx context.Context) (*[]model.Discount, error) {
	var discounts []model.Discount

	filter := bson.M{}

	cursor, err := discountRepository.collection.Find(ctx, filter)
	if err != nil {
		return &discounts, errors.Wrap(err, utils.ErrorExecuteQuery.Error())
	}

	if err = cursor.All(ctx, &discounts); err != nil {
		return &discounts, errors.Wrap(err, utils.ErrorDecode.Error())
	}

	return &discounts, nil
}

func (discountRepository *discountRepository) CreateDiscount(ctx context.Context, discount *model.Discount) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)

	defer cancel()

	result, err := discountRepository.collection.InsertOne(ctx, discount)
	if err != nil {
		return utils.EmptyString, errors.Wrap(err, utils.ErrorExecuteQuery.Error())
	}

	oid, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return utils.EmptyString, errors.Wrap(errors.New("error convert hex to oid"), utils.ErrorConvert.Error())
	}

	return oid.Hex(), nil
}

func (discountRepository *discountRepository) UpdateDiscount(ctx context.Context, discount *model.Discount) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)

	defer cancel()

	oid, err := primitive.ObjectIDFromHex(discount.ID)
	if err != nil {
		return errors.Wrap(err, utils.ErrorConvert.Error())
	}

	filter := bson.M{"_id": oid}

	discountByte, err := bson.Marshal(discount)
	if err != nil {
		return errors.Wrap(err, utils.ErrorMarshal.Error())
	}

	var object bson.M

	err = bson.Unmarshal(discountByte, &object)
	if err != nil {
		return errors.Wrap(err, utils.ErrorUnmarshal.Error())
	}

	delete(object, "_id")

	update := bson.M{
		"$set": object,
	}

	result, err := discountRepository.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return errors.Wrap(err, utils.ErrorExecuteQuery.Error())
	}

	if result.MatchedCount == 0 {
		return errors.Wrap(errors.New("not found"), utils.ErrorExecuteQuery.Error())
	}

	return nil
}

func (discountRepository *discountRepository) DeleteDiscount(ctx context.Context, uuid string) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)

	defer cancel()

	oid, err := primitive.ObjectIDFromHex(uuid)
	if err != nil {
		return errors.Wrap(errors.New("error convert hex to oid"), utils.ErrorConvert.Error())
	}

	filter := bson.M{"_id": oid}

	result, err := discountRepository.collection.DeleteOne(ctx, filter)
	if err != nil {
		return errors.Wrap(err, utils.ErrorExecuteQuery.Error())
	}

	if result.DeletedCount == 0 {
		return errors.Wrap(errors.New("not found"), utils.ErrorExecuteQuery.Error())
	}

	return nil
}
