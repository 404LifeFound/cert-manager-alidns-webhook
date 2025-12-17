package utils

import (
	"strings"
)

// UnFqdn converts the fqdn into a name removing the trailing dot.
func UnFqdn(name string) string {
	n := len(name)
	if n != 0 && name[n-1] == '.' {
		return name[:n-1]
	}
	return name
}

func ExtractRR(fqdn, domain string) string {
	name := UnFqdn(fqdn)
	if idx := strings.Index(name, "."+domain); idx != -1 {
		return name[:idx]
	}

	return name
}
