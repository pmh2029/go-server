package interfaces

import "context"

type UploadUsecase interface {
	FileUpload(
		ctx context.Context,
		file interface{},
	) (string, error)
}
