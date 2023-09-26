package main

import (
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"github.com/naumovrus/finance-transaction-api/internal/cache"
	"github.com/naumovrus/finance-transaction-api/internal/config"
	"github.com/naumovrus/finance-transaction-api/internal/handler"
	"github.com/naumovrus/finance-transaction-api/internal/httpserver"
	"github.com/naumovrus/finance-transaction-api/internal/repository"
	"github.com/naumovrus/finance-transaction-api/internal/service"
	"github.com/sirupsen/logrus"
)

func main() {

	//init logger
	logrus.SetFormatter(new(logrus.JSONFormatter))

	// load .env files
	if err := godotenv.Load(); err != nil {
		log.Fatalf("error load .env files: %s", err.Error())
	}
	cfg := config.LoadConfig()

	db, err := repository.NewPostgresDB(repository.DBConfig{
		Username: cfg.Username,
		Host:     cfg.Host,
		Port:     cfg.PortDb,
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   cfg.Dbname,
		SSLMode:  cfg.Sslmode,
	})
	if err != nil {
		logrus.Fatalf("unable to intitialize db: %s", err.Error())
	}

	repos := repository.NewRepository(db)
	service := service.NewService(repos)

	redisCache := cache.NewRedisCache("redis:6379", 1, 1000, service)
	err = redisCache.SetCachedData()
	if err != nil {
		log.Fatalf("error occured while get data from redis: %s", err)
	}

	handlers := handler.NewHandler(service, redisCache)
	srv := new(httpserver.Server)

	if err := srv.Run("8080", handlers.InitRoutes()); err != nil {
		logrus.Fatalf("error occured while running http server: %s", err.Error())
	}
}

/*
примитивная транзакционная система
пользователи, баланс, отправка

при перезагрузке сервиса не терялась история баланса
не терять транзакцию в случае краша сервиса
*/
