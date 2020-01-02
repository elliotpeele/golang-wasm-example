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

package cmd

import (
	"io"

	"github.com/op/go-logging"
)

func setupLogging(out io.Writer, debug bool) {
	backend := logging.NewBackendFormatter(
		logging.NewLogBackend(out, "", 0),
		logging.MustStringFormatter("%{color}%{time:15:04:05.000} %{shortpkg}.%{shortfunc} ▶ %{level:.8s} %{id:03x}%{color:reset} %{message}"),
	)
	level := logging.AddModuleLevel(backend)
	if debug {
		level.SetLevel(logging.DEBUG, "")
	} else {
		level.SetLevel(logging.INFO, "")
	}
	logging.SetBackend(backend)
}
