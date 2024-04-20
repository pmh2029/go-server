package handlers

import (
	"go-server/internal/pkg/domains/interfaces"
	"go-server/internal/pkg/domains/models/dtos"
	"go-server/internal/pkg/domains/models/entities"
	"go-server/internal/pkg/repositories"
	"go-server/internal/pkg/usecases"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type userHandler struct {
	userUsecase interfaces.UserUsecase
	logger      *logrus.Logger
}

func NewUserHandler(logger *logrus.Logger, db *gorm.DB) *userHandler {
	userRepo := repositories.NewUserRepository(db, logger)
	userUsecase := usecases.NewUserUsecase(userRepo, logger)
	return &userHandler{
		userUsecase,
		logger,
	}
}

func (h *userHandler) Register(c *gin.Context) {
	req := dtos.RegisterRequestDto{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, dtos.BaseResponse{
			Code: 1,
			Error: &dtos.ErrorResponse{
				ErrorDetails: err.Error(),
			},
		})
		return
	}

	user := entities.User{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	}

	user, err = h.userUsecase.Create(c, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dtos.BaseResponse{
			Code: 1,
			Error: &dtos.ErrorResponse{
				ErrorMessage: err.Error(),
				ErrorDetails: err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, dtos.BaseResponse{
		Code: 0,
		Data: dtos.RegisterResponseDto{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
		},
	})
}
