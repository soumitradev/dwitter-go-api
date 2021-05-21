package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"regexp"
	"time"

	"cloud.google.com/go/storage"
	firebase "firebase.google.com/go/v4"
	"github.com/golang/gddo/httputil/header"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

var bucket *storage.BucketHandle

func initCDN() {
	opt := option.WithCredentialsFile("./cdn_key.json")

	config := &firebase.Config{
		StorageBucket: "dwitter-72e9d.appspot.com",
	}

	app, err := firebase.NewApp(context.Background(), config, opt)
	if err != nil {
		panic(fmt.Errorf("error initializing app: %v", err))
	}

	fireClient, err := app.Storage(context.Background())
	if err != nil {
		panic(fmt.Errorf("error initializing client: %v", err))
	}

	bucket, err = fireClient.DefaultBucket()
	if err != nil {
		panic(fmt.Errorf("error initializing bucket: %v", err))
	}
}

func DeleteFile(link string) error {
	re1, err := regexp.Compile(`^https://storage\.googleapis\.com/download/storage/v1/b/dwitter\-72e9d\.appspot\.com/o/\w+\.\w+\?generation\=\d+\&alt\=media$`)
	if err != nil {
		panic(err)
	}
	matched := re1.MatchString(link)
	if matched {
		re2, err := regexp.Compile(`^https://storage\.googleapis\.com/download/storage/v1/b/dwitter\-72e9d\.appspot\.com/o/`)
		if err != nil {
			panic(err)
		}
		nlink := re2.ReplaceAllString(link, "")
		re3, err := regexp.Compile(`\?generation\=\d+\&alt\=media$`)
		if err != nil {
			panic(err)
		}
		finlink := re3.ReplaceAllString(nlink, "")

		o := bucket.Object(finlink)
		if err := o.Delete(ctx); err != nil {
			if err.Error() == "storage: object doesn't exist" {
				return errors.New("media not found")
			}
			return fmt.Errorf("Object(%q).Delete: %v", finlink, err)
		}

		return nil
	} else {
		return errors.New("media link invalid")
	}
}

// Handle login requests
func uploadFile(w http.ResponseWriter, r *http.Request) {

	supportedFormats := map[string]bool{
		"image/bmp":       true, // BMP
		"image/gif":       true, // GIF
		"image/jpeg":      true, // JPEG
		"image/webp":      true, // WEBP
		"image/png":       true, // PNG
		"video/mp4":       true, // MP4
		"video/x-msvideo": true, // AVI
		"video/ogg":       true, // OGG
		"video/webm":      true, // WEBM
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

	formdata := r.MultipartForm
	files := formdata.File["files"]

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
		fmt.Println(files[i].Header)
		file, err := files[i].Open()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer file.Close()

		ctx, cancel := context.WithTimeout(ctx, time.Second*50)
		defer cancel()

		// Upload an object with storage.Writer.

		randID := genID(30)
		found := true
		for found {
			query := &storage.Query{Prefix: randID}
			it := bucket.Objects(ctx, query)
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
				randID = genID(30)
			}
		}

		obj := bucket.Object(randID + filepath.Ext(files[i].Filename))
		wc := obj.NewWriter(ctx)
		if _, err = io.Copy(wc, file); err != nil {
			panic(fmt.Errorf("io.Copy: %v", err))
		}
		if err := wc.Close(); err != nil {
			panic(fmt.Errorf("io.Copy: %v", err))
		}
		links = append(links, wc.Attrs().MediaLink)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(links)
}
