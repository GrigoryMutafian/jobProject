package main

import (
	"errors"
	"fmt"
	"jobProject/internal/db"
	"jobProject/internal/handlers"
	"jobProject/internal/repository"
	"jobProject/internal/usecase"
	"log"
	"net/http"
)

func main() {
	if err := db.InitDB(); err != nil {
		fmt.Println(err)
		return
	}

	defer db.DB.Close()

	subRepo := &repository.PostgresSubs{DB: db.DB}
	subUC := usecase.NewSubUsecase(subRepo)

	if err := handlers.Init(subUC); err != nil {
		log.Fatalf("handlers init: %v", err)
	}

	http.HandleFunc("/CreateSub", handlers.CreateSub)
	http.HandleFunc("/ReadSubByID", handlers.ReadSubByID)
	log.Println("listening on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("server error: %v", err)
	}
}
