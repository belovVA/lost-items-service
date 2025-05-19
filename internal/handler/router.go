package handler

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"lost-items-service/internal/middleware"
)

const (
	ErrBodyRequest     = "Invalid Request Body"
	ErrRequestFields   = "Invalid Request Fields"
	ErrInvalidRole     = "invalid role in Request"
	ErrUUIDParsing     = "invalid ID format"
	ErrQueryParameters = "invalid query parameters"
	ErrInvalidStatus   = "invalid moderation status"
)

// Roles
const (
	ModeratorRole = "moderator"
	AdminRole     = "admin"
	UserRole      = "user"
)

// Logging
const (
	UserIDKey       = "userId"
	ErrorKey        = "error"
	ErrorServiceMsg = "service error"
	AnnIDKey        = "annId"
)

type Service interface {
	AuthService
	UserService
	AnnouncementService
	ImageService
}

type Router struct {
	service Service
}

func NewRouter(service Service, jwtSecret string, logger *slog.Logger) *chi.Mux {
	r := chi.NewRouter()

	router := &Router{service: service}
	r.Route("/api/v1", func(api chi.Router) {
		api.Use(middleware.NewValidator().Middleware)
		api.Use(middleware.ContextLoggerMiddleware(logger))

		api.Post("/register", http.HandlerFunc(router.registerHandler))
		api.Post("/login", http.HandlerFunc(router.loginHandler))

		api.Group(func(protected chi.Router) {
			protected.Use(middleware.NewJWT(jwtSecret).Authenticate)

			protected.Get("/user", router.infoOwnHandler)
			protected.Patch("/user/{userId}", router.updateUserHandler)
			protected.Delete("/user/{userId}", router.deleteUserHandler)
			protected.With(middleware.RequireRoles(AdminRole)).Get("/user/{userId}", router.infoUserHandler)
			protected.With(middleware.RequireRoles(AdminRole)).Post("/users", router.infoAllUsersHandler)

			protected.Post("/announcement/create", router.createAnnHandler)
			protected.Post("/announcement/images/add", router.AddimagesHandler)
		})
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

func (r *Router) infoOwnHandler(w http.ResponseWriter, req *http.Request) {
	h := NewUserHandler(r.service)
	h.InfoOwnUser(w, req)
}

func (r *Router) infoUserHandler(w http.ResponseWriter, req *http.Request) {
	h := NewUserHandler(r.service)
	h.InfoUser(w, req)
}

func (r *Router) updateUserHandler(w http.ResponseWriter, req *http.Request) {
	h := NewUserHandler(r.service)
	h.UpdateUser(w, req)
}

func (r *Router) deleteUserHandler(w http.ResponseWriter, req *http.Request) {
	h := NewUserHandler(r.service)
	h.DeleteUser(w, req)
}

func (r *Router) infoAllUsersHandler(w http.ResponseWriter, req *http.Request) {
	h := NewUserHandler(r.service)
	h.UsersInfo(w, req)
}

func (r *Router) createAnnHandler(w http.ResponseWriter, req *http.Request) {
	h := NewAnnHandler(r.service)
	h.CreateAnnouncement(w, req)
}

func (r *Router) AddimagesHandler(w http.ResponseWriter, req *http.Request) {
	h := NewImageHandler(r.service)
	h.AddImagesToAnnouncement(w, req)
}
