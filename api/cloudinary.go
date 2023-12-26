package api

import (
	"context"
	"time"

	"github.com/cloudinary/cloudinary-go"
	"github.com/cloudinary/cloudinary-go/api/uploader"
	"github.com/olartbaraq/spectrumshelf/utils"
)

// type CloudinaryValues struct {
// 	config *utils.Config
// }

// func NewCloudinaryValues(config *utils.Config) *CloudinaryValues {
// 	return &CloudinaryValues{
// 		config: config,
// 	}
// }

var config *utils.Config

func ImageUploadHelper(filename interface{}) (string, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	//create cloudinary instance
	cldInstance, err := cloudinary.NewFromParams(config.CloudName, config.CloudinaryApiKey, config.CloudinaryApiSecret)
	//log.Println("are we here?", config.CloudName, config.CloudinaryApiKey, config.CloudinaryApiSecret)
	if err != nil {
		println("Error creating cloudinary instance", err)
		return "", err
	}

	//upload file
	uploadParam, err := cldInstance.Upload.Upload(ctx, filename, uploader.UploadParams{Folder: config.CloudinaryFolder})
	if err != nil {
		return "", err
	}
	return uploadParam.SecureURL, nil
}
