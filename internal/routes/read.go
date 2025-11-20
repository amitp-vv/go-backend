package routes

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/amitp-vv/go-backend/internal/models"

	"github.com/gorilla/mux"
)

func RegisterReadRoutes(r *mux.Router) {
	readRouter := r.PathPrefix("/read").Subrouter()
	readRouter.HandleFunc("/distributions", getDistributions).Methods(http.MethodGet)
	readRouter.HandleFunc("/claims", getClaims).Methods(http.MethodGet)
}

// GET /read/distributions
func getDistributions(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("pageSize"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 200 {
		pageSize = 50
	}
	offset := (page - 1) * pageSize

	var items []models.Distribution
	var total int64
	models.DB.Model(&models.Distribution{}).Count(&total)
	models.DB.Order("id DESC").Limit(pageSize).Offset(offset).Find(&items)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"items":    items,
		"page":     page,
		"pageSize": pageSize,
		"total":    total,
	})
}

// GET /read/claims
func getClaims(w http.ResponseWriter, r *http.Request) {
	user := r.URL.Query().Get("user")
	if len(user) != 42 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid_user"})
		return
	}
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("pageSize"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 200 {
		pageSize = 50
	}
	offset := (page - 1) * pageSize

	var items []models.Claim
	var total int64
	models.DB.Model(&models.Claim{}).Where("user_id = ?", user).Count(&total)
	models.DB.Where("user_id = ?", user).Order("id DESC").Limit(pageSize).Offset(offset).Find(&items)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"items":    items,
		"page":     page,
		"pageSize": pageSize,
		"total":    total,
	})
}
