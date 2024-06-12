package handlers

import (
	"errors"
	"go-server/internal/pkg/domains/interfaces"
	"go-server/internal/pkg/domains/models/dtos"
	"go-server/internal/pkg/domains/models/entities"
	"go-server/internal/pkg/repositories"
	"go-server/internal/pkg/usecases"
	"go-server/pkg/shared/database"
	"go-server/pkg/shared/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type placeHandler struct {
	placeUsecase interfaces.PlaceUsecase
	db           *gorm.DB
	logger       *logrus.Logger
}

func NewPlaceHandler(
	logger *logrus.Logger,
	db *gorm.DB,
) *placeHandler {
	placeRepo := repositories.NewPlaceRepository(db, logger)
	categoryRepo := repositories.NewCategoryRepository(db, logger)
	placeCategoryRepo := repositories.NewPlaceCategoryRepository(db, logger)

	placeUsecase := usecases.NewPlaceUsecase(
		placeRepo,
		placeCategoryRepo,
		categoryRepo,
		logger,
	)

	return &placeHandler{
		placeUsecase,
		db,
		logger,
	}
}

func (h *placeHandler) CreatePlace(c *gin.Context) {
	req := dtos.CreatePlaceRequestDto{}
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

	place, err, errDetail := h.placeUsecase.Create(c, h.db, req)
	if err != nil {
		if errors.Is(err, usecases.CreatePlaceCategoriesIsNull) {
			c.JSON(http.StatusOK, dtos.BaseResponse{
				Code:    1,
				Message: err,
				Error: &dtos.ErrorResponse{
					ErrorDetails: errDetail,
				},
			})
			return
		}

		if errors.Is(err, usecases.CreatePlaceCategoriesNotFound) {
			c.JSON(http.StatusOK, dtos.BaseResponse{
				Code:    2,
				Message: err,
				Error: &dtos.ErrorResponse{
					ErrorDetails: errDetail,
				},
			})
			return
		}

		if errors.Is(err, usecases.CreatePlaceImagesIsNull) {
			c.JSON(http.StatusOK, dtos.BaseResponse{
				Code:    3,
				Message: err,
				Error: &dtos.ErrorResponse{
					ErrorDetails: errDetail,
				},
			})
			return
		}

		c.JSON(http.StatusInternalServerError, dtos.BaseResponse{
			Message: InternalServerError,
			Error: &dtos.ErrorResponse{
				ErrorDetails: errDetail,
			},
		})
		return
	}

	c.JSON(http.StatusOK, dtos.BaseResponse{
		Code:    0,
		Message: "Created success",
		Data: gin.H{
			"place": place,
		},
	})
}

func (h *placeHandler) ListPlacePaginate(c *gin.Context) {
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

	if keyword, ok := c.GetQuery("keyword"); ok {
		conditions["keyword"] = keyword
	}

	categoryIDQuery, ok := c.GetQuery("category_id")
	if ok {
		categoryID, err := strconv.Atoi(categoryIDQuery)
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
		conditions["category_id"] = categoryID
	}

	places, total, err := h.placeUsecase.FindListPaginate(c, pageData, conditions)
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
			"places":       places,
			"page":         pageData["page"],
			"per_page":     pageData["per_page"],
			"total_record": total,
			"total_page":   utils.CalcTotalPage(total, pageData["per_page"]),
		},
	})
}

func (h *placeHandler) UpdatePlace(c *gin.Context) {
	placeIDParam := c.Param("place_id")
	placeID, err := strconv.Atoi(placeIDParam)
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

	req := dtos.UpdatePlaceRequestDto{}
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

	place, err, errDetail := h.placeUsecase.Update(c, h.db, placeID, req)
	if err != nil {
		if errors.Is(err, usecases.CreatePlaceCategoriesIsNull) {
			c.JSON(http.StatusOK, dtos.BaseResponse{
				Code:    1,
				Message: err,
				Error: &dtos.ErrorResponse{
					ErrorDetails: errDetail,
				},
			})
			return
		}

		if errors.Is(err, usecases.CreatePlaceCategoriesNotFound) {
			c.JSON(http.StatusOK, dtos.BaseResponse{
				Code:    2,
				Message: err,
				Error: &dtos.ErrorResponse{
					ErrorDetails: errDetail,
				},
			})
			return
		}

		if errors.Is(err, usecases.CreatePlaceImagesIsNull) {
			c.JSON(http.StatusOK, dtos.BaseResponse{
				Code:    3,
				Message: err,
				Error: &dtos.ErrorResponse{
					ErrorDetails: errDetail,
				},
			})
			return
		}

		c.JSON(http.StatusInternalServerError, dtos.BaseResponse{
			Message: InternalServerError,
			Error: &dtos.ErrorResponse{
				ErrorDetails: errDetail,
			},
		})
		return
	}

	c.JSON(http.StatusOK, dtos.BaseResponse{
		Code:    0,
		Message: "Updated success",
		Data: gin.H{
			"place": place,
		},
	})
}

