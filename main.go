package main

import (
	"context"
	"fmt"

	"github.com/fastly/compute-sdk-go/fsthttp"
	"github.com/starshine-sys/pkgo"
)

// The entry point for your application.
//
// Use this function to define your main request handling logic. It could be
// used to route based on the request properties (such as method or path), send
// the request to a backend, make completely new requests, and/or generate
// synthetic responses.

func main() {
	fsthttp.ServeFunc(func(ctx context.Context, w fsthttp.ResponseWriter, r *fsthttp.Request) {

		// Filter requests that have unexpected methods.
		if r.Method != "HEAD" && r.Method != "GET" {
			w.WriteHeader(fsthttp.StatusMethodNotAllowed)
			fmt.Fprintf(w, "This method is not allowed\n")
			return
		}

		// If request is to the `/` path...
		if r.URL.Path == "/Fronting" {
			// Create a new request.

			// Add request headers.
			// req.Header.Set("Custom-Header", "Welcome to Compute@Edge!")
			// req.Header.Set(
			//   "Another-Custom-Header",
			//   "Recommended reading: https://developer.fastly.com/learning/compute"
			// )

			// Override cache TTL.
			// req.CacheOptions.TTL = 60

			// Forward the request to a backend named "TheOrigin".

			currentfronter := GetCurrentFronter()
			if currentfronter != "" {
				fronter := fmt.Sprintf(`
				<!DOCTYPE html>
				<html lang="en">
				<head>
					<meta charset="UTF-8">
					<meta http-equiv="X-UA-Compatible" content="IE=edge">
					<meta name="viewport" content="width=device-width, initial-scale=1.0">
					<title>Document</title>
				</head>
				<body>
				<p>"The current fronter is: %f!"</p>
					
				</body>
				</html>`, currentfronter)
				fmt.Fprintln(w, fronter)
				w.Header().Set("Content-Type", "text/html; charset=utf-8")

			}

			// Remove response headers.
			// resp.Header.Del("Yet-Another-Custom-Header")

			// Copy all headers from the response.
			// w.Header().Reset(resp.Header.Clone())

			// Log to a Fastly endpoint.
			// NOTE: You will need to import "github.com/fastly/compute-sdk-go/rtlog"
			// for this to work
			// endpoint := rtlog.Open("my_endpoint")
			// fmt.Fprintln(endpoint, "Hello from the edge!")

			// Send a default synthetic response.

			return
		} else if r.URL.Path == "/" {

			if r.Method != "HEAD" && r.Method != "GET" {
				w.WriteHeader(fsthttp.StatusMethodNotAllowed)
				fmt.Fprintf(w, "This method is not allowed\n")
				return
			}
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
			<p>"But nothing happened!"</p>
				
			</body>
			</html>`)

		}
		// TODO: Put a magicarp picture above!

		// Catch all other requests and return a 404.
		w.WriteHeader(fsthttp.StatusNotFound)
		fmt.Fprintf(w, "The page you requested could not be found\n")
	})
}

func GetCurrentFronter() string {
	// DO NOT COMMIT UNTIL YOU VERIFY PRIVACY SETTINGS
	sysID := "REDACTED"

	pk := pkgo.New("")
	front, err := pk.Fronters(sysID)
	if err != nil {
		// Change this later
		print("There has been an error!")
		return ""
	}
	frontingmembername := front.Members[0].Name
	return frontingmembername

}
