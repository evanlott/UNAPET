package pages

import "pest/functions"
import "os"

const WWW_ROOT string = "/var/www/htdocs/unapet/workspace/"

func EnsureNotNull(s string) {

	if s == "" {
		functions.ErrorResponse("Bad request.")
		os.Exit(0)
	}
}
