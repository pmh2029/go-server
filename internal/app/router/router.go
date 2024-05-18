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
	categoryHandler := handlers.NewCategoryHandler(r.Logger, r.DB)
	placeHandler := handlers.NewPlaceHandler(r.Logger, r.DB)
	tripHandler := handlers.NewTripHandler(r.Logger, r.DB)

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
			authApi.POST("/admin/login", userHandler.AdminLogin)
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
	adminApi := r.Engine.Group("/api")
	adminApi.Use(middleware.CheckAuthentication(), middleware.CheckRole())

	{
		userApi := privateApi.Group("/app/user")
		{
			userApi.PATCH("/:user_id", userHandler.Update)
			userApi.GET("/info", userHandler.DetailUser)
		}

		bannerApi := adminApi.Group("/banner")
		{
			bannerApi.POST("/", bannerHandler.CreateBanner)
			bannerApi.GET("/", bannerHandler.ListBanner)
			bannerApi.PATCH("/:banner_id", bannerHandler.Update)
			bannerApi.GET("/:banner_id", bannerHandler.DetailBanner)
			bannerApi.DELETE("/:banner_id", bannerHandler.DeleteBanner)
		}

		bannerAppApi := privateApi.Group("/app/banner")
		{
			bannerAppApi.GET("/", bannerHandler.ListBanner)
			bannerAppApi.GET("/:banner_id", bannerHandler.DetailBanner)
		}

		categoryApi := adminApi.Group("/category")
		{
			categoryApi.POST("/", categoryHandler.CreateCategory)
			categoryApi.GET("/", categoryHandler.ListCategory)
			categoryApi.PATCH("/:category_id", categoryHandler.Update)
			categoryApi.GET("/:category_id", categoryHandler.DetailCategory)
			categoryApi.DELETE("/:category_id", categoryHandler.DeleteCategory)
		}

		categoryAppApi := privateApi.Group("/app/category")
		{
			categoryAppApi.GET("/", categoryHandler.ListCategory)
			categoryAppApi.GET("/:category_id", categoryHandler.DetailCategory)
		}

		placeApi := adminApi.Group("/place")
		{
			placeApi.POST("/", placeHandler.CreatePlace)
			placeApi.GET("/", placeHandler.ListPlacePaginate)
			placeApi.PATCH("/:place_id", placeHandler.UpdatePlace)
			placeApi.GET("/:place_id", placeHandler.DetailPlace)
			placeApi.DELETE("/:place_id", placeHandler.DeletePlace)
			placeApi.GET("/all_places", placeHandler.ListAllPlace)
		}

		placeAppApi := privateApi.Group("/app/place")
		{
			placeAppApi.GET("/", placeHandler.ListPlacePaginate)
			placeAppApi.GET("/:place_id", placeHandler.DetailPlace)
			placeAppApi.GET("/all_places", placeHandler.ListAllPlace)
		}

		tripApi := privateApi.Group("/app/trip")
		{
			tripApi.POST("/", tripHandler.CreateTrip)
			tripApi.GET("/:user_id", tripHandler.ListTrip)
		}
	}
}
