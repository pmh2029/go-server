package handlers

import (
	"go-server/internal/pkg/domains/interfaces"
	"go-server/internal/pkg/domains/models/dtos"
	"go-server/internal/pkg/domains/models/entities"
	"go-server/internal/pkg/repositories"
	"go-server/internal/pkg/usecases"
	"go-server/pkg/shared/utils"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"gorm.io/gorm"
)

type userHandler struct {
	userUsecase  interfaces.UserUsecase
	logger       *logrus.Logger
	googleConfig *oauth2.Config
}

func NewUserHandler(logger *logrus.Logger, db *gorm.DB) *userHandler {
	userRepo := repositories.NewUserRepository(db, logger)
	userUsecase := usecases.NewUserUsecase(userRepo, logger)

	googleConfig := utils.SetupConfig()
	return &userHandler{
		userUsecase,
		logger,
		googleConfig,
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

func (h *userHandler) GoogleLogin(c *gin.Context) {
	// Create oauthState cookie
	oauthState := utils.GenerateStateOauthCookie(c)
	/*
		AuthCodeURL receive state that is a token to protect the user
		from CSRF attacks. You must always provide a non-empty string
		and validate that it matches the the state query parameter
		on your redirect callback.
	*/

	h.googleConfig.AuthCodeURL(oauthState)
	c.Redirect(http.StatusTemporaryRedirect, h.googleConfig.AuthCodeURL(oauthState))
}

func (h *userHandler) GoogleCallback(c *gin.Context) {
	oauthstate, err := c.Cookie("oauthstate")
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

	state := c.Request.FormValue("state")
	code := c.Request.FormValue("code")

	c.Request.Header.Add("content-type", "application/json")
	if state != oauthstate {
		c.JSON(http.StatusTemporaryRedirect, dtos.BaseResponse{
			Code: 2,
			Error: &dtos.ErrorResponse{
				ErrorMessage: "Invalid state",
				ErrorDetails: "Invalid state",
			},
		})
		return
	}

	token, err := h.googleConfig.Exchange(c, code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dtos.BaseResponse{
			Code: 3,
			Error: &dtos.ErrorResponse{
				ErrorMessage: err.Error(),
				ErrorDetails: err.Error(),
			},
		})
		return
	}

	response, err := http.Get(utils.OauthGoogleUrlAPI + token.AccessToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dtos.BaseResponse{
			Code: 4,
			Error: &dtos.ErrorResponse{
				ErrorMessage: err.Error(),
				ErrorDetails: err.Error(),
			},
		})
		return
	}

	defer response.Body.Close()
	contents, err := io.ReadAll(response.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dtos.BaseResponse{
			Code: 5,
			Error: &dtos.ErrorResponse{
				ErrorMessage: err.Error(),
				ErrorDetails: err.Error(),
			},
		})
		return
	}

	// send back response to browser
	log.Println(string(contents))
}
