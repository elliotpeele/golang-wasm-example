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

package api

import (
	"context"

	"github.com/elliotpeele/golang-wasm-example/api/pb"
	"github.com/op/go-logging"
)

var logger = logging.MustGetLogger("golang-wasm-example.api")

// Server provides an interface for gRPC
type Server struct {
}

// New creates an instance of a server
func New() (*Server, error) {
	return &Server{}, nil
}

// ListUsers lists all users
func (s *Server) ListUsers(ctx context.Context, e *pb.Empty) (*pb.UsersResponse, error) {
	return &pb.UsersResponse{}, nil
}

// ListProjects lists all projects
func (s *Server) ListProjects(ctx context.Context, e *pb.Empty) (*pb.ProjectsResponse, error) {
	return &pb.ProjectsResponse{}, nil
}
