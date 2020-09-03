package utils

import "regexp"

var (
	domainRegexp = regexp.MustCompile("^[a-zA-Z0-9][-a-zA-Z0-9]{0,62}(\\.[a-zA-Z0-9][-a-zA-Z0-9]{0,62})+$")
	schemaRegexp = regexp.MustCompile("http(s)?|w(s){1,2}")
)

func CheckDomain(domain string) bool {
	return domainRegexp.MatchString(domain)
}

func CheckSchema(schema string) bool {
	return schemaRegexp.MatchString(schema)
}
