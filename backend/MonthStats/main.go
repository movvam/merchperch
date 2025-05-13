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
	"fmt"
	"log"
	"net/http"
	"time"

	spotifyauth "github.com/zmb3/spotify/v2/auth"

	"github.com/zmb3/spotify/v2"
)

// redirectURI is the OAuth redirect URI for the application.
// You must register an application at Spotify's developer portal
// and enter this value.
const redirectURI = "http://localhost:8080/callback"

var scopes = []string{spotifyauth.ScopeUserFollowRead, spotifyauth.ScopeUserLibraryRead, spotifyauth.ScopeUserTopRead, spotifyauth.ScopeUserReadRecentlyPlayed}

var (
	auth             = spotifyauth.New(spotifyauth.WithRedirectURL(redirectURI), spotifyauth.WithScopes(scopes...))
	ch               = make(chan *spotify.Client)
	state            = "abc123"
	artistsWithShops []spotify.FullArtist
)

// type ShopArtist struct {
// 	NAME       string          `json:"name"`
// 	ID         string          `json:"id"`
// 	SHOP_URL   string          `json:"shop_url"`
// 	PHOTO_URL  string          `json:"photo_url"`
// 	POPULARITY spotify.Numeric `json:"popularity"`
// }

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
	// numArtistsGot := 0
	currentTopArtistsResp, err := client.CurrentUsersTopArtists(context.Background(), spotify.Timerange(spotify.Range("short_term")))
	if err != nil {
		log.Fatal(err)
	}
	currentTopArtists := currentTopArtistsResp.Artists

	currentTopTracksResp, err := client.CurrentUsersTopTracks(context.Background(), spotify.Timerange(spotify.Range("short_term")))
	if err != nil {
		log.Fatal(err)
	}

	for _, topTrack := range currentTopTracksResp.Tracks {
		fmt.Println("track: ", topTrack.Name, "artist: ", topTrack.Artists)
	}

	currentUsersRecentlyPlayedTracks, err := client.PlayerRecentlyPlayedOpt(context.Background(), &spotify.RecentlyPlayedOptions{AfterEpochMs: getOneWeekAgoTimeStamp()})
	if err != nil {
		log.Fatal(err)
	}

	printArtists(currentTopArtists, currentUsersRecentlyPlayedTracks)

	// fmt.Println(numArtistsGot)

	// // Figure out which artists have connected their Shopify shops with Spotify
	// wg := sync.WaitGroup{}
	// wg.Add(len(followedArtists))
	// for _, artist := range followedArtists {
	// 	go makeShopExistenceCheckRequest(&wg, artist)
	// }

	// wg.Wait()

	// fmt.Println("Artists with shops:", artistsWithShops)
	// writeShops()
}

// func buildShopArtist(artist spotify.FullArtist) ShopArtist {
// 	var shopArtist ShopArtist
// 	var dataStr string

// 	id := artist.ID.String()
// 	name := artist.Name
// 	photoUrl := ""
// 	if len(artist.Images) > 0 {
// 		photoUrl = artist.Images[0].URL
// 	}
// 	shopUrl := fmt.Sprintf("https://shop.spotify.com/en/artist/%v/store", id)
// 	dataStr += fmt.Sprintf("%v %v %v\n", name, id, shopUrl)

// 	shopArtist.ID = id
// 	shopArtist.NAME = name
// 	shopArtist.SHOP_URL = shopUrl
// 	shopArtist.PHOTO_URL = photoUrl
// 	shopArtist.POPULARITY = artist.Popularity
// 	return shopArtist
// }

func getOneWeekAgoTimeStamp() int64 {
	oneWeekAgo := time.Now().AddDate(0, 0, -7)
	return oneWeekAgo.Unix()
}

func printArtists(artists []spotify.FullArtist, tracks []spotify.RecentlyPlayedItem) {
	for _, artist := range artists {
		fmt.Println("top artist: ", artist.Name, "genre: ", artist.Genres)
	}

	freqMap := make(map[string]int)
	for _, element := range tracks {
		freqMap[string(element.Track.ID.String())]++
	}

	// for _, track := range tracks {
	// 	fmt.Println("track: ", track.Track.Name, "artist: ", track.Track.Artists[0].Name, "plays: ", freqMap[string(track.Track.ID.String())])
	// }
	// tracks.frequency()
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
