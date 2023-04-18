package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/go-micro/cli/debug/trace/jaeger"
	ot "github.com/go-micro/plugins/v4/wrapper/trace/opentracing"
	"github.com/opentracing/opentracing-go"
	opentracinglog "github.com/opentracing/opentracing-go/log"
	cfg "github.com/uber/jaeger-client-go"
	"go-micro.dev/v4"
	"go-micro.dev/v4/client"
	"go-micro.dev/v4/logger"
	"go-micro.dev/v4/registry"
	pb "opentracing/go-micro-grpc/client-client/proto"
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
		jaeger.Name("go-micro-client"),
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
	)
	srv.Init()

	// Create client
	c := pb.NewServerServerService(service, srv.Client())

	for {
		// Call service
		rsp, err := c.Call(context.Background(), &pb.CallRequest{Name: "John"})
		if err != nil {
			logger.Fatal(err)
		}

		logger.Info(rsp)

		time.Sleep(1 * time.Second)
	}
}

func rpcCallWrap(tracer opentracing.Tracer) client.CallWrapper {
	return func(cf client.CallFunc) client.CallFunc {
		return func(ctx context.Context, node *registry.Node, req client.Request, rsp interface{}, opts client.CallOptions) error {
			if tracer == nil {
				tracer = opentracing.GlobalTracer()
			}
			name := fmt.Sprintf("go-micro.client::%s.%s", req.Service(), req.Endpoint())
			ctx, span, err := ot.StartSpanFromContext(ctx, tracer, name)
			if err != nil {
				return err
			}
			defer span.Finish()
			if err = cf(ctx, node, req, rsp, opts); err != nil {
				span.LogFields(opentracinglog.String("error", err.Error()))
				span.SetTag("error", true)
			}
			return err
		}
	}
}
