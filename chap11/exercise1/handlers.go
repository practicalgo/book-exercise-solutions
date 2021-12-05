package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func packageRegHandler(
	w http.ResponseWriter,
	r *http.Request,
	config appConfig,
) {
	d := pkgRegisterResponse{}
	err := r.ParseMultipartForm(5000)
	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusBadRequest,
		)
		return
	}
	mForm := r.MultipartForm
	fHeader := mForm.File["filedata"][0]

	packageName := mForm.Value["name"][0]
	packageVersion := mForm.Value["version"][0]

	d.ID = fmt.Sprintf(
		"%s-%s-%s",
		packageName,
		packageVersion,
		fHeader.Filename,
	)
	nBytes, err := uploadData(config, d.ID, fHeader)
	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	config.logger.Printf(
		"Package uploaded: %s. Bytes written: %d\n",
		d.ID,
		nBytes,
	)
	jsonData, err := json.Marshal(d)
	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, string(jsonData))
}

func packageGetHandler(
	w http.ResponseWriter,
	r *http.Request,
	config appConfig,
) {
	queryParams := r.URL.Query()
	packageID := queryParams.Get("id")
	download := queryParams.Get("download")

	exists, err := config.packageBucket.Exists(
		r.Context(), packageID,
	)
	if err != nil || !exists {
		http.Error(w, "invalid package ID", http.StatusBadRequest)
		return
	}

	if len(download) == 0 {
		url, err := config.packageBucket.SignedURL(
			r.Context(),
			packageID,
			nil,
		)
		if err != nil {
			http.Error(
				w,
				err.Error(),
				http.StatusInternalServerError,
			)
			return
		}

		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	}

	reader, err := config.packageBucket.NewReader(r.Context(), packageID, nil)
	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}
	defer reader.Close()
	w.Header().Add("Content-Type", reader.ContentType())
	w.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", packageID))
	io.Copy(w, reader)
}

func packageHandler(
	w http.ResponseWriter,
	r *http.Request,
	config appConfig,
) {
	switch r.Method {
	case "POST":
		packageRegHandler(w, r, config)
	case "GET":
		packageGetHandler(w, r, config)
	default:
		http.Error(w,
			"Cannot process request",
			http.StatusMethodNotAllowed,
		)
	}
}
