package handlers

import (
	"encoding/json"
	"fmt"
	"log"

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nurmeden/courses-service/internal/app/model"
	"github.com/nurmeden/courses-service/internal/app/usecase"
)

type CourseController struct {
	courseService usecase.CourseUsecase
}

func NewCourseController(courseService *usecase.CourseUsecase) *CourseController {
	return &CourseController{
		courseService: *courseService,
	}
}

func (cc *CourseController) CreateCourse(c *gin.Context) {
	var courseInput model.CourseInput

	// Извлечение данных о курсе из тела запроса
	if err := c.ShouldBindJSON(&courseInput); err != nil {
		fmt.Println("ok")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	course, err := cc.courseService.CreateCourse(&courseInput)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": course})
}

func (cc *CourseController) GetCourse(c *gin.Context) {
	courseID := c.Param("id")

	course, err := cc.courseService.GetCourseByID(courseID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": course})
}

func (cc *CourseController) UpdateCourse(c *gin.Context) {
	courseID := c.Param("id")

	var courseUpdateInput model.CourseUpdateInput

	if err := c.ShouldBindJSON(&courseUpdateInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	course, err := cc.courseService.UpdateCourse(courseID, &courseUpdateInput)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": course})
}

func (cc *CourseController) DeleteCourse(c *gin.Context) {
	courseID := c.Param("id")

	err := cc.courseService.DeleteCourse(courseID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Курс успешно удален"})
}

func (cc *CourseController) GetCoursesByStudentID(c *gin.Context) {
	studentID := c.Param("id")
	fmt.Printf("studentID: %v\n", studentID)
	course, err := cc.courseService.GetCoursesByStudentID(studentID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": course})
}

func (cc *CourseController) GetCourseStudentsHandler(c *gin.Context) {
	courseID := c.Param("id")

	resp, err := http.Get("http://localhost:8000/api/students/" + courseID + "/students")
	fmt.Printf("resp: %v\n", resp)
	if err != nil {
		// Обработка ошибки
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get student courses"})
		return
	}
	defer resp.Body.Close()

	fmt.Printf("resp: %v\n", resp.Body)

	var student *model.Student
	err = json.NewDecoder(resp.Body).Decode(&student)
	if err != nil {
		// Обработка ошибки
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode student courses"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"students": student})
}
