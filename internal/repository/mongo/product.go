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
		return product, errors.Wrap(err, utils.ErrorConvert.Error())
	}

	filter := bson.M{"_id": oid}

	result := productRepository.collection.FindOne(ctx, filter)
	if result.Err() != nil {
		return product, errors.Wrap(result.Err(), utils.ErrorExecuteQuery.Error())
	}

	if err = result.Decode(&product); err != nil {
		return product, errors.Wrap(err, utils.ErrorDecode.Error())
	}

	return product, nil
}

func (productRepository *productRepository) GetProductByTitle(ctx context.Context, title string) (*model.Product, error) {
	var product *model.Product

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)

	defer cancel()

	filter := bson.M{"title": title}

	result := productRepository.collection.FindOne(ctx, filter)
	if result.Err() != nil {
		return product, errors.Wrap(result.Err(), utils.ErrorExecuteQuery.Error())
	}

	if err := result.Decode(&product); err != nil {
		return product, errors.Wrap(err, utils.ErrorDecode.Error())
	}

	return product, nil
}

func (productRepository *productRepository) GetAllProducts(ctx context.Context) (*[]model.Product, error) {
	var products []model.Product

	filter := bson.M{}

	cursor, err := productRepository.collection.Find(ctx, filter)
	if err != nil {
		return &products, errors.Wrap(err, utils.ErrorExecuteQuery.Error())
	}

	if err = cursor.All(ctx, &products); err != nil {
		return &products, errors.Wrap(err, utils.ErrorDecode.Error())
	}

	return &products, nil
}

func (productRepository *productRepository) CreateProduct(ctx context.Context, product *model.Product) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)

	defer cancel()

	result, err := productRepository.collection.InsertOne(ctx, product)
	if err != nil {
		return utils.EmptyString, errors.Wrap(err, utils.ErrorExecuteQuery.Error())
	}

	oid, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return utils.EmptyString, errors.Wrap(errors.New("error convert hex to oid"), utils.ErrorConvert.Error())
	}

	return oid.Hex(), nil
}

func (productRepository *productRepository) UpdateProduct(ctx context.Context, product *model.Product) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)

	defer cancel()

	oid, err := primitive.ObjectIDFromHex(product.ID)
	if err != nil {
		return errors.Wrap(err, utils.ErrorConvert.Error())
	}

	filter := bson.M{"_id": oid}

	productByte, err := bson.Marshal(product)
	if err != nil {
		return errors.Wrap(err, utils.ErrorMarshal.Error())
	}

	var object bson.M

	err = bson.Unmarshal(productByte, &object)
	if err != nil {
		return errors.Wrap(err, utils.ErrorUnmarshal.Error())
	}

	delete(object, "_id")

	update := bson.M{
		"$set": object,
	}

	result, err := productRepository.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return errors.Wrap(err, utils.ErrorExecuteQuery.Error())
	}

	if result.MatchedCount == 0 {
		return errors.Wrap(errors.New("not found"), utils.ErrorExecuteQuery.Error())
	}

	return nil
}

func (productRepository *productRepository) DeleteProduct(ctx context.Context, uuid string) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)

	defer cancel()

	oid, err := primitive.ObjectIDFromHex(uuid)
	if err != nil {
		return errors.Wrap(errors.New("error convert hex to oid"), utils.ErrorConvert.Error())
	}

	filter := bson.M{"_id": oid}

	result, err := productRepository.collection.DeleteOne(ctx, filter)
	if err != nil {
		return errors.Wrap(err, utils.ErrorExecuteQuery.Error())
	}

	if result.DeletedCount == 0 {
		return errors.Wrap(errors.New("not found"), utils.ErrorExecuteQuery.Error())
	}

	return nil
}
