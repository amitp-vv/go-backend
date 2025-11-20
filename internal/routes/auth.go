package routes

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/amitp-vv/go-backend/internal/middleware"
	"github.com/amitp-vv/go-backend/internal/models"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	// "github.com/ethereum/go-ethereum/crypto" // For signature verification
)

func RegisterAuthRoutes(r *mux.Router) {
	authRouter := r.PathPrefix("/auth").Subrouter()
	authRouter.Use(middleware.RateLimit)
	authRouter.HandleFunc("/nonce", nonceHandler).Methods(http.MethodPost)
	authRouter.HandleFunc("/login", loginHandler).Methods(http.MethodPost)
}

// POST /auth/nonce
func nonceHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Address string `json:"address"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || len(req.Address) != 42 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid_request"})
		return
	}
	address := strings.ToLower(req.Address)
	nonce := randomHex(16)
	var wallet models.Wallet
	result := models.DB.Where("id = ?", address).First(&wallet)
	if result.Error != nil {
		wallet = models.Wallet{ID: address, UserID: address, Nonce: nonce}
		models.DB.Create(&wallet)
	} else if wallet.Nonce == "" {
		wallet.Nonce = nonce
		models.DB.Save(&wallet)
	}
	json.NewEncoder(w).Encode(map[string]string{"address": address, "nonce": wallet.Nonce})
}

// POST /auth/login
func loginHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Address   string `json:"address"`
		Signature string `json:"signature"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || len(req.Address) != 42 || len(req.Signature) < 10 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid_request"})
		return
	}
	address := strings.ToLower(req.Address)
	var wallet models.Wallet
	result := models.DB.Where("id = ?", address).First(&wallet)
	if result.Error != nil || wallet.Nonce == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "nonce_missing"})
		return
	}
	// TODO: Verify signature using go-ethereum/crypto
	// recovered := verifyMessage(message, req.Signature)
	// if strings.ToLower(recovered) != address {
	//     w.WriteHeader(http.StatusUnauthorized)
	//     json.NewEncoder(w).Encode(map[string]string{"error": "signature_mismatch"})
	//     return
	// }

	wallet.Nonce = ""
	models.DB.Save(&wallet)
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "devsecret"
	}
	claims := middleware.AuthPayload{Address: address}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(secret))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "tx_failed"})
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"token": tokenStr})
}

// Helper to generate random hex string
func randomHex(n int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, n)
	rand.Read(b)
	hex := ""
	for _, v := range b {
		hex += "0123456789abcdef"[v%16 : v%16+1]
	}
	return hex
}
