package usecases

import (
	"context"
	"errors"
	"go-server/internal/pkg/domains/interfaces"
	"go-server/internal/pkg/domains/models/dtos"
	"go-server/internal/pkg/domains/models/entities"
	"go-server/pkg/shared/database"
	"go-server/pkg/shared/utils"
	"strings"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var (
	CreatePlaceCategoriesIsNull   = errors.New("Category must be selected")
	CreatePlaceCategoriesNotFound = errors.New("Category not found")
	CreatePlaceImagesIsNull       = errors.New("Image must be selected")

	UpdatePlaceIDNotFound         = errors.New("Place not found")
	UpdatePlaceCategoriesIsNull   = errors.New("Category must be selected")
	UpdatePlaceCategoriesNotFound = errors.New("Category not found")
	UpdatePlaceImagesIsNull       = errors.New("Image must be selected")

	DetailPlaceIDNotFound = errors.New("Place not found")

	DeletePlaceIDNotFound = errors.New("Place not found")
)

type placeUsecase struct {
	placeRepo         interfaces.PlaceRepository
	placeCategoryRepo interfaces.PlaceCategoryRepository
	categoryRepo      interfaces.CategoryRepository
	logger            *logrus.Logger
}

func NewPlaceUsecase(
	placeRepo interfaces.PlaceRepository,
	placeCategoryRepo interfaces.PlaceCategoryRepository,
	categoryRepo interfaces.CategoryRepository,
	logger *logrus.Logger,
) interfaces.PlaceUsecase {
	return &placeUsecase{
		placeRepo,
		placeCategoryRepo,
		categoryRepo,
		logger,
	}
}

func (u *placeUsecase) Create(
	ctx context.Context,
	db *gorm.DB,
	req dtos.CreatePlaceRequestDto,
) (entities.Place, error, map[string]interface{}) {
	if len(req.Categories) == 0 {
		return entities.Place{}, CreatePlaceCategoriesIsNull, map[string]interface{}{
			"categories": req.Categories,
		}
	}

	categoriesInDb, err := u.categoryRepo.PluckIDByConditions(ctx, map[string]interface{}{
		"id": req.Categories,
	})
	if err != nil {
		return entities.Place{}, err, map[string]interface{}{
			"error": err,
		}
	}

	if len(categoriesInDb) == 0 || len(utils.CheckSliceDiff(req.Categories, categoriesInDb)) > 0 {
		return entities.Place{}, CreatePlaceCategoriesNotFound, map[string]interface{}{
			"categories":       req.Categories,
			"categories_in_db": categoriesInDb,
		}
	}

	if len(req.Images) == 0 {
		return entities.Place{}, CreatePlaceImagesIsNull, map[string]interface{}{
			"images": req.Images,
		}
	}

	place := entities.Place{
		Name:        req.Name,
		Address:     req.Address,
		Latitude:    req.Latitude,
		Longitude:   req.Longitude,
		Description: req.Description,
		Images:      "|" + strings.Join(req.Images, "|") + "|",
		Price:       req.Price,
	}
	if err := database.Transaction(ctx, db, func(tx *gorm.DB) error {
		place, err = u.placeRepo.CreateWithTx(tx, place)
		if err != nil {
			return err
		}

		var placeCategories []entities.PlaceCategory
		for _, category := range req.Categories {
			placeCategories = append(placeCategories, entities.PlaceCategory{
				PlaceID:    place.ID,
				CategoryID: category,
			})
		}

		err = u.placeCategoryRepo.BatchCreateWithTx(tx, placeCategories)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return entities.Place{}, err, map[string]interface{}{
			"error": err,
		}
	}

	return place, nil, map[string]interface{}{}
}

func (u *placeUsecase) FindListPaginate(
	ctx context.Context,
	pageData map[string]int,
	conditions map[string]interface{},
) ([]entities.Place, int64, error) {
	return u.placeRepo.FindListPaginate(ctx, pageData, conditions)
}

func (u *placeUsecase) Update(
	ctx context.Context,
	db *gorm.DB,
	placeID int,
	req dtos.UpdatePlaceRequestDto,
) (entities.Place, error, map[string]interface{}) {
	place, err := u.placeRepo.TakeByConditions(ctx, map[string]interface{}{
		"id": placeID,
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entities.Place{}, UpdatePlaceIDNotFound, map[string]interface{}{}
		}
	}

	if len(req.Categories) == 0 {
		return entities.Place{}, UpdatePlaceCategoriesIsNull, map[string]interface{}{
			"categories": req.Categories,
		}
	}

	categoriesInDb, err := u.categoryRepo.PluckIDByConditions(ctx, map[string]interface{}{
		"id": req.Categories,
	})
	if err != nil {
		return entities.Place{}, err, map[string]interface{}{
			"error": err,
		}
	}

	if len(categoriesInDb) == 0 || len(utils.CheckSliceDiff(req.Categories, categoriesInDb)) > 0 {
		return entities.Place{}, UpdatePlaceCategoriesNotFound, map[string]interface{}{
			"categories":       req.Categories,
			"categories_in_db": categoriesInDb,
		}
	}

	if len(req.Images) == 0 {
		return entities.Place{}, UpdatePlaceImagesIsNull, map[string]interface{}{
			"images": req.Images,
		}
	}

	newPlace := entities.Place{
		Name:        req.Name,
		Address:     req.Address,
		Latitude:    req.Latitude,
		Longitude:   req.Longitude,
		Description: req.Description,
		Images:      "|" + strings.Join(req.Images, "|") + "|",
		Price:       req.Price,
	}
	if err := database.Transaction(ctx, db, func(tx *gorm.DB) error {
		place, err = u.placeRepo.UpdateWithTx(tx, place, newPlace)
		if err != nil {
			return err
		}

		err = u.placeCategoryRepo.DeleteByConditionsWithTx(tx, map[string]interface{}{
			"place_id": placeID,
		})
		if err != nil {
			return err
		}

		var placeCategories []entities.PlaceCategory
		for _, category := range req.Categories {
			placeCategories = append(placeCategories, entities.PlaceCategory{
				PlaceID:    place.ID,
				CategoryID: category,
			})
		}

		err = u.placeCategoryRepo.BatchCreateWithTx(tx, placeCategories)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return entities.Place{}, err, map[string]interface{}{
			"error": err,
		}
	}

	return place, nil, map[string]interface{}{}
}

func (u *placeUsecase) TakeByConditionsWithPreload(
	ctx context.Context,
	conditions map[string]interface{},
) (entities.Place, error) {
	place, err := u.placeRepo.TakeByConditionsWithPreload(ctx, conditions)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entities.Place{}, DetailPlaceIDNotFound
		}
		return entities.Place{}, err
	}

	return place, err
}

func (u *placeUsecase) Delete(
	ctx context.Context,
	db *gorm.DB,
	conditions map[string]interface{},
) error {
	place, err := u.placeRepo.TakeByConditions(ctx, conditions)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return DeletePlaceIDNotFound
		}
		return err
	}

	if err := database.Transaction(ctx, db, func(tx *gorm.DB) error {
		err = u.placeRepo.DeleteByConditionsWithTx(tx, map[string]interface{}{
			"id": place.ID,
		})
		if err != nil {
			return err
		}

		err = u.placeCategoryRepo.DeleteByConditionsWithTx(tx, map[string]interface{}{
			"place_id": place.ID,
		})
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (u *placeUsecase) FindByConditions(
	ctx context.Context,
	conditions map[string]interface{},
) ([]entities.Place, error) {
	return u.placeRepo.FindByConditions(ctx, conditions)
}
