package main

/*
 * Copyright 2018 IAB Tech Lab & OpenRTB Group
 * Copyright 2018 Google LLC
 *
 * Author: Curtis Light, Google
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strconv"

	"github.com/cmlight/adstxtparse/appstoreparse"
	"github.com/cmlight/adstxtparse/urlcanonical"
)

const metadataFormat = `Parsed metadata:
  Developer URL: %s
  Bundle ID:     %s
  Store ID:      %s

`

const urlFormat = `Derived app-ads.txt URLs:
  Registerable Domain URL: %s
  Subdomain URL:           %s

`

func startSampleServer(sampleFilePort int) {
	http.HandleFunc("/sample_app_store.html", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "sample_app_store.html")
	})

	// Synchronously creating listener to ensure server ready at function return.
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(sampleFilePort))
	if err != nil {
		log.Fatalf("Could not listen on port for HTTP server: %v", err)
	}
	go http.Serve(listener, nil)
}

// This very simple sample app fetches the HTML content from a single app store
// URL passed by flag from the command line.  The app may be helpful in
// evaluating whether an authorized seller verifier can correctly parse the app
// store listing metadata according to the app-ads.txt specification.
func main() {
	sampleFilePort := flag.Int("sample_file_server_port", 0,
		"If non-zero, runs a sample HTTP server that hosts sample_app_store.html on localhost.")
	appStoreURL := flag.String("app_store_url", "", "App store listing URL")

	flag.Parse()

	if *sampleFilePort > 0 {
		startSampleServer(*sampleFilePort)
	}

	if *appStoreURL == "" {
		log.Fatalln("Flag --app_store_url required")
	}

	resp, err := http.Get(*appStoreURL)
	if err != nil {
		log.Fatalf("Error fetching app store listing: %v\n", err)
	}

	if resp.StatusCode != 200 {
		log.Fatalf("Did not receive HTTP 200 response, instead was %s.", resp.Status)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading app store listing: %v\n", err)
	}

	meta, err := appstoreparse.ParseAppStorePageHTML(string(body))
	if err != nil {
		log.Fatalf("Error parsing app store listing: %v\n", err)
	}

	fmt.Printf(metadataFormat, meta.DeveloperURL, meta.BundleID, meta.StoreID)

	if meta.DeveloperURL == "" {
		fmt.Println("No developer URL found to parse.")
	} else {
		reg, sub, err := urlcanonical.TranslateDeveloperURLToAppAdsTxtPaths(meta.DeveloperURL)
		if err != nil {
			log.Fatalf("Error translating developer URLs: %v\n", err)
		}

		fmt.Printf(urlFormat, reg, sub)
	}
}
