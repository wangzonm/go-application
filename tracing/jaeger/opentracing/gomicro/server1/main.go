package main

import (
	"os"
	"server1/handler"
	pb "server1/proto"

	ot "github.com/go-micro/plugins/v4/wrapper/trace/opentracing"
	"go-micro.dev/v4"
	"go-micro.dev/v4/logger"

	"github.com/go-micro/cli/debug/trace/jaeger"
	cfg "github.com/uber/jaeger-client-go"
)

var (
	service = "server1"
	version = "latest"
)

func main() {
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
	if err := pb.RegisterServer1Handler(srv.Server(), new(handler.Server1)); err != nil {
		logger.Fatal(err)
	}
	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
