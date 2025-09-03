package auth_controller

import (
	"fmt"
	"log"

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

	userdn := "uid=" + ldap.EscapeDN(username) + "," + ldapdn
	err = l.Bind(userdn, password)
	if err != nil {
		return err
	}

	filter := fmt.Sprintf("(uid=%s)", ldap.EscapeFilter(username))

	searchRequest := ldap.NewSearchRequest(
		ldapdn,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		filter,
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
