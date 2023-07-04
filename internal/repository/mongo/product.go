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

type productRepository struct {
	collection *mongo.Collection
}

func NewProductRepository(storage *mongo.Database, collection string) repository.ProductRepository {
	return &productRepository{
		collection: storage.Collection(collection),
	}
}

func (productRepository *productRepository) GetProduct(ctx context.Context, uuid string) (*model.Product, error) {
	var product *model.Product

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)

	defer cancel()

	oid, err := primitive.ObjectIDFromHex(uuid)
	if err != nil {
		return product, errors.Wrap(utils.ErrorConvert, "error convert hex to oid")
	}

	filter := bson.M{"_id": oid}

	result := productRepository.collection.FindOne(ctx, filter)
	if result.Err() != nil {
		if errors.Is(result.Err(), mongo.ErrNoDocuments) {
			return product, errors.Wrap(utils.ErrorExecuteQuery, err.Error())
		}

		return product, errors.Wrap(utils.ErrorExecuteQuery, err.Error())
	}

	if err = result.Decode(&product); err != nil {
		return product, errors.Wrap(utils.ErrorDecode, err.Error())
	}

	return product, nil
}

func (productRepository *productRepository) GetAllProducts(ctx context.Context) (*[]model.Product, error) {
	var products []model.Product

	filter := bson.M{}

	cursor, err := productRepository.collection.Find(ctx, filter)
	if err != nil {
		return &products, errors.Wrap(utils.ErrorExecuteQuery, err.Error())
	}

	if err = cursor.All(ctx, &products); err != nil {
		return &products, errors.Wrap(utils.ErrorDecode, err.Error())
	}

	return &products, nil
}

func (productRepository *productRepository) CreateProduct(ctx context.Context, product *model.Product) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)

	defer cancel()

	result, err := productRepository.collection.InsertOne(ctx, product)
	if err != nil {
		return utils.EmptyString, errors.Wrap(utils.ErrorExecuteQuery, err.Error())
	}

	oid, ok := result.InsertedID.(primitive.ObjectID)
	if ok {
		return oid.Hex(), nil
	}

	return utils.EmptyString, errors.Wrap(utils.ErrorConvert, "error convert oid to hex")
}

func (productRepository *productRepository) UpdateProduct(ctx context.Context, product *model.Product) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)

	defer cancel()

	oid, err := primitive.ObjectIDFromHex(product.ID)
	if err != nil {
		return errors.Wrap(utils.ErrorConvert, "error convert hex to oid")
	}

	filter := bson.M{"_id": oid}

	productByte, err := bson.Marshal(product)
	if err != nil {
		return errors.Wrap(utils.ErrorMarshal, err.Error())
	}

	var object bson.M

	err = bson.Unmarshal(productByte, &object)
	if err != nil {
		return errors.Wrap(utils.ErrorUnmarshal, err.Error())
	}

	delete(object, "_id")

	update := bson.M{
		"$set": object,
	}

	result, err := productRepository.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return errors.Wrap(utils.ErrorExecuteQuery, err.Error())
	}

	if result.MatchedCount == 0 {
		return errors.Wrap(utils.ErrorExecuteQuery, "not found")
	}

	return nil
}

func (productRepository *productRepository) DeleteProduct(ctx context.Context, uuid string) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)

	defer cancel()

	oid, err := primitive.ObjectIDFromHex(uuid)
	if err != nil {
		return errors.Wrap(utils.ErrorConvert, "error convert hex to oid")
	}

	filter := bson.M{"_id": oid}

	result, err := productRepository.collection.DeleteOne(ctx, filter)
	if err != nil {
		return errors.Wrap(utils.ErrorExecuteQuery, err.Error())
	}

	if result.DeletedCount == 0 {
		return errors.Wrap(utils.ErrorExecuteQuery, "not found")
	}

	return nil
}
