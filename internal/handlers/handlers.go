package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/viditagrawal56/url-shortner/internal/config"
)

type API struct {
	cfg *config.Config
}

func NewRouter(cfg *config.Config) http.Handler {
	api := &API{
		cfg: cfg,
	}

	r := chi.NewRouter()

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