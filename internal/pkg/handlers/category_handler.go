package handlers

import (
	"errors"
	"go-server/internal/pkg/domains/interfaces"
	"go-server/internal/pkg/domains/models/dtos"
	"go-server/internal/pkg/repositories"
	"go-server/internal/pkg/usecases"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type categoryHandler struct {
	categoryUsecase interfaces.CategoryUsecase
	logger          *logrus.Logger
}

func NewCategoryHandler(
	logger *logrus.Logger,
	db *gorm.DB,
) *categoryHandler {
	categoryRepo := repositories.NewCategoryRepository(db, logger)
	categoryUsecase := usecases.NewCategoryRepository(categoryRepo, logger)

	return &categoryHandler{
		categoryUsecase,
		logger,
	}
}

func (h *categoryHandler) CreateCategory(c *gin.Context) {
	req := dtos.CreateCategoryRequestDto{}
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

	category, err := h.categoryUsecase.Create(c, req)
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
		Data: gin.H{
			"category": category,
		},
	})
}

func (h *categoryHandler) ListCategory(c *gin.Context) {
	categories, err := h.categoryUsecase.FindByConditions(c, map[string]interface{}{})
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
			"categories": categories,
		},
	})
}

func (h *categoryHandler) Update(c *gin.Context) {
	categoryIDParam := c.Param("category_id")

	categoryID, err := strconv.Atoi(categoryIDParam)
	if err != nil {
		c.JSON(http.StatusOK, dtos.BaseResponse{
			Code:    400,
			Message: BadRequest,
			Error: &dtos.ErrorResponse{
				ErrorDetails: err.Error(),
			},
		})
		return
	}

	req := dtos.UpdateCategoryRequestDto{}
	err = c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusOK, dtos.BaseResponse{
			Code:    400,
			Message: BadRequest,
			Error: &dtos.ErrorResponse{
				ErrorDetails: err.Error(),
			},
		})
		return
	}

	category, err := h.categoryUsecase.Update(c, req, map[string]interface{}{
		"id": categoryID,
	})

	if err != nil {
		if errors.Is(err, usecases.UpdateCategoryIDNotFound) {
			c.JSON(http.StatusOK, dtos.BaseResponse{
				Code:    1,
				Message: "Category not found",
				Error: &dtos.ErrorResponse{
					ErrorDetails: err,
				},
			})
			return
		}
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
		Message: "Updated success",
		Data: gin.H{
			"category": category,
		},
	})
}

func (h *categoryHandler) DetailCategory(c *gin.Context) {
	categoryIDParam := c.Param("category_id")

	categoryID, err := strconv.Atoi(categoryIDParam)
	if err != nil {
		c.JSON(http.StatusOK, dtos.BaseResponse{
			Code:    400,
			Message: BadRequest,
			Error: &dtos.ErrorResponse{
				ErrorDetails: err.Error(),
			},
		})
		return
	}

	category, err := h.categoryUsecase.TakeByConditions(c, map[string]interface{}{
		"id": categoryID,
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusOK, dtos.BaseResponse{
				Code:    1,
				Message: "Category not found",
				Error: &dtos.ErrorResponse{
					ErrorDetails: err,
				},
			})
			return
		}
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
			"category": category,
		},
	})
}

func (h *categoryHandler) DeleteCategory(c *gin.Context) {
	categoryIDParam := c.Param("category_id")

	categoryID, err := strconv.Atoi(categoryIDParam)
	if err != nil {
		c.JSON(http.StatusOK, dtos.BaseResponse{
			Code:    400,
			Message: BadRequest,
			Error: &dtos.ErrorResponse{
				ErrorDetails: err.Error(),
			},
		})
		return
	}

	err = h.categoryUsecase.DeleteByConditions(c, map[string]interface{}{
		"id": categoryID,
	})
	if err != nil {
		if errors.Is(err, usecases.DeleteCategoryIDNotFound) {
			c.JSON(http.StatusOK, dtos.BaseResponse{
				Code:    1,
				Message: "Category not found",
				Error: &dtos.ErrorResponse{
					ErrorDetails: err,
				},
			})
			return
		}
		if errors.Is(err, usecases.DeleteCategoryCategoryHasPlace) {
			c.JSON(http.StatusOK, dtos.BaseResponse{
				Code:    2,
				Message: err.Error(),
				Error: &dtos.ErrorResponse{
					ErrorDetails: err,
				},
			})
			return
		}
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
	})
}
