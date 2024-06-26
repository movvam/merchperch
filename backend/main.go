// This example demonstrates how to authenticate with Spotify using the authorization code flow.
// In order to run this example yourself, you'll need to:
//
//  1. Register an application at: https://developer.spotify.com/my-applications/
//     - Use "http://localhost:8080/callback" as the redirect URI
//  2. Set the SPOTIFY_ID environment variable to the client ID you got in step 1.
//  3. Set the SPOTIFY_SECRET environment variable to the client secret from step 1.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	spotifyauth "github.com/zmb3/spotify/v2/auth"

	"github.com/zmb3/spotify/v2"
)

// redirectURI is the OAuth redirect URI for the application.
// You must register an application at Spotify's developer portal
// and enter this value.
const redirectURI = "http://localhost:8080/callback"

var (
	auth             = spotifyauth.New(spotifyauth.WithRedirectURL(redirectURI), spotifyauth.WithScopes(spotifyauth.ScopeUserFollowRead, spotifyauth.ScopeUserLibraryRead))
	ch               = make(chan *spotify.Client)
	state            = "abc123"
	artistsWithShops []spotify.FullArtist
)

type ShopArtist struct {
	NAME       string          `json:"name"`
	ID         string          `json:"id"`
	SHOP_URL   string          `json:"shop_url"`
	PHOTO_URL  string          `json:"photo_url"`
	POPULARITY spotify.Numeric `json:"popularity"`
}

func main() {
	// first start an HTTP server
	http.HandleFunc("/callback", completeAuth)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Got request for:", r.URL.String())
	})
	go func() {
		err := http.ListenAndServe(":8080", nil)
		if err != nil {
			log.Fatal(err)
		}
	}()

	url := auth.AuthURL(state)
	fmt.Println("Please log in to Spotify by visiting the following page in your browser:", url)

	// wait for auth to complete
	client := <-ch

	// Get list of artists the user follows
	numArtistsGot := 0
	followedArtistsResp, err := client.CurrentUsersFollowedArtists(context.Background(), spotify.Limit(50))
	if err != nil {
		log.Fatal(err)
	}
	followedArtists := followedArtistsResp.Artists
	numArtistsGot += len(followedArtists)
	for numArtistsGot < int(followedArtistsResp.Total) {
		nextArtists, err := client.CurrentUsersFollowedArtists(
			context.Background(),
			spotify.Limit(50),
			spotify.After(followedArtists[numArtistsGot-1].ID.String()))
		followedArtists = append(followedArtists, nextArtists.Artists...)
		numArtistsGot = len(followedArtists)
		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println(followedArtistsResp.Total)
	fmt.Println(numArtistsGot)

	// Figure out which artists have connected their Shopify shops with Spotify
	wg := sync.WaitGroup{}
	wg.Add(len(followedArtists))
	for _, artist := range followedArtists {
		go makeShopExistenceCheckRequest(&wg, artist)
	}

	wg.Wait()

	fmt.Println("Artists with shops:", artistsWithShops)
	writeShops()
}

func buildShopArtist(artist spotify.FullArtist) ShopArtist {
	var shopArtist ShopArtist
	var dataStr string

	id := artist.ID.String()
	name := artist.Name
	photoUrl := ""
	if len(artist.Images) > 0 {
		photoUrl = artist.Images[0].URL
	}
	shopUrl := fmt.Sprintf("https://shop.spotify.com/en/artist/%v/store", id)
	dataStr += fmt.Sprintf("%v %v %v\n", name, id, shopUrl)

	shopArtist.ID = id
	shopArtist.NAME = name
	shopArtist.SHOP_URL = shopUrl
	shopArtist.PHOTO_URL = photoUrl
	shopArtist.POPULARITY = artist.Popularity
	return shopArtist
}

func writeShops() {
	var dataStr string
	var artistsJson []ShopArtist

	// Build the list of ShopArtist objects
	for _, artist := range artistsWithShops {
		shopArtist := buildShopArtist(artist)
		artistsJson = append(artistsJson, shopArtist)
	}
	// Write to JSON file
	data := []byte(dataStr)
	err := os.WriteFile("followed_artists.txt", data, 0644)
	if err != nil {
		log.Panic(err)
	}
	jsonToWrite, _ := json.Marshal(artistsJson)
	err = os.WriteFile("output.json", jsonToWrite, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func makeShopExistenceCheckRequest(wg *sync.WaitGroup, artist spotify.FullArtist) {
	resp, err := http.Get(fmt.Sprintf("https://generic.wg.spotify.com/shopify-merch/public/v0/artist/%v/storefront-token", artist.ID.String()))
	if err != nil {
		log.Fatalln(err)
	}
	if resp.StatusCode == 200 {
		artistsWithShops = append(artistsWithShops, artist)
	}
	wg.Done()
}

func completeAuth(w http.ResponseWriter, r *http.Request) {
	tok, err := auth.Token(r.Context(), state, r)
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusForbidden)
		log.Fatal(err)
	}
	if st := r.FormValue("state"); st != state {
		http.NotFound(w, r)
		log.Fatalf("State mismatch: %s != %s\n", st, state)
	}

	// use the token to get an authenticated client
	client := spotify.New(auth.Client(r.Context(), tok))
	fmt.Fprintf(w, "Login Completed!")
	ch <- client
}
