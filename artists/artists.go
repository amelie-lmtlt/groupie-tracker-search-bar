package groupietracker

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
)

const (
	BASE_API_ADDRESS  = "https://groupietrackers.herokuapp.com/api/"
	ARTIST_ENDPOINT   = BASE_API_ADDRESS + "artists"
	DATE_ENDPOINT     = BASE_API_ADDRESS + "dates"
	LOCATION_ENDPOINT = BASE_API_ADDRESS + "locations"
	RELATION_ENDPOINT = BASE_API_ADDRESS + "relation"
)

type Artists struct {
	ID           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
	ConcertsURL  string   `json:"relations"`
}

type Date struct {
	ID    int `json:"id"`
	Dates int `json:"dates"`
}

type Location struct {
	ID        int      `json:"id"`
	Locations []string `json:"locations"`
}

type Relation struct {
	ID             int                 `json:"id"`
	ConcertDetails map[string][]string `json:"datesLocations"`
}

type ArtistPage struct {
	ID             int
	Name           string
	Image          string
	Members        []string
	CreationDate   int
	FirstAlbum     string
	ConcertDetails map[string][]string
}

type ArtistsResponse struct {
	ArtistsData   []Artists
	ArtistsNumber int
	PageSize      int
	PagesNumber   int
}

func (a *Artists) GetDates() (Relation, error) {
	resp, err := http.Get(a.ConcertsURL)
	if err != nil {
		return Relation{}, errors.New("Could not fetch concert data")
	}
	defer resp.Body.Close()

	var result Relation
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return Relation{}, errors.New("Didn't get expected format for concert info json")
	}
	return result, nil
}

func GetRelationData(id int) ([]byte, error) {
	resp, err := http.Get(RELATION_ENDPOINT + "/" + strconv.Itoa(id))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func GetArtistData(id int) ([]byte, error) {
	resp, err := http.Get(ARTIST_ENDPOINT + "/" + strconv.Itoa(id))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func GetArtistsData() ([]byte, error) {
	response, err := http.Get(ARTIST_ENDPOINT)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
