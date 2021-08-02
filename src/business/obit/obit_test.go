package obit

import (
	"context"
	"github.com/obada-foundation/node/business/tests"
	"github.com/obada-foundation/sdkgo"
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

type ObitTests struct {
	service *Service
}

func TestService(t *testing.T) {
	test := tests.NewIntegration(t)
	t.Cleanup(test.Teardown)

	sdk, err := sdkgo.NewSdk(test.Logger, false)

	if err != nil {
		t.Fatal(err.Error())
	}

	service := NewObitService(sdk, test.Logger, test.DB, nil, nil)

	tests := ObitTests{
		service: service,
	}

	t.Run("generateID", tests.generateID)
	t.Run("checksum", tests.checksum)
}

func (os ObitTests) generateID(t *testing.T) {
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

	for _, tc := range testCases {
		args := tc.args

		got, err := os.service.GenerateID(args.serialNumberHash, args.manufacturer, args.partNumber)

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

func (os ObitTests) checksum(t *testing.T) {
	testCases := []struct{
		arg sdkgo.ObitDto
		want string
	}{
		{
			arg: sdkgo.ObitDto{
				ObitIDDto: sdkgo.ObitIDDto{
					SerialNumberHash: "cae6b797ae2627d96689fed03adc28311d5f2175253c3a0e375301e225ddf44d",
					Manufacturer: "SONY",
					PartNumber: "PN123456S",
				},
				OwnerDid: "did:obada:owner:123456",
				ModifiedOn: 1624387537,
			},
			want: "2eb12c48ad2f073c49b95fcf2190cec40548c69fdc6f49135dee0753020f1624",
		},
	}

	for _, tc := range testCases {
		got, err := os.service.Checksum(context.Background(), tc.arg)

		if err != nil {
			t.Error(err.Error())
		}

		if got != tc.want {
			t.Errorf(
				"service.Checksum(%v+) = %s want %s",
				tc.arg, got, tc.want,
			)
		}
	}
}