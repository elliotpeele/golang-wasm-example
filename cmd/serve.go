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
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/elliotpeele/golang-wasm-example/api"
	"github.com/elliotpeele/golang-wasm-example/api/pb"
	"github.com/elliotpeele/golang-wasm-example/frontend/assets"
	"github.com/spf13/cobra"
	"github.com/tarndt/wasmws"
	"google.golang.org/grpc"
)

// serveCmd represents the index command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "",
	Long:  "",
	RunE: func(cmd *cobra.Command, args []string) error {
		port, err := cmd.PersistentFlags().GetString("port")
		if err != nil {
			return err
		}

		appCtx, appCancel := context.WithCancel(context.Background())
		defer appCancel()

		//Setup HTTP / Websocket server
		router := http.NewServeMux()
		wsl := wasmws.NewWebSocketListener(appCtx)
		router.HandleFunc("/grpc-proxy", wsl.ServeHTTP)
		router.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
			if req.URL.Path == "/" {
				req.URL.Path = "/index.html"
			}
			assets.ServeFile(w, req, req.URL.Path)
		})
		httpServer := &http.Server{Addr: ":" + port, Handler: router}
		//Run HTTP server
		go func() {
			defer appCancel()
			logger.Errorf("HTTP Listen and Server failed; Details: %s", httpServer.ListenAndServe())
		}()

		//gRPC setup
		s := grpc.NewServer(grpc.UnaryInterceptor(unaryInterceptor))
		srv, err := api.New()
		if err != nil {
			return err
		}
		pb.RegisterWASMExampleServer(s, srv)

		//Run gRPC server
		go func() {
			defer appCancel()
			if err := s.Serve(wsl); err != nil {
				logger.Errorf("Failed to serve gRPC connections; Details: %s", err)
			}
		}()

		//Handle signals
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
		go func() {
			logger.Infof("Received shutdown signal: %s", <-sigs)
			appCancel()
		}()

		//Shutdown
		<-appCtx.Done()
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), time.Second*2)
		defer shutdownCancel()

		grpcShutdown := make(chan struct{}, 1)
		go func() {
			s.GracefulStop()
			grpcShutdown <- struct{}{}
		}()

		httpServer.Shutdown(shutdownCtx)
		select {
		case <-grpcShutdown:
		case <-shutdownCtx.Done():
			s.Stop()
		}

		return nil
	},
}

func unaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()

	h, err := handler(ctx, req)

	//logging
	logger.Infof("request - Method:%s\tDuration:%s\tError:%v\n",
		info.FullMethod,
		time.Since(start),
		err)

	return h, err
}

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.PersistentFlags().String("port", "8080", "port on which to listen")
}
