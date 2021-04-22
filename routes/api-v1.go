package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/pantame/server/api/auth"
	"github.com/pantame/server/api/file"
	"github.com/pantame/server/api/session"
	"github.com/pantame/server/api/user"
	"github.com/pantame/server/api/verification"
)

func Api1(router *fiber.App) {
	v1 := router.Group("/api/v1")


	// SESSIONS
	v1.Post("/sessions", session.NewSession)
	v1.Post("/sessions/validations", session.Validate)

	// MANDATORY LOGIN BELOW
	v1.Use(auth.Auth)

	v1.Get("/sessions", session.GetAllActiveSessionsByUserID)
	v1.Delete("/sessions", session.Logout)

	v1.Post("/verifications", verification.NewVerificationCode)

	// USERS
	v1.Put("/users/name", user.UpdateNameUser)
	v1.Get("/users", user.GetUser)
	v1.Get("/users/access-passes", user.GetAccessPasses)
	v1.Post("/users/access-passes", user.NewAccessPass)
	v1.Delete("/users/access-passes", user.RemoveAccessPass)

	// FILES
	files := v1.Group("/files")
	files.Post("/metadata", file.PostMetadata)
	files.Get("/metadata/all", file.GetAllMetadata)
	files.Get("/metadata/:uuid", file.GetMetadata)
	files.Post("/parts", file.PostPart)
	files.Get("/parts/:uuid/:n", file.GetPart)
}
