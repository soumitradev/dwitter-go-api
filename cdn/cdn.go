// Package cdn provides some functions to interface with the firebase bucket
package cdn

import (
	"bytes"
	"context"
	"dwitter_go_graphql/auth"
	"dwitter_go_graphql/common"
	"dwitter_go_graphql/util"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"image/png"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"time"

	"cloud.google.com/go/storage"
	firebase "firebase.google.com/go/v4"
	"github.com/disintegration/imaging"
	"github.com/golang/gddo/httputil/header"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

func init() {
	// Connect to bucket
	opt := option.WithCredentialsFile("./cdn_key.json")
	config := &firebase.Config{
		StorageBucket: "dwitter-72e9d.appspot.com",
	}

	app, err := firebase.NewApp(common.BaseCtx, config, opt)
	if err != nil {
		panic(fmt.Errorf("error initializing app: %v", err))
	}

	fireClient, err := app.Storage(common.BaseCtx)
	if err != nil {
		panic(fmt.Errorf("error initializing client: %v", err))
	}

	common.Bucket, err = fireClient.DefaultBucket()
	if err != nil {
		panic(fmt.Errorf("error initializing common.Bucket: %v", err))
	}
}

// Convert a link to just it's stem, like "media/:id"
func LinkToLocation(link string) (string, error) {
	// Ckeck if link matches the format of a link
	linkRegex, err := regexp.Compile(`^https://storage\.googleapis\.com/download/storage/v1/b/dwitter\-72e9d\.appspot\.com/o/\w+\%2F\w+\.\w+\?.+$`)
	if err != nil {
		return "", err
	}
	matched := linkRegex.MatchString(link)
	if matched {
		// Reduce the link to its stem
		prefixRegex, err := regexp.Compile(`^https://storage\.googleapis\.com/download/storage/v1/b/dwitter\-72e9d\.appspot\.com/o/`)
		if err != nil {
			return "", err
		}
		noPrefix := prefixRegex.ReplaceAllString(link, "")
		queryRegex, err := regexp.Compile(`\?.+$`)
		if err != nil {
			return "", err
		}
		stem := queryRegex.ReplaceAllString(noPrefix, "")

		slashRegex, err := regexp.Compile(`\%2F`)
		if err != nil {
			return "", err
		}
		finalLink := slashRegex.ReplaceAllString(stem, "/")

		return finalLink, nil
	} else {
		return "", errors.New("media link invalid")
	}
}

// Delete an object given a stem, and allow deletion of its thumbnail
func DeleteLocation(location string, deleteThumb bool) error {
	// Delete corresponding thumbnail
	if deleteThumb {
		mediaCheckRegex, err := regexp.Compile(`media\/`)
		if err != nil {
			return err
		}
		thumbLocNoFormat := mediaCheckRegex.ReplaceAllString(location, "thumb/")

		extensionRegex, err := regexp.Compile(`\.\w+$`)
		if err != nil {
			return err
		}
		thumbLoc := extensionRegex.ReplaceAllString(thumbLocNoFormat, ".png")

		o := common.Bucket.Object(thumbLoc)
		if err := o.Delete(common.BaseCtx); err != nil {
			if err.Error() == "storage: object doesn't exist" {
				return errors.New("thumbnail for media not found")
			}
			return fmt.Errorf("Object(%q).Delete: %v", location, err)
		}
	}

	// Delete object
	o := common.Bucket.Object(location)
	if err := o.Delete(common.BaseCtx); err != nil {
		if err.Error() == "storage: object doesn't exist" {
			return errors.New("media not found")
		}
		return fmt.Errorf("Object(%q).Delete: %v", location, err)
	}
	return nil
}

// Handle media upload requests
func UploadMediaHandler(w http.ResponseWriter, r *http.Request) {
	// Check authentication
	authHeader := r.Header.Get("authorization")
	_, err := auth.Authenticate(authHeader)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
	} else {
		supportedFormats := map[string]bool{
			"image/gif":  true, // GIF
			"image/jpeg": true, // JPEG
			"image/png":  true, // PNG
			"video/mp4":  true, // MP4
		}
		isImage := map[string]bool{
			"image/gif":  true,  // GIF
			"image/jpeg": true,  // JPEG
			"image/png":  true,  // PNG
			"video/mp4":  false, // MP4
		}

		// Check if content type is "multipart/form-data"
		if r.Header.Get("Content-Type") != "" {
			value, _ := header.ParseValueAndParams(r.Header, "Content-Type")
			if value != "multipart/form-data" {
				msg := "Content-Type header is not multipart/form-data"
				http.Error(w, msg, http.StatusUnsupportedMediaType)
				return
			}
		}

		// Limit size to 8*8MB = 64MB
		err := r.ParseMultipartForm(64 << 20)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		formData := r.MultipartForm
		files := formData.File["files"]

		// Enforce limits
		if len(files) > 8 {
			msg := "Too many files. Limit is 8 files."
			http.Error(w, msg, http.StatusBadRequest)
			return
		}

		for _, file := range files {
			if file.Size > (8 << 20) {
				msg := "File too large. Limit is 8 files, 8MB each."
				http.Error(w, msg, http.StatusRequestEntityTooLarge)
				return
			}
			if !supportedFormats[file.Header.Get("Content-Type")] {
				msg := "Format unsupported."
				http.Error(w, msg, http.StatusUnsupportedMediaType)
				return
			}
		}

		links := []string{}

		for i := range files {
			// Open file in request
			file, err := files[i].Open()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer file.Close()

			operationCtx, cancel := context.WithTimeout(common.BaseCtx, time.Second*50)
			defer cancel()

			// Upload an object with storage.Writer.

			// Get a unique name for object
			randID := util.GenID(30)
			found := true
			for found {
				query := &storage.Query{Prefix: randID}
				it := common.Bucket.Objects(operationCtx, query)
				numObj := 0
				for {
					_, err := it.Next()
					if err == iterator.Done {
						if numObj == 0 {
							found = false
							break
						}
					}
					if err != nil {
						log.Fatal(err)
					}
					numObj += 1
					break
				}
				if numObj != 0 {
					randID = util.GenID(30)
				}
			}

			// Write file to cloud
			obj := common.Bucket.Object("media/" + randID + filepath.Ext(files[i].Filename))
			writer := obj.NewWriter(operationCtx)
			if _, err = io.Copy(writer, file); err != nil {
				panic(fmt.Errorf("io.Copy: %v", err))
			}
			if err := writer.Close(); err != nil {
				panic(fmt.Errorf("io.Copy: %v", err))
			}

			// Read it back
			reader, err := obj.NewReader(operationCtx)
			if err != nil {
				panic(fmt.Errorf("reader: %v", err))
			}
			defer reader.Close()
			data, err := ioutil.ReadAll(reader)
			if err != nil {
				panic(fmt.Errorf("ioutil.ReadAll: %v", err))
			}

			var thumb *image.NRGBA

			// Make thumbnail based on image or video
			if isImage[files[i].Header.Get("Content-Type")] {
				// Make thumbnail
				image, _, err := image.Decode(bytes.NewReader(data))
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				thumb = imaging.Thumbnail(image, 640, 360, imaging.Lanczos)
			} else {
				tempVidPath := "./tmp/" + util.GenID(40) + ".mp4"
				tempVid, err := os.Create(tempVidPath)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				_, err = tempVid.Write(data)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				path, err := generateVideoThumbnail(tempVidPath)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				picDat, err := imaging.Open(path)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				thumb = imaging.Thumbnail(picDat, 640, 360, imaging.Lanczos)

				err = os.Remove(path)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				err = os.Remove(tempVidPath)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			}

			// Save thumbnail
			thumbObj := common.Bucket.Object("thumb/" + randID + ".png")
			thumbWriter := thumbObj.NewWriter(operationCtx)
			if err = png.Encode(thumbWriter, thumb); err != nil {
				panic(fmt.Errorf("png.Encode: %v", err))
			}
			if err := thumbWriter.Close(); err != nil {
				panic(fmt.Errorf("png.Encode: %v", err))
			}

			mediaLink := writer.Attrs().MediaLink
			common.MediaCreatedButNotUsed[mediaLink] = true

			go destroyObjectAfterExpire(10, mediaLink)

			links = append(links, mediaLink)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(links)
	}
}

// Handle pfp upload requests
func UploadPFPHandler(w http.ResponseWriter, r *http.Request) {
	// Check authentication
	authHeader := r.Header.Get("authorization")
	username, err := auth.Authenticate(authHeader)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
	} else {
		supportedFormats := map[string]bool{
			"image/bmp":  true, // BMP
			"image/gif":  true, // GIF
			"image/jpeg": true, // JPEG
			"image/webp": true, // WEBP
			"image/png":  true, // PNG
		}

		// Check if content type is "multipart/form-data"
		if r.Header.Get("Content-Type") != "" {
			value, _ := header.ParseValueAndParams(r.Header, "Content-Type")
			if value != "multipart/form-data" {
				msg := "Content-Type header is not multipart/form-data"
				http.Error(w, msg, http.StatusUnsupportedMediaType)
				return
			}
		}

		// Limit size to 8MB
		err := r.ParseMultipartForm(8 << 20)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		formData := r.MultipartForm
		files := formData.File["files"]

		// Enforce limits
		if len(files) > 1 {
			msg := "Too many files. Limit is 1 file."
			http.Error(w, msg, http.StatusBadRequest)
			return
		}

		file := files[0]
		if file.Size > (8 << 20) {
			msg := "File too large. Limit is 8MB."
			http.Error(w, msg, http.StatusRequestEntityTooLarge)
			return
		}
		if !supportedFormats[file.Header.Get("Content-Type")] {
			msg := "Format unsupported."
			http.Error(w, msg, http.StatusUnsupportedMediaType)
			return
		}

		links := []string{}

		for i := range files {
			// Open file
			file, err := files[i].Open()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer file.Close()

			operationCtx, cancel := context.WithTimeout(common.BaseCtx, time.Second*50)
			defer cancel()

			// Upload to cloud
			randID := util.GenID(30)
			found := true
			for found {
				query := &storage.Query{Prefix: randID}
				it := common.Bucket.Objects(operationCtx, query)
				numObj := 0
				for {
					_, err := it.Next()
					if err == iterator.Done {
						if numObj == 0 {
							found = false
							break
						}
					}
					if err != nil {
						log.Fatal(err)
					}
					numObj += 1
					break
				}
				if numObj != 0 {
					randID = util.GenID(30)
				}
			}

			obj := common.Bucket.Object("pfp/pfp_" + username + filepath.Ext(files[i].Filename))
			writer := obj.NewWriter(operationCtx)
			if _, err = io.Copy(writer, file); err != nil {
				panic(fmt.Errorf("io.Copy: %v", err))
			}
			if err := writer.Close(); err != nil {
				panic(fmt.Errorf("io.Copy: %v", err))
			}
			links = append(links, writer.Attrs().MediaLink)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(links)
	}
}

// Generate a thumbnail from a video in the tmp directory
func generateVideoThumbnail(filepath string) (string, error) {
	// command line args, path, and command
	command := "ffmpeg"
	frameExtractionTime := "0:00:00.000"
	vframes := "1"
	qv := "2"
	output := "./tmp/" + time.Now().Format(time.Kitchen) + util.GenID(40) + ".png"

	cmd := exec.Command(command,
		"-ss", frameExtractionTime,
		"-i", filepath, // to read from
		"-vframes", vframes,
		"-q:v", qv,
		output)

	// run the command and don't wait for it to finish. waiting exec is run
	// ignore errors for examples-sake
	err := cmd.Run()
	return output, err
}

// Destroy an object when it expires
func destroyObjectAfterExpire(minutes int, id string) {
	time.Sleep(time.Minute * time.Duration(minutes))
	if common.MediaCreatedButNotUsed[id] {
		loc, err := LinkToLocation(id)
		if err != nil {
			fmt.Printf("Error finding media: %v", err)
		}
		err = DeleteLocation(loc, true)
		if err != nil {
			fmt.Printf("Error auto-deleting media: %v", err)
		}
	}
}
