package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/nurmeden/courses-service/internal/app/handlers"
	"github.com/nurmeden/courses-service/internal/app/repository"
	"github.com/nurmeden/courses-service/internal/app/usecase"
	"github.com/nurmeden/courses-service/internal/database"
	"github.com/sirupsen/logrus"
)

func main() {
	logfile, err := os.OpenFile("logrus.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer logfile.Close()
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(logfile)
	logrus.SetLevel(logrus.DebugLevel)

	client := database.SetupDatabase()
	if client == nil {
		fmt.Println("hero")
		return
	}
	defer client.Disconnect(context.Background())
	fmt.Printf("client: %v\n", client)
	courseRepo, _ := repository.NewCourseRepository(client, "coursesdb", "courses")

	courseUsecase := usecase.NewCourseUsecase(courseRepo)

	studentHandler := handlers.NewCourseController(courseUsecase)
	r := gin.Default()
	courseAPI := r.Group("/api/courses")
	{
		courseAPI.POST("/", studentHandler.CreateCourse)
		courseAPI.GET("/:id", studentHandler.GetCourse)
		courseAPI.GET("/:id/courses", studentHandler.GetCoursesByStudentID)
		courseAPI.GET("/:id/students", studentHandler.GetCourseStudentsHandler)
		courseAPI.PUT("/:id", studentHandler.UpdateCourse)
		courseAPI.DELETE("/:id", studentHandler.DeleteCourse)
	}
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}

}
