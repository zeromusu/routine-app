package app

import (
	"log"
	"routine-app-server/internal/config"
	"routine-app-server/internal/db"
	"routine-app-server/internal/domain"
	"routine-app-server/internal/infrastructure/persistence"
	"routine-app-server/internal/interfaces/handler"
	"routine-app-server/internal/router"
	"routine-app-server/internal/usecase"
)

func Run() {
	// Load Config
	cfg := config.LoadConfig()

	// Initialize Database
	database, err := db.InitDB(cfg.DB)
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}

	// Database Migration
	err = database.AutoMigrate(&domain.Routine{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	routineRepo := persistence.NewRoutinePersistence(database)
	routineUseCase := usecase.NewRoutineUseCase(routineRepo)
	routineHandler := handler.NewRoutineHandler(routineUseCase)

	appHandlers := router.AppHandlers{
		Routine: routineHandler,
	}

	r := router.NewRouter(appHandlers)
	r.Run(":" + cfg.App.AppPort)
}
