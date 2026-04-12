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
	// 1. Секретный ключ для JWT (обычно берется из .env, но для теста зададим так)
	os.Setenv("JWT_SECRET", "my_super_secret_key_123")

	// 2. Подключение к PostgreSQL
	// ВАЖНО: Замени данные (user, password, dbname) на свои настройки Postgres!
	dsn := "host=localhost user=postgres password=твои_пароль dbname=postgres port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Автомиграция (создаст таблицу users в базе)
	db.AutoMigrate(&entity.User{})

	// 3. Инициализация слоев (Repo -> Usecase)
	userRepo := repo.NewUserRepo(db)
	userUseCase := usecase.NewUserUseCase(userRepo)

	// 4. Настройка роутера Gin
	router := gin.Default()
	v1Group := router.Group("/v1")
	
	// Подключаем наши роуты
	v1.NewUserRoutes(v1Group, userUseCase)

	// 5. Запуск сервера
	log.Println("Server is running on port 8090...")
	router.Run(":8090")
}