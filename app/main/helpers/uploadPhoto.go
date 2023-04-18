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

func UploadPhoto(c echo.Context, formUsername string, updata ...bool) (newPhoto *multipart.FileHeader, err error) {
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

	if newPhoto == nil {
		return nil, nil
	}

	if len(updata) == 0 {
		oldFilename := fmt.Sprint(formUsername, rand.Intn(1000), newPhoto.Filename)
		newPhoto.Filename = fmt.Sprint("members/", oldFilename)
	} else if len(updata) == 1 {
		newPhoto.Filename = fmt.Sprint("members/", newPhoto.Filename)
	} else { // too agrs [0,1,..]=bool
		return nil, c.JSON(http.StatusBadRequest, echo.Map{
			"message": "Upload Photo too agrs [0]=bool",
		})
	}

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
