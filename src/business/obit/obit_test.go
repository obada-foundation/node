package obit

import (
	"context"
	"github.com/obada-foundation/node/business/tests"
	"testing"

	"github.com/obada-foundation/sdkgo"
)

func TestObit(t *testing.T) {
	logger, teardown := tests.NewUnit(t)
	sdk, err := sdkgo.NewSdk(logger, true)

	if err != nil {
		t.Fatal("Cannot construct sdk")
	}

	service := NewObitService(sdk, logger, nil)

	t.Log("Given need to check obit service functionality")
	var dto sdkgo.ObitDto
	service.Create(context.Background(), dto)

	t.Cleanup(teardown)
}
