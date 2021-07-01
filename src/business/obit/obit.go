package obit

import (
	"github.com/awslabs/amazon-qldb-driver-go/qldbdriver"
	"github.com/obada-foundation/sdkgo"
	"log"
)

// Service provider an API to manage obits
type Service struct {
	logger *log.Logger
	sdk    *sdkgo.Sdk
	qldb   *qldbdriver.QLDBDriver
}

// NewObitService creates new version of Obit service
func NewObitService(sdk *sdkgo.Sdk, logger *log.Logger, qldb *qldbdriver.QLDBDriver) *Service {
	return &Service{
		logger: logger,
		sdk:    sdk,
		qldb:   qldb,
	}
}

// Create method creates a new Obit
func (os Service) Create(dto sdkgo.ObitDto) error {
	obit, err := os.sdk.NewObit(dto)

	if err != nil {
		return err
	}

	os.logger.Printf("%v", obit)

	return nil
}

// Update method updates a new Obit
func (os Service) Update() {

}

// Delete method search Obits by given criteria ??
func (os Service) Delete() {

}

// Search method search Obits by given criteria
func (os Service) Search() {

}

// History the history of Obit changes
func (os Service) History() {

}
