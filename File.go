package main

import (
	"path/filepath"
	"strings"
)

func IsFileImage(filePath string) bool {
	ext := strings.ToLower(filepath.Ext(filePath))
	return ext == ".jpg" || ext == ".jpeg" || ext == ".png"
}

func IsFileAlreadyProcessed(filePath string) bool {
	return strings.Contains(filepath.Base(filePath), "_output")
}

func GetOutputPath(inputPath string) string {
	fileName := filepath.Base(inputPath)
	fileNameWithoutExt := strings.TrimSuffix(fileName, filepath.Ext(fileName))
	fileExt := filepath.Ext(fileName)

	outputFileName := fileNameWithoutExt + "_output" + fileExt

	outputPath := filepath.Join(filepath.Dir(inputPath), outputFileName)

	return outputPath
}
