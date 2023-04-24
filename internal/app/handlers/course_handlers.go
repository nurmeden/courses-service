package handlers

import (
	"fmt"

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nurmeden/courses-service/internal/app/model"
	"github.com/nurmeden/courses-service/internal/app/usecase"
)

// CourseController контроллер для операций с курсами
type CourseController struct {
	courseService usecase.CourseUsecase
}

// NewCourseController создает новый экземпляр контроллера CourseController
func NewCourseController(courseService *usecase.CourseUsecase) *CourseController {
	return &CourseController{
		courseService: *courseService,
	}
}

// CreateCourse обрабатывает запрос на создание нового курса
func (cc *CourseController) CreateCourse(c *gin.Context) {
	var courseInput model.CourseInput

	// Извлечение данных о курсе из тела запроса
	if err := c.ShouldBindJSON(&courseInput); err != nil {
		fmt.Println("ok")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Создание нового курса
	course, err := cc.courseService.CreateCourse(&courseInput)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": course})
}

// GetCourse обрабатывает запрос на получение информации о курсе по его ID
func (cc *CourseController) GetCourse(c *gin.Context) {
	courseID := c.Param("id")

	// Получение информации о курсе по его ID
	course, err := cc.courseService.GetCourseByID(courseID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": course})
}

// UpdateCourse обрабатывает запрос на обновление информации о курсе
func (cc *CourseController) UpdateCourse(c *gin.Context) {
	courseID := c.Param("id")

	var courseUpdateInput model.CourseUpdateInput

	// Извлечение данных для обновления курса из тела запроса
	if err := c.ShouldBindJSON(&courseUpdateInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Обновление информации о курсе
	course, err := cc.courseService.UpdateCourse(courseID, &courseUpdateInput)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": course})
}

// DeleteCourse обрабатывает запрос на удаление курса по его ID
func (cc *CourseController) DeleteCourse(c *gin.Context) {
	courseID := c.Param("id")

	// Удаление курса
	err := cc.courseService.DeleteCourse(courseID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Курс успешно удален"})
}

// func (cc *CourseController) GetCoursesByStudentID(c *gin.Context) {
// 	studentID := c.Param("id")
// 	courses, err := cc.courseService.GetCoursesByStudentID(studentID)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"data": courses})
// }
