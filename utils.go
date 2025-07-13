package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/bootdotdev/learn-file-storage-s3-golang-starter/internal/database"
)

func generatePresignedURL(s3Client *s3.Client, bucket, key string, expireTime time.Duration) (string, error) {
	presignClient := s3.NewPresignClient(s3Client)
	presignObject, err := presignClient.PresignGetObject(context.Background(), &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}, s3.WithPresignExpires(expireTime))
	if err != nil {
		return "", err
	}
	return presignObject.URL, nil
}

func (cfg *apiConfig) dbVideoToSignedVideo(video database.Video) (database.Video, error) {
	if video.VideoURL == nil {
		return video, nil
	}

	videoParts := strings.Split(*video.VideoURL, ",")
	if len(videoParts) != 2 {
		return video, fmt.Errorf("invalid video URL")
	}
	bucket := videoParts[0]
	videoUrl := videoParts[1]

	presignedUrl, err := generatePresignedURL(cfg.s3Client, bucket, videoUrl, 1*time.Hour)
	if err != nil {
		return video, err
	}
	video.VideoURL = &presignedUrl

	return video, nil
}
