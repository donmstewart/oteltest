package oteltest_test

import (
	"context"
	"errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

	"github.com/donmstewart/oteltest/pkg/v1/oteltest"
)

const (
	LogErrorMsg       = "log error message"
	ErrorAttributeKey = "ErrorMsg"
)

var _ = Describe("Tracer", func() {
	var (
		ctx   context.Context
		err   error
		span  trace.Span
		spans []*oteltest.Span
		sr    *oteltest.SpanRecorder
		tp    *oteltest.TracerProvider
	)
	Describe("Test the creation of Spans & their data", func() {
		BeforeEach(func() {
			// Create the span Context object
			ctx = context.Background()
			// Create an error whost message we will add as an attribute to the OTEL Spam
			err = errors.New(LogErrorMsg)
			// Setup OTEL Tracing test code
			sr = new(oteltest.SpanRecorder)
			// Create the OTEL test harness TracerProvider
			tp = oteltest.NewTracerProvider(oteltest.WithSpanRecorder(sr))
			// We need an OTEL span for our log functionality test cases
			ctx, span = tp.Tracer("Log Utils").Start(ctx, "Log Utils")
			span.AddEvent("CreateSpan", trace.WithAttributes(
				attribute.String(ErrorAttributeKey, err.Error()),
			))

			span.RecordError(err)
			span.SetStatus(codes.Error, "Span Status set for "+LogErrorMsg)
			// Tear down the OTEL Span
			defer span.End()
		})
		It("should be possible to retrieve Span Events & Attributes", func() {
			spans = sr.Completed()
			By("extracting the first span we can verify the name of the span")
			Expect(spans[0].Name()).To(Equal("Log Utils"))
			By("extracting the first Event verify its name")
			Expect(spans[0].Events()[0].Name).To(Equal("CreateSpan"))
			By("validating the number of attributes added to the span")
			Expect(spans[0].Events()[0].Attributes).Should(HaveLen(1))
			By("validating the attributes has an entry in the map with the given key name and value")
			Expect(spans[0].Events()[0].Attributes[attribute.Key(ErrorAttributeKey)].AsString()).
				Should(Equal(LogErrorMsg))
		})
		It("should be possible to retrieve Span Events & Attributes", func() {
			spans = sr.Completed()
			By("checking the span we should see it is in error")
			Expect(spans[0].StatusCode()).Should(Equal(codes.Error))
			By("checking the span name we should see it is an exception")
			Expect(spans[0].Events()[1].Name).Should(Equal("exception"))
			By("validating the Spans attribute array has an entry in the map with the key 'exception.type' " +
				"and the value '*errors.errorString'")
			Expect(spans[0].Events()[1].Attributes[attribute.Key("exception.type")].AsString()).
				Should(Equal("*errors.errorString"))
			By("checking the attributes type we should be able to determine it is 'attribute.String'")
			v := spans[0].Events()[1].Attributes[attribute.Key("exception.message")]
			Expect(v.Type().String()).Should(Equal(attribute.STRING.String()))
			By("validating the attributes has an entry in the map with the key 'exception.type' " +
				"and the value '*errors.errorString'")
			Expect(spans[0].Events()[1].Attributes[attribute.Key("exception.message")].AsString()).
				Should(Equal("log error message"))
		})
	})
})
