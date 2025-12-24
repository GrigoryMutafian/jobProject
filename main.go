package main

import (
	"errors"
	"jobProject/internal/db"
	"jobProject/internal/handlers"
	"jobProject/internal/logger"
	"jobProject/internal/repository"
	"jobProject/internal/usecase"
	"log"
	"log/slog"
	"net/http"
	"os"

	httpSwagger "github.com/swaggo/http-swagger"
)

func main() {
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "INFO"
	}

	slog.Info("Starting application", "log_level", logLevel)

	logger.InitLogger(logLevel)

	if err := db.InitDB(); err != nil {
		slog.Error("Failed to initialize database", "error", err)
		os.Exit(1)
	}
	slog.Info("Database initialized successfully")

	subRepo := &repository.PostgresSubs{DB: db.DB}
	subUC := usecase.NewSubUsecase(subRepo)

	if err := handlers.Init(subUC); err != nil {
		slog.Error("Failed to initialize handlers", "error", err)
		os.Exit(1)
	}
	slog.Info("Handlers initialized successfully")

	http.HandleFunc("/CreateColumn", handlers.CreateColumn)
	http.HandleFunc("/ReadSubByID", handlers.ReadSubByID)
	http.HandleFunc("/PatchColumnByID", handlers.PatchColumnByID)
	http.HandleFunc("/DeleteColumnByID", handlers.DeleteColumnByID)
	http.HandleFunc("/TotalPriceByPeriod", handlers.TotalPriceByPeriod)
	http.HandleFunc("/swagger/", httpSwagger.WrapHandler)
	http.HandleFunc("/ListSubscriptions", handlers.ListSubscriptions)

	log.Println("listening on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("server error: %v", err)
	}
}
