package httpserver

import (
	"net/http"
	"os"
	"time"

	"dkmbackend/internal/config"
	"dkmbackend/internal/httpserver/handlers"
	"dkmbackend/internal/httpserver/middleware"
	"dkmbackend/internal/services"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func NewRouter(cfg config.Config, blog *services.BlogService, product *services.ProductService, career *services.CareerService, auth *services.AuthService) http.Handler {
	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   cfg.AllowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// static for uploads
	r.Handle("/uploads/*", http.StripPrefix("/uploads/", http.FileServer(http.Dir("uploads"))))

	// health
	r.Get("/api/health", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(200); _, _ = w.Write([]byte("OK")) })

	// auth
	authH := handlers.NewAuthHandler(auth)
	r.Post("/api/auth/login", authH.Login)
	r.Get("/api/auth/me", authH.Me)

	// upload
	base := os.Getenv("BACKEND_URL")
	if base == "" {
		base = "http://localhost:" + itoa(cfg.Port)
	}
	uploadH := handlers.NewUploadHandler(base)
	r.Post("/api/upload", uploadH.Upload)

	// applications
	appH := handlers.NewApplicationHandler()
	// public POST - rate limited to prevent abuse (in-memory)
	r.With(middleware.RateLimit(5, time.Minute)).Post("/api/applications", appH.Submit)

	// blogs
	blogH := handlers.NewBlogHandler(blog)
	r.Get("/api/blogs", blogH.List)
	r.Get("/api/blogs/{param}", blogH.Get) // id or slug
	r.Group(func(rt chi.Router) {
		rt.Use(middleware.Auth(auth))
		rt.Post("/api/blogs", blogH.Create)
		rt.Put("/api/blogs/{id}", blogH.Update)
		rt.Delete("/api/blogs/{id}", blogH.Delete)
	})

	// products
	prodH := handlers.NewProductHandler(product)
	r.Get("/api/products", prodH.List)
	r.Get("/api/products/{id}", prodH.Get)
	r.Group(func(rt chi.Router) {
		rt.Use(middleware.Auth(auth))
		rt.Post("/api/products", prodH.Create)
		rt.Put("/api/products/{id}", prodH.Update)
		rt.Delete("/api/products/{id}", prodH.Delete)
	})

	// careers
	carH := handlers.NewCareerHandler(career)
	r.Get("/api/careers/public", carH.PublicList)
	r.Group(func(rt chi.Router) {
		rt.Use(middleware.Auth(auth))
		rt.Get("/api/careers", carH.List)
		rt.Get("/api/careers/{id}", carH.Get)
		rt.Post("/api/careers", carH.Create)
		rt.Put("/api/careers/{id}", carH.Update)
		rt.Delete("/api/careers/{id}", carH.Delete)
	})

	// admin application routes
	r.Group(func(rt chi.Router) {
		rt.Use(middleware.Auth(auth))
		rt.Get("/api/admin/applications", appH.AdminList)
		rt.Get("/api/admin/applications/{id}", appH.AdminGet)
		rt.Delete("/api/admin/applications/{id}", appH.AdminDelete)
		rt.Get("/api/admin/applications/cv/{filename}", appH.AdminGetCV)
	})

	return r
}

// tiny int to string to avoid fmt import here
func itoa(i int) string {
	if i == 0 {
		return "0"
	}
	neg := false
	if i < 0 {
		neg = true
		i = -i
	}
	var b [20]byte
	bp := len(b)
	for i > 0 {
		bp--
		b[bp] = byte('0' + i%10)
		i /= 10
	}
	if neg {
		bp--
		b[bp] = '-'
	}
	return string(b[bp:])
}
