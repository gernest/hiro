package models

import "sync"

// reservedNames special names that aren"t allowed as usernames.
// Source : https://ldpreload.com/blog/names-to-reserve
var reservedNames = map[string]bool{
	// Hostnames with special/reserved meaning.
	"autoconfig":    true, // Thunderbird autoconfig
	"autodiscover":  true, // MS Outlook/Exchange autoconfig
	"broadcasthost": true, // Network broadcast hostname
	"isatap":        true, // IPv6 tunnel autodiscovery
	"localdomain":   true, // Loopback
	"localhost":     true, // Loopback
	"wpad":          true, // Proxy autodiscovery

	// Common protocol hostnames.
	"ftp":     true,
	"imap":    true,
	"mail":    true,
	"news":    true,
	"pop":     true,
	"pop3":    true,
	"smtp":    true,
	"usenet":  true,
	"uucp":    true,
	"webmail": true,
	"www":     true,

	// Email addresses known used by certificate authorities during
	// verification.
	"admin":            true,
	"administrator":    true,
	"hostmaster":       true,
	"info":             true,
	"is":               true,
	"it":               true,
	"mis":              true,
	"postmaster":       true,
	"root":             true,
	"ssladmin":         true,
	"ssladministrator": true,
	"sslwebmaster":     true,
	"sysadmin":         true,
	"webmaster":        true,

	// RFC-2142-defined names not already covered.
	"abuse":     true,
	"marketing": true,
	"noc":       true,
	"sales":     true,
	"security":  true,
	"support":   true,

	// Common no-reply email addresses.
	"mailer-daemon": true,
	"nobody":        true,
	"noreply":       true,
	"no-reply":      true,

	// Sensitive filenames.
	"clientaccesspolicy.xml": true, // Silverlight cross-domain policy file.
	"crossdomain.xml":        true, // Flash cross-domain policy file.
	"favicon.ico":            true,
	"humans.txt":             true,
	"keybase.txt":            true, // Keybase ownership-verification URL.
	"robots.txt":             true,
	".htaccess":              true,
	".htpasswd":              true,

	// Other names which could be problems depending on URL/subdomain
	// structure.
	"account":     true,
	"accounts":    true,
	"blog":        true,
	"buy":         true,
	"clients":     true,
	"contact":     true,
	"contactus":   true,
	"contact-us":  true,
	"copyright":   true,
	"dashboard":   true,
	"doc":         true,
	"docs":        true,
	"download":    true,
	"downloads":   true,
	"enquiry":     true,
	"faq":         true,
	"help":        true,
	"inquiry":     true,
	"license":     true,
	"login":       true,
	"logout":      true,
	"me":          true,
	"myaccount":   true,
	"payments":    true,
	"plans":       true,
	"portfolio":   true,
	"preferences": true,
	"pricing":     true,
	"privacy":     true,
	"profile":     true,
	"register":    true,
	"secure":      true,
	"settings":    true,
	"signin":      true,
	"signup":      true,
	"ssl":         true,
	"status":      true,
	"subscribe":   true,
	"terms":       true,
	"tos":         true,
	"user":        true,
	"users":       true,
	"weblog":      true,
	"work":        true,

	//bq specific
	"bq":      true,
	"gernest": true,
}
var mu sync.RWMutex

// IsReserved returns true if the name is reserved
func IsReserved(name string) bool {
	mu.RLock()
	v := reservedNames[name]
	mu.RUnlock()
	return v
}
