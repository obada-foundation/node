package obit

import "github.com/amzn/ion-go/ion"

type Obit struct {
	ObitDID string `ion:"ObitDID" json:"obit_did"`
	Usn string `ion:"Usn" json:"usn"`
	ObitDIDVersions string `ion:"ObitDIDVersions" json:"obit_did_versions"`
	OwnerDID string `ion:"OwnerDID" json:"owner_did"`
	Obd_DID string `ion:"Obd_DID" json:"obd_did"`
	ObitStatus string `ion:"ObitStatus" json:"obit_status"`
	Manufacturer string `ion:"Manufacturer" json:"manufacturer"`
	PartNumber string `ion:"PartNumber" json:"part_number"`
	SerialNumberHash string `ion:"SerialNumberHash" json:"serial_number_hash"`
	MetaData []MetaDataRecord `ion:"MetaData" json:"metadata"`
	DocLinks []DocumentHashLink `ion:"DocLinks" json:"doc_links"`
	StructuredData []StructuredDataRecord `ion:"StructuredData" json:"structured_data"`
	ModifiedAt ion.Timestamp `ion:"ModifiedAt" json:"modified_at"`
	RootHash string `ion:"RootHash" json:"root_hash"`
}

type MetaDataRecord struct {
	Key string `ion:"Key" json:"key"`
	Value string `ion:"Value" json:"value"`
}

type DocumentHashLink struct {
	Name string `ion:"Name" json:"name"`
	Hashlink string `ion:"Hashlink" json:"hashlink"`
}

type StructuredDataRecord struct {
	Key string `ion:"Key" json:"key"`
	Value string `ion:"Value" json:"value"`
}

// NewObit is what we require from clients when adding a Obit.
type NewObit struct {

}

type UpdateObit struct {

}