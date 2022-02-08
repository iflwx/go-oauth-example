package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)


const clientID = "<gitlab client id>"
const clientSecret = "<gitlab client secret>"
const gitlabServer = "<gitlab server url>"

func main() {
	fs := http.FileServer(http.Dir("public"))
	http.Handle("/", fs)

	// We will be using `httpClient` to make external HTTP requests later in our code
	httpClient := http.Client{}

	// Create a new redirect route route
	http.HandleFunc("/oauth/redirect", func(w http.ResponseWriter, r *http.Request) {
		// First, we need to get the value of the `code` query param
		err := r.ParseForm()
		if err != nil {
			fmt.Fprintf(os.Stdout, "could not parse query: %v\n", err)
			w.WriteHeader(http.StatusBadRequest)
		}
		code := r.FormValue("code")
        
        fmt.Fprintf(os.Stdout, "code: %v\n", code)

		// Next, lets for the HTTP request to call the github oauth enpoint
		// to get our access token
		reqURL := fmt.Sprintf("%s/oauth/token?client_id=%s&client_secret=%s&code=%s&grant_type=%s&redirect_uri=%s", gitlabServer, clientID, clientSecret, code, "authorization_code", "http://localhost:8080/oauth/redirect")
		req, err := http.NewRequest(http.MethodPost, reqURL, nil)
		if err != nil {
			fmt.Fprintf(os.Stdout, "could not create HTTP request: %v", err)
			w.WriteHeader(http.StatusBadRequest)
		}
		// We set this header since we want the response
		// as JSON
		req.Header.Set("accept", "application/json")

		// Send out the HTTP request
		res, err := httpClient.Do(req)
		if err != nil {
			fmt.Fprintf(os.Stdout, "could not send HTTP request: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		defer res.Body.Close()
        
		// Parse the request body into the `OAuthAccessResponse` struct
		var t OAuthAccessResponse        
		if err := json.NewDecoder(res.Body).Decode(&t); err != nil {
			fmt.Fprintf(os.Stdout, "could not parse JSON response: %v", err)
			w.WriteHeader(http.StatusBadRequest)
		}
        
        fmt.Fprintf(os.Stdout, "http response body parse succ! access_token: %s, token_type:%s, expires_in:%d, refresh_token:%s, created_at:%d\n", t.AccessToken, t.TokenType, t.ExpiresIn, t.RefreshToken, t.CreatedAt)

		// Finally, send a response to redirect the user to the "welcome" page
		// with the access token
		w.Header().Set("Location", "/welcome.html?access_token="+t.AccessToken)
		w.WriteHeader(http.StatusFound)
	})

	http.ListenAndServe(":8080", nil)
}

type OAuthAccessResponse struct {
	AccessToken string `json:"access_token"`
    TokenType string `json:"token_type"`
    ExpiresIn int `json:"expires_in"`
    RefreshToken string `json:"refresh_token"`
    CreatedAt int `json:"created_at`
}
