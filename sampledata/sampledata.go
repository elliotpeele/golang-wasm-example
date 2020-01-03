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

package sampledata

import (
	"bytes"
	"encoding/json"
	"path"

	"github.com/elliotpeele/golang-wasm-example/sampledata/models"
)

//go:generate go build -o gen ./datagen
//go:generate mkdir -p generated
//go:generate ./gen --users 1000 --projects 450
//go:generate go-bindata --prefix "sampledata/generated" -pkg $GOPACKAGE -nometadata -nomemcopy -md5checksum -o generated_data.go generated/...

// Users returns a list of user instances
func Users() ([]models.User, error) {
	var users []models.User
	return users, load("users.json", &users)
}

// Projects returns a  list of project instances
func Projects() ([]models.Project, error) {
	var projects []models.Project
	return projects, load("projects.json", &projects)
}

func load(fn string, obj interface{}) error {
	name := path.Join("generated", fn)
	blob, err := Asset(name)
	if err != nil {
		return err
	}
	r := bytes.NewReader(blob)
	return json.NewDecoder(r).Decode(obj)
}
