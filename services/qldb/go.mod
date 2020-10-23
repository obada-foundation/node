module github.com/obada-protocol/server-gateway/services/qldb

go 1.14

require (
	github.com/ardanlabs/conf v1.3.2
	github.com/aws/aws-sdk-go v1.34.28
	github.com/awslabs/amazon-qldb-driver-go v0.0.0-20200909015545-4ac582ff3b8f
	github.com/dimfeld/httptreemux/v5 v5.2.2
	github.com/go-playground/locales v0.13.0
	github.com/go-playground/universal-translator v0.17.0
	github.com/leodido/go-urn v1.2.0 // indirect
	github.com/pkg/errors v0.9.1
	go.opentelemetry.io/contrib/instrumentation/net/http v0.11.0
	go.opentelemetry.io/otel v0.11.0
	go.opentelemetry.io/otel/exporters/trace/zipkin v0.11.0
	go.opentelemetry.io/otel/sdk v0.11.0
	gopkg.in/go-playground/assert.v1 v1.2.1 // indirect
	gopkg.in/go-playground/validator.v9 v9.31.0
)
