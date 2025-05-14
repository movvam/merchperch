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
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"

	spotifyauth "github.com/zmb3/spotify/v2/auth"

	"github.com/zmb3/spotify/v2"

	"github.com/machinebox/graphql"
)

// redirectURI is the OAuth redirect URI for the application.
// You must register an application at Spotify's developer portal
// and enter this value.
const redirectURI = "http://localhost:8080/callback"

type ShopArtist struct {
	NAME             string           `json:"name"`
	ID               string           `json:"id"`
	SPOTIFY_SHOP_URL string           `json:"spotify_shop_url"`
	SHOPIFY_URL      string           `json:"shopify_url"`
	PHOTO_URL        string           `json:"photo_url"`
	POPULARITY       spotify.Numeric  `json:"popularity"`
	PRODUCTS         []ShopifyProduct `json:"products"`
}

type ProductImage struct {
	URL string `json:"url"`
}

type ProductVariantPrice struct {
	Amount       string `json:"amount"`
	CurrencyCode string `json:"currencyCode"`
}

type ProductVariant struct {
	Price ProductVariantPrice `json:"price"`
}

type ShopifyProductWithArtist struct {
	PRODUCT ShopifyProduct `json:"product"`
	ARTIST  string         `json:"artist"`
}

type ShopifyProduct struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	ProductType string `json:"productType"`
	Handle      string `json:"handle"`

	Images struct {
		Edges []struct {
			Node ProductImage `json:"node"`
		} `json:"edges"`
	} `json:"images"`

	Variants struct {
		Edges []struct {
			Node ProductVariant `json:"node"`
		} `json:"edges"`
	} `json:"variants"`
}

type ShopifyProductEdge struct {
	Cursor string         `json:"cursor"`
	Node   ShopifyProduct `json:"node"`
}

type TopLevelResponse struct {
	Data ShopifyProductResponse `json:"data"`
}

type ShopifyProductResponse struct {
	Products struct {
		PageInfo struct {
			HasNextPage bool `json:"hasNextPage"`
		} `json:"pageInfo"`
		Edges []ShopifyProductEdge `json:"edges"`
	} `json:"products"`
}

var (
	auth                 = spotifyauth.New(spotifyauth.WithRedirectURL(redirectURI), spotifyauth.WithScopes(spotifyauth.ScopeUserFollowRead, spotifyauth.ScopeUserLibraryRead))
	ch                   = make(chan *spotify.Client)
	state                = "abc123"
	artistsThatHaveShops []spotify.FullArtist
	shopArtists          []ShopArtist
	shopKeys             = make(map[string]string)
	products             []ShopifyProductWithArtist
)

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

	// print out number of followed artists
	fmt.Println(followedArtistsResp.Total)

	// Figure out which artists have connected their Shopify shops with Spotify
	wg := sync.WaitGroup{}
	wg.Add(len(followedArtists))
	for _, artist := range followedArtists {
		go makeShopExistenceCheckRequest(&wg, artist)
	}

	wg.Wait()

	fmt.Println("Artists with shops:", artistsThatHaveShops)
	writeShops()
}

