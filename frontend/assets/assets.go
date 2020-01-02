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

package assets

import (
	"bytes"
	"fmt"
	"net/http"
	"path"
	"time"

	"github.com/op/go-logging"
)

//go:generate go-bindata --prefix "assets/dist" -pkg $GOPACKAGE -nometadata -nomemcopy -md5checksum -o assets_generated.go dist/...

var logger = logging.MustGetLogger("golang-wasm-example.frontend.assets")

// ServeFile serves compiled asset files
func ServeFile(w http.ResponseWriter, req *http.Request, html string) {
	logger.Debugf("serving %s", html)
	name := path.Join("dist", html)
	blob, err := Asset(name)
	if err != nil {
		http.NotFound(w, req)
		return
	}
	info, err := AssetInfo(name)
	if err != nil {
		http.NotFound(w, req)
		return
	}
	base := path.Base(name)
	// TODO: insert cache control headers
	if size := info.Size(); size != 0 {
		w.Header().Set("Content-Length", fmt.Sprintf("%d", size))
	}
	http.ServeContent(w, req, base, time.Time{}, bytes.NewReader(blob))
}
