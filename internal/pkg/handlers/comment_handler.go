package handlers

import (
	"errors"
	"go-server/internal/pkg/domains/models/dtos"
	"go-server/internal/pkg/domains/models/entities"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type commentHandler struct {
	db     *gorm.DB
	logger *logrus.Logger
}

func NewCommentHandler(db *gorm.DB, logger *logrus.Logger) *commentHandler {
	return &commentHandler{
		db:     db,
		logger: logger,
	}
}

func (h *commentHandler) Create(c *gin.Context) {
	req := dtos.CreateCommentRequestDto{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, dtos.BaseResponse{
			Code:    400,
			Message: "Bad Request",
			Error: &dtos.ErrorResponse{
				ErrorDetails: err,
			},
		})
		return
	}

	newComment := entities.Comment{
		PlaceID: req.PlaceID,
		UserID:  req.UserID,
		Comment: req.Comment,
		Rate:    req.Rate,
	}

	err = h.db.Create(&newComment).Error
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
			"comment": gin.H{
				"id": newComment.ID,
			},
		},
	})
}

func (h *commentHandler) Update(c *gin.Context) {
	commentIDParam := c.Param("comment_id")
	commentID, err := strconv.Atoi(commentIDParam)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dtos.BaseResponse{
			Code:    0,
			Message: InternalServerError,
			Error: &dtos.ErrorResponse{
				ErrorDetails: err,
			},
		})
		return
	}

	req := dtos.UpdateCommentRequestDto{}
	err = c.ShouldBindJSON(&req)
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

	var comment entities.Comment
	err = h.db.Where("id = ?", commentID).Find(&comment).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusOK, dtos.BaseResponse{
				Code:    1,
				Message: "Comment not found",
				Error: &dtos.ErrorResponse{
					ErrorDetails: err,
				},
			})
			return
		}
		c.JSON(http.StatusInternalServerError, dtos.BaseResponse{
			Code:    0,
			Message: InternalServerError,
			Error: &dtos.ErrorResponse{
				ErrorDetails: err,
			},
		})
		return
	}

	err = h.db.Model(&comment).Where("id = ?", commentID).Updates(map[string]interface{}{
		"comment": req.Comment,
		"rate":    req.Rate,
	}).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, dtos.BaseResponse{
			Code:    0,
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
	})
}

func (h *commentHandler) Delete(c *gin.Context) {
	commentIDParam := c.Param("comment_id")
	commentID, err := strconv.Atoi(commentIDParam)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dtos.BaseResponse{
			Code:    0,
			Message: InternalServerError,
			Error: &dtos.ErrorResponse{
				ErrorDetails: err,
			},
		})
		return
	}

	err = h.db.Where("id = ?", commentID).Delete(&entities.Comment{}).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, dtos.BaseResponse{
			Code:    0,
			Message: InternalServerError,
			Error: &dtos.ErrorResponse{
				ErrorDetails: err,
			},
		})
		return
	}

	c.JSON(http.StatusOK, dtos.BaseResponse{
		Code:    0,
		Message: "Deleted success",
	})
}


