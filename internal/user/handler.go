package user

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	service *Service
	logger  *log.Logger
}

func NewHandler(s *Service, l *log.Logger) *Handler {
	return &Handler{service: s, logger: l}
}

func (h *Handler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Post("/register", h.Register)
	r.Get("/", h.ListAll)
	r.Get("/{id}", h.GetByID)
	r.Patch("/verify", h.Verify)
	return r
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Printf("Handler Error: failed to decode JSON: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, err := h.service.Register(r.Context(), &req)
	if err != nil {
		h.logger.Printf("Handler Error: registration failed: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func (h *Handler) ListAll(w http.ResponseWriter, r *http.Request) {
	users, err := h.service.repo.GetAll(r.Context())
	if err != nil {
		h.logger.Printf("Handler Error: failed to fetch users: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")

	user, err := h.service.FindByID(r.Context(), idStr)
	if err != nil {
		h.logger.Printf("Handler: User %s not found", idStr)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *Handler) Verify(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email string `json:"email"`
		Code  string `json:"code"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if err := h.service.repo.VerifyEmail(r.Context(), req.Email, req.Code); err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}
