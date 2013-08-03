/*Save and Get Files. Only .png format supported right now.
*TODO: save and serve files by content type
 */
package hamster

import (
	"bufio"
	"bytes"
	"encoding/json"
	//"fmt"
	"image"
	"image/png"
	"io"
	"labix.org/v2/mgo/bson"
	"net/http"
)

type File struct {
	FileName string
}

//Saves file read from request body
//POST:/api/v1/files/:fileName
func (s *Server) SaveFile(w http.ResponseWriter, r *http.Request) {
	s.logger.SetPrefix("SaveFile: ")

	//get file params and file data reader
	file_name := s.getFileName(w, r)
	fileReader := bufio.NewReader(r.Body)
	defer r.Body.Close()

	//get meta data
	meta_data_json := r.Header.Get("X-Meta-Data")

	//parse meta data
	var metadata map[string]interface{}
	json.Unmarshal([]byte(meta_data_json), &metadata)

	//get session
	session := s.db.GetSession()
	defer session.Close()
	db := session.DB("")

	//create file
	file, err := db.GridFS("fs").Create(file_name)
	if err != nil {
		s.internalError(r, w, err, "could not create file")
	}

	//copy incoming data to file
	_, err = io.Copy(file, fileReader)
	if err != nil {
		s.internalError(r, w, err, "could not copy file")
	}

	//set content type and meta deta
	file.SetContentType("image/png")
	file.SetMeta(metadata)

	//encode file id and serve
	file_id := encodeBase64Token(file.Id().(bson.ObjectId).Hex())
	response := SaveFileResponse{FileId: file_id, FileName: file.Name()}

	err = file.Close()
	if err != nil {
		s.internalError(r, w, err, "could not close file")
	}

	//respond

	s.serveJson(w, &response)

}

//Gets file from GridFS and writes to response body
//GET:/api/v1/files/:fileName/:fileId handler
func (s *Server) GetFile(w http.ResponseWriter, r *http.Request) {
	s.logger.SetPrefix("GetFile: ")

	//get file params
	_, file_id := s.getFileParams(w, r)

	//get session
	session := s.db.GetSession()
	defer session.Close()
	db := session.DB("")

	//open file from GridFS
	file, err := db.GridFS("fs").OpenId(bson.ObjectIdHex(file_id))
	if err != nil {
		s.internalError(r, w, err, "could not open file")
	}

	//copy buffer
	var buf bytes.Buffer
	_, err = io.Copy(&buf, file)
	if err != nil {
		s.internalError(r, w, err, "could not copy buffer")
	}

	//decode buffer
	img, _, err := image.Decode(&buf)

	if err != nil {
		s.internalError(r, w, err, "could not decode image")
	}

	contentType := file.ContentType()
	//fmt.Printf("content type: %v \n", contentType)

	err = file.Close()
	if err != nil {
		s.internalError(r, w, err, "could not close file")
	}

	//set content type and write to response body
	w.Header().Set("Content-Type", contentType)
	//TODO encode by mime type
	png.Encode(w, img)

}
