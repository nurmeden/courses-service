package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/joho/godotenv"
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
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(logfile)
	logger.SetLevel(logrus.DebugLevel)

	if err := godotenv.Load(".env"); err != nil {
		log.Print("No .env file found")
	}
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})

	defer redisClient.Close()

	_, err = redisClient.Ping().Result()
	if err != nil {
		fmt.Println(err.Error())
		logger.Fatal("Ошибка подключения к Redis:", err)
	}

	dbName := os.Getenv("DATABASE_NAME")
	mongoURI := os.Getenv("MONGODB_URI")
	collectionName := os.Getenv("COLLECTION_NAME")

	client, err := database.SetupDatabase(context.Background(), mongoURI)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer client.Disconnect(context.Background())
	fmt.Printf("client: %v\n", client)
	courseRepo, _ := repository.NewCourseRepository(client, dbName, collectionName, redisClient, logger)

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
