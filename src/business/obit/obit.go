package obit

import (
	sdk "github.com/obada-foundation/sdk-go"
)

type Service struct {
	sdk *sdk.Sdk
}

func NewObitService(sdk *sdk.Sdk) *Service {
	return &Service{
		sdk: sdk,
	}
}

func (os Service) Create() {

}

func (os Service) Update() {

}

func (os Service) Delete() {

}

func (os Service) Search() {

}

func (os Service) History() {

}
