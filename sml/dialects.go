package sml

import (
	"fmt"

	log "github.com/maagodata/maago-commons/logger"
)

// This file contains common utility functions that are dialect specific.

// GetQuotes returns the dialect specific quotes.
func GetQuotes(dialect string) (startQuote string, endQuote string, err error) {

	supportedDialects := []string{
		DialectDuckDB,
	}

	if !IsStringInSlice(dialect, supportedDialects) {
		err = fmt.Errorf("unsupported dialect: %s", dialect)
		log.Error(err)
		return
	}

	startQuote = "`"
	endQuote = startQuote

	return
}
