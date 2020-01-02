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
	"context"
	"fmt"
	"log"
	"runtime"
	"syscall/js"
	"time"

	"google.golang.org/grpc"
	"github.com/tarndt/wasmws"

	"github.com/elliotpeele/golang-wasm-example/api/pb"
	"github.com/elliotpeele/golang-wasm-example/frontend/dom"
)

var (
	version    = "devel"
	buildDate  string
	commitHash string
)

func versionInfo() string {
	return fmt.Sprintf(`
golang-wasm-example frontend:
 version     : %s
 build date  : %s
 git hash    : %s
 go version  : %s
 go compiler : %s
 platform    : %s/%s
`, version, buildDate, commitHash, runtime.Version(), runtime.Compiler, runtime.GOOS, runtime.GOARCH)
}

func mkInterfaceSlice(items ...interface{}) []interface{} {
	data := make([]interface{}, len(items))
	for i, item := range items {
		data[i] = item
	}
	return data
}

type client struct {
	ctx  context.Context
	conn *grpc.ClientConn
	cli  pb.WASMExampleClient
}

func main() {
	log.Println("WASM Go Initialized")
	log.Println(versionInfo())

	ch := make(chan struct{})

	//App context setup
	appCtx, appCancel := context.WithCancel(context.Background())
	defer appCancel()

	//Dial setup
	const dialTO = time.Second * 30
	dialCtx, dialCancel := context.WithTimeout(appCtx, dialTO)
	defer dialCancel()

	loc, err := dom.ParseLocation()
	if err != nil {
		log.Fatalf("failed to parse location information: %s", err)
	}
	log.Printf("location info: %+v", loc)

	//Connect to remote gRPC server
	websocketURL := fmt.Sprintf("ws://%s/grpc-proxy", loc.Host)
	conn, err := grpc.DialContext(dialCtx, "passthrough:///"+websocketURL, grpc.WithContextDialer(wasmws.GRPCDialer), grpc.WithDisableRetry(), grpc.WithMaxMsgSize(1*1024*1024*1024), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not gRPC dial: %s; Details: %s", websocketURL, err)
	}

	cli := &client{
		ctx:  appCtx,
		conn: conn,
		cli:  pb.NewWASMExampleClient(conn),
	}
	cli.registerCallbacks()

	// default display of users
	cli.listUsers(js.Value{}, nil)

	<-ch
}

func (c *client) registerCallbacks() {
	dom.RegisterFunc("listUsers", c.listUsers)
	dom.RegisterFunc("listProjects", c.listProjects)
}

func (c *client) listUsers(this js.Value, i []js.Value) interface{} {
	log.Println("listUsers called")
	go func() {
		resp, err := c.cli.ListUsers(c.ctx, &pb.Empty{})
		if err != nil {
			log.Printf("error listing users: %s", err)
			return
		}
		dp := &datapager{
			headers: []string{
				"id",
				"updated at",
				"first name",
				"last name",
				"email",
			},
		}
		for _, u := range resp.Users {
			dp.items = append(dp.items, mkInterfaceSlice(
				u.Id,
				time.Unix(u.UpdatedAt.Seconds, 0).UTC().String(),
				u.FirstName,
				u.LastName,
				u.Email,
			))
		}
		dp.registerCallbacks()
		dp.applyFilter("")
		dp.renderPage(0)
	}()
	return nil
}

func (c *client) listProjects(this js.Value, i []js.Value) interface{} {
	log.Println("listProjects called")
	go func() {
		resp, err := c.cli.ListProjects(c.ctx, &pb.Empty{})
		if err != nil {
			log.Printf("error listing projects: %s", err)
			return
		}
		dp := &datapager{
			headers: []string{
				"id",
				"name",
				"updated at",
			},
		}
		for _, p := range resp.Projects {
			dp.items = append(dp.items, mkInterfaceSlice(
				p.Id,
				p.Name,
				time.Unix(p.UpdatedAt.Seconds, 0).UTC().String(),
			))
		}
		dp.registerCallbacks()
		dp.applyFilter("")
		dp.renderPage(0)
	}()
	return nil
}