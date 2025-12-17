package handlers

import (
	"encoding/json"
	"fmt"
	"jobProject/internal/api"
	"jobProject/internal/conv"
	"jobProject/internal/model"
	"jobProject/internal/usecase"
	"log"
	"net/http"
	"regexp"
	"strconv"

	_ "jobProject/docs"
)

var subUC *usecase.SubUsecase

var uuidRegex = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)

func validateUUID(uuid string) bool {
	return uuidRegex.MatchString(uuid)
}

func Init(uc *usecase.SubUsecase) error {
	if uc == nil {
		return fmt.Errorf("nil usecase")
	}
	subUC = uc
	return nil
}

// @Summary Создать подписку
// @Description Создает новую запись о подписке
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param subscription body model.Subscription true "Данные подписки"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string "Некорректный JSON или параметры"
// @Failure 409 {object} map[string]string "Конфликт"
// @Router /CreateColumn [post]
func CreateColumn(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed: ", http.StatusMethodNotAllowed)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	defer r.Body.Close()

	var newSub model.Subscription

	err := json.NewDecoder(r.Body).Decode(&newSub)
	if err != nil {
		http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err := subUC.CreateColumnUC(r.Context(), newSub); err != nil {
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
	response := map[string]string{"name of added subscription is": *newSub.Service}
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "JSON encoding error: "+err.Error(), http.StatusInternalServerError)
	}
}

// @Summary Получить подписку по ID
// @Description Возвращает подписку по идентификатору ID
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id query int true "ID подписки"
// @Success 200 {object} map[string]model.SubscriptionDB
// @Failure 400 {object} map[string]string "Некорректный id или ошибка"
// @Failure 404 {object} map[string]string "Подписка не найдена"
// @Failure 409 {object} map[string]string "Конфликт"
// @Failure 500 {object} map[string]string "Внутренняя ошибка"
// @Router /ReadSubByID [get]
func ReadSubByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed: ", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "id input is clear", http.StatusBadRequest)
		return
	}

	idInt, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "pars error", http.StatusBadRequest)
		return
	}

	sub, err := subUC.ReadColumnUC(r.Context(), idInt)
	if err != nil {
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
	w.WriteHeader(http.StatusOK)
	text := fmt.Sprintf("column id: %d", idInt)
	response := map[string]model.SubscriptionDB{text: sub}
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "JSON encoding error: "+err.Error(), http.StatusInternalServerError)
		return
	}

}

// @Summary Частично обновить подписку по ID
// @Description Обновляет выбранные поля записи
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id query int true "ID подписки"
// @Param subscription body model.Subscription true "Патч-данные"
// @Success 200 {object} map[int]string "ID -> updated"
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string "Подписка не найдена"
// @Failure 409 {object} map[string]string "Конфликт"
// @Failure 500 {object} map[string]string
// @Router /PatchColumnByID [patch]
func PatchColumnByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		http.Error(w, "Method Not Allowed: ", http.StatusMethodNotAllowed)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	defer r.Body.Close()

	var patchBody model.Subscription

	err := json.NewDecoder(r.Body).Decode(&patchBody)
	if err != nil {
		http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "id input is clear", http.StatusBadRequest)
		return
	}

	idInt, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "pars error", http.StatusBadRequest)
		return
	}
	err = subUC.PatchColumnByID(r.Context(), idInt, patchBody)
	if err != nil {
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
	w.WriteHeader(http.StatusOK)
	response := map[int]string{idInt: "updated"}
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "JSON encoding error: "+err.Error(), http.StatusInternalServerError)
		return
	}

}

// @Summary Удалить подписку по ID
// @Description Удаляет запись о подписке
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id query int true "ID подписки"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string "Подписка не найдена"
// @Failure 409 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /DeleteColumnByID [delete]
func DeleteColumnByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method Not Allowed: ", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "id input is clear", http.StatusBadRequest)
		return
	}

	idInt, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "pars error", http.StatusBadRequest)
		return
	}

	err = subUC.DeleteColumnByID(r.Context(), idInt)
	if err != nil {
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
	w.WriteHeader(http.StatusOK)
	text := fmt.Sprintf("deleted column id: %d", idInt)
	response := map[string]string{text: "OK"}
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "JSON encoding error: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

// @Summary Получить сумму подписок за период
// @Description Считает суммарную стоимость подписок по id пользователя, названию подписки и периоду
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param user_id query string true "ID пользователя (uuid)"
// @Param service query string true "Название сервиса"
// @Param date_from query string true "Период начала подписки MM-YYYY"
// @Param date_to query string true "Период конца подписки MM-YYYY"
// @Success 200 {object} map[string]int
// @Failure 400 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /TotalPriceByPeriod [get]
func TotalPriceByPeriod(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := r.URL.Query().Get("user_id")
	service := r.URL.Query().Get("service")
	dateFrom := r.URL.Query().Get("date_from")
	dateTo := r.URL.Query().Get("date_to")

	if userID == "" || service == "" || dateFrom == "" || dateTo == "" {
		http.Error(w, "user_id, service, date_from, date_to required", http.StatusBadRequest)
		return
	}

	if !validateUUID(userID) {
		http.Error(w, "invalid user_id format: must be a valid UUID (e.g., 70601fee-2bf1-4721-ae6f-7636e79a0cbb)", http.StatusBadRequest)
		return
	}

	fromTime, err := conv.ParseMMYYYY(dateFrom)
	if err != nil {
		http.Error(w, "wrong date_from format", http.StatusBadRequest)
		return
	}
	toTime, err := conv.ParseMMYYYY(dateTo)
	if err != nil {
		http.Error(w, "wrong date_to format", http.StatusBadRequest)
		return
	}

	total, err := subUC.TotalPriceByPeriod(r.Context(), userID, service, fromTime, toTime)
	if err != nil {
		if usecase.IsValidationErr(err) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]int{"total": total})
}

func ListSubscriptions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		http.Error(w, "user_id parameter is required", http.StatusBadRequest)
		return
	}

	if !validateUUID(userID) {
		http.Error(w, "invalid user_id format: must be a valid UUID (70601fee-2bf1-4721-ae6f-7636e79a0cbb)", http.StatusBadRequest)
		return
	}

	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	params := api.PaginationParams{
		Page:  1,
		Limit: 10,
	}

	if pageStr != "" {
		page, err := strconv.Atoi(pageStr)
		if err != nil || page < 1 {
			http.Error(w, "invalid page parameter: must be a positive integer", http.StatusBadRequest)
			return
		}
		params.Page = page
	}

	if limitStr != "" {
		limit, err := strconv.Atoi(limitStr)
		if err != nil || limit < 1 {
			http.Error(w, "invalid limit parameter: must be a positive integer", http.StatusBadRequest)
			return
		}
		if limit > 100 {
			http.Error(w, "invalid limit parameter: maximum value is 100", http.StatusBadRequest)
			return
		}
		params.Limit = limit
	}

	params.Validate()

	response, err := subUC.ListSubscriptions(r.Context(), userID, params)
	if err != nil {
		log.Printf("error listing subscriptions: %v", err)
		http.Error(w, fmt.Sprintf("internal error: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("error encoding response: %v", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
}
