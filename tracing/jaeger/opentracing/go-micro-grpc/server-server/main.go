package main

import (
	"opentracing/go-micro-grpc/server-server/handler"
	pb "opentracing/go-micro-grpc/server-server/proto"
	"os"

	"github.com/go-micro/cli/debug/trace/jaeger"
	ot "github.com/go-micro/plugins/v4/wrapper/trace/opentracing"
	cfg "github.com/uber/jaeger-client-go"
	"go-micro.dev/v4"
	"go-micro.dev/v4/logger"
)

var (
	service = "go-micro-server"
	version = "latest"
)

func main() {
	_ = os.Setenv("JAEGER_AGENT_HOST", "127.0.0.1")
	_ = os.Setenv("JAEGER_AGENT_PORT", "6831")
	_ = os.Setenv("JAEGER_SAMPLER_TYPE", cfg.SamplerTypeConst)
	_ = os.Setenv("JAEGER_SAMPLER_PARAM", "1")
	// Create tracer
	tracer, closer, err := jaeger.NewTracer(
		jaeger.Name(service),
		jaeger.FromEnv(true),
		jaeger.GlobalTracer(true),
	)
	if err != nil {
		logger.Fatal(err)
	}
	defer closer.Close()

	// Create service
	srv := micro.NewService(
		micro.WrapCall(ot.NewCallWrapper(tracer)),
		micro.WrapClient(ot.NewClientWrapper(tracer)),
		micro.WrapHandler(ot.NewHandlerWrapper(tracer)),
		micro.WrapSubscriber(ot.NewSubscriberWrapper(tracer)),
	)
	srv.Init(
		micro.Name(service),
		micro.Version(version),
	)

	// Register handler
	if err := pb.RegisterServerServerHandler(srv.Server(), new(handler.ServerServer)); err != nil {
		logger.Fatal(err)
	}
	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
