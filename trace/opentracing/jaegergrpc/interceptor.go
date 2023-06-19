package jaegergrpc

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type textMapReader struct {
	metadata.MD
}

func (t textMapReader) ForeachKey(handler func(key, val string) error) error {
	for key, val := range t.MD {
		for _, v := range val {
			if err := handler(key, v); err != nil {
				return err
			}
		}
	}
	return nil
}

func UnaryServerInterceptor(o ...ServerOption) grpc.UnaryServerInterceptor {
	opts := serverOptions{
		requestIgnorer: IgnoreNone,
	}
	for _, o := range o {
		o(&opts)
	}
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		if opts.requestIgnorer(info) {
			return handler(ctx, req)
		}
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			md = metadata.New(nil)
		} else {
			md = md.Copy()
		}
		carrier := textMapReader{md}
		tracer := opentracing.GlobalTracer()
		spanCtx, _ := tracer.Extract(opentracing.TextMap, carrier)
		span := tracer.StartSpan(info.FullMethod, opentracing.ChildOf(spanCtx))
		defer span.Finish()
		ext.SpanKindRPCServer.Set(span)
		ctx = opentracing.ContextWithSpan(ctx, span)
		ret, err := handler(ctx, req)
		if err != nil {
			ext.LogError(span, err)
		}
		return ret, err
	}
}

type serverOptions struct {
	requestIgnorer RequestIgnorerFunc
}

// ServerOption sets options for server-side tracing.
type ServerOption func(*serverOptions)

// RequestIgnorerFunc is the type of a function for use in
// WithServerRequestIgnorer.
type RequestIgnorerFunc func(*grpc.UnaryServerInfo) bool

// WithServerRequestIgnorer returns a ServerOption which sets r as the
// function to use to determine whether or not a server request should
// be ignored. If r is nil, all requests will be reported.
func WithServerRequestIgnorer(r RequestIgnorerFunc) ServerOption {
	if r == nil {
		r = IgnoreNone
	}
	return func(o *serverOptions) {
		o.requestIgnorer = r
	}
}

// IgnoreNone is a RequestIgnorerFunc which ignores no requests.
func IgnoreNone(*grpc.UnaryServerInfo) bool {
	return false
}

type textMapWriter struct {
	metadata.MD
}

func (t textMapWriter) Set(key, val string) {
	t.MD[key] = append(t.MD[key], val)
}

func UnaryClientInterceptor() grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		span, _ := opentracing.StartSpanFromContext(ctx, method)
		defer span.Finish()
		ext.SpanKindRPCClient.Set(span)

		md, ok := metadata.FromOutgoingContext(ctx)
		if !ok {
			md = metadata.New(nil)
		} else {
			md = md.Copy()
		}

		tracer := opentracing.GlobalTracer()
		carrier := textMapWriter{md}
		tracer.Inject(span.Context(), opentracing.TextMap, carrier)
		ctx = metadata.NewOutgoingContext(ctx, md)
		err := invoker(ctx, method, req, reply, cc, opts...)
		if err != nil {
			ext.LogError(span, err)
		}
		return err
	}
}
