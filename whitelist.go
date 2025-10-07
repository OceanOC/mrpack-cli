package main

import "strings"

var Whitelist = [...]string{"github.com", "cdn.modrinth.com", "githubusercontent.com", "gitlab.com"}

func (mp ModPack) NotInWhitelist(domain string) bool {
	for _, wl := range Whitelist {
		if strings.Contains(domain, wl) {
			return false
		}
	}
	return true
}
