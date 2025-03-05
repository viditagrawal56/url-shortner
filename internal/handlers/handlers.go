package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/viditagrawal56/url-shortner/internal/auth"
	"github.com/viditagrawal56/url-shortner/internal/config"
	"github.com/viditagrawal56/url-shortner/internal/db"
	"github.com/viditagrawal56/url-shortner/internal/models"
)

type API struct {
	cfg  *config.Config
	auth *auth.Service
}

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func NewRouter(database *db.Database, cfg *config.Config) http.Handler {
	authService := auth.New(database, cfg)
	api := &API{
		cfg:  cfg,
		auth: authService,
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	//CORS middleware
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	//Public routes
	r.Group(func(r chi.Router) {
		r.Post("/register", api.HandleUserRegistration)
	})
	r.Get("/", api.HandleHome)

	return r
}

func (a *API) HandleHome(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := map[string]string{"message": "Hello there"}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}

func (a *API) HandleUserRegistration(w http.ResponseWriter, r *http.Request) {
	var credentials models.Credentials
	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Validate input
	if credentials.Email == "" || credentials.Password == "" {
		respondWithError(w, http.StatusBadRequest, "Email and password are required")
		return
	}

	// Check if user already exists
	err := a.auth.RegisterUser(credentials.Email, credentials.Password)
	if err != nil {
		switch {
		case errors.Is(err, auth.ErrUserExists):
			respondWithError(w, http.StatusConflict, "User already exists")
		case errors.Is(err, auth.ErrInvalidInput):
			respondWithError(w, http.StatusBadRequest, "Invalid input")
		default:
			log.Printf("Error registering user: %v", err)
			respondWithError(w, http.StatusInternalServerError, "Failed to register user")
		}
		return
	}

	respondWithJSON(w, http.StatusCreated, Response{
		Success: true,
		Message: "User registered successfully",
	})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, Response{
		Success: false,
		Message: message,
	})
}
