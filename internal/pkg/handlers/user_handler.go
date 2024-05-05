package handlers

import (
	"errors"
	"go-server/internal/pkg/domains/interfaces"
	"go-server/internal/pkg/domains/models/dtos"
	"go-server/internal/pkg/domains/models/entities"
	"go-server/internal/pkg/repositories"
	"go-server/internal/pkg/usecases"
	"go-server/pkg/shared/auth"
	"go-server/pkg/shared/utils"
	"io"
	"log"
	"net/http"
	"strconv"

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
		c.JSON(http.StatusOK, dtos.BaseResponse{
			Code:    400,
			Message: "Bad Request",
			Error: &dtos.ErrorResponse{
				ErrorDetails: err,
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
		if errors.Is(err, usecases.RegisterEmailExisted) {
			c.JSON(http.StatusOK, dtos.BaseResponse{
				Code:    1,
				Message: "Email existed",
				Error: &dtos.ErrorResponse{
					ErrorDetails: err.Error(),
				},
			})
			return
		}
		if errors.Is(err, usecases.RegisterUsernameExisted) {
			c.JSON(http.StatusOK, dtos.BaseResponse{
				Code:    2,
				Message: "Username existed",
				Error: &dtos.ErrorResponse{
					ErrorDetails: err.Error(),
				},
			})
			return
		}
		c.JSON(http.StatusInternalServerError, dtos.BaseResponse{
			Code:    0,
			Message: "Internal Server Error",
			Error: &dtos.ErrorResponse{
				ErrorDetails: err.Error(),
			},
		})
		return
	}

	accessToken, err := auth.GenerateHS256JWT(map[string]interface{}{
		"user_id": user.ID,
		"sub":     user.Username,
		"email":   user.Email,
	})
	if err != nil {
		c.JSON(http.StatusOK, dtos.BaseResponse{
			Code:    0,
			Message: "Internal Server Error",
			Error: &dtos.ErrorResponse{
				ErrorDetails: err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, dtos.BaseResponse{
		Code:    0,
		Message: "OK",
		Data: dtos.RegisterResponseDto{
			User: entities.User{
				ID:       user.ID,
				Username: user.Username,
				Email:    user.Email,
			},
			AccessToken: accessToken,
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
				ErrorDetails: err.Error(),
			},
		})
		return
	}

	// send back response to browser
	log.Println(string(contents))
}

func (h *userHandler) Login(c *gin.Context) {
	req := dtos.LoginRequestDto{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, dtos.BaseResponse{
			Code: 400,
			Error: &dtos.ErrorResponse{
				ErrorDetails: err.Error(),
			},
		})
		return
	}

	user, accessToken, err := h.userUsecase.Login(c, req)
	if err != nil {
		if errors.Is(err, usecases.EmailNotFound) {
			c.JSON(http.StatusOK, dtos.BaseResponse{
				Code:    1,
				Message: "Email not found",
				Error: &dtos.ErrorResponse{
					ErrorDetails: err.Error(),
				},
			})
			return
		}
		if errors.Is(err, usecases.WrongPassword) {
			c.JSON(http.StatusOK, dtos.BaseResponse{
				Code:    2,
				Message: "Wrong password",
				Error: &dtos.ErrorResponse{
					ErrorDetails: err.Error(),
				},
			})
			return
		}
		c.JSON(http.StatusInternalServerError, dtos.BaseResponse{
			Code: 0,
			Error: &dtos.ErrorResponse{
				ErrorDetails: err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, dtos.BaseResponse{
		Code: 0,
		Data: dtos.LoginResponseDto{
			User: entities.User{
				ID:       user.ID,
				Username: user.Username,
				Email:    user.Email,
			},
			AccessToken: accessToken,
		},
		Message: "Login success",
	})
}

func (h *userHandler) Update(c *gin.Context) {
	userIDParam := c.Param("user_id")

	userID, err := strconv.Atoi(userIDParam)
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

	req := dtos.UpdateUserRequestDto{}
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

	user, err := h.userUsecase.Update(c, map[string]interface{}{
		"id": userID,
	}, req)
	if err != nil {
		if errors.Is(err, usecases.UpdateUserIDNotFound) {
			c.JSON(http.StatusOK, dtos.BaseResponse{
				Code:    1,
				Message: "User not found",
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
		Data: dtos.UpdateUserResponseDto{
			User: user,
		},
	})
}

func (h *userHandler) DetailUser(c *gin.Context) {
	userID := c.MustGet("user_id")

	if userID == nil {
		c.JSON(http.StatusOK, dtos.BaseResponse{
			Code:    1,
			Message: "Unauthorized",
		})
		return
	}

	user, err := h.userUsecase.TakeByConditions(c, map[string]interface{}{
		"id": userID,
	})
	if err != nil {
		if errors.Is(err, usecases.DetailUserIDNotFound) {
			c.JSON(http.StatusOK, dtos.BaseResponse{
				Code:    2,
				Message: "User not found",
				Error:   &dtos.ErrorResponse{ErrorDetails: err},
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
			"user": user,
		},
	})
}
