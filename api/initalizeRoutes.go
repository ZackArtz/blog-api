package api

import (
	"github.com/gofiber/cors"
	"github.com/gofiber/fiber"
	"github.com/gofiber/fiber/middleware"
)

func (s *Server) initializeRoutes() {
	s.Router = fiber.New()

	s.Router.Use(cors.New())
	s.Router.Use(middleware.Logger())

	s.Router.Get("/", func(ctx *fiber.Ctx) {
		ctx.JSON(struct {
			Message string `json:"message"`
		}{
			Message: "hello!",
		})
	})

	s.Router.Post("/api/login", s.Login)

	s.Router.Get("/api/users", s.GetAllUsers)
	s.Router.Post("/api/users", s.CreateUser)
	s.Router.Get("/api/users/:id", s.GetUser)
	s.Router.Put("/api/users/:id", s.UpdateUser)
	s.Router.Delete("/api/users/:id", s.DeleteUser)

	s.Router.Get("/api/articles", s.GetArticles)
	s.Router.Post("/api/articles", s.CreateArticle)
	s.Router.Get("/api/articles/:id", s.GetArticleByID)
	s.Router.Get("/api/articles/slug/:slug", s.GetArticleBySlug)
	s.Router.Post("/api/articles/:id", s.UpdateArticle)
	s.Router.Delete("/api/articles/:id", s.DeleteArticle)
}
