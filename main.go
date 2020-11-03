package main

import (
	"fmt"
	"net"
	"os"
	"time"

	"github.com/hashicorp/go-hclog"
	btdprotos "github.com/iantal/btd/protos/btd"
	"github.com/iantal/dta/internal/files"
	"github.com/iantal/dta/internal/server"
	protos "github.com/iantal/dta/protos/dta"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" // postgres
	"github.com/spf13/viper"

	gpprotos "github.com/iantal/dta/protos/gradle-parser"
	mcdprotos "github.com/iantal/mcd/protos/mcd"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func gRPCConnection(host string) *grpc.ClientConn {
	conn, err := grpc.Dial(
		host,
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(1024*1000*3000)),
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithTimeout(60*time.Second),
	)
	if err != nil {
		fmt.Printf("Error connecting to %s", host)
		panic(err)
	}
	return conn
}

func main() {
	viper.AutomaticEnv()
	log := hclog.Default()

	// create a new gRPC server, use WithInsecure to allow http connections
	gs := grpc.NewServer()

	bp := fmt.Sprintf("%v", viper.Get("BASE_PATH"))
	rmHost := fmt.Sprintf("%v", viper.Get("RM_HOST"))
	btdHost := fmt.Sprintf("%v", viper.Get("BTD_HOST"))
	gradleParserHost := fmt.Sprintf("%v", viper.Get("GP_HOST"))
	mcdHost := fmt.Sprintf("%v", viper.Get("MCD_HOST"))

	stor, err := files.NewLocal(bp, 1024*1000*1000*5)
	if err != nil {
		log.Error("Unable to create storage", "error", err)
		os.Exit(1)
	}

	// setup GRPC for BTD
	connBTD := gRPCConnection(btdHost)
	defer connBTD.Close()
	cc := btdprotos.NewUsedBuildToolsClient(connBTD)

	// setup GRPC for GRADLE-PARSER
	connGP := gRPCConnection(gradleParserHost)
	defer connGP.Close()
	gpc := gpprotos.NewGradleParseServiceClient(connGP)

	// setup GRPC for MCD
	connMCD := gRPCConnection(mcdHost)
	defer connMCD.Close()
	mcd := mcdprotos.NewDownloaderClient(connMCD)

	user := viper.Get("POSTGRES_USER")
	password := viper.Get("POSTGRES_PASSWORD")
	database := viper.Get("POSTGRES_DB")
	host := viper.Get("POSTGRES_HOST")
	port := viper.Get("POSTGRES_PORT")
	connection := fmt.Sprintf("host=%v port=%v user=%v dbname=%v password=%v sslmode=disable", host, port, user, database, password)

	db, err := gorm.Open("postgres", connection)
	defer db.Close()
	if err != nil {
		panic("Failed to connect to database!")
	}

	err = db.DB().Ping()
	if err != nil {
		panic("Ping failed!")
	}

	c := server.NewCommitExplorer(log, db, bp, cc, gpc, mcd, rmHost, stor)

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
