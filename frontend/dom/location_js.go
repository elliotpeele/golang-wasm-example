// Copyright (c) 2019 Elliot Peele <elliot@bentlogic.net>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package dom

import "syscall/js"

// Location represents the javascript location object
// https://www.w3docs.com/snippets/javascript/how-to-get-current-url-in-javascript.html
type Location struct {
	Href     string
	Host     string
	Hostname string
	Protocol string
	Pathname string
	Search   string
	Hash     string
}

// ParseLocation parses the javascript location object into a Location struct instance
func ParseLocation() (*Location, error) {
	loc := js.Global().Get("location")
	return &Location{
		Href:     loc.Get("href").String(),
		Host:     loc.Get("host").String(),
		Hostname: loc.Get("hostname").String(),
		Protocol: loc.Get("protocol").String(),
		Pathname: loc.Get("pathname").String(),
		Search:   loc.Get("search").String(),
		Hash:     loc.Get("hash").String(),
	}, nil
}
