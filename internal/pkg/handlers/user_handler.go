package handlers

import (
	"errors"
	"go-server/internal/pkg/domains/interfaces"
	"go-server/internal/pkg/domains/models/dtos"
	"go-server/internal/pkg/domains/models/entities"
	"go-server/internal/pkg/repositories"
	"go-server/internal/pkg/services"
	"go-server/internal/pkg/usecases"
	"go-server/pkg/shared/auth"
	"go-server/pkg/shared/database"
	"go-server/pkg/shared/utils"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"gorm.io/gorm"
)

type userHandler struct {
	userUsecase  interfaces.UserUsecase
	logger       *logrus.Logger
	googleConfig *oauth2.Config
	mailService  services.MailServiceInterface
	db           *gorm.DB
}

func NewUserHandler(logger *logrus.Logger, db *gorm.DB) *userHandler {
	userRepo := repositories.NewUserRepository(db, logger)
	userUsecase := usecases.NewUserUsecase(userRepo, logger)
	mailService := services.NewMailService()

	googleConfig := utils.SetupConfig()
	return &userHandler{
		userUsecase,
		logger,
		googleConfig,
		mailService,
		db,
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
		Active:   true,
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

	tokenID := uuid.New().String()
	accessToken, err := auth.GenerateHS256JWT(map[string]interface{}{
		"user_id":  user.ID,
		"sub":      user.Username,
		"email":    user.Email,
		"is_admin": user.IsAdmin,
		"token_id": tokenID,
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

	err = h.db.Create(&entities.UserToken{
		UserID:  user.ID,
		TokenID: tokenID,
	}).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, dtos.BaseResponse{
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

	user, accessToken, tokenID, err := h.userUsecase.Login(c, req)
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
		if !user.Active {
			c.JSON(http.StatusOK, dtos.BaseResponse{
				Code:    3,
				Message: "User not active",
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

	userToken := entities.UserToken{
		UserID:  user.ID,
		TokenID: tokenID,
	}
	err = h.db.Create(&userToken).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, dtos.BaseResponse{
			Code:    0,
			Message: InternalServerError,
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

func (h *userHandler) AdminLogin(c *gin.Context) {
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

	user, accessToken, tokenID, err := h.userUsecase.Login(c, req)
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
	if !user.IsAdmin {
		c.JSON(http.StatusForbidden, dtos.BaseResponse{
			Code:    3,
			Message: "Forbidden",
			Error: &dtos.ErrorResponse{
				ErrorDetails: "Role invalid",
			},
		})
		return
	}

	err = h.db.Create(&entities.UserToken{
		UserID:  user.ID,
		TokenID: tokenID,
	}).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, dtos.BaseResponse{
			Code:    0,
			Message: InternalServerError,
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

func (h *userHandler) ForgotPassword(c *gin.Context) {
	req := dtos.ForgotPasswordRequestDto{}

	err := c.ShouldBindJSON(&req)
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

	user, err := h.userUsecase.TakeByConditions(c, map[string]interface{}{
		"email": req.Email,
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusOK, dtos.BaseResponse{
				Code:    1,
				Message: "Email not found",
				Error: &dtos.ErrorResponse{
					ErrorDetails: err.Error(),
				},
			})
			return
		}

		c.JSON(http.StatusInternalServerError, dtos.BaseResponse{
			Message: InternalServerError,
			Error: &dtos.ErrorResponse{
				ErrorDetails: err.Error(),
			},
		})
		return
	}

	newPassword, err := utils.GeneratePassword(10)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dtos.BaseResponse{
			Message: InternalServerError,
			Error: &dtos.ErrorResponse{
				ErrorDetails: err.Error(),
			},
		})
		return
	}

	hashPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dtos.BaseResponse{
			Message: InternalServerError,
			Error: &dtos.ErrorResponse{
				ErrorDetails: err.Error(),
			},
		})
		return
	}

	err = h.db.Model(&user).Where("id = ?", user.ID).UpdateColumn("password", hashPassword).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, dtos.BaseResponse{
			Message: InternalServerError,
			Error: &dtos.ErrorResponse{
				ErrorDetails: err.Error(),
			},
		})
		return
	}

	data := map[string]interface{}{
		"to":       req.ReceiveEmail,
		"username": user.Username,
		"password": newPassword,
		"year":     time.Now().Year(),
	}

	err = h.mailService.SendOneTimePasswordMail("new_password_template.html", data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dtos.BaseResponse{
			Message: InternalServerError,
			Error: &dtos.ErrorResponse{
				ErrorDetails: err.Error(),
			},
		})
		return
	}

	go func(userID int) {
		err = h.db.Where("user_id = ?", userID).Delete(&entities.UserToken{}).Error
		if err != nil {
			return
		}
	}(user.ID)

	c.JSON(http.StatusOK, dtos.BaseResponse{
		Code:    0,
		Message: "OK",
	})
}

func (h *userHandler) ChangePassword(c *gin.Context) {
	userID := c.MustGet("user_id")
	tokenID := c.MustGet("token_id").(string)
	if userID == nil || tokenID == "" {
		c.JSON(http.StatusOK, dtos.BaseResponse{
			Code:    1,
			Message: "Unauthorized",
		})
		return
	}

	var req dtos.ChangePasswordRequestDto
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dtos.BaseResponse{
			Message: BadRequest,
			Error: &dtos.ErrorResponse{
				ErrorDetails: err.Error(),
			},
		})
		return
	}

	user, err := h.userUsecase.TakeByConditions(c, map[string]interface{}{
		"id": userID,
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusOK, dtos.BaseResponse{
				Code:    1,
				Message: "User not found",
				Error: &dtos.ErrorResponse{
					ErrorDetails: "invalid id",
				},
			})
			return
		}
		c.JSON(http.StatusInternalServerError, dtos.BaseResponse{
			Message: InternalServerError,
			Error: &dtos.ErrorResponse{
				ErrorDetails: err.Error(),
			},
		})
		return
	}

	if !utils.CheckPasswordHash(req.OldPassword, user.Password) {
		c.JSON(http.StatusOK, dtos.BaseResponse{
			Code:    2,
			Message: "Wrong password",
			Error: &dtos.ErrorResponse{
				ErrorDetails: "WrongPassword",
			},
		})
		return
	}
	if req.NewPassword != req.ConfirmPassword {
		c.JSON(http.StatusOK, dtos.BaseResponse{
			Code:    3,
			Message: "Password not match",
			Error: &dtos.ErrorResponse{
				ErrorDetails: "PasswordNotMatch",
			},
		})
		return
	}

	hashedPass, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dtos.BaseResponse{
			Message: InternalServerError,
			Error: &dtos.ErrorResponse{
				ErrorDetails: err.Error(),
			},
		})
		return
	}

	user.Password = hashedPass
	err = h.db.Model(&user).Where("id = ?", user.ID).UpdateColumn("password", user.Password).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, dtos.BaseResponse{
			Message: InternalServerError,
			Error: &dtos.ErrorResponse{
				ErrorDetails: err.Error(),
			},
		})
		return
	}

	go func(tokenID string) {
		err = h.db.Where("token_id = ?", tokenID).Delete(&entities.UserToken{}).Error
		if err != nil {
			return
		}
	}(tokenID)

	c.JSON(http.StatusOK, dtos.BaseResponse{
		Code:    0,
		Message: "Change password success",
	})
}

