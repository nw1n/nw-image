package main

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"

	"github.com/nfnt/resize"
)

func RunImageOperation(operationType string, imageFilePath string) (image.Image, error) {
	img, err := LoadImage(imageFilePath)
	if err != nil {
		fmt.Printf("Error opening image: %s\n", err)
		return nil, err
	}

	switch operationType {
	case "grayscale":
		grayImg := ToGrayscale(img)

		outputPath := GetOutputPath(imageFilePath)
		err = SaveImage(outputPath, grayImg)
		if err != nil {
			fmt.Printf("Error saving grayscale image: %s\n", err)
			return nil, err
		}

		fmt.Printf("Grayscale image saved at: %s\n", outputPath)
		return grayImg, nil

	case "resize":
		width := 480
		if len(os.Args) > 3 {
			width = 0
			fmt.Sscanf(os.Args[3], "%d", &width)

			if width == 0 {
				fmt.Printf("Error: Invalid width: %s\n", os.Args[3])
				return nil, nil
			}
		}
		resizedImg := ResizeImage(img, width)

		outputPath := GetOutputPath(imageFilePath)
		err = SaveImage(outputPath, resizedImg)
		if err != nil {
			fmt.Printf("Error saving resized image: %s\n", err)
			return nil, err
		}

		fmt.Printf("Resized image saved at: %s\n", outputPath)
		return resizedImg, nil

	default:
		fmt.Printf("Error: Unknown operation type: %s\n", operationType)
		return nil, nil
	}
}

func ResizeImage(img image.Image, width int) image.Image {
	bounds := img.Bounds()
	aspectRatio := float64(bounds.Dy()) / float64(bounds.Dx())
	newHeight := int(float64(width) * aspectRatio)

	return resize.Resize(uint(width), uint(newHeight), img, resize.Lanczos3)
}

func LoadImage(filePath string) (image.Image, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	return img, err
}

func SaveImage(filePath string, img image.Image) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	switch filepath.Ext(filePath) {
	case ".jpg", ".jpeg":
		err = jpeg.Encode(file, img, nil)
	case ".png":
		err = png.Encode(file, img)
	default:
		err = fmt.Errorf("Unsupported file format: %s", filepath.Ext(filePath))
	}

	return err
}

func ToGrayscale(img image.Image) image.Image {
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
