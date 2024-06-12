package handlers

import (
	"encoding/json"
	"errors"
	"go-server/internal/pkg/domains/models/dtos"
	"go-server/internal/pkg/domains/models/entities"
	"go-server/pkg/shared/database"
	"go-server/pkg/shared/utils"
	"net/http"
	"strconv"
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

		for i, place := range day.Places {
			var placeInDB entities.Place
			err = h.db.Where("id = ?", place.ID).First(&placeInDB).Error
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
			day.Places[i].Place = placeInDB
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
			Users:    req.Users,
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

	longitudeQuery, longitudeOk := c.GetQuery("longitude")
	latitudeQuery, latitudeOk := c.GetQuery("latitude")
	if !longitudeOk || !latitudeOk {
		c.JSON(http.StatusOK, dtos.BaseResponse{
			Code:    400,
			Message: InternalServerError,
			Error: &dtos.ErrorResponse{
				ErrorDetails: "Must provide longitude and latitude",
			},
		})
		return
	}
	longitude, err := strconv.ParseFloat(longitudeQuery, 64)
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

	latitude, err := strconv.ParseFloat(latitudeQuery, 64)
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
	err = h.db.Preload("Days").Where("owner = ?", userID).Order("updated_at DESC").Find(&trips).Error
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

	for _, trip := range trips {
		for _, day := range trip.Days {
			for i, place := range day.PlacesJson {
				day.PlacesJson[i].Distance = utils.Haversine(place.Latitude, place.Longitude, latitude, longitude)
			}
		}
	}

	c.JSON(http.StatusOK, dtos.BaseResponse{
		Code:    0,
		Message: "OK",
		Data: gin.H{
			"trips": trips,
		},
	})
}
