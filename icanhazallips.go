// https://gist.github.com/miguelmota/7b765edff00dc676215d6174f3f30216
package main

import (
	"errors"
	"log"
	"net"
	"net/http"
	"strings"
)

func getIP(r *http.Request) (string, error) {
	ips := r.Header.Get("X-Forwarded-For")
	splitIps := strings.Split(ips, ",")

	if len(splitIps) > 0 {
		netIP := net.ParseIP(splitIps[len(splitIps)-1])
		if netIP != nil {
			return netIP.String(), nil
		}
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return "", err
	}

	netIP := net.ParseIP(ip)
	if netIP != nil {
		ip := netIP.String()
		if ip == "::1" {
			return "127.0.0.1", nil
		}
		return ip, nil
	}

	return "", errors.New("IP not found")
}

func handler(w http.ResponseWriter, r *http.Request) {
	ip, err := getIP(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(ip))
}

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServeTLS(":8080", "/certs/server.cert", "/certs/server.key", nil))
}
