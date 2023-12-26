package api

import (
	"context"
	"time"

	"github.com/cloudinary/cloudinary-go"
	"github.com/cloudinary/cloudinary-go/api/uploader"
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
	// cldInstance, err := cloudinary.NewFromParams(config.CloudName, config.CloudinaryApiKey, config.CloudinaryApiSecret)
	// log.Println("are we here?", config.CloudName, config.CloudinaryApiKey, config.CloudinaryApiSecret)
	// if err != nil {
	// 	println("Error creating cloudinary instance", err)
	// 	return "", err
	// }

	cldInstance, err := cloudinary.NewFromParams(
		"dxijsd5cc",
		"222494914895738",
		"k0ocMFYnHUkqfsq4YZaLivBdqKk")
	if err != nil {
		return "", err
	}

	//upload file
	uploadParam, err := cldInstance.Upload.Upload(ctx, filename, uploader.UploadParams{Folder: "ra_nkan-cloudinary"})
	if err != nil {
		return "", err
	}
	return uploadParam.SecureURL, nil
}
