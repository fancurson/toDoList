package main

import (
	"fmt"
	"os"

	todo "github.com/fancurson/toDoList"
	"github.com/fancurson/toDoList/pkg/handler"
	"github.com/fancurson/toDoList/pkg/repository"
	"github.com/fancurson/toDoList/pkg/service"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func InitConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}

func main() {

	must(InitConfig(), "error while reading configs")

	must(godotenv.Load(), "opening .env file error")

	db, err := repository.NewPostgresDB(repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
	})
	must(err, "Connecting to database errror")

	repos := repository.NewRepository(db)
	serv := service.NewService(repos)
	handler := handler.NewHandler(serv)

	srv := new(todo.Server)

	fmt.Println("Starting server on :3000...")
	if err := srv.Run(viper.GetString("port"), handler.InitRouters()); err != nil {
		logrus.Fatalf("Error occured while starting http server: %w", err)
	}
}

func must(err error, message string) {
	if err != nil {
		logrus.Fatalf("%s: %v", message, err.Error())
	}
}
