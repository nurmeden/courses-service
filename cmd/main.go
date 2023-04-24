package main

import (
	"context"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/nurmeden/courses-service/internal/app/handlers"
	"github.com/nurmeden/courses-service/internal/app/repository"
	"github.com/nurmeden/courses-service/internal/app/usecase"
	"github.com/nurmeden/courses-service/internal/database"
)

func main() {
	client := database.SetupDatabase()
	if client == nil {
		fmt.Println("hero")
		return
	}
	defer client.Disconnect(context.Background())
	fmt.Printf("client: %v\n", client)
	courseRepo, _ := repository.NewCourseRepository(client, "taskdb", "courses")

	courseUsecase := usecase.NewCourseUsecase(courseRepo)

	studentHandler := handlers.NewCourseController(courseUsecase)
	r := gin.Default()
	courseAPI := r.Group("/api/courses")
	{
		courseAPI.POST("/", studentHandler.CreateCourse)
		courseAPI.GET("/:id", studentHandler.GetCourse)
		// courseAPI.GET("/:id", studentHandler.)
		courseAPI.PUT("/:id", studentHandler.UpdateCourse)
		courseAPI.DELETE("/:id", studentHandler.DeleteCourse)
	}
	// courseStudentsAPI := r.Group("/api/courses/:id/students")
	// {
	// 	courseStudentsAPI.GET("/", getCourseStudents)
	// 	courseStudentsAPI.POST("/", addStudentToCourse)
	// 	courseStudentsAPI.DELETE("/:studentID", removeStudentFromCourse)
	// }
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}

}
