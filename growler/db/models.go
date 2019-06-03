package db

import "gopkg.in/mgo.v2/bson"

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

// GridFile is the struct which defines the file
// in GridFS
type GridFile struct {
	Name string `json:"Name"`
	Size int    `json:"Size"`
}
