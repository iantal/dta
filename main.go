package main

import (
	"fmt"
	"net"
	"os"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/iantal/dta/internal/files"
	"github.com/iantal/dta/internal/server"
	protos "github.com/iantal/dta/protos/dta"
	btdprotos "github.com/iantal/btd/protos/btd"

	"github.com/spf13/viper"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	viper.AutomaticEnv()
	log := hclog.Default()

	// create a new gRPC server, use WithInsecure to allow http connections
	gs := grpc.NewServer()

	bp := fmt.Sprintf("%v", viper.Get("BASE_PATH"))
	rmHost := fmt.Sprintf("%v", viper.Get("RM_HOST"))
	btdHost := fmt.Sprintf("%v", viper.Get("BTD_HOST"))

	stor, err := files.NewLocal(bp, 1024*1000*1000*5)
	if err != nil {
		log.Error("Unable to create storage", "error", err)
		os.Exit(1)
	}

	conn, err := grpc.Dial(
		btdHost,
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(1024*1000*3000)),
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithTimeout(60*time.Second),
	)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	cc := btdprotos.NewUsedBuildToolsClient(conn)

	c := server.NewCommitExplorer(log, bp, cc, rmHost, stor)

	// register the currency server
	protos.RegisterCommitExplorerServer(gs, c)

	// register the reflection service which allows clients to determine the methods
	// for this gRPC service
	reflection.Register(gs)

	// create a TCP socket for inbound server connections
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", 8006))
	if err != nil {
		log.Error("Unable to create listener", "error", err)
		os.Exit(1)
	}

	log.Info("Starting server", "bind_address", l.Addr().String())
	// listen for requests
	gs.Serve(l)
}
