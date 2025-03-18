package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	consulAPI "github.com/hashicorp/consul/api"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"src/internal/conf"
	"src/internal/repository"
	"src/internal/service"
	pb "src/protos/Tipster"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	// Name is the name of the compiled software.
	Name = "domain.service.tipster"
	// Version is the version of the compiled software.
	Version string
	// configPath is the config file path.
	configPath string

	id, _ = os.Hostname()
)

func init() {
	flag.StringVar(&configPath, "conf", "../../configs/config.yaml", "config path, eg: -conf config.yaml")
}

func newApp(logger log.Logger, hs *http.Server, gs *grpc.Server) *kratos.App {
	return kratos.New(
		kratos.ID(id),
		kratos.Name(Name),
		kratos.Version(Version),
		kratos.Metadata(map[string]string{}),
		kratos.Logger(logger),
		kratos.Server(
			hs,
			gs,
		),
	)
}

type GlobalConfig struct {
	AppSettings struct {
		ShutdownTimeoutInSecond int `json:"ShutdownTimeoutInSecond"`
	} `json:"AppSettings"`
	DomainService struct {
		Tipster struct {
			HostName string `json:"HostName"`
			Port     string `json:"Port"`
		} `json:"Tipster"`
	} `json:"DomainService"`
}

type MongoDBConfig struct {
	MongoDbConnection struct {
		Uri      string `json:"uri"`
		Database string `json:"database"`
	} `json:"MongoDbConnection"`
}

func main() {
	flag.Parse()
	logger := log.With(log.NewStdLogger(os.Stdout),
		"ts", log.DefaultTimestamp,
		"caller", log.DefaultCaller,
		"service.id", id,
		"service.name", Name,
		"service.version", Version,
		"trace.id", tracing.TraceID(),
		"span.id", tracing.SpanID(),
	)

	c := config.New(
		config.WithSource(
			file.NewSource(configPath),
		),
	)

	defer c.Close()

	if err := c.Load(); err != nil {
		panic(err)
	}

	var bc conf.Bootstrap
	if err := c.Scan(&bc); err != nil {
		panic(err)
	}

	// Initialize Consul client with config address
	consulConfig := consulAPI.DefaultConfig()
	consulConfig.Address = bc.Consul.Address

	consulClient, err := consulAPI.NewClient(consulConfig)
	if err != nil {
		logger.Log(log.LevelError, "msg", "failed to create consul client", "error", err)
		panic(err)
	}

	// Get configuration from Consul
	kv := consulClient.KV()

	// Read MongoDB configuration
	mongodbPair, _, err := kv.Get("tipster", nil)
	if err != nil {
		logger.Log(log.LevelError, "msg", "failed to get mongodb config from consul", "error", err)
		panic(err)
	}
	if mongodbPair == nil {
		logger.Log(log.LevelError, "msg", "mongodb config not found in consul")
		panic("mongodb config not found in consul")
	}

	// Read GRPC configuration
	grpcPair, _, err := kv.Get("Global", nil)
	if err != nil {
		logger.Log(log.LevelError, "msg", "failed to get grpc config from consul", "error", err)
		panic(err)
	}
	if grpcPair == nil {
		logger.Log(log.LevelError, "msg", "grpc config not found in consul")
		panic("grpc config not found in consul")
	}

	// Parse MongoDB config
	var mongodbCfg MongoDBConfig
	if err := json.Unmarshal(mongodbPair.Value, &mongodbCfg); err != nil {
		logger.Log(log.LevelError, "msg", "failed to parse mongodb consul config", "error", err)
		panic(err)
	}

	// // Parse GRPC config
	var grpcCfg GlobalConfig
	if err := json.Unmarshal(grpcPair.Value, &grpcCfg); err != nil {
		logger.Log(log.LevelError, "msg", "failed to parse grpc consul config", "error", err)
		panic(err)
	}

	// Create the Bootstrap config
	bc.Mongodb = &conf.MongoDbConnection{
		Uri:      mongodbCfg.MongoDbConnection.Uri,
		Database: mongodbCfg.MongoDbConnection.Database,
	}

	// Convert Port from string to int32
	port, err := strconv.ParseInt(grpcCfg.DomainService.Tipster.Port, 10, 32)
	if err != nil {
		logger.Log(log.LevelError, "msg", "failed to parse port number", "error", err)
		panic(err)
	}

	host := grpcCfg.DomainService.Tipster.HostName

	bc.GrpcServer = &conf.GRPCServer{
		Port: int32(port),
		Host: host,
	}

	// Initialize MongoDB client

	mongoClient, err := mongo.Connect(context.Background(), options.Client().ApplyURI(bc.Mongodb.Uri))
	// mongoClient, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://root:pass.123@localhost:22097/"))
	if err != nil {
		logger.Log(log.LevelError, "msg", "failed to connect to mongodb", "error", err)
		panic(err)
	}
	defer mongoClient.Disconnect(context.Background())

	socialLogger := log.With(logger,
		"module", "social",
		"service", "social.v1.SocialService",
	)

	// Initialize repository
	db := mongoClient.Database(bc.Mongodb.Database)
	// db := mongoClient.Database("tipster")
	socialRepo := repository.NewSocialRepository(db, socialLogger)
	// HTTP Server
	httpSrv := http.NewServer(
		http.Address(":8000"),
	)

	// gRPC Server
	grpcAddr := fmt.Sprintf("%s:%d", bc.GrpcServer.Host, bc.GrpcServer.Port)
	// grpcAddr := fmt.Sprintf("%s:%d", "0.0.0.0", 9000)
	grpcSrv := grpc.NewServer(
		grpc.Address(grpcAddr),
		grpc.Middleware(
			recovery.Recovery(),
		),
	)

	logger.Log(log.LevelInfo, "msg", "starting gRPC server", "address", grpcAddr)

	solcialSvc := service.NewSocialServiceService(
		socialRepo,
		socialLogger,
	)
	pb.RegisterSocialServiceServer(grpcSrv, solcialSvc)
	app := newApp(logger, httpSrv, grpcSrv)
	if err := app.Run(); err != nil {
		panic(err)
	}
}
