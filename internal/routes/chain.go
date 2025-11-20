package routes

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	// "go-backend/internal/services/chain" // Uncomment and implement this service
)

func RegisterChainRoutes(r *mux.Router) {
	chainRouter := r.PathPrefix("/chain").Subrouter()
	chainRouter.HandleFunc("/claimable", claimableHandler).Methods(http.MethodGet)
	// Add more chain-related routes here
}

// GET /chain/claimable
func claimableHandler(w http.ResponseWriter, r *http.Request) {
	distributionIdStr := r.URL.Query().Get("distributionId")
	user := r.URL.Query().Get("user")

	distributionId, err := strconv.Atoi(distributionIdStr)
	if err != nil || distributionId < 0 || len(user) != 42 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid_request"})
		return
	}

	// TODO: Call your service logic here, e.g.:
	// amount, err := chain.ReadClaimableAmount(distributionId, user)
	// if err != nil {
	//     w.WriteHeader(http.StatusInternalServerError)
	//     json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
	//     return
	// }
	amount := "mocked_amount" // Replace with actual logic

	json.NewEncoder(w).Encode(map[string]interface{}{"amount": amount})
}
