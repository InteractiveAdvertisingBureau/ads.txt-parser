# app-ads.txt URL parsing reference implementation

This repository contains a reference implementation of app-ads.txt parsing utilities, along with test cases for verifying compatibility.  These packages are written in the Go software language, although the code and unit tests can be easily ported to another language if needed.

Currently included are two packages:

* `appstoreparse`: an app store listing page HTML parser for extracting HTML `<meta>` tags containing app metadata
* `urlcanonical`: a URL canonicalization package for converting app developer URLs to the corresponding app-ads.txt paths to crawl

The package includes a crawl sample app written in Go, useful for testing compatibility of an individual app store URL from the command line.

## Set up environment

If not already installed, download and install the Go distribution from https://golang.org/

## Running sample app

The sample app can be run against any compliant app store URL.  For example, this command line demonstrates crawling an app in Google Play Store.

```
go run examples/appadstxtcrawl/sample_app.go --app_store_url=https://play.google.com/store/apps/details?id=com.google.android.apps.maps
```

Sample output:

```
Parsed metadata:
  Developer URL: http://maps.google.com/about/
  Bundle ID:     com.google.android.apps.maps
  Store ID:      com.google.android.apps.maps

Derived app-ads.txt URLs:
  Registerable Domain URL: http://google.com/app-ads.txt
  Subdomain URL:           http://maps.google.com/app-ads.txt

```

If the script receives a response other than HTTP 200 from the web server, it will terminate with an error message.

Note: if the app store URL does not provide the required HTML meta tags, the
script will display an empty result such as the following:

```
Parsed metadata:
  Developer URL: 
  Bundle ID:     
  Store ID:      

No developer URL found to parse.
```

## Using a local store HTML page

To run the crawl sample app against an internal sample web server simulating an app listing page HTML file (`sample_app_store.html`), specify the desired port using the `--sample_file_server_port` flag.  The app will crawl the app store URL, output the parse results, and then immediately terminate.  Note: the web server will look for files relative to your shell's current path, so change directory to the local HTML file location prior to running.

```
cd examples/appadstxtcrawl
go run sample_app.go --app_store_url=http://localhost:8081/sample_app_store.html --sample_file_server_port=8081
```

Sample output:

```
Parsed metadata:
  Developer URL: https://sample.path.to/page
  Bundle ID:     com.example.myapp
  Store ID:      SKU12345

Derived app-ads.txt URLs:
  Registerable Domain URL: https://path.to/app-ads.txt
  Subdomain URL:           https://sample.path.to/app-ads.txt
```
