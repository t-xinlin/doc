package examples

import (
	"context"
	"github.com/openzipkin/zipkin-go"
)

func doSomeWork(context.Context) {}

func ExampleNewTracer() {
	tracer := GetTracer("demoService", "127.0.0.1:9898")
	// tracer can now be used to create spans.
	span := tracer.StartSpan("some_operation")
	// ... do some work ...
	span.Finish()

	childSpan := tracer.StartSpan("some_operation2", zipkin.Parent(span.Context()))
	// ... do some work ...
	childSpan.Finish()

	span.Finish()

	// Output:
}
