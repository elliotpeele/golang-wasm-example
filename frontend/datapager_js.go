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

package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"syscall/js"

	"github.com/elliotpeele/golang-wasm-example/frontend/dom"
)

const (
	pageSize        = 20
	selectablePages = 5
)

type datapager struct {
	headers       []string
	items         [][]interface{}
	curPage       int
	filteredItems []int
}

func (d *datapager) renderPage(offset int) {
	headerTmpl := dom.GetElementByID("headerNode")
	rowTmpl := dom.GetElementByID("rowNode")

	// Clear header and content
	header := dom.GetElementByID("tableHeader").WithInnerHTML("")
	content := dom.GetElementByID("tableContent").WithInnerHTML("")

	// Populate header
	tr := dom.CreateElement("tr")
	for _, item := range d.headers {
		tr.Append(headerTmpl.Clone().WithID("").WithInnerHTML(item))
	}
	header.Append(tr)

	// Populate table content
	start := offset * pageSize
	for i := 0; i < pageSize && len(d.filteredItems) > start+i; i++ {
		row := rowTmpl.Clone().WithID("")
		for _, item := range d.items[d.filteredItems[start+i]] {
			if element, ok := item.(*dom.Element); ok {
				row.Append(dom.CreateElement("td").Append(element))
			} else {
				row.Append(dom.CreateElement("td").WithInnerHTML(item))
			}
		}
		content.Append(row)
	}

	d.renderPaginator(offset)
}

func (d *datapager) renderPaginator(offset int) {
	activePageTmpl := dom.GetElementByID("activePage")
	inactivePageTmpl := dom.GetElementByID("inactivePage")
	pagerLeftDisabled := dom.GetElementByID("pagerLeftDisabled").Clone().WithID(fmt.Sprintf("pld%d", offset))
	pagerLeftEnabled := dom.GetElementByID("pagerLeftEnabled").Clone().WithID(fmt.Sprintf("ple%d", offset))
	pagerRightDisabled := dom.GetElementByID("pagerRightDisabled").Clone().WithID(fmt.Sprintf("prd%d", offset))
	pagerRightEnabled := dom.GetElementByID("pagerRightEnabled").Clone().WithID(fmt.Sprintf("pre%d", offset))

	pager := dom.GetElementByID("pager").WithInnerHTML("")

	maxPages := selectablePages
	if d.pageCount() < selectablePages {
		maxPages = d.pageCount()
	}

	start := offset - (maxPages / 2)
	if start < 0 {
		start = 0
	}

	// render the maxPages with the appropriate page active
	if offset != 0 {
		pager.Append(pagerLeftEnabled)
	} else {
		pager.Append(pagerLeftDisabled)
	}
	for i := start; i < start+maxPages; i++ {
		if offset == i {
			active := activePageTmpl.Clone().WithID(fmt.Sprintf("%d", i))
			active.GetElementByTagName("button").WithInnerHTML(offset + 1)
			pager.Append(active)
		} else {
			inactive := inactivePageTmpl.Clone().WithID(fmt.Sprintf("%d", i))
			inactive.GetElementByTagName("button").WithID(fmt.Sprintf("btn%d", i)).WithInnerHTML(i + 1)
			pager.Append(inactive)
		}
	}
	if offset != d.pageCount() {
		pager.Append(pagerRightEnabled)
	} else {
		pager.Append(pagerRightDisabled)
	}
}

func (d *datapager) applyFilter(fltr string) {
	d.filteredItems = nil
	for i, row := range d.items {
		for _, item := range row {
			// can't filter on dom elements
			if _, ok := item.(*dom.Element); ok {
				continue
			}
			// assume everything else is a string
			// FIXME: should let data provider indicate which columns are filterable
			if strings.Contains(item.(string), fltr) {
				d.filteredItems = append(d.filteredItems, i)
				// Only add each row once
				break
			}
		}
	}
}

func (d *datapager) pageCount() int {
	if len(d.filteredItems)%pageSize == 0 {
		return len(d.filteredItems) / pageSize
	}
	return (len(d.filteredItems) / pageSize) + 1
}

func (d *datapager) nextPage() {
	d.curPage++
	d.renderPage(d.curPage)
}

func (d *datapager) prevPage() {
	d.curPage--
	d.renderPage(d.curPage)
}

func (d *datapager) page(offset int) {
	d.curPage = offset
	d.renderPage(d.curPage)
}

func (d *datapager) registerCallbacks() {
	dom.RegisterFunc("nextPage", func(this js.Value, i []js.Value) interface{} {
		d.nextPage()
		return nil
	})
	dom.RegisterFunc("prevPage", func(this js.Value, i []js.Value) interface{} {
		d.prevPage()
		return nil
	})
	dom.RegisterFunc("page", func(this js.Value, i []js.Value) interface{} {
		idx, err := strconv.Atoi(i[0].String())
		if err != nil {
			log.Printf("error accessing page %s", i[0].String())
			return nil
		}
		d.page(idx - 1)
		return nil
	})
	dom.GetElementByID("srch-term").RemoveAllEventListeners().AddEventListener("input", func(this js.Value, i []js.Value) interface{} {
		log.Printf("event listener triggered")
		fltrStr := dom.GetElementByID("srch-term").Get("value").String()
		d.applyFilter(fltrStr)
		// Reset current page pointer if filter makes collection smaller than current location
		if d.pageCount() < d.curPage {
			d.curPage = 0
		}
		d.renderPage(d.curPage)
		return nil
	})
}
