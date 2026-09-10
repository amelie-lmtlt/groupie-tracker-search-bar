package main

import (
	"encoding/json"
	gt "groupietracker/artists"
	"html/template"
	"mime"
	"net/http"
	"os"
	"strconv"
	"strings"
)

const (
	ISE_MESSAGE     = "500 Internal Server Error"
	SERVICE_MESSAGE = "503 Service Unavailable"
)

func HandlerHome(w http.ResponseWriter, r *http.Request) {
	var t *template.Template
	var err error

	if r.URL.Path == "/" {
		t, err = template.ParseFiles("static/templates/index.html")
	} else if r.URL.Path == "/content.html" {
		t, err = template.ParseFiles("static/templates/content.html")
	} else {
		http.Error(w, "404 Not Found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(w, "error when loading template: "+err.Error(), http.StatusInternalServerError)
		return
	}

	err = t.Execute(w, nil)
	if err != nil {
		http.Error(w, "error when executing template: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func ArtistEndpoint(w http.ResponseWriter, r *http.Request) {
	// Number of artists per page
	pagination := 10

	paginationSelect, err := strconv.Atoi(r.URL.Query().Get("n"))
	if err == nil {
		if paginationSelect == 0 {
			pagination = 5
		} else if paginationSelect == 2 {
			pagination = 25
		}
	}

	// Get all artists
	artists, err := gt.GetArtistsData()
	if err != nil {
		http.Error(w, SERVICE_MESSAGE+": "+err.Error(), http.StatusServiceUnavailable)
		return
	}

	var artistsArray []gt.Artists
	err = json.Unmarshal(artists, &artistsArray)
	if err != nil {
		http.Error(w, ISE_MESSAGE, http.StatusInternalServerError)
		return
	}

	// Total number of artists
	artistsNumber := len(artistsArray)

	// Total number of pages
	pageNumber := artistsNumber / pagination
	if artistsNumber%pagination != 0 {
		pageNumber++
	}

	// Get requested page
	page := 1

	pageSelect, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err == nil && pageSelect > 0 {
		page = pageSelect
	}

	// Prevent invalid pages
	if page > pageNumber {
		page = pageNumber
	}

	// Calculate slice indexes
	begin := (page - 1) * pagination
	end := min(begin+pagination, artistsNumber)

	if begin < artistsNumber {
		artistsArray = artistsArray[begin:end]
	} else {
		artistsArray = []gt.Artists{}
	}

	// Create response
	artistsObj := gt.ArtistsResponse{
		ArtistsData:   artistsArray,
		ArtistsNumber: artistsNumber,
		PageSize:      pagination,
		PagesNumber:   pageNumber,
	}

	response, err := json.Marshal(artistsObj)
	if err != nil {
		http.Error(w, ISE_MESSAGE, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
func ArtistSearch(w http.ResponseWriter, r *http.Request) {
	artists, err := gt.GetArtistsData()
	if err != nil {
		http.Error(w, SERVICE_MESSAGE+": "+err.Error(), http.StatusServiceUnavailable)
		return
	}
	var artistsArray []gt.Artists
	err = json.Unmarshal(artists, &artistsArray)
	if err != nil {
		http.Error(w, ISE_MESSAGE, http.StatusInternalServerError)
		return
	}

	query := r.PathValue("query")
	if query == "" {
		http.Error(w, "Bad Request: query string must be provided", http.StatusBadRequest)
		return
	}
	filteredArtists := []gt.Artists{}
	for _, artist := range artistsArray {
		if strings.Index(strings.ToLower(artist.Name), strings.ToLower(query)) != -1 {
			filteredArtists = append(filteredArtists, artist)
		}
	}

	artistsObj, _ := json.Marshal(filteredArtists)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", strconv.Itoa(len(artistsObj)))
	w.Write(artistsObj)
}

func HandlerArtistPage(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Bad Request: id must be provided and must be integer", http.StatusBadRequest)
		return
	}

	if id < 1 {
		http.Error(w, "Bad Request: id must be at least 1", http.StatusBadRequest)
		return
	}

	artistData, err := gt.GetArtistData(id)
	if err != nil {
		http.Error(w, SERVICE_MESSAGE+": "+err.Error(), http.StatusServiceUnavailable)
		return
	}
	var artist gt.Artists
	err = json.Unmarshal(artistData, &artist)
	if err != nil {
		http.Error(w, ISE_MESSAGE, http.StatusInternalServerError)
		return
	}
	if artist.ID != id {
		http.Error(w, "404 Not Found", http.StatusNotFound)
		return
	}

	relationData, err := gt.GetRelationData(id)
	if err != nil {
		http.Error(w, SERVICE_MESSAGE+": "+err.Error(), http.StatusServiceUnavailable)
		return
	}
	var relation gt.Relation
	err = json.Unmarshal(relationData, &relation)

	if err != nil {
		http.Error(w, ISE_MESSAGE, http.StatusInternalServerError)
		return
	}
	page := gt.ArtistPage{
		ID:             artist.ID,
		Name:           artist.Name,
		Image:          artist.Image,
		Members:        artist.Members,
		CreationDate:   artist.CreationDate,
		FirstAlbum:     artist.FirstAlbum,
		ConcertDetails: relation.ConcertDetails,
	}
	t, err := template.ParseFiles("static/templates/artist.html")
	if err != nil {
		http.Error(w, "error when loading template: "+err.Error(), http.StatusInternalServerError)
		return
	}
	err = t.Execute(w, page)
	if err != nil {
		http.Error(w, "error when executing template: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func RelationEndpoint(w http.ResponseWriter, r *http.Request) {
	if r.PathValue("id") == "" {
		http.Error(w, "400 Bad Request: artist ID must be present", http.StatusBadRequest)
		return
	}
	artistID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || artistID < 1 {
		http.Error(w, "400 Bad Request: artist ID must be a positive integer", http.StatusBadRequest)
		return
	}

	relation, err := gt.GetRelationData(artistID)
	if err != nil {
		http.Error(w, SERVICE_MESSAGE+": "+err.Error(), http.StatusServiceUnavailable)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", strconv.Itoa(len(relation)))
	w.Write(relation)
}

func StaticPublicRessourceServer(w http.ResponseWriter, r *http.Request) {
	requestedFile, err := os.OpenInRoot("static/", r.URL.Path[1:])
	if err != nil {
		http.Error(w, "404 Not Found", http.StatusNotFound)
		return
	}
	statInfo, err := requestedFile.Stat()
	if err != nil {
		http.Error(w, ISE_MESSAGE, http.StatusInternalServerError)
		return
	}
	extension := strings.LastIndex(requestedFile.Name(), ".")
	contentType := "application/octet-stream"
	if extension > 0 && mime.TypeByExtension(requestedFile.Name()[extension:]) != "" {
		contentType = mime.TypeByExtension(requestedFile.Name()[extension:])
	}
	w.Header().Set("Content-Type", contentType)
	http.ServeContent(w, r, "", statInfo.ModTime(), requestedFile)
}

func AllowMethod(h http.HandlerFunc, methods []string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		for _, m := range methods {
			if m == r.Method {
				h(w, r)
				return
			}
		}
		w.Header().Set("Allow", strings.Join(methods, ", "))
		http.Error(w, "405 Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
