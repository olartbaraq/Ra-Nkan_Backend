package api

import (
	"context"
	"log"
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

func ImageUploadHelper(filename interface{}) (string, error) {

	otherConfig, err := utils.LoadOtherConfig("..")
	if err != nil {
		log.Fatal("Could not load env config in cloudinary", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	//create cloudinary instance
	cldInstance, err := cloudinary.NewFromParams(otherConfig.CloudName, otherConfig.CloudApiKey, otherConfig.CloudApiSecret)
	log.Println("are we here?", otherConfig.CloudName, otherConfig.CloudApiKey, otherConfig.CloudApiSecret)
	if err != nil {
		//println("Error creating cloudinary instance", err)
		return "", err
	}

	//upload file
	uploadParam, err := cldInstance.Upload.Upload(ctx, filename, uploader.UploadParams{Folder: otherConfig.CloudUploadFolder})
	if err != nil {
		return "", err
	}
	return uploadParam.SecureURL, nil
}
