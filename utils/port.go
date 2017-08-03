// +build !go1.7

package utils

import "net/url"

func urlPort(u *url.URL) string {
	return u.Port()
}
