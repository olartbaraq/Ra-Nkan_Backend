package api

import (
	"context"
	"time"

	"github.com/cloudinary/cloudinary-go"
	"github.com/cloudinary/cloudinary-go/api/uploader"
	config "github.com/olartbaraq/spectrumshelf/configs"
)

// type CloudinaryValues struct {
// 	config *utils.Config
// }

// func NewCloudinaryValues(config *utils.Config) *CloudinaryValues {
// 	return &CloudinaryValues{
// 		config: config,
// 	}
// }

func ImageUploadHelper(filename interface{}) (string, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	//create cloudinary instance
	cldInstance, err := cloudinary.NewFromParams(config.EnvCloudName(), config.EnvCloudAPIKey(), config.EnvCloudAPISecret())
	//log.Println("are we here?", config.EnvCloudName(), config.EnvCloudAPIKey(), config.EnvCloudAPISecret())
	if err != nil {
		//println("Error creating cloudinary instance", err)
		return "", err
	}

	//upload file
	uploadParam, err := cldInstance.Upload.Upload(ctx, filename, uploader.UploadParams{Folder: config.EnvCloudUploadFolder()})
	if err != nil {
		return "", err
	}
	return uploadParam.SecureURL, nil
}
