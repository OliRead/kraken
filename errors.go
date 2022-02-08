package kraken

import (
	"errors"
)

var (
	// ErrGeneral represents all the EGeneral errors from the Kraken API
	ErrGeneral = errors.New("EGeneral")
	// ErrAPI represents all the EAPI errors from the Kraken API
	ErrAPI = errors.New("EAPI")
	// ErrQuery represents all the EQuery errors from the Kraken API
	ErrQuery = errors.New("EQuery")
	// ErrOrder represents all the EOrder errors from the Kraken API
	ErrOrder = errors.New("EOrder")
	// ErrTrade represents all the ETrade errors from the Kraken API
	ErrTrade = errors.New("ETrade")
	// ErrFunding represents all the EFunding errors from the Kraken API
	ErrFunding = errors.New("EFunding")
	// ErrService represents all the EService errors from the Kraken API
	ErrService = errors.New("EService")
	// ErrSession represents all the ESession errors from the Kraken API
	ErrSession = errors.New("ESession")

	// ErrAPIUnknown an unknown error was returned from the API
	ErrAPIUnknown = errors.New("unknown API error")
	// ErrDryRun dry run has been specified so action cannot be completed
	ErrDryRun = errors.New("dry run")
	// ErrParse error during parsing of an API response
	ErrParse = errors.New("parse error")
	// ErrNetwork error occoured during the transportation of a message
	ErrNetwork = errors.New("network error")
)
