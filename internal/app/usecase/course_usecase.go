package usecase

import (
	"github.com/nurmeden/courses-service/internal/app/model"
	"github.com/nurmeden/courses-service/internal/app/repository"
)

type CourseUsecases interface {
	GetCourses() ([]model.Course, error)
	GetCourseByID(id string) (*model.Course, error)
	CreateCourse(courseInput *model.CourseInput) (*model.Course, error)
	UpdateCourse(id string, courseUpdateInput *model.CourseUpdateInput) (*model.Course, error)
	DeleteCourse(id string) error
	AddStudentToCourse(courseID string, studentID string) (*model.Course, error)
	RemoveStudentFromCourse(courseID string, studentID string) (*model.Course, error)
}

// CourseService представляет сервис для работы с курсами
type CourseUsecase struct {
	courseRepo *repository.CourseRepository
}

// NewCourseService создает новый экземпляр сервиса курсов
func NewCourseUsecase(courseRepo *repository.CourseRepository) *CourseUsecase {
	return &CourseUsecase{
		courseRepo: courseRepo,
	}
}

// CreateCourse создает новый курс на основе входных параметров
func (cs *CourseUsecase) CreateCourse(courseInput *model.CourseInput) (*model.Course, error) {
	course, err := model.NewCourse(courseInput)
	if err != nil {
		return nil, err
	}
	_, err = cs.courseRepo.CreateCourse(course)
	if err != nil {
		return nil, err
	}

	return course, nil
}

// UpdateCourse обновляет информацию о курсе на основе входных параметров
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

// DeleteCourse удаляет курс по его идентификатору
func (cs *CourseUsecase) DeleteCourse(courseID string) error {
	err := cs.courseRepo.DeleteCourse(courseID)
	if err != nil {
		return err
	}

	return nil
}

// GetCourseByID возвращает курс по его идентификатору
func (cs *CourseUsecase) GetCourseByID(courseID string) (*model.Course, error) {
	course, err := cs.courseRepo.GetCourseByID(courseID)
	if err != nil {
		return nil, err
	}

	return course, nil
}

// GetAllCourses возвращает все курсы
func (cs *CourseUsecase) GetAllCourses() ([]*model.Course, error) {
	courses, err := cs.courseRepo.GetAllCourses()
	if err != nil {
		return nil, err
	}

	return courses, nil
}

// func (cs *CourseUsecase) AddStudentToCourse(courseID string, studentID string) (*model.Course, error) {
// 	course, err := cs.courseRepo.GetCourseByID(courseID)
// 	if err != nil {
// 		return nil, err
// 	}

// 	student, err := studentRepository.Read()
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Добавляем идентификатор студента в список студентов курса
// 	course.Students = append(course.Students, studentID)

// 	// Обновляем информацию о курсе в репозитории
// 	err = courseRepository.UpdateCourse(course)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return course, nil
// }
