package network

import (
	"errors"
	"net"
	"net/url"
)

func GetValidURL(addr string) (string, error) {
	u, err := url.Parse(addr)
	if err != nil {
		host, port, err := net.SplitHostPort(addr)
		if err != nil {
			return "", err
		}

		if net.ParseIP(host) == nil {
			return "", errors.New("address is invalid")
		}

		return "http://" + host + ":" + port, nil
	}
	u.Scheme = "http"
	return u.String(), nil
}

// The function `GetLocalIP` retrieves the local IP address of the machine it is running on.
func GetLocalIP() (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return "", err
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}
			return ip.String(), nil
		}
	}
	return "", errors.New("Warning, failed to get local ip, washtub features may not working well")
}
