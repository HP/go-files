package main

import (
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"

	"github.com/bootdotdev/learn-file-storage-s3-golang-starter/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerUploadThumbnail(w http.ResponseWriter, r *http.Request) {
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

	fmt.Println("uploading thumbnail for video", videoID, "by user", userID)

	// TODO: implement the upload here
	const maxMemory = 10 * 1024 * 1024    // 10mb
	err = r.ParseMultipartForm(maxMemory) // 10mb
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't parse multipart form", err)
		return
	}

	formFile, fileHeader, err := r.FormFile("thumbnail")
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't get thumbnail file", err)
		return
	}
	defer formFile.Close()
	mediaType := fileHeader.Header.Get("Content-Type")

	fileExtensions := map[string]string{
		"image/jpeg": ".jpg",
		"image/png":  ".png",
	}

	mimeType, _, err := mime.ParseMediaType(mediaType)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Unsupported media type", err)
		return
	}
	if mimeType != "image/jpeg" && mimeType != "image/png" && mimeType != "image/gif" && mimeType != "image/webp" && mimeType != "image/svg+xml" {
		respondWithError(w, http.StatusBadRequest, "Unsupported media type", nil)
		return
	}

	fileExtension, ok := fileExtensions[mimeType]
	if !ok {
		respondWithError(w, http.StatusBadRequest, "Unsupported media type", nil)
		return
	}
	filePath := cfg.assetsRoot + "/" + videoID.String() + fileExtension
	file, err := os.Create(filePath)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create file", err)
		return
	}
	defer file.Close()

	_, err = io.Copy(file, formFile)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't copy file", err)
		return
	}

	video, err := cfg.db.GetVideo(videoID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't find video", err)
		return
	}
	if video.UserID != userID {
		respondWithError(w, http.StatusForbidden, "You are not authorized to upload a thumbnail for this video", nil)
		return
	}

	cleanFilePath := filePath
	if len(cleanFilePath) > 0 && cleanFilePath[0] == '.' {
		cleanFilePath = cleanFilePath[1:]
	}
	thumbnailURL := "http://localhost:" + cfg.port + cleanFilePath

	nextVideo := video
	nextVideo.ThumbnailURL = &thumbnailURL

	err = cfg.db.UpdateVideo(nextVideo)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't update video", err)
		return
	}

	respondWithJSON(w, http.StatusOK, nextVideo)
}
