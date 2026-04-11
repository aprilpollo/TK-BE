package main

import (
	"aprilpollo/internal/adapters/config"
	googleAdapter "aprilpollo/internal/adapters/google"
	"aprilpollo/internal/adapters/middleware"
	"aprilpollo/internal/adapters/routes"
	"aprilpollo/internal/adapters/routes/handler"
	"aprilpollo/internal/adapters/storage/cache"
	"aprilpollo/internal/adapters/storage/orm"
	"aprilpollo/internal/adapters/storage/repository"
	"aprilpollo/internal/core/services"
	"aprilpollo/internal/utils"
	"fmt"
	"log"
	"runtime/debug"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	fiberLogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	cfg, err := config.GetConfig()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("✔ [INFO] Loading Configuration")

	db, err := orm.NewGormDB(cfg.Database, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	fmt.Println("✔ [INFO] Database Connection")

	redis, err := cache.NewRedisClient(cfg.Redis)
	if err != nil {
		log.Fatal(err)
	}
	defer redis.Close()

	fmt.Println("✔ [INFO] Redis Connection")

	// --- Repositories (output adapters) ---
	oauthRepo := repository.NewOauthRepository(db.GetDB())
	orgRepo := repository.NewOrganizationRepository(db.GetDB())
	userRepo := repository.NewUserRepository(db.GetDB())

	// --- Services (core / use cases) ---
	googleVerifier := googleAdapter.NewGoogleVerifier(cfg.Oauth.GoogleProvider.ClientID)
	oauthSvc := services.NewOauthService(oauthRepo, googleVerifier, utils.JWTConfig{
		SecretKey:  cfg.JWT.SecretKey,
		Issuer:     cfg.JWT.Issuer,
		Subject:    cfg.JWT.Subject,
		ExpireDays: cfg.JWT.JwtExpireDaysCount,
	})

	orgSvc := services.NewOrganizationService(orgRepo)
	userSvc := services.NewUserService(userRepo, orgRepo)

	// --- Middleware ---
	jwtMiddleware := middleware.JWTProtected(cfg.JWT.SecretKey)

	// --- Handlers (input adapters) ---
	oauthHandler := handler.NewOauthHandler(oauthSvc)
	orgHandler := handler.NewOrganizationHandler(orgSvc)
	userHandler := handler.NewUserHandler(userSvc, orgSvc)

	// --- Fiber app ---
	app := fiber.New(fiber.Config{
		AppName: cfg.App.AppName,
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.App.AllowedCredentialOrigins,
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization,X-Organization-ID",
		AllowCredentials: false,
	}))

	app.Use(fiberLogger.New(fiberLogger.Config{
		Format: "${time} | ${status} | ${latency} | ${ip} | ${method} | ${path}\n",
	}))

	app.Use(recover.New(recover.Config{
		EnableStackTrace: true,
		StackTraceHandler: func(c *fiber.Ctx, e interface{}) {
			log.Printf("[PANIC] path=%s method=%s error=%v\n%s",
				c.Path(),
				c.Method(),
				e,
				debug.Stack(),
			)
			c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "internal server error",
			})
		},
	}))

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"app":     cfg.App.AppName,
			"version": cfg.App.AppVersion,
		})
	})

	routes.RegisterOauthRoutes(app, oauthHandler)
	routes.RegisterUserRoutes(app, userHandler, jwtMiddleware)
	routes.RegisterOrganizationRoutes(app, orgHandler, jwtMiddleware)

	if err := app.Listen(fmt.Sprintf(":%s", cfg.App.ApiPort)); err != nil {
		log.Println(err)
	}
}
