package usecases

import (
	"context"
	"go-server/internal/pkg/domains/interfaces"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/sirupsen/logrus"
)

type uploadUsecase struct {
	cld    *cloudinary.Cloudinary
	logger *logrus.Logger
}

func NewUploadUsecase(
	cld *cloudinary.Cloudinary,
	logger *logrus.Logger,
) interfaces.UploadUsecase {
	return &uploadUsecase{
		cld:    cld,
		logger: logger,
	}
}

func (u *uploadUsecase) FileUpload(
	ctx context.Context,
	file interface{},
) (string, error) {
	uploadParam, err := u.cld.Upload.Upload(ctx, file, uploader.UploadParams{
		Folder: os.Getenv("CLOUDINARY_UPLOAD_FOLDER"),
	})

	if err != nil {
		return "", err
	}

	return uploadParam.URL, nil
}