func (h *placeHandler) DetailPlace(c *gin.Context) {
	placeIDParam := c.Param("place_id")
	placeID, err := strconv.Atoi(placeIDParam)
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

	place, err := h.placeUsecase.TakeByConditionsWithPreload(c, map[string]interface{}{
		"id": placeID,
	})
	if err != nil {
		if errors.Is(err, usecases.DetailPlaceIDNotFound) {
			c.JSON(http.StatusOK, dtos.BaseResponse{
				Code:    1,
				Message: err.Error(),
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

	c.JSON(http.StatusOK, dtos.BaseResponse{
		Code:    0,
		Message: "OK",
		Data: gin.H{
			"place": place,
		},
	})
}

func (h *placeHandler) DeletePlace(c *gin.Context) {
	placeIDParam := c.Param("place_id")
	placeID, err := strconv.Atoi(placeIDParam)
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

	err = h.placeUsecase.Delete(c, h.db, map[string]interface{}{
		"id": placeID,
	})
	if err != nil {
		if errors.Is(err, usecases.DeletePlaceIDNotFound) {
			c.JSON(http.StatusOK, dtos.BaseResponse{
				Code:    1,
				Message: err.Error(),
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

	c.JSON(http.StatusOK, dtos.BaseResponse{
		Code:    0,
		Message: "Deleted success",
	})
}

func (h *placeHandler) ListAllPlace(c *gin.Context) {
	conditions := make(map[string]interface{})

	if keyword, ok := c.GetQuery("keyword"); ok {
		conditions["keyword"] = keyword
	}

	categoryIDQuery, ok := c.GetQuery("category_id")
	if ok {
		categoryID, err := strconv.Atoi(categoryIDQuery)
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
		conditions["category_id"] = categoryID
	}

	places, err := h.placeUsecase.FindByConditions(c, conditions)
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
			"places": places,
		},
	})
}

func (h *placeHandler) ListComment(c *gin.Context) {
	placeIDParam := c.Param("place_id")
	placeID, err := strconv.Atoi(placeIDParam)
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

	pageData := make(map[string]int)
	conditions := make(map[string]interface{})
	conditions["place_id"] = placeID
	pageQuery, ok := c.GetQuery("page")
	if ok {
		page, err := strconv.Atoi(pageQuery)
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
				Message: BadRequest,
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

	if rate, ok := c.GetQuery("rate"); ok {
		conditions["rate"] = rate
	}

	orderCondition := "updated_at DESC"
	if order, ok := c.GetQuery("order"); ok {
		orderInt, err := strconv.Atoi(order)
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
		if orderInt == 2 {
			orderCondition = "updated_at ASC"
		}
	}

	var comments []entities.Comment
	var count int64
	err = h.db.Model(&entities.Comment{}).Where(conditions).Count(&count).Error
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

	err = h.db.Scopes(database.Pagination(pageData)).Preload("User").Where(conditions).Order(orderCondition).Find(&comments).Error
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
			"comments":     comments,
			"page":         pageData["page"],
			"per_page":     pageData["per_page"],
			"total_record": count,
			"total_page":   utils.CalcTotalPage(count, pageData["per_page"]),
		},
	})
}

type placeResponse struct {
	entities.Place
	Note      string `json:"note"`
	VisitTime int    `json:"visit_time"`
	StartTime int    `json:"start_time"`
	Vehicle   int    `json:"vehicle"`
}

func (h *placeHandler) ListSuggestPlace(c *gin.Context) {
	longitude, longitudeOk := c.GetQuery("longitude")
	latitude, latitudeOk := c.GetQuery("latitude")
	keyword, keywordOk := c.GetQuery("keyword")

	var places []entities.Place
	if !longitudeOk || !latitudeOk {
		c.JSON(http.StatusOK, dtos.BaseResponse{
			Code:    400,
			Message: InternalServerError,
			Error: &dtos.ErrorResponse{
				ErrorDetails: "Must provide longitude and latitude",
			},
		})
		return

	} else {
		if !keywordOk {
			query := `
			SELECT *,
			(6371 * acos(cos(radians(?)) * cos(radians(latitude)) * cos(radians(longitude) - radians(?)) + sin(radians(?)) * sin(radians(latitude)))) AS distance
			FROM places
			ORDER BY distance ASC, rate DESC;
			`
			err := h.db.Raw(query, latitude, longitude, latitude).Scan(&places).Error
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
		} else {
			query := `
			SELECT *,
			(6371 * acos(cos(radians(?)) * cos(radians(latitude)) * cos(radians(longitude) - radians(?)) + sin(radians(?)) * sin(radians(latitude)))) AS distance
			FROM places
			WHERE LOWER(name) LIKE LOWER(?)
			ORDER BY distance ASC, rate DESC;
			`
			err := h.db.Raw(query, latitude, longitude, latitude, "%"+keyword+"%").Scan(&places).Error
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
		}
	}

	var placeResponses []placeResponse
	for _, place := range places {
		placeResponses = append(placeResponses, placeResponse{
			Place:     place,
			Note:      "",
			VisitTime: 30,
			StartTime: 0,
			Vehicle:   2,
		})
	}

	c.JSON(http.StatusOK, dtos.BaseResponse{
		Code:    0,
		Message: "OK",
		Data: gin.H{
			"places": placeResponses,
		},
	})
}
