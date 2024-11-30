package storages

import (
	"context"
	"greenenvironment/configs"
	"greenenvironment/constant"
	"greenenvironment/helper"
	"mime/multipart"
	"net/http"
	"strings"
	"time"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/labstack/echo/v4"
)

type StorageInterface interface {
	ImageValidation(files *multipart.FileHeader) (multipart.File, error)
	UploadImageToCloudinary(file interface{}, folderPath string) (string, error)
	UploadFileHandler(c echo.Context) error
}

type storage struct {
	conf configs.CloudinaryConfig
}

type ImageResponse struct {
	ImageUrl string `json:"image_url"`
}

func NewStorage(conf configs.CloudinaryConfig) StorageInterface {
	return &storage{conf: conf}
}

func (s *storage) ImageValidation(file *multipart.FileHeader) (multipart.File, error) {
	var response multipart.File

	if file.Size > 2*1024*1024 {
		return nil, constant.ErrSizeFile
	}
	fileType := file.Header.Get("Content-Type")

	if !strings.HasPrefix(fileType, "image/") {
		return nil, constant.ErrContentTypeFile
	}

	src, _ := file.Open()
	defer src.Close()

	response = src

	return response, nil
}

func (s *storage) UploadImageToCloudinary(file interface{}, folderPath string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cld, err := cloudinary.NewFromParams(s.conf.CloudName, s.conf.ApiKeyStorage, s.conf.ApiSecretStorage)

	if err != nil {
		return "", err
	}

	resp, err := cld.Upload.Upload(ctx, file, uploader.UploadParams{Folder: folderPath})
	if err != nil {
		return "", err
	}

	return resp.SecureURL, nil
}

// Upload Image
// @Summary      Upload Image
// @Description  Upload an image to Cloudinary and return the image URL
// @Tags         Upload
// @Accept       multipart/form-data
// @Produce      json
// @Param        image  formData  file  true  "Image file to upload"
// @Success      201    {object}  helper.Response{data=ImageResponse} "Upload success"
// @Failure      400    {object}  helper.Response{data=string} "Invalid input or validation error"
// @Failure      500    {object}  helper.Response{data=string} "Internal server error"
// @Router       /media/upload [post]
func (s *storage) UploadFileHandler(c echo.Context) error {
	file, err := c.FormFile("image")

	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, "Image is required", nil))
	}

	src, err := s.ImageValidation(file)

	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, err.Error(), nil))
	}

	imageURL, err := s.UploadImageToCloudinary(src, "ecomate/")

	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, "Image upload failed", nil))
	}

	var response = new(ImageResponse)
	response.ImageUrl = imageURL
	return c.JSON(http.StatusCreated, helper.ObjectFormatResponse(true, "Upload success", response))
}
