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
	"reflect"
	"testing"
)

// Demonstrates parsing a complete, valid HTML5 doc with compliant meta tags.
const validBasicHTML = `
<!DOCTYPE html>
<html lang="en">
	<head>
		<title>
			Example app listing
		</title>
		<meta charset="UTF-8">
		<meta name="appstore:developer_url" content="https://www.path.to/page">
		<meta name="appstore:bundle_id" content="com.example.myapp">
		<meta name="appstore:store_id" content="SKU12345">
	</head>
	<body>
		<h1>
			Example app
		</h1>
	</body>
</html>`

// Confirms that parsing still works correctly with HTML tags and attributes in
// uppercase
const validBasicHTMLUppercase = `
<!DOCTYPE HTML>
<HTML LANG="en">
	<HEAD>
		<TITLE>
			Example app listing
		</TITLE>
		<META CHARSET="UTF-8">
		<META NAME="appstore:developer_url" CONTENT="https://www.path.to/page">
		<META NAME="appstore:bundle_id" CONTENT="com.example.myapp">
		<META NAME="appstore:store_id" CONTENT="SKU12345">
	</HEAD>
	<BODY>
		<H1>
			Example app
		</H1>
	</BODY>
</HTML>`

// Demonstrates parsing a valid HTML5 doc with compliant meta tags, truncated
// at the closing </head> tag.
const validTruncatedHTML = `
<!DOCTYPE html>
<html lang="en">
	<head>
		<title>
			Example app listing
		</title>
		<meta charset="UTF-8">
		<meta name="appstore:developer_url" content="https://www.path.to/page">
		<meta name="appstore:bundle_id" content="com.example.myapp">
		<meta name="appstore:store_id" content="SKU12345">
	</head>`

// Demonstrates <meta> tags in HTML <body> tag which will be ignored, as the
// tags must reside within the <head> tag.
const metaInBodyNonCompliant = `
<!DOCTYPE html>
<html lang="en">
	<head>
		<title>
			Example app listing
		</title>
		<meta charset="UTF-8">
	</head>
	<body>
		<meta name="appstore:developer_url" content="https://www.path.to/page">
		<meta name="appstore:bundle_id" content="com.example.myapp">
		<meta name="appstore:store_id" content="SKU12345">
		<h1>
			Example app
		</h1>
	</body>
</html>`

var testScenarios = []struct {
	input        string
	wantMetadata AppStoreMetadata
}{
	{
		validBasicHTML,
		AppStoreMetadata{
			DeveloperURL: "https://www.path.to/page",
			BundleID:     "com.example.myapp",
			StoreID:      "SKU12345",
		},
	},
	{
		validBasicHTMLUppercase,
		AppStoreMetadata{
			DeveloperURL: "https://www.path.to/page",
			BundleID:     "com.example.myapp",
			StoreID:      "SKU12345",
		},
	},
	{
		validTruncatedHTML,
		AppStoreMetadata{
			DeveloperURL: "https://www.path.to/page",
			BundleID:     "com.example.myapp",
			StoreID:      "SKU12345",
		},
	},
	{
		metaInBodyNonCompliant,
		AppStoreMetadata{},
	},
	{
		"", // Empty HTML document
		AppStoreMetadata{},
	},
}

func TestParseAppStorePageHTML(t *testing.T) {
	for _, scenario := range testScenarios {
		t.Run(scenario.input, func(t *testing.T) {
			haveMetadata, err := ParseAppStorePageHTML(scenario.input)
			if !reflect.DeepEqual(*haveMetadata, scenario.wantMetadata) {
				t.Errorf("have [%v] want [%v]", haveMetadata, scenario.wantMetadata)
			}

			if err != nil {
				t.Error("Received unexpected error")
			}
		})
	}
}
