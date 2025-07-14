# Files storage and CDN with S3 and CloudFront

This is a [Boot.dev course](https://www.boot.dev/courses/learn-file-servers-s3-cloudfront-golang) project that I have worked through. A really useful course, it walks you through all the necessary AWS steps to setup a file storage and CDN with S3 and CloudFront. When the steps are provided this is quite easy.

I chose to use the Go language, typescript is also available at boot.dev.

What I liked:

- The course is well structured and has good pace.
- AWS steps are well explained and easy to follow. Links to the relevant dashboards are provided.
- Solutions for both public files (CDN) and private files (presigned URLs) are provided.

What I didn't like:

- Not much really. Overall a great course.

How to make it more advanced:

- [Signed urls](https://www.boot.dev/courses/learn-file-servers-s3-cloudfront-golang/lessons/signed-urls) with CloudFront.
- Upload of videos with presigned put requests directly from the client to S3.
- Video processing in the background, or via Lambda functions triggered by S3 events.

My referral link: https://www.boot.dev?bannerlord=hanspetter

## learn-file-storage-s3-golang-starter (Tubely)

This repo contains the starter code for the Tubely application - the #1 tool for engagement bait - for the "Learn File Servers and CDNs with S3 and CloudFront" [course](https://www.boot.dev/courses/learn-file-servers-s3-cloudfront-golang) on [boot.dev](https://www.boot.dev)

## Quickstart

*This is to be used as a *reference\* in case you need it, you should follow the instructions in the course rather than trying to do everything here.

## 1. Install dependencies

- [Go](https://golang.org/doc/install)
- `go mod download` to download all dependencies
- [FFMPEG](https://ffmpeg.org/download.html) - both `ffmpeg` and `ffprobe` are required to be in your `PATH`.

```bash
# linux
sudo apt update
sudo apt install ffmpeg

# mac
brew update
brew install ffmpeg
```

- [SQLite 3](https://www.sqlite.org/download.html) only required for you to manually inspect the database.

```bash
# linux
sudo apt update
sudo apt install sqlite3

# mac
brew update
brew install sqlite3
```

- [AWS CLI](https://docs.aws.amazon.com/cli/latest/userguide/getting-started-install.html)

## 2. Download sample images and videos

```bash
./samplesdownload.sh
# samples/ dir will be created
# with sample images and videos
```

## 3. Configure environment variables

Copy the `.env.example` file to `.env` and fill in the values.

```bash
cp .env.example .env
```

You'll need to update values in the `.env` file to match your configuration, but _you won't need to do anything here until the course tells you to_.

## 3. Run the server

```bash
go run .
```

- You should see a new database file `tubely.db` created in the root directory.
- You should see a new `assets` directory created in the root directory, this is where the images will be stored.
- You should see a link in your console to open the local web page.
