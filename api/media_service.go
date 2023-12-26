package api

import (
	"log"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type mediaUpload interface {
	FileUpload(file File) (string, error)
	RemoteUpload(url Url) (string, error)
}

type media struct{}

func NewMediaUpload() mediaUpload {
	return &media{}
}

func (*media) FileUpload(file File) (string, error) {
	//validate
	if V, ok := binding.Validator.Engine().(*validator.Validate); ok {
		//print("When trying to validate", ok)
		V.Struct(file)
	}

	//upload
	uploadUrl, err := ImageUploadHelper(file.File)
	if err != nil {
		log.Println("ImageFileHelper error:", err)
		return "", err
	}
	return uploadUrl, nil
}

func (*media) RemoteUpload(url Url) (string, error) {

	//upload
	uploadUrl, errUrl := ImageUploadHelper(url.Url)
	if errUrl != nil {
		return "", errUrl
	}
	return uploadUrl, nil
}
