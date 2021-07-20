package helper

type ObitDef struct {
	SerialNumberHash string `json:"serial_number_hash"`
	Usn string `json:"usn"`
	DID string `json:"obit_did"`
	Usn58 string `json:"usn_base58"`
}

type KV struct {
	Key string `ion:"Key" json:"key"`
	Value string `ion:"Value" json:"value"`
}

type LocalObit struct {
	SerialNumber string `json:"serial_number"`
	Manufacturer string `json:"manufacturer"`
	PartNumber string `json:"part_number"`
	ObitStatus string `json:"obit_status"`
	Owner string `json:"owner"`
	ModifiedOn int64 `json:"modified_on"`
	Metadata []KV `json:"metadata"`
	StructuredData []KV `json:"structured_data"`
	Documents interface{} `json:"documents"`
}
