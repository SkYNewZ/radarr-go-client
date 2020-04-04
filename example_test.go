package radarr_test

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/SkYNewZ/radarr"
)

// Instanciate a standard client
func ExampleNew_basic() {
	client, err := radarr.New("https://my.radarr-instance.fr", "radarr-api-key", nil)
	if err != nil {
		log.Fatalln(err)
	}

	movie, err := client.Movies.GetMovie(217)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("%s", movie.Title)

	// Output:
	// Frozen II
}

// Instanciate a client with a custom HTTP client
func ExampleNew_advanced() {
	client, err := radarr.New("https://my.radarr-instance.fr", "radarr-api-key", &http.Client{
		Timeout: time.Second * 10,
	})
	if err != nil {
		log.Fatalln(err)
	}

	movie, err := client.Movies.GetMovie(217)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("%s", movie.Title)

	// Output:
	// Frozen II
}
