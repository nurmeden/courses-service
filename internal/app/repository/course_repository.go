package repository

import (
	"context"
	"errors"

	"github.com/go-redis/redis"
	"github.com/nurmeden/courses-service/internal/app/model"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CourseRepository struct {
	client     *mongo.Client
	collection *mongo.Collection
	cache      *redis.Client
	logger     *logrus.Logger
}

func NewCourseRepository(client *mongo.Client, dbName string, collectionName string, cache *redis.Client, logger *logrus.Logger) (*CourseRepository, error) {
	r := &CourseRepository{
		client: client,
		cache:  cache,
		logger: logger,
	}
	collection := client.Database(dbName).Collection(collectionName)
	r.collection = collection

	return r, nil
}

func (r *CourseRepository) GetCourseByID(id string) (*model.Course, error) {
	var course model.Course

	filter := bson.M{"id": id}

	err := r.collection.FindOne(context.Background(), filter).Decode(&course)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("course not found")
		}
		return nil, err
	}

	return &course, nil
}

func (r *CourseRepository) GetAllCourses() ([]*model.Course, error) {
	var courses []*model.Course

	cursor, err := r.collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	err = cursor.All(context.Background(), &courses)
	if err != nil {
		return nil, err
	}

	return courses, nil
}

func (r *CourseRepository) CreateCourse(course *model.Course) (*model.Course, error) {

	_, err := r.collection.InsertOne(context.Background(), course)
	if err != nil {
		return nil, err
	}

	return course, nil
}

func (r *CourseRepository) UpdateCourse(course *model.Course) (*model.Course, error) {
	filter := bson.M{"id": course.ID}

	result := r.collection.FindOneAndUpdate(context.Background(), filter, bson.M{"$set": course}, options.FindOneAndUpdate().SetReturnDocument(options.After))
	if result.Err() != nil {
		if errors.Is(result.Err(), mongo.ErrNoDocuments) {
			return nil, errors.New("course not found")
		}
		return nil, result.Err()
	}

	updatedCourse := &model.Course{}
	err := result.Decode(updatedCourse)
	if err != nil {
		return nil, err
	}

	return updatedCourse, nil
}

func (r *CourseRepository) DeleteCourse(id string) error {
	filter := bson.M{"_id": id}
	result, err := r.collection.DeleteOne(context.Background(), filter)
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("course not found")
	}

	return nil
}

func (r *CourseRepository) GetCoursesByStudentID(id string) (*model.Course, error) {
	var course model.Course

	filter := bson.M{"students": id}

	err := r.collection.FindOne(context.Background(), filter).Decode(&course)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("course not found")
		}
		return nil, err
	}

	return &course, nil
}
