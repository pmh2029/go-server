package handlers

import (
	"encoding/json"
	"errors"
	"go-server/internal/pkg/domains/models/dtos"
	"go-server/internal/pkg/domains/models/entities"
	"go-server/pkg/shared/database"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type tripHandler struct {
	db     *gorm.DB
	logger *logrus.Logger
}

func NewTripHandler(
	logger *logrus.Logger,
	db *gorm.DB,
) *tripHandler {
	return &tripHandler{
		db,
		logger,
	}
}

func (h *tripHandler) CreateTrip(c *gin.Context) {
	req := dtos.CreateTripRequestDto{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusOK, dtos.BaseResponse{
			Code:    400,
			Message: "Bad Request",
			Error: &dtos.ErrorResponse{
				ErrorDetails: err.Error(),
			},
		})
		return
	}

	var users []entities.User
	if len(req.Users) > 0 {
		err = h.db.Where("id IN (?) AND active = ?", req.Users, true).Find(&users).Error
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
		if len(req.Users) != len(users) {
			c.JSON(http.StatusOK, dtos.BaseResponse{
				Code:    1,
				Message: "User not found",
				Error: &dtos.ErrorResponse{
					ErrorDetails: err,
				},
			})
			return
		}
	}

	if len(req.Days) == 0 {
		c.JSON(http.StatusOK, dtos.BaseResponse{
			Code:    2,
			Message: "Days must be selected",
			Error: &dtos.ErrorResponse{
				ErrorDetails: err,
			},
		})
		return
	}

	var days []entities.Day
	for _, day := range req.Days {
		if len(day.Places) == 0 {
			c.JSON(http.StatusOK, dtos.BaseResponse{
				Code:    3,
				Message: "Places must be selected",
				Error: &dtos.ErrorResponse{
					ErrorDetails: err,
				},
			})
			return
		}

		for _, place := range day.Places {
			err = h.db.Where("id = ?", place.ID).First(&entities.Place{}).Error
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					c.JSON(http.StatusOK, dtos.BaseResponse{
						Code:    4,
						Message: "Place not found",
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
		}

		places, err := json.Marshal(&day.Places)
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

		newDay := entities.Day{
			Places: string(places),
		}

		days = append(days, newDay)
	}

	err = database.Transaction(c, h.db, func(tx *gorm.DB) error {
		trip := entities.Trip{
			Owner:    req.Owner,
			Name:     req.Name,
			FromDate: time.Unix(int64(req.FromDate), 0),
			ToDate:   time.Unix(int64(req.ToDate), 0),
		}
		err = tx.Create(&trip).Error
		if err != nil {
			return err
		}

		for i := range days {
			days[i].TripID = trip.ID
		}

		var userTrips []entities.UserTrip
		var strUserIDs []string
		if len(req.Users) > 0 {
			for _, userID := range req.Users {
				userTrips = append(userTrips, entities.UserTrip{
					UserID: userID,
					TripID: trip.ID,
				})
				strUserIDs = append(strUserIDs, strconv.Itoa(userID))
			}

			trip.UserIDs = "," + strings.Join(strUserIDs, ",") + ","
			err = tx.Model(&trip).Updates(trip).Error
			if err != nil {
				return err
			}
		}

		err = tx.CreateInBatches(&userTrips, len(req.Users)).Error
		if err != nil {
			return err
		}

		err = tx.Create(&days).Error
		if err != nil {
			return err
		}

		return nil
	})
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
		Message: "Created success",
	})
}

func (h *tripHandler) ListTrip(c *gin.Context) {
	userIDParam := c.Param("user_id")
	userID, err := strconv.Atoi(userIDParam)
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

	var trips []entities.Trip
	err = h.db.Preload("Days").Where("owner = ?", userID).Find(&trips).Error
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
		Message: "OK",
		Data: gin.H{
			"trips": trips,
		},
	})
}
