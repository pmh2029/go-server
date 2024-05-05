package router

import (
	"go-server/internal/pkg/domains/models/dtos"
	"go-server/internal/pkg/handlers"
	"go-server/pkg/shared/middleware"
	"go-server/pkg/shared/validator"
	"net/http"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Router struct {
	Engine *gin.Engine
	DB     *gorm.DB
	Logger *logrus.Logger
}

func (r *Router) InitializeRouter(logger *logrus.Logger) {
	r.Engine.Use(gin.Logger())
	r.Engine.Use(gin.Recovery())
	r.Engine.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"*"},
		AllowCredentials: true,
		ExposeHeaders:    []string{"Set-Cookie"},
		AllowWebSockets:  true,
		AllowFiles:       true,
	}))
	// r.Engine.Use(func(c *gin.Context) {
	// 	if c.Request.Method == "OPTIONS" {
	// 		c.AbortWithStatus(204)
	// 		return
	// 	}
	// 	c.Next()
	// })
	r.Logger = logger
}

func (r *Router) SetupHandler() {
	_ = validator.New()
	cld, err := cloudinary.NewFromParams(
		os.Getenv("CLOUDINARY_CLOUD_NAME"),
		os.Getenv("CLOUDINARY_API_KEY"),
		os.Getenv("CLOUDINARY_API_SECRET"),
	)
	if err != nil {
		return
	}

	r.Engine.Use(middleware.RequestID())

	userHandler := handlers.NewUserHandler(r.Logger, r.DB)
	uploadHandler := handlers.NewUploadHandler(cld, r.Logger)
	bannerHandler := handlers.NewBannerHandler(r.Logger, r.DB)

	// health check
	r.Engine.GET("/", func(c *gin.Context) {
		data := dtos.BaseResponse{
			Code:  0,
			Data:  gin.H{"message": "Health check OK!"},
			Error: nil,
		}
		c.JSON(http.StatusOK, data)
	})

	// router api
	publicApi := r.Engine.Group("/api")
	{
		// auth
		authApi := publicApi.Group("/auth")
		{
			authApi.POST("/register", userHandler.Register)
			authApi.POST("/login", userHandler.Login)
			authApi.GET("/google/login", userHandler.GoogleLogin)
			authApi.GET("/google/callback", userHandler.GoogleCallback)
		}

		//
		uploadApi := publicApi.Group("/upload")
		{
			uploadApi.POST("/", uploadHandler.FileUpload)
		}
	}

	privateApi := r.Engine.Group("/api")
	privateApi.Use(middleware.CheckAuthentication())

	{
		userApi := privateApi.Group("/user")
		{
			userApi.PATCH("/:user_id", userHandler.Update)
		}

		bannerApi := privateApi.Group("/banner")
		{
			bannerApi.POST("/", bannerHandler.CreateBanner)
			bannerApi.GET("/", bannerHandler.ListBanner)
			bannerApi.PATCH("/:banner_id", bannerHandler.Update)
			bannerApi.GET("/:banner_id", bannerHandler.DetailBanner)
			bannerApi.DELETE("/:banner_id", bannerHandler.DeleteBanner)
		}
	}
}
