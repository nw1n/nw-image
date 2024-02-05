package main

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	// Check if an image file is provided as an argument
	if len(os.Args) < 3 {
		fmt.Println("Aborting. Needs two arguemnts: type of operation and input file.")
		os.Exit(1)
	}

	// Get the operation type from the command-line arguments
	operationType := os.Args[1]
	imageFilePath := os.Args[2]

	// Check if the file exists
	if _, err := os.Stat(imageFilePath); os.IsNotExist(err) {
		fmt.Printf("Error: File '%s' not found.\n", imageFilePath)
		os.Exit(1)
	}

	// Open the image file
	img, err := loadImage(imageFilePath)
	if err != nil {
		fmt.Printf("Error opening image: %s\n", err)
		os.Exit(1)
	}

	// Perform the operation based on the operation type
	switch operationType {
	case "grayscale":
		// Convert the image to grayscale
		grayImg := toGrayscale(img)

		// Save the grayscale image to a new file
		outputPath := getOutputPath(imageFilePath)
		err = saveImage(outputPath, grayImg)
		if err != nil {
			fmt.Printf("Error saving grayscale image: %s\n", err)
			os.Exit(1)
		}

		fmt.Printf("Grayscale image saved at: %s\n", outputPath)
	default:
		fmt.Printf("Error: Unknown operation type: %s\n", operationType)
		os.Exit(1)
	}
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
	// Extract file name and extension
	fileName := filepath.Base(inputPath)
	fileNameWithoutExt := strings.TrimSuffix(fileName, filepath.Ext(fileName))

	// Create a new file name with "_grayscale" suffix
	outputFileName := fmt.Sprintf("%s_grayscale.jpg", fileNameWithoutExt)

	// Save the new file in the same directory as the original file
	outputPath := filepath.Join(filepath.Dir(inputPath), outputFileName)

	return outputPath
}
