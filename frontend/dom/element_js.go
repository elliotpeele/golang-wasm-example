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

// Element represents a dom element
type Element struct {
	js.Value
}

// FromValue returns an element from a js.Value
func FromValue(val js.Value) *Element {
	return &Element{
		val,
	}
}

// WithInnerHTML overwrites any text or sub nodes of this element
func (e *Element) WithInnerHTML(content interface{}) *Element {
	e.Set("innerHTML", content)
	return e
}

// GetInnerHTML returns the inner html data
func (e *Element) GetInnerHTML() *Element {
	return FromValue(e.Get("innerHTML"))
}

// WithID sets the id attribute of the element
func (e *Element) WithID(id string) *Element {
	e.Set("id", id)
	return e
}

// Append adds a child element to this element
func (e *Element) Append(e2 *Element) *Element {
	e.Call("appendChild", e2)
	return e
}

// Clone makes a copy of the current element
func (e *Element) Clone() *Element {
	return FromValue(e.Call("cloneNode", true))
}

// GetElementByTagName returns the first element with a matching tag name under the current element
func (e *Element) GetElementByTagName(tag string) *Element {
	return FromValue(e.Call("getElementsByTagName", tag).Index(0))
}

// GetElementsByTagName returns a slice of elements with the given tag name
func (e *Element) GetElementsByTagName(tag string) []*Element {
	items := e.Call("getElementsByTagName", tag)
	elements := make([]*Element, items.Length())
	for i := 0; i < items.Length(); i++ {
		elements[i] = FromValue(items.Index(i))
	}
	return elements
}

// ParentElement returns the parent element of the current element
func (e *Element) ParentElement() *Element {
	return FromValue(e.Get("parentElement"))
}

// AddEventListener attaches an event listener to the current element
func (e *Element) AddEventListener(event string, listener CallbackFunc) {
	e.Call("addEventListener", event, js.FuncOf(listener))
}

// RemoveAllEventListeners removes all event listeners from the current element
func (e *Element) RemoveAllEventListeners() *Element {
	newElement := e.Clone()
	e.ParentElement().ReplaceChild(newElement, e)
	return newElement
}

// ReplaceChild replaces a child node based on the child node
func (e *Element) ReplaceChild(newChild *Element, oldChild *Element) {
	e.Call("replaceChild", newChild, oldChild)
}