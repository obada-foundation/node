package sdkgo

import (
	"fmt"
	"github.com/obada-foundation/sdkgo/hash"
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/obada-foundation/sdkgo/properties"
)

const obadaReleaseDateUTC = 1624387536

// Sdk OBADA SDK
type Sdk struct {
	logger   *log.Logger
	debug    bool
	validate *validator.Validate
}

// NewSdk creates a new obada SDK instance
func NewSdk(logger *log.Logger, debug bool) (*Sdk, error) {
	v, err := initializeValidator()

	if err != nil {
		return nil, err
	}

	return &Sdk{
		logger:   logger,
		debug:    debug,
		validate: v,
	}, nil
}

func initializeValidator() (*validator.Validate, error) {
	v := validator.New()

	if err := v.RegisterValidation("min-modified-on", validateMinModifiedOn); err != nil {
		return v, err
	}

	return v, nil
}

func validateMinModifiedOn(fl validator.FieldLevel) bool {
	return fl.Field().Int() >= obadaReleaseDateUTC
}

// NewObit creates new obit
func (sdk *Sdk) NewObit(dto ObitDto) (Obit, error) {
	var o Obit
	o.debug = sdk.debug
	o.logger = sdk.logger

	err := sdk.validate.Struct(dto)
	if err != nil {
		return o, err
	}

	if sdk.debug {
		sdk.logger.Printf("NewObit(%v)", dto)
	}

	snProp, err := properties.NewStringProperty(
		"Making serialNumberHash",
		dto.SerialNumberHash,
		sdk.logger,
		sdk.debug,
	)

	if err != nil {
		return o, err
	}

	manufacturerProp, err := properties.NewStringProperty(
		"Making manufacturer hash",
		dto.Manufacturer,
		sdk.logger,
		sdk.debug,
	)

	if err != nil {
		return o, err
	}

	pnProp, err := properties.NewStringProperty(
		"Making partNumber hash",
		dto.PartNumber,
		sdk.logger,
		sdk.debug,
	)

	if err != nil {
		return o, err
	}

	obdDidProp, err := properties.NewStringProperty(
		"Making obdDid hash",
		dto.ObdDid,
		sdk.logger,
		sdk.debug,
	)

	if err != nil {
		return o, err
	}

	ownerDidProp, err := properties.NewStringProperty(
		"Making ownerDid hash",
		dto.OwnerDid,
		sdk.logger,
		sdk.debug,
	)

	if err != nil {
		return o, err
	}

	statusProp, err := properties.NewStatusProperty(dto.Status, sdk.logger, sdk.debug)

	if err != nil {
		return o, err
	}

	obitIDProp, err := properties.NewObitIDProperty(snProp, manufacturerProp, pnProp, sdk.logger, sdk.debug)

	if err != nil {
		return o, err
	}

	modifiedOn, err := properties.NewIntProperty(
		"Making modifiedOn hash",
		dto.ModifiedOn,
		sdk.logger,
		sdk.debug,
	)

	if err != nil {
		return o, err
	}

	metadataProp, err := properties.NewMapProperty(
		"Making matadata hash",
		dto.Matadata,
		sdk.logger,
		sdk.debug,
	)

	if err != nil {
		return o, err
	}

	strctDataProp, err := properties.NewMapProperty(
		"Making structuredData hash",
		dto.StructuredData,
		sdk.logger,
		sdk.debug,
	)

	if err != nil {
		return o, err
	}

	documentsProp, err := properties.NewMapProperty(
		"Making documents hash",
		dto.Documents,
		sdk.logger,
		sdk.debug,
	)

	if err != nil {
		return o, err
	}

	altIDsProp, err := properties.NewSliceStrProperty(
		"Making alternateIDs property",
		dto.AlternateIDS,
		sdk.logger,
		sdk.debug,
	)

	if err != nil {
		return o, err
	}

	o.obitID = obitIDProp
	o.serialNumberHash = snProp
	o.manufacturer = manufacturerProp
	o.partNumber = pnProp
	o.obdDid = obdDidProp
	o.ownerDid = ownerDidProp
	o.status = statusProp
	o.metadata = metadataProp
	o.structuredData = strctDataProp
	o.documents = documentsProp
	o.modifiedOn = modifiedOn
	o.alternateIDS = altIDsProp

	return o, nil
}

// GetObitID returns obit ID
func (o Obit) GetObitID() properties.ObitID {
	return o.obitID
}

// GetSerialNumberHash returns serial number hash Obit property
func (o Obit) GetSerialNumberHash() properties.StringProperty {
	return o.serialNumberHash
}

