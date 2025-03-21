package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/google/uuid"
	"github.com/viditagrawal56/url-shortner/internal/auth"
	"github.com/viditagrawal56/url-shortner/internal/config"
	"github.com/viditagrawal56/url-shortner/internal/db"
	"github.com/viditagrawal56/url-shortner/internal/models"
	"github.com/viditagrawal56/url-shortner/internal/urlShortner"
)

type API struct {
	cfg         *config.Config
	auth        *auth.Service
	urlShortner *urlShortner.URLShortnerService
}

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}

func NewRouter(database *db.Database, cfg *config.Config) http.Handler {
	authService := auth.New(database, cfg)
	urlShortnerService := urlShortner.NewUrlShortnerService(database.GetDB())

	api := &API{
		cfg:         cfg,
		auth:        authService,
		urlShortner: urlShortnerService,
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger, middleware.Recoverer)

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
		r.Post("/login", api.HandleUserLogin)
		r.Get("/{shortCode}", api.HandleRedirect)
	})

	//Protected routes
	r.Group(func(r chi.Router) {
		r.Use(api.auth.AuthMiddleware)
		r.Post("/urls", api.HandleCreateShortURL)
		r.Get("/urls", api.HandleGetUserShortURLs)
		r.Delete("/urls/{urlID}", api.HandleDeleteShortURL)
	})

	return r
}

func (a *API) HandleUserRegistration(w http.ResponseWriter, r *http.Request) {
	// Retrieve the credentials from the request body
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

	// Register the user
	err := a.auth.RegisterUser(credentials.Email, credentials.Password)
	if err != nil {
		switch {
		case errors.Is(err, auth.ErrUserExists):
			respondWithError(w, http.StatusConflict, "User with this email already exists")
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

func (a *API) HandleUserLogin(w http.ResponseWriter, r *http.Request) {
	// Retrieve the credentials from the request body
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

	// Start the login process
	token, err := a.auth.LoginUser(credentials.Email, credentials.Password)
	if err != nil {
		switch {
		case errors.Is(err, auth.ErrInvalidCredentials):
			respondWithError(w, http.StatusUnauthorized, "Invalid email or password")
		case errors.Is(err, auth.ErrInvalidInput):
			respondWithError(w, http.StatusBadRequest, "Invalid input")
		default:
			log.Printf("Error logging in user: %v", err)
			respondWithError(w, http.StatusInternalServerError, "Failed to login user")
		}
		return
	}

	respondWithJSON(w, http.StatusOK, Response{
		Success: true,
		Message: "User logged in successfully",
		Data: map[string]string{
			"token": token,
		},
	})
}

func (a *API) HandleCreateShortURL(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(auth.UserIDKey).(uuid.UUID)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "You are not authorized, Please login")
		return
	}

	var req models.CreateShortURLRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Request")
		return
	}

	// Validate the incoming URL
	if req.OriginalURL == "" || !isValidURL(req.OriginalURL) {
		respondWithError(w, http.StatusBadRequest, "Please enter a valid URL")
		return
	}

	shortURL, err := a.urlShortner.CreateShortURL(userID, req.OriginalURL, req.Options)
	if err != nil {
		log.Printf("Error creating short URL: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to create the short URL")
		return
	}

	respondWithJSON(w, http.StatusCreated, Response{
		Success: true,
		Data:    shortURL,
	})
}

func (a *API) HandleRedirect(w http.ResponseWriter, r *http.Request) {
	shortCode := chi.URLParam(r, "shortCode")
	if shortCode == "" {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	// TODO: get the email from the cookies afer figuring out the visitor auth
	email := ""

	shortURL, err := a.urlShortner.ResolveShortURL(shortCode, email)
	if err != nil {
		switch {
		case errors.Is(err, urlShortner.ErrURLNotFound):
			respondWithError(w, http.StatusNotFound, "URL not found")
		case errors.Is(err, urlShortner.ErrURLExpired):
			respondWithError(w, http.StatusGone, "URL has expired")
		case errors.Is(err, urlShortner.ErrURLNotYetValid):
			respondWithError(w, http.StatusForbidden, "URL is not yet valid")
		case errors.Is(err, urlShortner.ErrAuthenticationRequired):
			// TODO: Redirect to email verification page
			respondWithError(w, http.StatusNotFound, "Please verify your email")
		case errors.Is(err, urlShortner.ErrUnauthorizedAccess):
			respondWithError(w, http.StatusForbidden, "You are not authorized to access this URL")
		default:
			log.Printf("Error resolving short URL: %v", err)
			respondWithError(w, http.StatusInternalServerError, "An error occurred")
		}
		return
	}

	// Redirect to the original URL
	http.Redirect(w, r, shortURL.OriginalURL, http.StatusFound)
}

func (a *API) HandleGetUserShortURLs(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(auth.UserIDKey).(uuid.UUID)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	shortURLs, err := a.urlShortner.GetUserShortURLs(userID)
	if err != nil {
		log.Printf("error getting user's short URLs: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to get short URLs")
		return
	}

	respondWithJSON(w, http.StatusFound, Response{
		Success: true,
		Data:    shortURLs,
	})
}

func (a *API) HandleDeleteShortURL(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(auth.UserIDKey).(uuid.UUID)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	urlID, err := uuid.Parse(chi.URLParam(r, "urlID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid URL ID")
		return
	}

	err = a.urlShortner.DeleteShortURL(userID, urlID)
	if err != nil {
		switch {
		case errors.Is(err, urlShortner.ErrURLNotFound):
			respondWithError(w, http.StatusNotFound, "URL not found")
		default:
			log.Printf("Error deleting short URL: %v", err)
			respondWithError(w, http.StatusInternalServerError, "Failed to delete short URL")
		}
		return
	}
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

// TODO: Improve the url validation to make it more robust
func isValidURL(url string) bool {
	return strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://")
}
