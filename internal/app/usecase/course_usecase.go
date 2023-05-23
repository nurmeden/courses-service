package usecase

import (
	"github.com/nurmeden/courses-service/internal/app/model"
	"github.com/nurmeden/courses-service/internal/app/repository"
	"github.com/sirupsen/logrus"
)

type CourseUsecases interface {
	GetCourses() ([]model.Course, error)
	GetCourseByID(id string) (*model.Course, error)
	GetCoursesByStudentID(studentID string) (*model.Course, error)
	CreateCourse(courseInput *model.CourseInput) (*model.Course, error)
	UpdateCourse(id string, courseUpdateInput *model.CourseUpdateInput) (*model.Course, error)
	DeleteCourse(id string) error
	AddStudentToCourse(courseID string, studentID string) (*model.Course, error)
	RemoveStudentFromCourse(courseID string, studentID string) (*model.Course, error)
}

type CourseUsecase struct {
	courseRepo *repository.CourseRepository
	logger     *logrus.Logger
}

func NewCourseUsecase(courseRepo *repository.CourseRepository, logger *logrus.Logger) *CourseUsecase {
	return &CourseUsecase{
		courseRepo: courseRepo,
		logger:     logger,
	}
}

func (cs *CourseUsecase) CreateCourse(courseInput *model.CourseInput) (*model.Course, error) {
	course, err := model.NewCourse(courseInput)
	if err != nil {
		cs.logger.Error("")
		return nil, err
	}
	_, err = cs.courseRepo.CreateCourse(course)
	if err != nil {
		return nil, err
	}

	return course, nil
}

func (cs *CourseUsecase) UpdateCourse(courseID string, courseInput *model.CourseUpdateInput) (*model.Course, error) {
	course, err := cs.courseRepo.GetCourseByID(courseID)
	if err != nil {
		return nil, err
	}

	err = course.Update(courseInput)
	if err != nil {
		return nil, err
	}

	_, err = cs.courseRepo.UpdateCourse(course)
	if err != nil {
		return nil, err
	}

	return course, nil
}

func (cs *CourseUsecase) DeleteCourse(courseID string) error {
	err := cs.courseRepo.DeleteCourse(courseID)
	if err != nil {
		return err
	}

	return nil
}

func (cs *CourseUsecase) GetCourseByID(courseID string) (*model.Course, error) {
	course, err := cs.courseRepo.GetCourseByID(courseID)
	if err != nil {
		return nil, err
	}

	return course, nil
}

func (cs *CourseUsecase) GetAllCourses() ([]*model.Course, error) {
	courses, err := cs.courseRepo.GetAllCourses()
	if err != nil {
		return nil, err
	}

	return courses, nil
}

func (cs *CourseUsecase) GetCoursesByStudentID(studentID string) (*model.Course, error) {
	student, err := cs.courseRepo.GetCoursesByStudentID(studentID)
	if err != nil {
		return nil, err
	}

	return student, nil
}
