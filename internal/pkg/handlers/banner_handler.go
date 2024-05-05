package handlers

import (
	"go-server/internal/pkg/domains/interfaces"
	"go-server/internal/pkg/domains/models/dtos"
	"go-server/internal/pkg/repositories"
	"go-server/internal/pkg/usecases"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type bannerHandler struct {
	bannerUsecase interfaces.BannerUsecase
	logger        *logrus.Logger
}

func NewBannerHandler(
	logger *logrus.Logger,
	db *gorm.DB,
) *bannerHandler {
	bannerRepo := repositories.NewBannerRepository(db, logger)
	bannerUsecase := usecases.NewBannerUsecase(bannerRepo, logger)

	return &bannerHandler{
		bannerUsecase,
		logger,
	}
}

func (h *bannerHandler) CreateBanner(c *gin.Context) {
	req := dtos.CreateBannerRequestDto{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusOK, dtos.BaseResponse{
			Code:    400,
			Message: "Bad Request",
			Error: &dtos.ErrorResponse{
				ErrorDetails: err,
			},
		})
		return
	}

	banner, err := h.bannerUsecase.Create(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dtos.BaseResponse{
			Message: InternalServerError,
			Error: &dtos.ErrorResponse{
				ErrorDetails: err,
			},
		})
		return
	}

	c.JSON(http.StatusOK, dtos.BaseResponse{
		Code:    0,
		Message: "Created success",
		Data: dtos.CreateBannerResponseDto{
			Banner: banner,
		},
	})
}

func (h *bannerHandler) ListBanner(c *gin.Context) {
	banners, err := h.bannerUsecase.FindByConditions(c, map[string]interface{}{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, dtos.BaseResponse{
			Message: InternalServerError,
			Error: &dtos.ErrorResponse{
				ErrorDetails: err,
			},
		})
		return
	}

	c.JSON(http.StatusOK, dtos.BaseResponse{
		Code:    0,
		Message: "OK",
		Data: gin.H{
			"banners": banners,
		},
	})
}
