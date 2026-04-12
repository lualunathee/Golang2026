package main

import (
	"log"
	"os"

	"practice7/internal/entity"
	v1 "practice7/internal/controller/http/v1"
	"practice7/internal/usecase"
	"practice7/internal/usecase/repo"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	
	os.Setenv("JWT_SECRET", "my_super_secret_key_123")

	
	dsn := "host=localhost user=postgres password=postgres dbname=seventh port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	
	db.AutoMigrate(&entity.User{})

	userRepo := repo.NewUserRepo(db)
	userUseCase := usecase.NewUserUseCase(userRepo)

	
	router := gin.Default()
	v1Group := router.Group("/v1")
	
	
	v1.NewUserRoutes(v1Group, userUseCase)

	
	log.Println("Server is running on port 8090...")
	router.Run(":8090")
}