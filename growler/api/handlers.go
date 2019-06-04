package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/murlokito/growler/growler/db"
	"gopkg.in/mgo.v2"
)

// ResponseMessage defined to be used for serialization purposes
type ResponseMessage struct {
	Description string `json:"Description"`
}

// ResponseImageStream defined to be used for serialization purposes
type ResponseImageStream struct {
	Size int64  `json:"Size"`
	Name string `json:"Name"`
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func setHeadersImageStream(w http.ResponseWriter, fileSize int, fileName string) {

	/*
	 * Set headers to tell the caller how big the file is and stream it back
	 */
	w.Header().Set("Content-Length", strconv.Itoa(fileSize))
	w.Header().Set("Content-Disposition", "inline; filename=\""+fileName+"\"")

}

// Index is the handler for the '/' endpoint, which is to be used for
// tests only
func (ws *WebService) Index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	p := "Success!"
	respondWithJSON(w, http.StatusOK, p)
	return
}

// StreamImage is the handler for the '/images/stream/{id}' endpoint,
// which is to be used to get the stream of an image stored in GridFS
func (ws *WebService) StreamImage(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	var gridFile *mgo.GridFile
	var imageRecord *db.MongoImage

	payload := mux.Vars(r)

	if payload["id"] != "" {
		var id int

		id, err := strconv.Atoi(payload["id"])
		if err != nil {
			respondWithError(w, http.StatusBadRequest,
				err.Error())
			return
		}
		/*
		 * Get the record by ID
		 */

		if err := ws.MgoClient.DB.C("images").FindId(id).One(&imageRecord); err != nil {
			respondWithError(w, http.StatusBadRequest,
				fmt.Errorf("Error getting image by ID: %v", id).Error())
			return
		}

		/*
		 * Get the file from GridFS
		 */
		fileID := imageRecord.FileID

		gridFile, err := ws.MgoClient.DB.GridFS("images").OpenId(fileID)
		if err != nil {
			respondWithError(w, http.StatusBadRequest,
				err.Error())
			return
		}

		defer gridFile.Close()

	} else {
		respondWithError(w, http.StatusBadRequest,
			fmt.Errorf("'ID' not found in payload").Error())
		return
	}
	/*
		int(gridFile.Size())
		gridFile.Name()*/
	p := ResponseImageStream{
		Size: gridFile.Size(),
		Name: gridFile.Name(),
	}
	respondWithJSON(w, http.StatusOK, p)
	return
}

// GetImages is the handler for the '/images/get/{limit}' endpoint,
// which is to be used to get the stream of an image stored in GridFS
func (ws *WebService) GetImages(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	var imageCollection *db.MongoImages

	payload := mux.Vars(r)

	if payload["limit"] != "" {
		var limit int
		limit, err := strconv.Atoi(payload["limit"])
		if err != nil {
			respondWithError(w, http.StatusBadRequest,
				err.Error())
			return
		}

		/*
		 * Fetch a number of records up to a maximum equaling `limit`
		 */

		iter := ws.MgoClient.DB.C("images").Find(nil).Limit(limit).Iter()

		err = iter.All(&imageCollection)
		if err != nil {
			respondWithError(w, http.StatusBadRequest,
				fmt.Errorf("Error getting `images` collection").Error())
			return
		}

		respondWithJSON(w, http.StatusOK, imageCollection)
		return

	}
	respondWithError(w, http.StatusBadRequest,
		fmt.Errorf("'ID' not found in payload").Error())
	return

}
