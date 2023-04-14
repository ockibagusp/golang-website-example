package helpers

import (
	"fmt"
	"io"
	"math/rand"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
)

func UploadPhoto(c echo.Context, formUsername string) (newPhoto *multipart.FileHeader, err error) {
	// Set max file size
	err = c.Request().ParseMultipartForm(1024)
	if err != nil {
		return nil, c.JSON(http.StatusForbidden, echo.Map{
			"message": "File size is too big",
		})
	}

	// source
	newPhoto, err = c.FormFile("photo")
	if err != nil {
		return nil, c.JSON(http.StatusBadRequest, echo.Map{
			"message": "Bad Request",
		})
	}

	oldFilename := fmt.Sprint(formUsername, rand.Intn(1000), newPhoto.Filename)
	newPhoto.Filename = fmt.Sprint("members/", oldFilename)
	src, err := newPhoto.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	// Destination
	dst, err := os.Create(newPhoto.Filename)
	if err != nil {
		return nil, c.JSON(http.StatusBadRequest, echo.Map{
			"message": "Bad Request",
		})
	}
	defer dst.Close()

	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		return nil, c.JSON(http.StatusBadRequest, echo.Map{
			"message": "Bad Request",
		})
	}

	return
}