func (h *userHandler) Logout(c *gin.Context) {
	userID := c.MustGet("token_id")
	if userID == nil {
		c.JSON(http.StatusUnauthorized, dtos.BaseResponse{
			Code:    1,
			Message: "Unauthorized",
		})
		return
	}

	go func() {
		err := h.db.Where("token_id = ?", userID).Delete(&entities.UserToken{}).Error
		if err != nil {
			return
		}
	}()

	c.JSON(http.StatusOK, dtos.BaseResponse{
		Code:    0,
		Message: "Logout success",
	})
}

func (h *userHandler) UpdateStatus(c *gin.Context) {
	req := dtos.UpdateStatusRequestDto{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, dtos.BaseResponse{
			Code:    400,
			Message: BadRequest,
			Error: &dtos.ErrorResponse{
				ErrorDetails: err.Error(),
			},
		})
		return
	}

	var user entities.User
	err := h.db.Where("id = ?", req.UserID).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusOK, dtos.BaseResponse{
				Code:    1,
				Message: "User not found",
				Error: &dtos.ErrorResponse{
					ErrorDetails: "invalid id",
				},
			})
			return
		}
		c.JSON(http.StatusInternalServerError, dtos.BaseResponse{
			Message: InternalServerError,
			Error: &dtos.ErrorResponse{
				ErrorDetails: err.Error(),
			},
		})
		return
	}
	var isActive bool
	if req.Status == 1 {
		isActive = true
	} else {
		isActive = false
	}

	err = h.db.Model(&user).Where("id = ?", user.ID).UpdateColumn("active", isActive).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, dtos.BaseResponse{
			Message: InternalServerError,
			Error: &dtos.ErrorResponse{
				ErrorDetails: err.Error(),
			},
		})
		return
	}

	var message string
	if req.Status == 1 {
		message = "Activate account success"
	} else {
		go func(userID int) {
			err = h.db.Where("id = ?", userID).Delete(&entities.UserToken{}).Error
			if err != nil {
				return
			}
		}(user.ID)
		message = "Deactivate account success"
	}
	c.JSON(http.StatusOK, dtos.BaseResponse{
		Code:    0,
		Message: message,
	})
}

func (h *userHandler) ListUserPaginate(c *gin.Context) {
	pageData := make(map[string]int)
	conditions := make(map[string]interface{})
	pageQuery, ok := c.GetQuery("page")
	if ok {
		page, err := strconv.Atoi(pageQuery)
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
		pageData["page"] = page
	} else {
		pageData["page"] = 1
	}

	perPageQuery, ok := c.GetQuery("per_page")
	if ok {
		perPage, err := strconv.Atoi(perPageQuery)
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
		pageData["per_page"] = perPage
	} else {
		pageData["per_page"] = 20
	}

	if status, ok := c.GetQuery("status"); ok {
		isActive, err := strconv.Atoi(status)
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
		if isActive == 1 {
			conditions["active"] = true
		} else if isActive == 2 {
			conditions["active"] = false
		}
	}

	var count int64
	err := h.db.Model(&entities.User{}).Where(conditions).Count(&count).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, dtos.BaseResponse{
			Message: InternalServerError,
			Error: &dtos.ErrorResponse{
				ErrorDetails: err,
			},
		})
		return
	}

	var users []entities.User
	err = h.db.Scopes(database.Pagination(pageData)).Where(conditions).Order("updated_at DESC").Find(&users).Error
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
			"users":        users,
			"page":         pageData["page"],
			"per_page":     pageData["per_page"],
			"total_record": count,
			"total_page":   utils.CalcTotalPage(count, pageData["per_page"]),
		},
	})
}
