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

package client

import (
	"context"
	"fmt"
	"time"
	"log"
	"syscall/js"

	"google.golang.org/grpc"
	"github.com/elliotpeele/golang-wasm-example/api/pb"

	"github.com/tarndt/wasmws"

	"github.com/elliotpeele/golang-wasm-example/frontend/dom"
	"github.com/elliotpeele/golang-wasm-example/frontend/datatable"
)

// Client interfaces to the backend grpc/websocket server to bridge into the dom
type Client struct {
	ctx  context.Context
	conn *grpc.ClientConn
	cli  pb.WASMExampleClient
}

// New creates an instance of Client 
func New(ctx context.Context) (*Client, error) {
	//Dial setup
	const dialTO = time.Second * 30
	dialCtx, dialCancel := context.WithTimeout(ctx, dialTO)
	defer dialCancel()

	// Parse location information to get server connection info
	loc, err := dom.ParseLocation()
	if err != nil {
		return nil, fmt.Errorf("failed to parse location information: %s", err)
	}
	log.Printf("location info: %+v", loc)

	//Connect to remote gRPC server
	websocketURL := fmt.Sprintf("ws://%s/grpc-proxy", loc.Host)
	conn, err := grpc.DialContext(dialCtx, "passthrough:///"+websocketURL,
		grpc.WithContextDialer(wasmws.GRPCDialer),
		grpc.WithDisableRetry(),
		grpc.WithMaxMsgSize(1*1024*1024*1024),
		grpc.WithInsecure(),
	)
	if err != nil {
		return nil, fmt.Errorf("Could not gRPC dial: %s; Details: %s", websocketURL, err)
	}

	cli := &Client{
		ctx:  ctx,
		conn: conn,
		cli:  pb.NewWASMExampleClient(conn),
	}
	cli.registerCallbacks()
	return cli, nil
}

func wrapFunc(f func()) dom.CallbackFunc {
	return func(this js.Value, i []js.Value) interface{} {
		go func() {
			f()
		}()
		return nil
	}
}

func (c *Client) registerCallbacks() {
	dom.RegisterFunc("listUsers", wrapFunc(c.ListUsers))
	dom.RegisterFunc("listProjects", wrapFunc(c.ListProjects))
}

// ListUsers fetches users from the server and renders them in a table
func (c *Client) ListUsers() {
	log.Println("listUsers called")
	resp, err := c.cli.ListUsers(c.ctx, &pb.Empty{})
	if err != nil {
		log.Printf("error listing users: %s", err)
		return
	}
	dt := datatable.New().WithPagination().WithSearch()
	dt.Headers(
		"id",
		"updated at",
		"first name",
		"last name",
		"email",
	)
	for _, u := range resp.Users {
		dt.AppendRow(
			u.Id,
			time.Unix(u.UpdatedAt.Seconds, 0).UTC().String(),
			u.FirstName,
			u.LastName,
			u.Email,
		)
	}
	dt.Render()
}

// ListProjects fetches projects from the server and renderes them in a table
func (c *Client) ListProjects() {
	log.Println("listProjects called")
	resp, err := c.cli.ListProjects(c.ctx, &pb.Empty{})
	if err != nil {
		log.Printf("error listing projects: %s", err)
		return
	}
	dt := datatable.New().WithPagination().WithSearch()
	dt.Headers(
		"id",
		"name",
		"updated at",
		"button",
	)
	for _, p := range resp.Projects {
		btn := dom.GetElementByID("testButton").Clone().WithID(p.Id)
		dt.AppendRow(
			p.Id,
			p.Name,
			time.Unix(p.UpdatedAt.Seconds, 0).UTC().String(),
			btn,
		)
	}
	dt.Render()
}