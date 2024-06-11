package main

import (
	"go-server/internal/app/router"
	"go-server/internal/pkg/migrations"
	"go-server/pkg/shared/database"
	"go-server/pkg/shared/logging"
	"go-server/pkg/shared/logging/hooks"
	"os"

	"github.com/gin-gonic/gin"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	logger := logging.NewLogger()

	logger.AddHook(&hooks.RequestIDHook{})

	logger.Info("Init Database")
	db, err := database.NewDB(logger)
	if err != nil {
		logger.Fatalln("Failed to connect database.")
		panic(err)
	}
	logger.Info("Init Database Success")

	defer database.CloseDB(logger, db)

	logger.Info("Migrate Database")
	err = migrations.Migrate(db)
	if err != nil {
		logger.Fatalln("Failed to migrate database.")
		panic(err)
	}
	logger.Info("Migrate Database Success")

	engine := gin.New()
	router := &router.Router{
		Engine: engine,
		DB:     db,
	}
	router.InitializeRouter(logger)
	router.SetupHandler()

	err = engine.Run(":" + os.Getenv("API_PORT"))
	if err != nil {
		logger.Fatalln("Failed to run server.")
		panic(err)
	}
}
