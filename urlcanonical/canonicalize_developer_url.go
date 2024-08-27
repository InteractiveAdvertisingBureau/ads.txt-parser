// Package urlcanonical provides URL parsing features for ads.txt.
package urlcanonical

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
	"fmt"
	"net/url"
	"strings"

	"github.com/weppos/publicsuffix-go/net/publicsuffix"
)

// TranslateDeveloperURLToAppAdsTxtPaths returns a pair of strings indicating
// the app-ads.txt paths to crawl for a given developer URL.  If the URL
// includes a permissable subdomain, that value will be returned as the second
// return parameter; otherwise, that value will be an empty string.  Any error
// occurring on the parse attempt will be returned in the third parameter.
func TranslateDeveloperURLToAppAdsTxtPaths(input string) (string, string, error) {
	input = strings.ToLower(strings.TrimSpace(input))
	url, err := url.Parse(input)
	if err != nil {
		return "", "", fmt.Errorf("Unable to parse URL [%s]: %v", input, err)
	}

	originalURLScheme := url.Scheme
	if originalURLScheme != "https" && originalURLScheme != "http" {
		return "", "", fmt.Errorf("URL does not start with https/http: %s", input)
	}

	fullHostname := url.Hostname()
	if fullHostname == "" {
		return "", "", fmt.Errorf("Hostname not found in URL: %s", input)
	}

	registerableDomainPortion, err := publicsuffix.EffectiveTLDPlusOne(fullHostname)
	if err != nil {
		return "", "", fmt.Errorf("Unable to extract registerable domain from URL: %v", err)
	}

	var fullSubdomainPortion string
	if fullHostname != registerableDomainPortion {
		fullSubdomainPortion = fullHostname[0 : len(fullHostname)-len(registerableDomainPortion)-1]
	}

	// Retain only the last subdomain fragment from the full subdomain portion.
	lastIndexOfDot := strings.LastIndex(fullSubdomainPortion, ".")
	var relevantSubdomainPortion string
	if lastIndexOfDot == -1 {
		relevantSubdomainPortion = fullSubdomainPortion
	} else {
		relevantSubdomainPortion = fullSubdomainPortion[lastIndexOfDot+1 : len(fullSubdomainPortion)]
	}

	// Drop the "www" or "m" subdomain if present.
	if relevantSubdomainPortion == "www" || relevantSubdomainPortion == "m" {
		relevantSubdomainPortion = ""
	}

	registerableDomainResultURL := originalURLScheme + "://" + registerableDomainPortion + "/app-ads.txt"

	// Only supply a subdomain result URL if a valid subdomain was found.
	var subdomainResultURL string
	if len(relevantSubdomainPortion) > 0 {
		subdomainResultURL = originalURLScheme + "://" + relevantSubdomainPortion + "." + registerableDomainPortion + "/app-ads.txt"
	}

	return registerableDomainResultURL, subdomainResultURL, nil
}
