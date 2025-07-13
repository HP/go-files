package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/bootdotdev/learn-file-storage-s3-golang-starter/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerUploadVideo(w http.ResponseWriter, r *http.Request) {
	videoIDString := r.PathValue("videoID")
	videoID, err := uuid.Parse(videoIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid ID", err)
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find JWT", err)
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT", err)
		return
	}

	video, err := cfg.db.GetVideo(videoID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't find video", err)
		return
	}
	if video.UserID != userID {
		respondWithError(w, http.StatusForbidden, "You are not authorized to upload a video file for this video", nil)
		return
	}

	fmt.Println("uploading video", videoID, "by user", userID)

	const maxMemory = 1 * 1024 * 1024 * 1024 // 1gb
	err = r.ParseMultipartForm(maxMemory)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't parse multipart form", err)
		return
	}

	formFile, fileHeader, err := r.FormFile("video")
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't get video file", err)
		return
	}
	defer formFile.Close()
	mediaType := fileHeader.Header.Get("Content-Type")

	if mediaType != "video/mp4" {
		respondWithError(w, http.StatusBadRequest, "Unsupported media type", nil)
		return
	}

	tempFile, err := os.CreateTemp(os.TempDir(), "video-*.mp4")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create temp file", err)
		return
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()
	_, err = io.Copy(tempFile, formFile)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't copy video file", err)
		return
	}
	_, err = tempFile.Seek(0, io.SeekStart)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't seek to start of temp file", err)
		return
	}

	processingFileName, err := processVideoForFastStart(tempFile.Name())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't process video file", err)
		return
	}
	processingFile, err := os.Open(processingFileName)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't open processing file", err)
		return
	}
	defer processingFile.Close()
	defer os.Remove(processingFile.Name())

	randomBytes := make([]byte, 32)
	_, err = rand.Read(randomBytes)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't generate random string", err)
		return
	}

	aspectRatio, err := getVideoAspectRatio(tempFile.Name())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get video aspect ratio", err)
		return
	}

	var folderName string
	switch aspectRatio {
	case "16:9":
		folderName = "landscape"
	case "9:16":
		folderName = "portrait"
	default:
		folderName = "other"
	}

	fileName := folderName + "/" + hex.EncodeToString(randomBytes) + ".mp4"

	_, err = cfg.s3Client.PutObject(context.Background(), &s3.PutObjectInput{
		Bucket:      aws.String(cfg.s3Bucket),
		Key:         aws.String(fileName),
		ContentType: aws.String(mediaType),
		Body:        processingFile,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't upload video file", err)
		return
	}
	videoUrl := fmt.Sprintf("%s,%s", cfg.s3Bucket, fileName)

	nextVideo := video
	nextVideo.VideoURL = &videoUrl

	err = cfg.db.UpdateVideo(nextVideo)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't update video", err)
		return
	}

	signedVideo, err := cfg.dbVideoToSignedVideo(nextVideo)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't sign video URL", err)
		return
	}
	respondWithJSON(w, http.StatusOK, signedVideo)

}
