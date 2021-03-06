package routes

import (
	"encoding/json"
	"fmt"
	"github.com/gabriel-vasile/mimetype"
	"github.com/gorilla/mux"
	"go-api/model"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func buildDriveDirectLink(id string) string {
	return fmt.Sprintf(`https://drive.google.com/uc?export=download&id=%s`, id)
}

const MaxFileSize = 6 << 20 // 6 MB

func (app *App) RegisterFloorRoutes() {
	app.router.HandleFunc("/floors", app.CreateFloor).Methods("POST")
	app.router.HandleFunc("/floors/{id}", app.GetOneFloor).Methods("GET")
	app.router.HandleFunc("/floors", app.GetAllFloors).Methods("GET")
	//app.router.HandleFunc("/floors/{id}", app.UpdateFloor).Methods("PATCH")
	app.router.HandleFunc("/floors/{id}", app.DeleteFloor).Methods("DELETE")
}

var acceptedImages = map[string]bool{
	"image/png": true,
	"image/jpeg": true, // jpg are considered jpeg
}

func (app *App) CreateFloor(w http.ResponseWriter, r *http.Request) {
	var newFloor model.Floor
	r.Body = http.MaxBytesReader(w, r.Body, MaxFileSize+512)
	parseErr := r.ParseMultipartForm(MaxFileSize)
	if parseErr != nil {
		log.Println("App.CreateFloor - failed to parse message")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if r.MultipartForm == nil || r.MultipartForm.File == nil {
		log.Println("App.CreateFloor - expecting multipart form file")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	imageFile, _, err := r.FormFile("image")
	if err != nil {
		log.Println("App.CreateFloor - image is absent: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	mime, errMime := mimetype.DetectReader(imageFile)
	if errMime != nil {
		log.Println("App.CreateFloor - Error handling mime: " + errMime.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !acceptedImages[mime.String()] {
		log.Println("App.CreateFloor - The image must be of type jpg, jpeg or png: Mime was of type " + mime.String())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// MIME Reads part of the file, rewind to the start
	_, err = imageFile.Seek(0, io.SeekStart)
	if err != nil {
		log.Println("App.CreateFloor - Something went wrong with seeking back to the front")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	name := r.FormValue("name")
	if name == "" {
		log.Println("App.CreateFloor - name is absent")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	address := r.FormValue("address")
	if address == "" {
		log.Println("App.CreateFloor - address is absent")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id, parseErr := app.gDrive.UploadFloorPlan(name, imageFile)
	if parseErr != nil {
		log.Println("App.CreateFloor - Failed to upload image to Google Drive")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	newFloor.Name = name
	newFloor.DownloadURL = buildDriveDirectLink(id)
	newFloor.Address = address
	id, err = app.store.FloorProvider.CreateFloor(&newFloor)
	if err != nil {
		log.Printf("App.CreateBooking - error creating booking %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	newFloor.ID = id

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newFloor)
}

func (app *App) GetOneFloor(w http.ResponseWriter, r *http.Request) {
	floorID := mux.Vars(r)["id"]

	if floorID == "" {
		log.Printf("App.GetOneFloor - empty floor id")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	floor, err := app.store.FloorProvider.GetOneFloor(floorID)
	if err != nil {
		log.Printf("App.GetOneFloor - error getting floor from provider %v", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(floor)
}

func (app *App) GetAllFloors(w http.ResponseWriter, r *http.Request) {
	floors, err := app.store.FloorProvider.GetAllFloors()
	if err != nil {
		log.Printf("App.GetAllFloors - error getting all floors from provider %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(floors)
}

//func (app *App) UpdateFloor(w http.ResponseWriter, r *http.Request) {
//	floorID := mux.Vars(r)["id"]
//
//	if floorID == "" {
//		log.Printf("App.UpdateFloor - empty floor id")
//		w.WriteHeader(http.StatusBadRequest)
//		return
//	}
//
//	var updatedFloor model.Floor
//	reqBody, err := ioutil.ReadAll(r.Body)
//	if err != nil {
//		log.Printf("App.UpdateFloor - error reading request body %v", err)
//		w.WriteHeader(http.StatusBadRequest)
//		return
//	}
//	err = json.Unmarshal(reqBody, &updatedFloor)
//	if err != nil {
//		log.Printf("App.UpdateFloor - error unmarshaling request body %v", err)
//		w.WriteHeader(http.StatusBadRequest)
//		return
//	}
//
//	err = app.store.FloorProvider.UpdateFloor(floorID, &updatedFloor)
//	if err != nil {
//		log.Printf("App.UpdateFloor - error getting all floors from provider %v", err)
//		w.WriteHeader(http.StatusNotFound)
//		return
//	}
//}

func (app *App) DeleteFloor(w http.ResponseWriter, r *http.Request) {
	floorID := mux.Vars(r)["id"]

	if floorID == "" {
		log.Printf("App.DeleteFloor - empty floor id")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var deleteFloor model.DeleteFloor
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("App.DeleteFloor - error reading request body %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = json.Unmarshal(reqBody, &deleteFloor)
	if err != nil {
		log.Printf("App.DeleteFloor - error unmarshaling request body %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = app.store.FloorProvider.RemoveFloor(floorID, deleteFloor.ForceDelete)
	if err != nil {
		if strings.Contains(err.Error(), "invalid operation") {
			log.Printf("App.DeleteFloor - %v", err)
			w.WriteHeader(http.StatusForbidden)
		} else {
			log.Printf("App.DeleteFloor - error getting all floors from provider %v", err)
			w.WriteHeader(http.StatusNotFound)
		}
		return
	}
	w.WriteHeader(http.StatusOK)
}
