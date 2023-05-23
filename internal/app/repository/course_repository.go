package repository

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/go-redis/redis"
	"github.com/nurmeden/courses-service/internal/app/model"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println("Invalid id")
	}
	filter := bson.M{"_id": objectId}

	r.logger.Infof("Getting course by ID: %s", id)

	err = r.collection.FindOne(context.Background(), filter).Decode(&course)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			fmt.Printf("err.Error() in errors is: %v\n", err.Error())
			r.logger.Errorf("Course not found: %v", err)
			return nil, errors.New("course not found")
		}
		r.logger.Errorf("Failed to get course by ID: %v", err)
		return nil, err
	}

	r.logger.Infof("Got course by ID: %s", id)

	return &course, nil
}

func (r *CourseRepository) GetAllCourses() ([]*model.Course, error) {
	var courses []*model.Course

	cursor, err := r.collection.Find(context.Background(), bson.M{})
	if err != nil {
		r.logger.Errorf("Failed to get all courses: %v", err)
		return nil, err
	}
	defer cursor.Close(context.Background())

	err = cursor.All(context.Background(), &courses)
	if err != nil {
		r.logger.Errorf("Failed to decode courses: %v", err)
		return nil, err
	}

	r.logger.Infof("Found %d courses", len(courses))
	return courses, nil
}

func (r *CourseRepository) CreateCourse(course *model.Course) (*model.Course, error) {
	_, err := r.collection.InsertOne(context.Background(), course)
	if err != nil {
		r.logger.Errorf("Failed to create course: %v", err)
		return nil, err
	}

	r.logger.Infof("Created course: %v", course)
	return course, nil
}

func (r *CourseRepository) UpdateCourse(course *model.Course) (*model.Course, error) {
	objectID, err := primitive.ObjectIDFromHex(course.ID)
	if err != nil {
		log.Println("Invalid ID")
		return nil, err
	}

	filter := bson.M{"_id": objectID}
	update := bson.M{
		"$set": bson.M{
			// "name":        course.Name,
			// "description": course.Description,
			"students": course.Students,
		},
	}

	options := options.FindOneAndUpdate().SetReturnDocument(options.After)
	result := r.collection.FindOneAndUpdate(context.Background(), filter, update, options)
	if result.Err() != nil {
		if errors.Is(result.Err(), mongo.ErrNoDocuments) {
			r.logger.Errorf("Course not found: %v", result.Err())
			return nil, errors.New("course not found")
		}
		r.logger.Errorf("Failed to update course: %v", result.Err())
		return nil, result.Err()
	}

	updatedCourse := &model.Course{}
	err = result.Decode(updatedCourse)
	if err != nil {
		r.logger.Errorf("Failed to decode updated course: %v", err)
		return nil, err
	}

	r.logger.Infof("Updated course: %v", updatedCourse)
	return updatedCourse, nil
}

func (r *CourseRepository) DeleteCourse(id string) error {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println("Invalid id")
	}
	filter := bson.M{"_id": objectId}
	result, err := r.collection.DeleteOne(context.Background(), filter)
	if err != nil {
		r.logger.Errorf("Failed to delete course: %v", err)
		return err
	}

	if result.DeletedCount == 0 {
		r.logger.Errorf("Course not found")
		return errors.New("course not found")
	}

	r.logger.Infof("Deleted course with id %v", id)
	return nil
}

func (r *CourseRepository) GetCoursesByStudentID(id string) ([]*model.Course, error) {
	// objectID, err := primitive.ObjectIDFromHex(id)
	// if err != nil {
	// 	log.Println("Invalid ID")
	// 	return nil, err
	// }

	filter := bson.M{"students": bson.M{"$in": []string{id}}}

	cursor, err := r.collection.Find(context.Background(), filter)
	if err != nil {
		fmt.Printf("err.Error(): coursesFind %v\n", err.Error())
		r.logger.Errorf("Failed to get courses for student with ID %v: %v", id, err)
		return nil, err
	}
	defer cursor.Close(context.Background())

	courses := []*model.Course{}
	for cursor.Next(context.Background()) {
		var course model.Course
		if err := cursor.Decode(&course); err != nil {
			r.logger.Errorf("Error decoding course: %s", err)
			continue
		}
		courses = append(courses, &course)
	}

	if err := cursor.Err(); err != nil {
		r.logger.Errorf("Cursor error: %v", err)
		return nil, err
	}

	if len(courses) == 0 {
		r.logger.Infof("No courses found for student with ID %v", id)
		return nil, nil
	}
	fmt.Printf("courses: %v\n", courses)

	r.logger.Infof("Found courses for student with ID %v", id)
	return courses, nil
}
