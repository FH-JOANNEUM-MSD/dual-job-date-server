package handlers

import (
	"dual-job-date-server/internal/repository"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func UploadCompanyLogoHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	companyID, err := strconv.Atoi(vars["id"])
	if err != nil || companyID <= 0 {
		http.Error(w, "ungueltige company id", http.StatusBadRequest)
		return
	}

	if err := r.ParseMultipartForm(8 << 20); err != nil {
		http.Error(w, "multipart/form-data erwartet", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("logo")
	if err != nil {
		file, header, err = r.FormFile("file")
	}
	if err != nil {
		http.Error(w, "file feld 'logo' oder 'file' fehlt", http.StatusBadRequest)
		return
	}
	defer file.Close()

	fileData, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "datei konnte nicht gelesen werden", http.StatusBadRequest)
		return
	}

	contentType := header.Header.Get("Content-Type")
	if contentType == "" {
		contentType = http.DetectContentType(fileData)
	}

	result, err := repository.UploadCompanyLogo(companyID, header.Filename, contentType, fileData)
	if err != nil {
		http.Error(w, fmt.Sprintf("logo upload fehlgeschlagen: %v", err), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(result)
}