// GetPartNumber returns part number Obit property
func (o Obit) GetPartNumber() properties.StringProperty {
	return o.manufacturer
}

// GetManufacturer returns manufacturer Obit property
func (o Obit) GetManufacturer() properties.StringProperty {
	return o.manufacturer
}

// GetOwnerDID returns OBADA Obit owner DID
func (o Obit) GetOwnerDID() properties.StringProperty {
	return o.ownerDid
}

// GetObdDID returns OBADA Obit obd DID
func (o Obit) GetObdDID() properties.StringProperty {
	return o.obdDid
}

// GetMetadata returns Obit metadata
func (o Obit) GetMetadata() properties.KvProperty {
	return o.metadata
}

// GetStructuredData returns Obit structured data
func (o Obit) GetStructuredData() properties.KvProperty {
	return o.structuredData
}

// GetDocuments returns Obit documents
func (o Obit) GetDocuments() properties.KvProperty {
	return o.documents
}

// GetModifiedOn returns Obit modified on UNIX timestamp
func (o Obit) GetModifiedOn() properties.IntProperty {
	return o.modifiedOn
}

// GetAlternateIDS returns Obit alternatives identifiers
func (o Obit) GetAlternateIDS() properties.SliceStrProperty {
	return o.alternateIDS
}

// GetStatus returns Obit status
func (o Obit) GetStatus() properties.StatusProperty {
	return o.status
}

// GetRootHash returns obit root hash
func (o Obit) GetRootHash() (hash.Hash, error) {
	var rootHash hash.Hash

	if o.debug {
		o.logger.Println("\n\nObit root hash calculation")
	}

	sum := o.obitID.GetHash().GetDec() +
		o.serialNumberHash.GetHash().GetDec() +
		o.manufacturer.GetHash().GetDec() +
		o.partNumber.GetHash().GetDec() +
		o.ownerDid.GetHash().GetDec() +
		o.obdDid.GetHash().GetDec() +
		o.metadata.GetHash().GetDec() +
		o.structuredData.GetHash().GetDec() +
		o.documents.GetHash().GetDec() +
		o.modifiedOn.GetHash().GetDec() +
		o.alternateIDS.GetHash().GetDec() +
		o.status.GetHash().GetDec()

	if o.debug {
		o.logger.Println(fmt.Sprintf(
			"(%d + %d + %d + %d + %d + %d + %d + %d + %d + %d + %d) -> %d -> Dec2Hex(%d) -> %s",
			o.obitID.GetHash().GetDec(),
			o.serialNumberHash.GetHash().GetDec(),
			o.manufacturer.GetHash().GetDec(),
			o.partNumber.GetHash().GetDec(),
			o.ownerDid.GetHash().GetDec(),
			o.obdDid.GetHash().GetDec(),
			o.metadata.GetHash().GetDec(),
			o.structuredData.GetHash().GetDec(),
			o.documents.GetHash().GetDec(),
			o.modifiedOn.GetHash().GetDec(),
			o.status.GetHash().GetDec(),
			sum,
			sum,
			fmt.Sprintf("%x", sum),
		))
	}

	rootHash, err := hash.NewHash(fmt.Sprintf("%x", sum), o.logger, o.debug)

	if err != nil {
		return rootHash, err
	}

	if o.debug {
		o.logger.Printf("RootHash(%v) -> %q", o, rootHash.GetHash())
	}

	return rootHash, nil
}

// NewObitID creates new obit id
func (sdk *Sdk) NewObitID(dto ObitIDDto) (properties.ObitID, error) {
	var obitID properties.ObitID

	if sdk.debug {
		sdk.logger.Printf("NewObitID(%q, %q, %q)", dto.SerialNumberHash, dto.Manufacturer, dto.PartNumber)
	}

	err := sdk.validate.Struct(dto)
	if err != nil {
		return obitID, err
	}

	snProp, err := properties.NewStringProperty(
		"Making serialNumberHash",
		dto.SerialNumberHash,
		sdk.logger,
		sdk.debug,
	)

	if err != nil {
		return obitID, err
	}

	mnProp, err := properties.NewStringProperty(
		"Making manufacturer",
		dto.Manufacturer,
		sdk.logger,
		sdk.debug,
	)

	if err != nil {
		return obitID, err
	}

	pnProp, err := properties.NewStringProperty(
		"Making partNumber",
		dto.PartNumber,
		sdk.logger,
		sdk.debug,
	)

	if err != nil {
		return obitID, err
	}

	obitID, err = properties.NewObitIDProperty(snProp, mnProp, pnProp, sdk.logger, sdk.debug)

	if err != nil {
		return obitID, err
	}

	return obitID, nil
}
