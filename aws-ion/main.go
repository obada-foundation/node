package main

import (
	"encoding/json"
	"github.com/amzn/ion-go/ion"
	"github.com/arnaud-lb/php-go/php-go"
	"strconv"
	"strings"
)

// call php.Export() for its side effects
var _ = php.Export("ion", map[string]interface{}{
	"decode": IonDecode,
	"encode": IonEncode,
})

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

//
func IonEncode(jsonStr string) string {
	var o Obit

	err := json.Unmarshal([]byte(jsonStr), &o)

	if err != nil {
		return err.Error()
	}

	ionData, err := ion.MarshalText(o)

	if err != nil {
		return err.Error()
	}

	return string(ionData)
}

//
func IonDecode(bytes string) string {
	var o Obit

	strSplit := strings.Split(bytes, " ")

	var arr []byte

	// Converts string of bytes to array of bytes
	for _, v := range strSplit {
		// Cast string to integer
		i, err := strconv.ParseInt(v, 0, 64)

		if err != nil {
			return err.Error()
		}

		arr = append(arr, byte(i))
	}

	err := ion.Unmarshal(arr, &o)

	if err != nil {
		return err.Error()
	}

	json, err := json.Marshal(o)

	if err != nil {
		return err.Error()
	}

	return string(json)
}

func main() {

}
