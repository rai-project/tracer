package utils

import (
	"net/url"
	"strings"
)

func fixEndpoint(scheme, port, path, endpoint string) (string, error) {
	if !strings.HasPrefix(endpoint, "http://") && !strings.HasPrefix(endpoint, "https://") {
		endpoint = scheme + endpoint
	}
	u, err := url.Parse(endpoint)
	if err != nil {
		return "", err
	}
	if u.Scheme == "" {
		u.Scheme = scheme
	}
	if urlPort(u) == "" {
		u.Host = u.Host + ":" + port
	}
	if u.Path == "" {
		u.Path = path
	}
	return u.String(), nil
}

func FixEndpoints(scheme, port, path string) func(endpoints []string) []string {
	return func(endpoints []string) []string {
		res := []string{}
		for _, endpoint := range endpoints {
			endpoint, err := fixEndpoint(scheme, port, path, endpoint)
			if err != nil {
				continue
			}
			res = append(res, endpoint)
		}
		return res
	}
}
