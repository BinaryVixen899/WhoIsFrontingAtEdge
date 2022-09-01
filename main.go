package main

import (
	"context"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/fastly/compute-sdk-go/fsthttp"
	"github.com/fastly/compute-sdk-go/rtlog"
	"github.com/valyala/fastjson"
)

const BackendName = "Pluralkit"

// Our entry point
func main() {
	fsthttp.ServeFunc(func(ctx context.Context, w fsthttp.ResponseWriter, r *fsthttp.Request) {
		endpoint := rtlog.Open("API.Kitsune.Gay")
		// Filter requests that have unexpected methods.
		if r.Method != "HEAD" && r.Method != "GET" {
			w.WriteHeader(fsthttp.StatusMethodNotAllowed)
			fmt.Fprintf(endpoint, "This method is not allowed\n")
			return
		}

		// If request is to the `/` path...
		if r.URL.Path == "/Fronting" || r.URL.Path == "/fronting" {
			// Create a new request.

			currentfronter, err := GetCurrentFronter(&ctx, &w, endpoint)
			// Let's parse this into a string unless there's an error
			if err != nil {
				print("Getting the current fronter didn't work out as expected")
				w.WriteHeader(fsthttp.StatusInternalServerError)
				fmt.Fprintf(endpoint, err.Error())
				// At this point if it fails it's failed at the last step and let's assume this is on us

				// TODO: All the error handling logic, such as actually returning 500s
				// If you wrote 3 AM code, using returns to make sure you don't go on executing other code by accident is a great strategy!
				return
			}

			if currentfronter != nil {
				if r.Header.Get("Accept") == "application/json" {
					//TODO: Calculate the content length
					w.Header().Add("Content-Type", "application/json")

					// Serialize it into bytes
					currentfronterbytes, err := currentfronter.StringBytes()
					if err != nil {
						w.WriteHeader(fsthttp.StatusInternalServerError)
						print("There was an issue serializing into bytes", err.Error())
					}
					contentlengthint := binary.Size(currentfronterbytes)
					contentlength := string(contentlengthint)
					w.Header().Add("Content-Length", contentlength)
					w.Write(currentfronterbytes)
					return

				} else {
					// TODO: Figure out why String() wont' work here
					// TODO: Way better webpage
					currentfronter_string := fmt.Sprintf("%s", currentfronter)
					fmt.Fprintf(w, `
				<!DOCTYPE html>
				<html lang="en">
				<head>
					<meta charset="UTF-8">
					<meta http-equiv="X-UA-Compatible" content="IE=edge">
					<meta name="viewport" content="width=device-width, initial-scale=1.0">
					<title>Document</title>
				</head>
				<body>
				<p>"The current fronter is: %s!"</p>
					
				</body>
				</html>`, currentfronter_string)

				}
			}
			// TODO: Implement logging and send to
			// Log to a Fastly endpoint.

			// fmt.Fprintln(endpoint, "Hello from the edge!")
		} else if r.URL.Path == "/" {

			if r.Method != "HEAD" && r.Method != "GET" {
				w.WriteHeader(fsthttp.StatusMethodNotAllowed)
				fmt.Fprintf(endpoint, "This method is not allowed\n")
				return
			}

			// TODO: set status code here as well
			// Setting a default synthetic response
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			fmt.Fprintln(w, `
			<!DOCTYPE html>
			<html lang="en">
			<head>
				<meta charset="UTF-8">
				<meta http-equiv="X-UA-Compatible" content="IE=edge">
				<meta name="viewport" content="width=device-width, initial-scale=1.0">
				<title>Document</title>
			</head>
			<body>
			<img src="magikarp.png" alt ="Magikarp used Splash, ... but Nothing Happened!"
			<p>"But nothing happened!"</p>
				
			</body>
			</html>`)
			return
		}

		// Catch all other requests and return a 404.
		w.WriteHeader(fsthttp.StatusNotFound)
		fmt.Fprintf(endpoint, "The page you requested could not be found\n")
	})
}

func GetCurrentFronter(ctx *context.Context, w *fsthttp.ResponseWriter, endpoint *rtlog.Endpoint) (*fastjson.Value, error) {
	//TODO: Return some value for

	// TODO: Read this in from somewhere
	// sysID := "rzwbg "
	wrtr := *w
	req, err := fsthttp.NewRequest("GET", "https://api.pluralkit.me/v2/systems/rzwbg/fronters", nil)
	if err != nil {
		print("Oh no there has been an error when constructing the request!")
		//TODO: Customize this to write a different status
		wrtr.WriteHeader(fsthttp.StatusBadGateway)
		fmt.Fprintln(endpoint, err.Error())
		print(err.Error())
		return nil, err
	}
	req.Header.Set("accept", "*/*")
	//TODO: Figure out if a user agent gets sent by the caching layer, we want to be a good netizen after all!
	// req.Header.Set("user-agent", "curl/7.84.0")

	// Set the cache to pass
	req.Header.Set("Host", "api.pluralkit.me")
	req.CacheOptions.Pass = true

	// Print statement logging of the whole thing
	// TODO: Log this instead of print logging it
	fmt.Printf("Body: %v\n", req.Body)
	fmt.Printf("host: %v\n", req.Host)
	fmt.Printf("method: %v\n", req.Method)
	fmt.Printf("TLSInfo: %v\n", req.TLSInfo)
	fmt.Printf("url %v\n", req.URL)
	fmt.Println("Headers")
	for key, value := range req.Header {
		fmt.Println(key, value)
	}

	resp, err := req.Send(*ctx, BackendName)
	if err != nil {
		print("Oh no there has been an error when retrieving the primary fronter from pluralkit!")
		// TODO: Make sure that if the body isn't what we expect we throw an error
		wrtr.WriteHeader(fsthttp.StatusBadGateway)
		// TODO: Log here, we don't care about exposing this to the user
		fmt.Fprintln(endpoint, err.Error())
		print(err.Error())
		return nil, err
	}

	// Read the backend response body
	// This will read everything into memory so we only want to do it if we know that it is small
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		wrtr.WriteHeader(fsthttp.StatusInternalServerError)
		// TODO: Log here, we don't care about exposing this to the user
		print("There's been an error when reading the backend response body into json!")
		print(endpoint, err.Error())
		return nil, err
	}
	responsestring := string(body)
	println(responsestring)
	// Parse some Json
	var p fastjson.Parser
	response_body_json, err := p.ParseBytes(body)
	if err != nil {
		wrtr.WriteHeader(fsthttp.StatusInternalServerError)
		// TODO: Log here, we don't care about exposing this to the user
		print("There's been an error when parsing the response into JSON!")
		print(endpoint, err.Error())
		return nil, err
	}

	// Get the Fronting Member JSON
	fronting_member_json := response_body_json.Get("members", "0", "name")

	// Transform that JSON into a string
	return fronting_member_json, nil

}
