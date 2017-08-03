// +build go1.7,!go1.8

package utils

import (
	"net/url"
	"strings"
)

func urlPort(u *url.URL) string {
	hostport := u.Host
	colon := strings.IndexByte(hostport, ':')
	if colon == -1 {
		return ""
	}
	if i := strings.Index(hostport, "]:"); i != -1 {
		return hostport[i+len("]:"):]
	}
	if strings.Contains(hostport, "]") {
		return ""
	}
	return hostport[colon+len(":"):]
}