func queryShopifyStorefront(ctx context.Context, artistId string, storeURL string, token string) ([]ShopifyProductWithArtist, error) {
	client := graphql.NewClient(storeURL, graphql.WithHTTPClient(&http.Client{}))

	var allProducts []ShopifyProductWithArtist
	var cursor *string

	for {
		// Build query with optional cursor
		var b strings.Builder
		b.WriteString(`query { products(first: 250`)
		if cursor != nil {
			b.WriteString(fmt.Sprintf(`, after: "%s"`, *cursor))
		}
		b.WriteString(`) {
			pageInfo { hasNextPage }
			edges {
				cursor
				node {
					id
					title
					productType
					handle
					description
					images(first: 250) {
						edges {
							node {
								url
							}
						}
					}
					variants(first: 250) {
						edges {
							node {
								price {
									amount
									currencyCode
								}
							}
						}
					}
				}
			}
		}}`)
		query := b.String()

		// Create request
		req := graphql.NewRequest(query)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Shopify-Storefront-Access-Token", token)

		// Run query
		var resp ShopifyProductResponse
		if err := client.Run(ctx, req, &resp); err != nil {
			return nil, err
		}

		// Append titles and get next cursor
		edges := resp.Products.Edges
		for _, edge := range edges {
			allProducts = append(allProducts, ShopifyProductWithArtist{edge.Node, artistId})
		}
		fmt.Println("GOT PRODUCTS:", resp)
		fmt.Println("GOT PRODUCTS:", allProducts)
		// fmt.Println("HAS NEXT:", resp.Data.Products.PageInfo.HasNextPage)
		break
		// if !resp.Data.Products.PageInfo.HasNextPage {
		// 	break
		// }
		// lastCursor := edges[len(edges)-1].Cursor
		// cursor = &lastCursor
	}

	return allProducts, nil
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
	shopArtist.SPOTIFY_SHOP_URL = shopUrl
	shopArtist.PHOTO_URL = photoUrl
	shopArtist.POPULARITY = artist.Popularity
	return shopArtist
}

func writeShops() {
	// Write to JSON file
	jsonToWrite, _ := json.Marshal(shopArtists)
	err := os.WriteFile("../merchPerchUi/src/data/artistShops.json", jsonToWrite, 0644)
	if err != nil {
		log.Fatal(err)
	}

	jsonToWrite, _ = json.Marshal(products)
	err = os.WriteFile("../merchPerchUi/src/data/products.json", jsonToWrite, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func makeShopExistenceCheckRequest(wg *sync.WaitGroup, artist spotify.FullArtist) {
	resp, err := http.Get(fmt.Sprintf("https://generic.wg.spotify.com/shopify-merch/public/v0/artist/%v/storefront-token", artist.ID.String()))
	if err != nil {
		log.Fatalln("HTTP error:", err)
	}

	if resp.StatusCode == 200 {

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln("Failed to read response body:", err)
		}

		// Extract 32-char hex token
		tokenRe := regexp.MustCompile(`[a-f0-9]{32}`)
		token := tokenRe.Find(body)
		log.Println("token", string(token))
		shopKeys[artist.ID.String()] = string(token)

		// Extract the first .myshopify.com domain
		re := regexp.MustCompile(`[a-zA-Z0-9\-]+\.myshopify\.com`)
		storefrontURL := re.Find(body)
		if storefrontURL == nil {
			log.Println("No Shopify domain found")
			return
		}

		fmt.Println("Shopify URL:", string(storefrontURL))

		allProducts, err := queryShopifyStorefront(context.Background(), artist.ID.String(),
			fmt.Sprintf("https://%s/api/2025-04/graphql.json", string(storefrontURL)),
			string(token))
		if err != nil {
			// handle error
			log.Fatalln(fmt.Sprintf("Failed to query storefront %s response body: %s", storefrontURL, err))
		}
		var artistProducts []ShopifyProduct
		for _, shopifyProductWithArist := range allProducts {
			artistProducts = append(artistProducts, shopifyProductWithArist.PRODUCT)
			products = append(products, shopifyProductWithArist)
		}

		shopArtists = append(shopArtists, ShopArtist{
			ID:               artist.ID.String(),
			NAME:             artist.Name,
			SPOTIFY_SHOP_URL: fmt.Sprintf("https://shop.spotify.com/en/artist/%v/store", artist.ID.String()),
			SHOPIFY_URL:      string(storefrontURL),
			PHOTO_URL:        artist.Images[0].URL,
			POPULARITY:       artist.Popularity,
			PRODUCTS:         artistProducts,
		})
		artistsThatHaveShops = append(artistsThatHaveShops, artist)
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
