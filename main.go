package main

import (
	"github.com/KazakMaxim/auth-jwt-go/pkg/handler"
	"github.com/KazakMaxim/auth-jwt-go/pkg/repository"
	"github.com/KazakMaxim/auth-jwt-go/pkg/service"
	"github.com/KazakMaxim/auth-jwt-go/server"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))
	if err := InitConfig(); err != nil {
		logrus.Fatalf("Error initializing config: %s", err.Error())
	}

	db, err := repository.NewMongoDb()
	if err != nil {
		logrus.Fatalf("Failed to init db")
	}

	//Declaring all lower levels of dependencies
	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handler := handler.NewHandler(services)

	srv := new(server.Server)
	if err := srv.Run(viper.GetString("port"), handler.InitRoutes()); err != nil {
		logrus.Fatalf("Error occured while running http server: %s", err.Error())
	}
}

func InitConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
