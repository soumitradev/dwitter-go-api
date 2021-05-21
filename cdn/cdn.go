package cdn

import (
	"bytes"
	"context"
	"dwitter_go_graphql/auth"
	"dwitter_go_graphql/consts"
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

func InitCDN() {
	opt := option.WithCredentialsFile("./cdn_key.json")

	consts.MediaCreatedButNotUsed = make(map[string]bool)

	config := &firebase.Config{
		StorageBucket: "dwitter-72e9d.appspot.com",
	}

	app, err := firebase.NewApp(consts.BaseCtx, config, opt)
	if err != nil {
		panic(fmt.Errorf("error initializing app: %v", err))
	}

	fireClient, err := app.Storage(consts.BaseCtx)
	if err != nil {
		panic(fmt.Errorf("error initializing client: %v", err))
	}

	consts.Bucket, err = fireClient.DefaultBucket()
	if err != nil {
		panic(fmt.Errorf("error initializing consts.Bucket: %v", err))
	}
}

func LinkToLocation(link string) (string, error) {
	re1, err := regexp.Compile(`^https://storage\.googleapis\.com/download/storage/v1/b/dwitter\-72e9d\.appspot\.com/o/\w+\%2F\w+\.\w+\?.+$`)
	if err != nil {
		return "", err
	}
	matched := re1.MatchString(link)
	if matched {
		re2, err := regexp.Compile(`^https://storage\.googleapis\.com/download/storage/v1/b/dwitter\-72e9d\.appspot\.com/o/`)
		if err != nil {
			return "", err
		}
		nlink := re2.ReplaceAllString(link, "")
		re3, err := regexp.Compile(`\?.+$`)
		if err != nil {
			return "", err
		}
		stem := re3.ReplaceAllString(nlink, "")

		re4, err := regexp.Compile(`\%2F`)
		if err != nil {
			return "", err
		}
		finlink := re4.ReplaceAllString(stem, "/")

		return finlink, nil
	} else {
		return "", errors.New("media link invalid")
	}
}

func DeleteLocation(location string, deleteThumb bool) error {
	if deleteThumb {
		re1, err := regexp.Compile(`media\/`)
		if err != nil {
			return err
		}
		thumbLocNoFormat := re1.ReplaceAllString(location, "thumb/")

		re2, err := regexp.Compile(`\.\w+$`)
		if err != nil {
			return err
		}
		thumbLoc := re2.ReplaceAllString(thumbLocNoFormat, ".png")

		o := consts.Bucket.Object(thumbLoc)
		if err := o.Delete(consts.BaseCtx); err != nil {
			if err.Error() == "storage: object doesn't exist" {
				return errors.New("thumbnail for media not found")
			}
			return fmt.Errorf("Object(%q).Delete: %v", location, err)
		}
	}

	o := consts.Bucket.Object(location)
	if err := o.Delete(consts.BaseCtx); err != nil {
		if err.Error() == "storage: object doesn't exist" {
			return errors.New("media not found")
		}
		return fmt.Errorf("Object(%q).Delete: %v", location, err)
	}
	return nil
}

// Handle login requests
func UploadMedia(w http.ResponseWriter, r *http.Request) {
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
			file, err := files[i].Open()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer file.Close()

			myCtx, cancel := context.WithTimeout(consts.BaseCtx, time.Second*50)
			defer cancel()

			// Upload an object with storage.Writer.

			randID := util.GenID(30)
			found := true
			for found {
				query := &storage.Query{Prefix: randID}
				it := consts.Bucket.Objects(myCtx, query)
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
			obj := consts.Bucket.Object("media/" + randID + filepath.Ext(files[i].Filename))
			writer := obj.NewWriter(myCtx)
			if _, err = io.Copy(writer, file); err != nil {
				panic(fmt.Errorf("io.Copy: %v", err))
			}
			if err := writer.Close(); err != nil {
				panic(fmt.Errorf("io.Copy: %v", err))
			}

			// Read it back
			reader, err := obj.NewReader(myCtx)
			if err != nil {
				panic(fmt.Errorf("reader: %v", err))
			}
			defer reader.Close()
			data, err := ioutil.ReadAll(reader)
			if err != nil {
				panic(fmt.Errorf("ioutil.ReadAll: %v", err))
			}

			var thumb *image.NRGBA

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

				path, err := videoThumb(tempVidPath)
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
			thumbObj := consts.Bucket.Object("thumb/" + randID + ".png")
			thumbWriter := thumbObj.NewWriter(myCtx)
			if err = png.Encode(thumbWriter, thumb); err != nil {
				panic(fmt.Errorf("png.Encode: %v", err))
			}
			if err := thumbWriter.Close(); err != nil {
				panic(fmt.Errorf("png.Encode: %v", err))
			}

			mlink := writer.Attrs().MediaLink
			consts.MediaCreatedButNotUsed[mlink] = true

			go destroyObjectAfterExpire(10, mlink)

			links = append(links, mlink)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(links)
	}
}

// Handle login requests
func UploadPfp(w http.ResponseWriter, r *http.Request) {
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

		formdata := r.MultipartForm
		files := formdata.File["files"]

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
			file, err := files[i].Open()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer file.Close()

			myCtx, cancel := context.WithTimeout(consts.BaseCtx, time.Second*50)
			defer cancel()

			// Upload an object with storage.Writer.

			randID := util.GenID(30)
			found := true
			for found {
				query := &storage.Query{Prefix: randID}
				it := consts.Bucket.Objects(myCtx, query)
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

			obj := consts.Bucket.Object("pfp/pfp_" + username + filepath.Ext(files[i].Filename))
			wc := obj.NewWriter(myCtx)
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
}

func videoThumb(filepath string) (string, error) {
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

func destroyObjectAfterExpire(minutes int, id string) {
	time.Sleep(time.Minute * time.Duration(minutes))
	if consts.MediaCreatedButNotUsed[id] {
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
