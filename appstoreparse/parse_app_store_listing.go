package appstoreparse

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
	"strings"

	"golang.org/x/net/html"
)

// AppStoreMetadata holds the values parsed from the app listing HTML.
type AppStoreMetadata struct {

	// Value found in the appstore:developer_url <meta> tag.
	DeveloperURL string

	// Value found in the appstore:bundle_id <meta> tag.
	BundleID string

	// Value found in the appstore:store_id <meta> tag.
	StoreID string
}

// ParseAppStorePageHTML accepts a string containing an HTML doc, returning
// a struct containing the parsed "appstore:" meta tag values.
func ParseAppStorePageHTML(htmlContent string) (*AppStoreMetadata, error) {
	result := AppStoreMetadata{}
	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		return nil, fmt.Errorf("Error running HTML parser: %v", err)
	}
	var handleNode func(*html.Node, bool)
	handleNode = func(n *html.Node, parentNoteIsHeadTag bool) {
		// Ignore any elements contained within the <body> tag.
		if n.Type == html.ElementNode && n.Data == "body" {
			return
		}

		if n.Type == html.ElementNode && n.Data == "meta" && parentNoteIsHeadTag {
			var nameAttribute, contentAttribute string
			for _, a := range n.Attr {
				if a.Key == "name" {
					nameAttribute = a.Val
				}
				if a.Key == "content" {
					contentAttribute = a.Val
				}
			}
			if nameAttribute == "appstore:developer_url" {
				result.DeveloperURL = contentAttribute
			}
			if nameAttribute == "appstore:bundle_id" {
				result.BundleID = contentAttribute
			}
			if nameAttribute == "appstore:store_id" {
				result.StoreID = contentAttribute
			}
		}

		// See if this element is a <head> tag to check if <meta> tags are
		// immediate children.
		currentNodeIsHeadTag := n.Type == html.ElementNode && n.Data == "head"
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			handleNode(c, currentNodeIsHeadTag)
		}
	}
	handleNode(doc, false)
	return &result, nil
}
