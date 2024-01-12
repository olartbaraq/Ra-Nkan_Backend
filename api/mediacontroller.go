package api

// import (
// 	"log"
// 	"mime/multipart"
// 	"net/http"
// 	"net/url"
// 	"strings"

// 	"github.com/gin-gonic/gin"
// )

// func FileUpload() gin.HandlerFunc {
// 	return func(ctx *gin.Context) {
// 		form, err := ctx.MultipartForm()

// 		if err != nil {
// 			ctx.JSON(http.StatusInternalServerError, gin.H{
// 				"statusCode": http.StatusInternalServerError,
// 				"Error":      err.Error(),
// 				"data":       "Select a file to upload",
// 			})
// 			return
// 		}

// 		files := form.File["files[]"]
// 		uploadResults := make(chan UploadResult, len(files))

// 		for _, file := range files {
// 			go func(imageFile *multipart.FileHeader, results chan<- UploadResult) {
// 				log.Println("WE GOT TO FILE UPLOAD HERE >>>", imageFile.Filename)

// 				fileContent, err := imageFile.Open()
// 				if err != nil {
// 					results <- UploadResult{Err: err}
// 					return
// 				}
// 				defer fileContent.Close()

// 				uploadUrl, err := NewMediaUpload().FileUpload(File{File: fileContent})
// 				results <- UploadResult{Url: uploadUrl, Err: err}
// 			}(file, uploadResults)
// 		}

// 		var uploadUrls []string
// 		var uploadErrors []error

// 		for i := 0; i < len(files); i++ {
// 			result := <-uploadResults
// 			if result.Err != nil {
// 				uploadErrors = append(uploadErrors, result.Err)
// 			} else {
// 				uploadUrls = append(uploadUrls, result.Url)
// 			}
// 		}

// 		close(uploadResults)

// 		if len(uploadErrors) > 0 {
// 			ctx.JSON(http.StatusInternalServerError, gin.H{
// 				"statusCode": http.StatusInternalServerError,
// 				"Error":      "Error uploading files",
// 				"data":       uploadErrors,
// 			})
// 			return
// 		}

// 		ctx.JSON(
// 			http.StatusOK,
// 			gin.H{
// 				"statusCode": http.StatusOK,
// 				"message":    "files uploaded successfully",
// 				"data":       uploadUrls,
// 			},
// 		)
// 	}
// }

// // func RemoteUpload() gin.HandlerFunc {
// // 	return func(ctx *gin.Context) {

// // 		var url UrlImages

// // 		if err := ctx.ShouldBindJSON(&url); err != nil {
// // 			//fmt.Println(err.Error())
// // 			stringErr := string(err.Error())
// // 			//fmt.Println(stringErr)
// // 			if strings.Contains(stringErr, "isImageURL") {
// // 				ctx.JSON(http.StatusBadRequest, gin.H{
// // 					"Error": "image URL doesn't point to a real image",
// // 				})
// // 				return

// // 			}
// // 			ctx.JSON(http.StatusBadRequest, gin.H{
// // 				"Error": err.Error(),
// // 			})
// // 			return
// // 		}

// // 		uploadResults := make(chan UploadResult, len(url.ImageUrls))

// // 		for _, imageUrl := range url.ImageUrls {
// // 			go func(imageUrl string, results chan<- UploadResult) {
// // 				log.Println("WE GOT to REMOTE UPLOAD HERE >>>", imageUrl)

// // 				parsedUrl, err := NewUrlFromString(imageUrl)
// // 				if err != nil {
// // 					log.Println("WE COULDN'T PARSE STRING TO URL >>>", err)
// // 				}

// // 				uploadUrl, err := NewMediaUpload().RemoteUpload(*parsedUrl)
// // 				results <- UploadResult{Url: uploadUrl, Err: err}
// // 			}(imageUrl, uploadResults)
// // 		}

// // 		var uploadUrls []string
// // 		var uploadErrors []error

// // 		for i := 0; i < len(url.ImageUrls); i++ {
// // 			result := <-uploadResults
// // 			if result.Err != nil {
// // 				uploadErrors = append(uploadErrors, result.Err)
// // 			} else {
// // 				uploadUrls = append(uploadUrls, result.Url)
// // 			}
// // 		}

// // 		close(uploadResults)

// // 		if len(uploadErrors) > 0 {
// // 			ctx.JSON(http.StatusInternalServerError, gin.H{
// // 				"statusCode": http.StatusInternalServerError,
// // 				"Error":      "Error uploading files",
// // 				"data":       uploadErrors,
// // 			})
// // 			return
// // 		}

// // 		ctx.JSON(
// // 			http.StatusOK,
// // 			gin.H{
// // 				"statusCode": http.StatusOK,
// // 				"message":    "files uploaded successfully",
// // 				"data":       uploadUrls,
// // 			},
// // 		)

// // 	}
// // }
