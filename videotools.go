package main

import (
	"encoding/json"
	"fmt"
	"math"
	"os/exec"
)

func getVideoAspectRatio(filePath string) (string, error) {
	cmd := exec.Command("ffprobe", "-v", "error", "-print_format", "json", "-show_streams", filePath)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	var result map[string]interface{}
	err = json.Unmarshal(output, &result)
	if err != nil {
		return "", err
	}

	streams := result["streams"].([]interface{})

	for _, stream := range streams {
		streamMap := stream.(map[string]interface{})
		width := streamMap["width"].(float64)
		height := streamMap["height"].(float64)
		aspectRatio := width / height
		const tolerance = 0.03
		if math.Abs(aspectRatio-16.0/9.0) < tolerance {
			return "16:9", nil
		} else if math.Abs(aspectRatio-9.0/16.0) < tolerance {
			return "9:16", nil
		} else {
			return "other", nil
		}
	}
	return "other", nil
}

func processVideoForFastStart(filePath string) (string, error) {
	outputPath := filePath + ".processing"
	cmd := exec.Command("ffmpeg", "-i", filePath, "-c", "copy", "-movflags", "faststart", "-f", "mp4", outputPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("ffmpeg error: %v, output: %s", err, string(output))
	}
	return outputPath, nil
}
