package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
)

const (
	envVarName = "PORTNUMBER"
	portNumber = "8081"
)

func main() {
	greet()
	chosenPort := portNumber
	if len(os.Args) > 1 {
		_, portIsANumber := strconv.Atoi(os.Args[1])
		if portIsANumber != nil {
			fmt.Printf("error parsing port number: %s\n", portIsANumber.Error())
			return
		}
		chosenPort = os.Args[1]
	} else {
		portEnv, varExists := os.LookupEnv(envVarName)
		if varExists {
			_, portIsANumber := strconv.Atoi(portEnv)
			if portIsANumber != nil {
				fmt.Printf("error parsing port number: %s\n", portIsANumber.Error())
				return
			}
			chosenPort = portEnv
		}
	}
	http.HandleFunc("/", AllowMethod(HandlerHome, []string{http.MethodGet}))
	http.HandleFunc("/style/{file}", AllowMethod(StaticPublicRessourceServer, []string{http.MethodGet}))
	http.HandleFunc("/js/{file}", AllowMethod(StaticPublicRessourceServer, []string{http.MethodGet}))
	http.HandleFunc("/img/{file}", AllowMethod(StaticPublicRessourceServer, []string{http.MethodGet}))

	// Endpoints
	http.HandleFunc("/artist/", AllowMethod(ArtistEndpoint, []string{http.MethodGet}))
	http.HandleFunc("/artist/{id}", AllowMethod(HandlerArtistPage, []string{http.MethodGet}))
	//http.HandleFunc("/artist/search/{query}", AllowMethod(ArtistSearch, []string{http.MethodGet}))

	http.HandleFunc("/relation/{id}", AllowMethod(RelationEndpoint, []string{http.MethodGet}))
	log.Fatal(http.ListenAndServe(":"+chosenPort, nil))
}

func greet() {
	fmt.Println(`🎵🎵Hi🎵🎵`)
	fmt.Println("Listening on port", portNumber)
	fmt.Println("------------")
}
