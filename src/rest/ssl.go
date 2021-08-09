package rest

// sslMode defines ssl mode for rest server
type sslMode int8

const (
	// None defines to run http server only
	None sslMode = iota

	// Static defines to run both https and http server. Redirect http to https
	Static

	// Auto defines to run both https and http server. Redirect http to https. Https server with autocert support
	Auto
)

// SSLConfig holds all ssl params for rest server
type SSLConfig struct {
	SSLMode      sslMode
	Cert         string
	Key          string
	Port         int
	ACMELocation string
	ACMEEmail    string
}
