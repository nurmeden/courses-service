package repository

import (
	"context"
	"errors"

	"github.com/nurmeden/courses-service/internal/app/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// type CourseRepository interface {
// 	GetCourseByID(id string) (*model.Course, error)
// 	GetAllCourses() ([]*model.Course, error)
// 	CreateCourse(course *model.Course) (*model.Course, error)
// 	UpdateCourse(course *model.Course) (*model.Course, error)
// 	DeleteCourse(id string) error
// }

type CourseRepository struct {
	client     *mongo.Client
	collection *mongo.Collection
}

func NewCourseRepository(client *mongo.Client, dbName string, collectionName string) (*CourseRepository, error) {
	r := &CourseRepository{
		client: client,
	}
	// Получение коллекции студентов
	collection := client.Database(dbName).Collection(collectionName)
	r.collection = collection

	return r, nil
}

func (r *CourseRepository) GetCourseByID(id string) (*model.Course, error) {
	var course model.Course

	// Формирование фильтра по идентификатору
	filter := bson.M{"_id": id}

	// Поиск курса по идентификатору
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

	// Поиск всех курсов
	cursor, err := r.collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	// Декодирование результатов запроса в слайс курсов
	err = cursor.All(context.Background(), &courses)
	if err != nil {
		return nil, err
	}

	return courses, nil
}

func (r *CourseRepository) CreateCourse(course *model.Course) (*model.Course, error) {
	// Установка даты создания и обновления курса
	// now := time.Now()
	// course.CreatedAt = now
	// course.UpdatedAt = now

	// Вставка курса в базу данных
	_, err := r.collection.InsertOne(context.Background(), course)
	if err != nil {
		return nil, err
	}

	return course, nil
}

func (r *CourseRepository) UpdateCourse(course *model.Course) (*model.Course, error) {
	// Установка даты обновления курса
	// now := time.Now()
	// course.UpdatedAt = now

	// Формирование фильтра по идентификатору
	filter := bson.M{"_id": course.ID}

	// Обновление курса в базе данных
	result := r.collection.FindOneAndUpdate(context.Background(), filter, bson.M{"$set": course}, options.FindOneAndUpdate().SetReturnDocument(options.After))
	if result.Err() != nil {
		if errors.Is(result.Err(), mongo.ErrNoDocuments) {
			return nil, errors.New("course not found")
		}
		return nil, result.Err()
	}

	// Декодирование обновленного курса
	updatedCourse := &model.Course{}
	err := result.Decode(updatedCourse)
	if err != nil {
		return nil, err
	}

	return updatedCourse, nil
}

func (r *CourseRepository) DeleteCourse(id string) error {
	// Формирование фильтра по идентификатору
	filter := bson.M{"_id": id}
	// Удаление курса из базы данных
	result, err := r.collection.DeleteOne(context.Background(), filter)
	if err != nil {
		return err
	}

	// Проверка наличия удаленных документов
	if result.DeletedCount == 0 {
		return errors.New("course not found")
	}

	return nil
}
