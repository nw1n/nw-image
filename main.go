package main

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"os"
	"path/filepath"
	"strings"

	"github.com/nfnt/resize"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Aborting. Needs two arguemnts: type of operation and input file.")
		os.Exit(1)
	}

	operationType := os.Args[1]
	imageFilePath := os.Args[2]

	if _, err := os.Stat(imageFilePath); os.IsNotExist(err) {
		fmt.Printf("Error: File '%s' not found.\n", imageFilePath)
		os.Exit(1)
	}

	img, err := loadImage(imageFilePath)
	if err != nil {
		fmt.Printf("Error opening image: %s\n", err)
		os.Exit(1)
	}

	switch operationType {
	case "grayscale":
		grayImg := toGrayscale(img)

		outputPath := getOutputPath(imageFilePath)
		err = saveImage(outputPath, grayImg)
		if err != nil {
			fmt.Printf("Error saving grayscale image: %s\n", err)
			os.Exit(1)
		}

		fmt.Printf("Grayscale image saved at: %s\n", outputPath)
		os.Exit(0)

	case "resize":
		width := 480
		if len(os.Args) > 3 {
			width = 0
			fmt.Sscanf(os.Args[3], "%d", &width)

			if width == 0 {
				fmt.Printf("Error: Invalid width: %s\n", os.Args[3])
				os.Exit(1)
			}
		}
		resizedImg := resizeImage(img, width)

		outputPath := getOutputPath(imageFilePath)
		err = saveImage(outputPath, resizedImg)
		if err != nil {
			fmt.Printf("Error saving resized image: %s\n", err)
			os.Exit(1)
		}

		fmt.Printf("Resized image saved at: %s\n", outputPath)
		os.Exit(0)

	default:
		fmt.Printf("Error: Unknown operation type: %s\n", operationType)
		os.Exit(1)
	}
}

func resizeImage(img image.Image, width int) image.Image {
	bounds := img.Bounds()
	aspectRatio := float64(bounds.Dy()) / float64(bounds.Dx())
	newHeight := int(float64(width) * aspectRatio)

	return resize.Resize(uint(width), uint(newHeight), img, resize.Lanczos3)
}

func loadImage(filePath string) (image.Image, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	return img, err
}

func saveImage(filePath string, img image.Image) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	return jpeg.Encode(file, img, nil)
}

func toGrayscale(img image.Image) image.Image {
	bounds := img.Bounds()
	gray := image.NewGray(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			grayColor := color.GrayModel.Convert(img.At(x, y))
			gray.Set(x, y, grayColor)
		}
	}

	return gray
}

func getOutputPath(inputPath string) string {
	fileName := filepath.Base(inputPath)
	fileNameWithoutExt := strings.TrimSuffix(fileName, filepath.Ext(fileName))

	outputFileName := fmt.Sprintf("%s_transformed.jpg", fileNameWithoutExt)

	outputPath := filepath.Join(filepath.Dir(inputPath), outputFileName)

	return outputPath
}
