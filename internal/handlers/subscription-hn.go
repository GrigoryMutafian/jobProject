package handlers

import (
	"encoding/json"
	"fmt"
	"jobProject/internal/model"
	"jobProject/internal/usecase"
	"net/http"
)

var subUC *usecase.SubUsecase

func Init(uc *usecase.SubUsecase) error {
	if uc == nil {
		return fmt.Errorf("nil usecase")
	}
	subUC = uc
	return nil
}

func Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed: ", http.StatusMethodNotAllowed)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	defer r.Body.Close()

	var newSub model.Subscription

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(&newSub)
	if err != nil {
		http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err := subUC.Create(r.Context(), newSub); err != nil {
		switch {
		case usecase.IsValidationErr(err):
			http.Error(w, err.Error(), http.StatusBadRequest)
		case usecase.IsConflictErr(err):
			http.Error(w, err.Error(), http.StatusConflict)
		default:
			http.Error(w, "internal error: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	response := map[string]string{"name of added subscription is": newSub.Service}
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "JSON encoding error: "+err.Error(), http.StatusInternalServerError)
	}
}
