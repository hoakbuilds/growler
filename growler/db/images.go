package db

import (
	"bufio"
	"encoding/json"
	"errors"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"time"

	"gopkg.in/mgo.v2"
)

// writeToGridFile is used to write to GridFS
func writeToGridFile(file multipart.File, gridFile *mgo.GridFile) error {
	reader := bufio.NewReader(file)
	defer func() { file.Close() }()
	// make a buffer to keep chunks that are read
	buf := make([]byte, 1024)
	for {
		// read a chunk
		n, err := reader.Read(buf)
		if err != nil && err != io.EOF {
			return errors.New("Could not read the input file")
		}
		if n == 0 {
			break
		}
		// write a chunk
		if _, err := gridFile.Write(buf[:n]); err != nil {
			return errors.New("Could not write to GridFs for " + gridFile.Name())
		}
	}
	gridFile.Close()
	return nil
}

// PostImage is a method defined by the MongoClient structure, it is used
// by the crawler to post images
func (c *MongoClient) PostImage(image MongoImage) error {
	if gridFile, err := c.DB.GridFS("fs").Create(image.Name); err != nil {

		return err
	} else {
		gridFile.SetMeta(image)
		gridFile.SetName(image.Name)
		if err := writeToGridFile(file, gridFile); err != nil {

			return err
		}
	}
}

func decodeMongoImage(ur *http.Request) (MongoImage, error) {
	var u MongoImage
	if ur.Body == nil {
		return u, errors.New("no request body")
	}
	decoder := json.NewDecoder(ur.Body)
	err := decoder.Decode(&u)
	return u, err
}

func serveFromDB(w http.ResponseWriter, r *http.Request) {
	var gridfs *mgo.GridFS // Obtain GridFS via Database.GridFS(prefix)

	name := "somefile.pdf"
	f, err := gridfs.Open(name)
	if err != nil {
		log.Printf("Failed to open %s: %v", name, err)
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}
	defer f.Close()

	http.ServeContent(w, r, name, time.Now(), f) // Use proper last mod time
}

func ensureIndex(s *mgo.Session) {
	session := s.Copy()
	defer session.Close()

	c := session.DB("growler").C("images")

	index := mgo.Index{
		Key:        []string{"Name"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
	err := c.EnsureIndex(index)
	if err != nil {
		panic(err)
	}
}
