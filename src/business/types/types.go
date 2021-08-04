package types

// QLDBObit represent a node Obit
type QLDBObit struct {
	ObitDID          string            `ion:"ObitDID" db:"obit_id" json:"obit_did"`
	Usn              string            `ion:"Usn" json:"usn"`
	SerialNumberHash string            `ion:"SerialNumberHash" json:"serial_number_hash"`
	Manufacturer     string            `ion:"Manufacturer" json:"manufacturer"`
	PartNumber       string            `ion:"PartNumber" json:"part_number"`
	AlternateIDS     []string          `ion:"AlternateIDS" json:"alternate_ids"`
	OwnerDID         string            `ion:"OwnerDID" json:"owner_did"`
	ObdDID           string            `ion:"ObdDID" json:"obd_did"`
	Metadata         []KV              `ion:"MetaData" json:"metadata"`
	StructuredData   []KV              `ion:"StructuredData" json:"structured_data"`
	Documents        map[string]string `ion:"Documents" json:"documents"`
	ModifiedOn       int64             `ion:"ModifiedOn" json:"modified_on"`
	Status           string            `ion:"Status" json:"status"`
	Checksum         string            `ion:"Checksum" json:"checksum"`
}

// KV hosts key/value pairs for QLDBObit
type KV struct {
	Key   string `ion:"Key" json:"key"`
	Value string `ion:"Value" json:"value"`
}
