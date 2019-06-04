package db

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// MongoImage is the struct which defines the image
// in the mongo database
type MongoImage struct {
	ID          bson.ObjectId `bson:"_id"`
	Author      string        `bson:"author"`
	Caption     string        `bson:"caption"`
	ContentType string        `bson:"contentType"`
	DateTime    string        `bson:"dateTime"`
	FileID      bson.ObjectId `bson:"fileID"`
	FileSize    int64         `bson:"fileSize"`
	Height      int           `bson:"height"`
	Name        string        `bson:"name"`
	Width       int           `bson:"width"`
}

// MongoImages is defined here for serialization purposess
type MongoImages struct {
	Images []MongoImage
}

// MongoImageIndex is used by mongo to index the MongoImage
// structure
func MongoImageIndex() mgo.Index {
	return mgo.Index{
		Key:        []string{"Name"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
}

func newMongoImage(img MongoImage) *MongoImage {
	return &MongoImage{
		ID:          img.ID,
		Author:      img.Author,
		Caption:     img.Caption,
		ContentType: img.ContentType,
		DateTime:    img.DateTime,
		FileID:      img.FileID,
		FileSize:    img.FileSize,
		Height:      img.Height,
		Name:        img.Name,
		Width:       img.Width,
	}
}
