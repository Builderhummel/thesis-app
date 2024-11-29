package auth_controller

import (
	"log"
	"regexp"
	"strings"

	"github.com/go-ldap/ldap/v3"
)

func (auser *AuthUser) LDAP_authenticate(username, password string) error {
	ldapURL := Config.LDAPUrl // ldaps://ldap.example.com:636
	ldapdn := Config.LDAPDn   // ou=people,dc=tu-example,dc=de

	l, err := ldap.DialURL(ldapURL)
	if err != nil {
		return err
	}
	defer l.Close()

	userdn := "uid=" + username + "," + ldapdn
	err = l.Bind(userdn, password)
	if err != nil {
		return err
	}

	searchRequest := ldap.NewSearchRequest(
		ldapdn,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		"(uid="+username+")",
		[]string{},
		nil,
	)
	sr, err := l.Search(searchRequest)
	if err != nil {
		log.Fatal(err)
	}

	for _, entry := range sr.Entries {
		for _, attr := range entry.Attributes {
			if attr.Name == "uid" {
				auser.UID = attr.Values[0]
			} else if attr.Name == "cn" {
				auser.Name = attr.Values[0]
			} else if attr.Name == "mail" {
				auser.Email = attr.Values[0]
			}
		}
	}

	return nil
}

func SanitizeLDAPInput(input string) string {
	if input == "" {
		return ""
	}

	specialChars := map[string]string{
		",":  "\\,",
		"\\": "\\\\",
		"/":  "\\/",
		"#":  "\\#",
		"+":  "\\+",
		"<":  "\\<",
		">":  "\\>",
		";":  "\\;",
		"\"": "\\\"",
		"=":  "\\=",
	}

	sanitized := input

	for char, replacement := range specialChars {
		sanitized = strings.ReplaceAll(sanitized, char, replacement)
	}

	sanitized = regexp.MustCompile(`[\x00-\x1F\x7F]`).ReplaceAllString(sanitized, "")

	const maxLength = 255
	if len(sanitized) > maxLength {
		sanitized = sanitized[:maxLength]
	}

	return sanitized
}
