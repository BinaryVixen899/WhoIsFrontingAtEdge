package main

import (
	"context"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/fastly/compute-sdk-go/fsthttp"
	"github.com/valyala/fastjson"
)

const BackendName = "origin_0"

// Our entry point
func main() {
	fsthttp.ServeFunc(func(ctx context.Context, w fsthttp.ResponseWriter, r *fsthttp.Request) {

		// Filter requests that have unexpected methods.
		if r.Method != "HEAD" && r.Method != "GET" {
			w.WriteHeader(fsthttp.StatusMethodNotAllowed)
			fmt.Fprintf(w, "This method is not allowed\n")
			println("Yup, we're hitting the method not allowed thingy")
			return
		}

		// If request is to the `/` path...
		if r.URL.Path == "/Fronting" {
			// Create a new request.

			//Should probably do a pointer here
			currentfronter, err := GetCurrentFronter(ctx)
			// Let's parse this into a string unless there's an error
			if err != nil {
				println("Well that didn't work out as expected")
				print(err.Error())
				// TODO: All the error handling logic, such as actually returning 500s
			}

			if currentfronter != nil {
				if r.Header.Get("Accept") == "application/json" {
					//TODO: Calculate the content length
					w.Header().Add("Content-Type", "application/json")

					// Oh right, we have to serialize it into bytes duh
					currentfronterbytes, err := currentfronter.StringBytes()
					if err != nil {
						w.WriteHeader(fsthttp.StatusInternalServerError)
						print(err.Error())
					}
					// WYR: All we need to now is find the size and set content-length to it
					contentlengthint := binary.Size(currentfronterbytes)
					contentlength := string(contentlengthint)
					w.Header().Add("Content-Length", contentlength)
					w.Write(currentfronterbytes)

				} else {
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

			// Log to a Fastly endpoint.
			// NOTE: You will need to import "github.com/fastly/compute-sdk-go/rtlog"
			// for this to work
			// endpoint := rtlog.Open("my_endpoint")
			// fmt.Fprintln(endpoint, "Hello from the edge!")

			return
		} else if r.URL.Path == "/" {

			//We're using this so much we might as well just make it a method
			if r.Method != "HEAD" && r.Method != "GET" {
				w.WriteHeader(fsthttp.StatusMethodNotAllowed)
				fmt.Fprintf(w, "This method is not allowed\n")
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
			<p>"But nothing happened!"</p>
				
			</body>
			</html>`)
			// I'm pretty sure we need a return here

		}
		// TODO: Put a magicarp picture above!

		// Catch all other requests and return a 404.
		w.WriteHeader(fsthttp.StatusNotFound)
		fmt.Fprintf(w, "The page you requested could not be found\n")
	})
}

func GetCurrentFronter(ctx context.Context) (*fastjson.Value, error) {
	//TODO: Return some value for

	// TODO: Read this in from somewhere
	// sysID := "rzwbg "

	req, err := fsthttp.NewRequest("GET", "https://api.pluralkit.me/v2/systems/rzwbg/fronters", nil)
	if err != nil {
		print("Oh no there has been an error when constructing the request!")
		//TODO: Customize this to write a different status
		// w.WriteHeader(fsthttp.StatusBadGateway)
		//fmt.Fprintln(w, err.Error())
		print(err.Error())
		return nil, err
	}
	req.Header.Set("accept", "*/*")
	//TODO: Figure out the default user-agent
	// req.Header.Set("user-agent", "curl/7.84.0")

	// Set the cache to pass
	req.CacheOptions.Pass = true

	resp, err := req.Send(ctx, BackendName)
	if err != nil {
		print("Oh no there has been an error when retrieving the primary fronter from pluralkit!")
		//TODO: Customize this to write a different status
		// TODO: bubble this up
		// w.WriteHeader(fsthttp.StatusBadGateway)
		//fmt.Fprintln(w, err.Error())
		print(err.Error())
		return nil, err
	}

	// Read the backend response body
	// This will read everything into memory so we only want to do it if we know that it is small
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		print("There's been an error when reading the backend response body into json!")
		print(err.Error())
		return nil, err
	}
	responsestring := string(body)
	println(responsestring)
	// Parse some Json
	var p fastjson.Parser
	response_body_json, err := p.ParseBytes(body)
	if err != nil {
		print("There's been an error when parsing the response into JSON!")
		print(err.Error())
		return nil, err
	}

	// Get the Fronting Member JSON
	fronting_member_json := response_body_json.Get("members", "0", "name")

	// Transform that JSON into a string
	return fronting_member_json, nil

}
