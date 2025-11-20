package routes

import (
	"encoding/json"
	"net/http"

	"github.com/amitp-vv/go-backend/internal/middleware"

	"github.com/gorilla/mux"
)

func RegisterAdminRoutes(r *mux.Router) {
	adminRouter := r.PathPrefix("/admin").Subrouter()
	adminRouter.Use(middleware.RequireAuth)
	adminRouter.Use(middleware.RequireAdmin)

	adminRouter.HandleFunc("/whitelist/request", requestWhitelistUpdateTx).Methods(http.MethodPost)
	adminRouter.HandleFunc("/whitelist/confirm", confirmWhitelistUpdateTx).Methods(http.MethodPost)
	adminRouter.HandleFunc("/distributions", createDistributionTx).Methods(http.MethodPost)
}

// POST /admin/whitelist/request
func requestWhitelistUpdateTx(w http.ResponseWriter, r *http.Request) {
	var req struct {
		User   string `json:"user"`
		Status bool   `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || len(req.User) != 42 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid_request"})
		return
	}
	// TODO: Call your service logic here, e.g. chain.RequestWhitelistUpdateTx(req.User, req.Status)
	txHash := "mocked_tx_hash" // Replace with actual logic
	json.NewEncoder(w).Encode(map[string]string{"txHash": txHash})
}

// POST /admin/whitelist/confirm
func confirmWhitelistUpdateTx(w http.ResponseWriter, r *http.Request) {
	var req struct {
		User string `json:"user"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || len(req.User) != 42 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid_request"})
		return
	}
	// TODO: Call your service logic here, e.g. chain.ConfirmWhitelistUpdateTx(req.User)
	txHash := "mocked_tx_hash" // Replace with actual logic
	json.NewEncoder(w).Encode(map[string]string{"txHash": txHash})
}

// POST /admin/distributions
func createDistributionTx(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Amount string `json:"amount"`
	}
	// Simple regex check for digits only
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Amount == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid_request"})
		return
	}
	for _, c := range req.Amount {
		if c < '0' || c > '9' {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "invalid_amount"})
			return
		}
	}
	// TODO: Call your service logic here, e.g. chain.CreateDistributionTx(req.Amount)
	txHash := "mocked_tx_hash" // Replace with actual logic
	json.NewEncoder(w).Encode(map[string]string{"txHash": txHash})
}
