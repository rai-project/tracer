package utils

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"runtime"
	"strings"
	"sync"

	"github.com/Unknwon/com"
)

var (
	hostIp string
	once   sync.Once
)

func isLinuxSubsystem() bool {
	if runtime.GOOS != "linux" {
		return false
	}
	checkWSL := func(pth string) bool {
		if !com.IsFile(pth) {
			return false
		}
		byts, err := ioutil.ReadFile(pth)
		if err != nil {
			return false
		}
		ver := strings.ToLower(string(byts))
		return strings.Contains(ver, "microsoft") || strings.Contains(ver, "wsl")
	}
	return checkWSL("/proc/version") || checkWSL("/proc/sys/kernel/osrelease")
}

func getHostIP() string {
	if isLinuxSubsystem() {
		return "linux-wsl"
	}
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}

	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

func GetHostIP() string {
	once.Do(func() {
		hostIp = getHostIP()
	})
	return hostIp
}

func getExternalIpFrom(service string) (string, error) {
	rsp, err := http.Get(service)
	if err != nil {
		return "", err
	}
	defer rsp.Body.Close()

	buf, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return "", err
	}

	return string(bytes.TrimSpace(buf)), nil
}

func GetExternalIp() (string, error) {
	services := []string{
		"http://checkip.amazonaws.com",
		"http://myexternalip.com/raw",
		"http://icanhazip.com",
		"http://canihazip.com/s",
	}
	for _, service := range services {
		ip, err := getExternalIpFrom(service)
		if err == nil {
			return ip, err
		}
	}
	return "", errors.New("Cannot get external ip")
}

// GetLocalIP returns the non loopback local IP of the host
func GetLocalIP0() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", errors.New("Cannot get a list of network interfaces")
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}
		}
	}
	return "", errors.New("Cannot get local ip address")
}

func GetLocalIp() (string, error) {
	host, err := os.Hostname()
	if err != nil {
		return "", err
	}
	addrs, err := net.LookupIP(host)
	if err != nil {
		return "", err
	}
	for _, addr := range addrs {
		if ipv4 := addr.To4(); ipv4 != nil {
			return ipv4.String(), nil
		}
	}
	return "", errors.New("Cannot get local ip address")
}
