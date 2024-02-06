package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"

	"github.com/nfnt/resize"
)

var (
	version = "unknown"
	commit  = "unknown"
	date    = "unknown"
)

func main() {
	// Define a boolean flag "-v" for version
	showVersion := flag.Bool("v", false, "Prints the version")

	// Parse command line arguments
	flag.Parse()

	// If -v flag is provided, print version and exit
	if *showVersion {
		fmt.Printf("Version: %s\n", version)
		fmt.Printf("Git commit hash: %s\n", commit)
		fmt.Printf("Built at: %s\n", date)
		return
	}

	if len(os.Args) < 3 {
		fmt.Println("Aborting. Needs two arguments: type of operation and input file.")
		os.Exit(1)
	}

	operationType := os.Args[1]
	src := os.Args[2]

	if _, err := os.Stat(src); os.IsNotExist(err) {
		fmt.Printf("Error: File '%s' not found.\n", src)
		os.Exit(1)
	}

	// check if src is folder
	isFolder := false
	if info, err := os.Stat(src); err == nil && info.IsDir() {
		isFolder = true
	}

	if !isFolder {
		fmt.Printf("Processing file: %s\n", src)
		runImageOperation(operationType, src)
	}

	if isFolder {
		fmt.Printf("Processing folder: %s\n", src)

		files, err := filepath.Glob(filepath.Join(src, "*"))
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			os.Exit(1)
		}

		// print list of all files in folder
		for _, file := range files {
			fmt.Printf("Found file: %s\n", file)
			if !isFileImage(file) {
				fmt.Printf("Skipping non-image file: %s\n", file)
				continue
			}
			if isFileAlreadyProcessed(file) {
				fmt.Printf("Skipping already processed file: %s\n", file)
				continue
			}
			fmt.Printf("Processing file: %s\n", file)
			runImageOperation(operationType, file)
		}
	}
}

func runImageOperation(operationType string, imageFilePath string) (image.Image, error) {
	img, err := loadImage(imageFilePath)
	if err != nil {
		fmt.Printf("Error opening image: %s\n", err)
		return nil, err
	}

	switch operationType {
	case "grayscale":
		grayImg := toGrayscale(img)

		outputPath := getOutputPath(imageFilePath)
		err = saveImage(outputPath, grayImg)
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
		resizedImg := resizeImage(img, width)

		outputPath := getOutputPath(imageFilePath)
		err = saveImage(outputPath, resizedImg)
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

func isFileImage(filePath string) bool {
	ext := strings.ToLower(filepath.Ext(filePath))
	return ext == ".jpg" || ext == ".jpeg" || ext == ".png"
}

func isFileAlreadyProcessed(filePath string) bool {
	// check if file has "_output" in its name
	return strings.Contains(filepath.Base(filePath), "_output")
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
	fileExt := filepath.Ext(fileName)

	outputFileName := fileNameWithoutExt + "_output" + fileExt

	outputPath := filepath.Join(filepath.Dir(inputPath), outputFileName)

	return outputPath
}
