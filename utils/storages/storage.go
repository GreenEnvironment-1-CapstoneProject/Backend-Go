package storages

import (
	"context"
	"greenenvironment/configs"
	"greenenvironment/constant"
	"mime/multipart"
	"strings"
	"time"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

func ImageValidation(files []*multipart.FileHeader) ([]multipart.File, error) {
	var response []multipart.File

	for _, file := range files {
		if file.Size > 2*1024*1024 {
			return nil, constant.ErrSizeFile
		}
		fileType := file.Header.Get("Content-Type")

		if !strings.HasPrefix(fileType, "image/") {
			return nil, constant.ErrContentTypeFile
		}

		src, _ := file.Open()
		defer src.Close()

		response = append(response, src)
	}

	return response, nil
}

func UploadImageToCloudinary(file interface{}, folderPath string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	conf := configs.InitConfig().Cloudinary
	cld, err := cloudinary.NewFromParams(conf.CloudName, conf.ApiKeyStorage, conf.ApiSecretStorage)

	if err != nil {
		return "", err
	}

	resp, err := cld.Upload.Upload(ctx, file, uploader.UploadParams{Folder: folderPath})
	if err != nil {
		return "", err
	}

	return resp.SecureURL, nil
}
