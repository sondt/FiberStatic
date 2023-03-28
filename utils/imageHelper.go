package imageHelper

import (
	"github.com/disintegration/imaging"
	"github.com/gofiber/fiber/v2"
	"os"
	"strconv"
	"strings"
)

func ExtractSize(name string) (int, int) {
	//example: 400x300, w300, h300
	width := 0
	height := 0
	if strings.Contains(name, "x") {
		arr := strings.Split(name, "x")
		width, _ = strconv.Atoi(arr[0])
		height, _ = strconv.Atoi(arr[1])
	}

	if strings.Contains(name, "w") {
		width, _ = strconv.Atoi(strings.Replace(name, "w", "", -1))
	}

	if strings.Contains(name, "h") {
		height, _ = strconv.Atoi(strings.Replace(name, "h", "", -1))
	}

	return width, height
}

func CreateFolder(path string) error {
	err := os.MkdirAll(path, os.ModePerm)
	return err
}

func Resize(originalImage string, width, height int, distFolder string, imageName string) (string, error) {

	newFile := distFolder + "/" + imageName
	//check if file exists
	if _, err := os.Stat(newFile); err == nil {
		return newFile, nil
	}

	src, err := imaging.Open(originalImage)
	if err != nil {
		return "", err
	}

	if width > 0 && height > 0 {
		dst := imaging.Fit(src, width, height, imaging.Lanczos)
		err = imaging.Save(dst, newFile)
	} else if (width > 0 && height == 0) || (width == 0 && height > 0) {
		dst := imaging.Resize(src, width, height, imaging.Lanczos)
		err = imaging.Save(dst, newFile)
	}

	if err != nil {
		return "", err
	}

	return newFile, err
}

var publicFolder = "public"
var resizedFolder = "resized"

// ProcessRewriteImage process rewrite image
func ProcessRewriteImage(c *fiber.Ctx) error {
	resizeType := c.Params("type")
	resizePath := strings.ToLower(c.Path())

	arr := strings.Split(resizePath, "/")
	folders := arr[0 : len(arr)-1]
	resizePath = resizedFolder + strings.Join(folders, "/")
	fileName := arr[len(arr)-1]
	if strings.Contains(fileName, "?") {
		arr = strings.Split(fileName, "?")
		fileName = arr[0]
	}

	err := CreateFolder(resizePath)
	if err != nil {
		return c.SendStatus(500)
	}
	width, height := ExtractSize(resizeType)
	originalPath := strings.Replace(c.Path(), resizeType, publicFolder, 1)
	newFile, err := Resize("."+originalPath, width, height, resizePath, fileName)
	if err != nil {
		return c.SendStatus(500)
	}

	return c.SendFile(newFile)
}
