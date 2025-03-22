package handlers

import (
	"encoding/json"
	"errors"
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/google/uuid"
	"github.com/viditagrawal56/url-shortner/internal/auth"
	"github.com/viditagrawal56/url-shortner/internal/config"
	"github.com/viditagrawal56/url-shortner/internal/db"
	"github.com/viditagrawal56/url-shortner/internal/email"
	"github.com/viditagrawal56/url-shortner/internal/models"
	"github.com/viditagrawal56/url-shortner/internal/urlShortner"
)

type API struct {
	cfg          *config.Config
	auth         *auth.Service
	urlShortner  *urlShortner.URLShortnerService
	emailService *email.Service
}

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}

func NewRouter(database *db.Database, cfg *config.Config) http.Handler {
	authService := auth.New(database, cfg)
	urlShortnerService := urlShortner.NewUrlShortnerService(database.GetDB())
	emailConfig := &email.EmailConfig{
		SMTPHost:     cfg.Email.SMTPHost,
		SMTPPort:     cfg.Email.SMTPPort,
		SMTPUsername: cfg.Email.SMTPUsername,
		SMTPPassword: cfg.Email.SMTPPassword,
		FromEmail:    cfg.Email.FromEmail,
	}
	emailService := email.NewEmailService(emailConfig)

	api := &API{
		cfg:          cfg,
		auth:         authService,
		urlShortner:  urlShortnerService,
		emailService: emailService,
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
		r.Get("/auth/{shortCode}", api.HandleAuthRequest)
		r.Post("/auth/{shortCode}", api.HandleEmailSubmission)
		r.Get("/verify", api.HandleMagicLinkVerification)
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

	// Get the short URL details
	shortURL, err := a.urlShortner.GetShortURLByCode(shortCode)
	if err != nil {
		if errors.Is(err, urlShortner.ErrURLNotFound) {
			respondWithError(w, http.StatusNotFound, "URL not found")
		} else {
			respondWithError(w, http.StatusInternalServerError, "An error occurred")
		}
		return
	}

	// Check if URL requires authentication
	if shortURL.RequiresAuth {
		http.Redirect(w, r, "/auth/"+shortCode, http.StatusFound)
		return
	}

	// If no auth required, check time validity
	now := time.Now()
	if shortURL.ValidFrom != nil && now.Before(*shortURL.ValidFrom) {
		respondWithError(w, http.StatusForbidden, "URL is not yet valid")
		return
	}

	if shortURL.ValidUntil != nil && now.After(*shortURL.ValidUntil) {
		respondWithError(w, http.StatusGone, "URL has expired")
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

// Handler for showing the email input form
func (a *API) HandleAuthRequest(w http.ResponseWriter, r *http.Request) {
	shortCode := chi.URLParam(r, "shortCode")
	if shortCode == "" {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	// Check if the short URL exists
	_, err := a.urlShortner.GetShortURLByCode(shortCode)
	if err != nil {
		if errors.Is(err, urlShortner.ErrURLNotFound) {
			respondWithError(w, http.StatusNotFound, "URL not found")
		} else {
			respondWithError(w, http.StatusInternalServerError, "An error occurred")
		}
		return
	}

	// Render email input form
	tmpl := `
<!DOCTYPE html>
<html>
<head>
    <title>URL Authentication</title>
    <style>
        body { font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto; padding: 20px; }
        .container { background: #f9f9f9; border-radius: 5px; padding: 20px; box-shadow: 0 2px 4px rgba(0,0,0,0.1); }
        h1 { color: #333; }
        .form-group { margin-bottom: 15px; }
        label { display: block; margin-bottom: 5px; }
        input[type="email"] { width: 100%; padding: 8px; border: 1px solid #ddd; border-radius: 4px; }
        button { background: #4CAF50; color: white; border: none; padding: 10px 15px; border-radius: 4px; cursor: pointer; }
        .error { color: red; margin-top: 10px; }
    </style>
</head>
<body>
    <div class="container">
        <h1>Email Verification Required</h1>
        <p>This URL requires verification. Please enter your email address to receive a magic link.</p>
        
        <form method="POST" action="/auth/{{.ShortCode}}">
            <div class="form-group">
                <label for="email">Email Address:</label>
                <input type="email" id="email" name="email" required>
            </div>
            <button type="submit">Send Magic Link</button>
        </form>
    </div>
</body>
</html>
`
	t, err := template.New("auth").Parse(tmpl)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Template error")
		return
	}

	data := struct {
		ShortCode string
	}{
		ShortCode: shortCode,
	}

	w.Header().Set("Content-Type", "text/html")
	t.Execute(w, data)
}

// Handler for email submission
func (a *API) HandleEmailSubmission(w http.ResponseWriter, r *http.Request) {
	shortCode := chi.URLParam(r, "shortCode")
	if shortCode == "" {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	// Parse form
	if err := r.ParseForm(); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid form data")
		return
	}

	email := r.FormValue("email")
	if email == "" {
		http.Redirect(w, r, "/auth/"+shortCode, http.StatusFound)
		return
	}

	// Generate token
	token, err := a.auth.GenerateMagicLinkToken(email, shortCode)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to generate verification token")
		return
	}

	// Send email with magic link
	err = a.emailService.SendMagicLink(email, shortCode, token, a.cfg.Server.BaseURL)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to send verification email")
		return
	}

	// Show confirmation page
	tmpl := `
<!DOCTYPE html>
<html>
<head>
    <title>Email Sent</title>
    <style>
        body { font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto; padding: 20px; }
        .container { background: #f9f9f9; border-radius: 5px; padding: 20px; box-shadow: 0 2px 4px rgba(0,0,0,0.1); }
        h1 { color: #333; }
        .success { color: #4CAF50; }
    </style>
</head>
<body>
    <div class="container">
        <h1>Email Sent!</h1>
        <p class="success">We've sent a magic link to <strong>{{.Email}}</strong>.</p>
        <p>Please check your inbox and click the link to access the URL. The link will expire in 15 minutes.</p>
    </div>
</body>
</html>
`
	t, err := template.New("confirmation").Parse(tmpl)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Template error")
		return
	}

	data := struct {
		Email string
	}{
		Email: email,
	}

	w.Header().Set("Content-Type", "text/html")
	t.Execute(w, data)
}

