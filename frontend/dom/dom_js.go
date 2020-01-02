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

// CallbackFunc defines the signature of a callback
type CallbackFunc func(this js.Value, args []js.Value) interface{}

// RegisterFunc registers a java script function
func RegisterFunc(name string, fn CallbackFunc) {
	js.Global().Set(name, js.FuncOf(fn))
}

// Document returns the dom document node
func Document() *Element {
	val := js.Global().Get("document")
	return &Element{
		val,
	}
}

// GetElementByID calls into the dom to getElementById
// https://www.w3schools.com/jsref/met_document_getelementbyid.asp
func GetElementByID(name string) *Element {
	val := Document().Call("getElementById", name)
	return &Element{
		val,
	}
}

// CreateElement calls into the dom to create a new element
func CreateElement(name string) *Element {
	val := Document().Call("createElement", name)
	return &Element{
		val,
	}
}
