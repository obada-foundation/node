package helper

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/obada-foundation/sdkgo/properties"
	"log"

	"github.com/obada-foundation/sdkgo"
)

// Service provider an API to manage helper calls
type Service struct {
	logger *log.Logger
	sdk    *sdkgo.Sdk
}

func NewService(sdk *sdkgo.Sdk, logger *log.Logger) *Service {
	return &Service{
		logger: logger,
		sdk:    sdk,
	}
}

func hashStr(str string) (string, error) {
	h := sha256.New()

	if _, err := h.Write([]byte(str)); err != nil {
		return "", fmt.Errorf("cannot wite bytes %v to hasher: %w", []byte(str), err)
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}

func (s Service) ToDto(lo LocalObit) (sdkgo.ObitDto, error) {
	var dto sdkgo.ObitDto

	snHash, err := hashStr(lo.SerialNumber)

	if err != nil {
		return dto, err
	}

	dto.SerialNumberHash = snHash
	dto.Manufacturer = lo.Manufacturer
	dto.PartNumber = lo.PartNumber
	dto.OwnerDid = lo.Owner
	dto.Status = lo.ObitStatus
	dto.ModifiedOn = lo.ModifiedOn

	kvs := func(kvs []KV) []properties.KV {
		var newKVs []properties.KV

		for _, kv := range kvs {
			newKVs = append(newKVs, properties.KV(kv))
		}

		return newKVs
	}

	dto.Matadata = kvs(lo.Metadata)
	dto.StructuredData = kvs(lo.StructuredData)

	return dto, nil
}

// GenRootHash generates obit root hash
func (s Service) GenRootHash(lo LocalObit) (string, error) {
	dto, err := s.ToDto(lo)

	if err != nil {
		return "", err
	}

	o, err := s.sdk.NewObit(dto)

	if err != nil {
		return "", err
	}

	h, err := o.GetRootHash(nil)

	if err != nil {
		return "", err
	}

	return h.GetHash(), nil
}