func (a *API) HandleMagicLinkVerification(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		respondWithError(w, http.StatusBadRequest, "Missing token")
		return
	}
	// Verify token
	email, shortCode, err := a.auth.VerifyMagicLinkToken(token)
	if err != nil {
		// Show error page
		tmpl := `
<!DOCTYPE html>
<html>
<head>
    <title>Invalid or Expired Link</title>
    <style>
        body { font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto; padding: 20px; }
        .container { background: #f9f9f9; border-radius: 5px; padding: 20px; box-shadow: 0 2px 4px rgba(0,0,0,0.1); }
        h1 { color: #333; }
        .error { color: red; }
    </style>
</head>
<body>
    <div class="container">
        <h1>Invalid or Expired Link</h1>
        <p class="error">The magic link you used is invalid or has expired.</p>
        <p>Please request a new link to access the URL.</p>
    </div>
</body>
</html>
`
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusBadRequest)
		template.Must(template.New("error").Parse(tmpl)).Execute(w, nil)
		return
	}

	// Resolve the short URL with the verified email
	shortURL, err := a.urlShortner.ResolveShortURL(shortCode, email)
	if err != nil {
		switch {
		case errors.Is(err, urlShortner.ErrURLNotFound):
			respondWithError(w, http.StatusNotFound, "URL not found")
		case errors.Is(err, urlShortner.ErrURLExpired):
			respondWithError(w, http.StatusGone, "URL has expired")
		case errors.Is(err, urlShortner.ErrURLNotYetValid):
			respondWithError(w, http.StatusForbidden, "URL is not yet valid")
		case errors.Is(err, urlShortner.ErrUnauthorizedAccess):
			// Show unauthorized access page
			tmpl := `
<!DOCTYPE html>
<html>
<head>
    <title>Access Denied</title>
    <style>
        body { font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto; padding: 20px; }
        .container { background: #f9f9f9; border-radius: 5px; padding: 20px; box-shadow: 0 2px 4px rgba(0,0,0,0.1); }
        h1 { color: #333; }
        .error { color: red; }
    </style>
</head>
<body>
    <div class="container">
        <h1>Access Denied</h1>
        <p class="error">Your email address ({{.Email}}) is not authorized to access this URL.</p>
        <p>Please contact the URL owner for access.</p>
    </div>
</body>
</html>
`
			w.Header().Set("Content-Type", "text/html")
			w.WriteHeader(http.StatusForbidden)
			template.Must(template.New("unauthorized").Parse(tmpl)).Execute(w, struct{ Email string }{Email: email})
			return
		default:
			respondWithError(w, http.StatusInternalServerError, "An error occurred")
		}
		return
	}

	// Redirect to the original URL
	http.Redirect(w, r, shortURL.OriginalURL, http.StatusFound)
}

// TODO: Improve the url validation to make it more robust
func isValidURL(url string) bool {
	return strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://")
}
