package obit

import (
	sdk "github.com/obada-foundation/sdkgo"
)

// Service provider an API to manage obits
type Service struct {
	sdk *sdk.Sdk
}

// NewObitService creates new version of Obit service
func NewObitService(sdk *sdk.Sdk) *Service {
	return &Service{
		sdk: sdk,
	}
}

// Create method creates a new Obit
func (os Service) Create() {

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
