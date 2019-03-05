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
	"strings"
	"testing"
)

var testScenarios = []struct {
	input               string // Input developer URL as found in app store
	wantRegDomainOutput string // "Registerable" domain app-ads.txt URL
	wantSubdomainOutput string // Subdomain app-ads.txt URL
	wantErrorPrefix     string // Error message, if any
}{
	// Test cases from app-ads.txt spec.
	{
		"https://example.com/test",
		"https://example.com/app-ads.txt",
		"",
		"",
	},
	{
		"https://www.example.com/test",
		"https://example.com/app-ads.txt",
		"",
		"",
	},
	{
		"https://m.example.com/test",
		"https://example.com/app-ads.txt",
		"",
		"",
	},
	{
		"https://subdomain.example.com/test",
		"https://example.com/app-ads.txt",
		"https://subdomain.example.com/app-ads.txt",
		"",
	},
	{
		"https://another.subdomain.example.com/test",
		"https://example.com/app-ads.txt",
		"https://subdomain.example.com/app-ads.txt",
		"",
	},
	{
		"https://yet.another.subdomain.example.com/test",
		"https://example.com/app-ads.txt",
		"https://subdomain.example.com/app-ads.txt",
		"",
	},
	{
		"https://subdomain.www.example.com/test",
		"https://example.com/app-ads.txt",
		"",
		"",
	},
	{
		"https://another.subdomain.example.co.uk/test",
		"https://example.co.uk/app-ads.txt",
		"https://subdomain.example.co.uk/app-ads.txt",
		"",
	},
	{
		"https://another.subdomain.example.uk/test",
		"https://example.uk/app-ads.txt",
		"https://subdomain.example.uk/app-ads.txt",
		"",
	},

	// Test cases for error conditions from invalid input.
	{
		"",
		"",
		"",
		"URL does not start with https/http: ",
	},
	{
		"https://example/test",
		"",
		"",
		"Unable to extract registerable domain from URL:",
	},
	{
		"malformed",
		"",
		"",
		"URL does not start with https/http: malformed",
	},
	{
		"example.com/test",
		"",
		"",
		"URL does not start with https/http: example.com/test",
	},
	{
		"ftp://example.com/test",
		"",
		"",
		"URL does not start with https/http: ftp://example.com/test",
	},

	// Test additional valid variants.
	{
		"https://subdomain.example.com",
		"https://example.com/app-ads.txt",
		"https://subdomain.example.com/app-ads.txt",
		"",
	},
	{
		"https://subdomain.example.com/test?a=b&c=d#description",
		"https://example.com/app-ads.txt",
		"https://subdomain.example.com/app-ads.txt",
		"",
	},

	//
	{
		"https://subdomain.example.com/test?a=b&c=d#description",
		"https://example.com/app-ads.txt",
		"https://subdomain.example.com/app-ads.txt",
		"",
	},

	// Test case normalization.
	{
		"https://SubDomain.Example.Com",
		"https://example.com/app-ads.txt",
		"https://subdomain.example.com/app-ads.txt",
		"",
	},

	// Test trimming.
	{
		"  https://subdomain.example.com  ",
		"https://example.com/app-ads.txt",
		"https://subdomain.example.com/app-ads.txt",
		"",
	},

	// Test Unicode encoding.
	{
		"https://北京.點看.cn/test",
		"https://點看.cn/app-ads.txt",
		"https://北京.點看.cn/app-ads.txt",
		"",
	},

	// Test Punycode encoding.
	{
		"https://xn--1lq90i.xn--c1yn36f.cn/test",
		"https://xn--c1yn36f.cn/app-ads.txt",
		"https://xn--1lq90i.xn--c1yn36f.cn/app-ads.txt",
		"",
	},

	// Test ignoring authentication details in URL.
	{
		"https://user:pass@subdomain.example.com",
		"https://example.com/app-ads.txt",
		"https://subdomain.example.com/app-ads.txt",
		"",
	},
}

func TestTranslateDeveloperURLToAppAdsTxtPaths(t *testing.T) {
	for _, scenario := range testScenarios {
		t.Run(scenario.input, func(t *testing.T) {
			reg, sub, err := TranslateDeveloperURLToAppAdsTxtPaths(scenario.input)
			if reg != scenario.wantRegDomainOutput {
				t.Errorf("have [%s] want [%s]", reg, scenario.wantRegDomainOutput)
			}
			if sub != scenario.wantSubdomainOutput {
				t.Errorf("have [%s] want [%s]", sub, scenario.wantSubdomainOutput)
			}

			if err == nil {
				if scenario.wantErrorPrefix != "" {
					t.Errorf("Expected error with prefix [%s].", scenario.wantErrorPrefix)
				}
			} else {
				if scenario.wantErrorPrefix == "" {
					t.Errorf("Expected no error, but found [%v].", err)
				} else if !strings.HasPrefix(err.Error(), scenario.wantErrorPrefix) {
					t.Errorf("have [%v] want prefix [%s]", err, scenario.wantErrorPrefix)
				}
			}
		})
	}
}
