package handler

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"lost-items-service/internal/middleware"
)

const (
	ErrBodyRequest   = "Invalid Request Body"
	ErrRequestFields = "Invalid Request Fields"
	ErrInvalidRole   = "invalid role in Request"
	ErrUUIDParsing   = "invalid ID format"
)

const (
	ModeratorRole = "moderator"
	AdminRole     = "admin"
	UserRole      = "user"
)

type Service interface {
	AuthService
}

type Router struct {
	service Service
}

func NewRouter(service Service, jwtSecret string, logger *slog.Logger) *chi.Mux {
	r := chi.NewRouter()
	router := &Router{service: service}

	r.Use(middleware.NewValidator().Middleware)
	r.Use(middleware.ContextLoggerMiddleware(logger))
	r.Post("/register", http.HandlerFunc(router.registerHandler))
	r.Post("/login", http.HandlerFunc(router.loginHandler))

	r.Group(func(protected chi.Router) {
		protected.Use(middleware.NewJWT(jwtSecret).Authenticate)
	})
	return r
}

func getValidator(r *http.Request) *validator.Validate {
	if v, ok := r.Context().Value("validator").(*validator.Validate); ok {
		return v
	}
	return validator.New()
}

func getLogger(r *http.Request) *slog.Logger {
	if l, ok := r.Context().Value("logger").(*slog.Logger); ok {
		return l
	}
	return slog.Default() // fallback на глобальный
}

func (r *Router) registerHandler(w http.ResponseWriter, req *http.Request) {
	h := NewAuthHandler(r.service)
	h.Register(w, req)
}

func (r *Router) loginHandler(w http.ResponseWriter, req *http.Request) {
	h := NewAuthHandler(r.service)
	h.Login(w, req)
}
