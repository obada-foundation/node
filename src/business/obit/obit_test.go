package obit

import (
	"bytes"
	"github.com/obada-foundation/sdkgo"
	"log"
	"testing"
)

type testCase struct {
	args args
	want  ID
}

type args struct {
	serialNumberHash string
	manufacturer     string
	partNumber       string
}

func createService(t *testing.T) (*Service, *log.Logger) {
	sdk, err := sdkgo.NewSdk(nil, false)

	if err != nil {
		t.Fatal(err.Error())
	}

	var logStr bytes.Buffer

	logger := log.New(&logStr, "", 0)

	service := NewObitService(sdk, logger, nil, nil, nil)

	return service, logger
}

func TestService_GenerateID(t *testing.T) {
	testCases := []testCase{
		{
			args: args{
				serialNumberHash: "dc0fb8e9835790195bf4a8e5e122fe608e548f46f88410cc6792927bedbb6d55", // sha256("serial_number")
				manufacturer:     "manufacturer",
				partNumber:       "part number",
			},
			want: ID{
				ID:  "bb00c8da8424d0af25cbef87968f3784bc829671ff208c5dc9505ab2976a369f",
				DID: "did:obada:bb00c8da8424d0af25cbef87968f3784bc829671ff208c5dc9505ab2976a369f",
			},
		},
		{
			args: args{
				serialNumberHash: "cae6b797ae2627d96689fed03adc28311d5f2175253c3a0e375301e225ddf44d", // sha256("SN123456X")
				manufacturer:     "SONY",
				partNumber:       "PN123456S",
			},
			want: ID{
				ID: "d7cf869423d12f623f5611e48d6f6665bbc4a270b6e09da2f4c32bcb1b949ecd",
				DID: "did:obada:d7cf869423d12f623f5611e48d6f6665bbc4a270b6e09da2f4c32bcb1b949ecd",
			},
		},
	}

	service, _ := createService(t)

	for _, tc := range testCases {
		args := tc.args

		got, err := service.GenerateID(args.serialNumberHash, args.manufacturer, args.partNumber)

		if err != nil {
			t.Error(err.Error())
		}

		if got != tc.want {
			t.Errorf(
				"service.GenerateID(%q, %q, %q) = %v+ want %v+",
				args.serialNumberHash, args.manufacturer, args.partNumber, got, tc.want,
			)
		}
	}
}
