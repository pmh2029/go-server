package handlers

import (
	"go-server/internal/pkg/domains/interfaces"
	"go-server/internal/pkg/domains/models/dtos"
	"go-server/internal/pkg/usecases"
	"net/http"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type uploadHandler struct {
	uploadUsecase interfaces.UploadUsecase
	logger        *logrus.Logger
}

func NewUploadHandler(
	cld *cloudinary.Cloudinary,
	logger *logrus.Logger,
) *uploadHandler {
	uploadUsecase := usecases.NewUploadUsecase(cld, logger)
	return &uploadHandler{
		uploadUsecase,
		logger,
	}
}

func (h *uploadHandler) FileUpload(c *gin.Context) {
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusOK, dtos.BaseResponse{
			Code:    400,
			Message: BadRequest,
			Error: &dtos.ErrorResponse{
				ErrorDetails: err,
			},
		})
		return
	}

	url, err := h.uploadUsecase.FileUpload(c, file)
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
		Code: 0,
		Data: gin.H{
			"file": url,
		},
		Message: "OK",
	})
}
