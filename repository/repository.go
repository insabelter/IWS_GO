package repository

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/insabelter/IWS_GO/models"
)

var (
	ErrFeedbackNotFound = errors.New("feedback not found")
)

type repository struct {
	db *mongo.Database
}

func NewRepository(db *mongo.Database) Repository {
	return &repository{db: db}
}

func (r repository) GetFeedback(ctx context.Context, ID string) (models.Feedback, error) {
	var out models.Feedback
	err := r.db.
		Collection("feedback").
		FindOne(ctx, bson.M{"id": ID}).
		Decode(&out)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return models.Feedback{}, ErrFeedbackNotFound
		}
		return models.Feedback{}, err
	}
	return out, nil
}

func (r repository) CreateFeedback(ctx context.Context, feedback models.Feedback) (models.Feedback, error) {
	out, err := r.db.
		Collection("feedback").
		InsertOne(ctx, feedback)
	if err != nil {
		return models.Feedback{}, err
	}
	insertedID := out.InsertedID.(primitive.ObjectID).String()
	fmt.Println(insertedID)
	return feedback, nil
}

func (r repository) DeleteFeedback(ctx context.Context, ID string) error {
	out, err := r.db.
		Collection("feedback").
		DeleteOne(ctx, bson.M{"id": ID})
	if err != nil {
		return err
	}
	if out.DeletedCount == 0 {
		return ErrFeedbackNotFound
	}
	return nil
}
