package repository

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/insabelter/IWS_GO/models"
)

var (
	ErrFeedbackNotFound = errors.New("feedback not found")
)

type MongoDBRepository struct {
	db *mongo.Database
}

func NewMongoDBRepository(db *mongo.Database) Repository {
	return &MongoDBRepository{db: db}
}

func (r MongoDBRepository) GetAllFeedbacks(ctx context.Context) ([]models.Feedback, error) {
	var out []models.Feedback
	// Create cursor to iterate through the documents
	cursor, err := r.db.
		Collection("feedback").
		Find(ctx, bson.M{})
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return []models.Feedback{}, ErrFeedbackNotFound
		}
		return []models.Feedback{}, err
	}
	// Close the cursor when the function ends
	defer cursor.Close(ctx)
	// Get all the documents from the cursor and decode them into the out variable
	err = cursor.All(ctx, &out)
	if err != nil {
		return []models.Feedback{}, err
	}
	// If there are no documents, return an empty array
	if out == nil {
		return []models.Feedback{}, nil
	}
	return out, nil
}

func (r MongoDBRepository) GetFeedback(ctx context.Context, ID string) (models.Feedback, error) {
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

func (r MongoDBRepository) CreateFeedback(ctx context.Context, feedback models.Feedback) (models.Feedback, error) {
	_, err := r.db.
		Collection("feedback").
		InsertOne(ctx, feedback)
	if err != nil {
		return models.Feedback{}, err
	}
	return feedback, nil
}

func (r MongoDBRepository) DeleteFeedback(ctx context.Context, ID string) error {
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
