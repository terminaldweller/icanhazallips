// https://gist.github.com/miguelmota/7b765edff00dc676215d6174f3f30216
package main

import (
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	defaultTimeOut = 10.
)

var (
	errMalformedAdr  = errors.New("malformed address")
	errIPNotFound    = errors.New("ip not found")
	errBadFloatValue = errors.New("bad float value")
	errBadConfig     = errors.New("bad config")
)

func getDefaultOptions() map[string]float64 {
	return map[string]float64{
		"APP_CONTEXT_TIMEOUT":     defaultTimeOut,
		"APP_READ_HEADER_TIMEOUT": defaultTimeOut,
		"APP_READ_TIMEOUT":        defaultTimeOut,
		"APP_WRITE_TIMEOUT":       defaultTimeOut,
		"APP_IDLE_TIMEOUT":        defaultTimeOut,
	}
}

func getIP(request *http.Request) (string, error) {
	ips := request.Header.Get("X-Forwarded-For")
	splitIps := strings.Split(ips, ",")

	log.Println(request.RemoteAddr)

	if len(splitIps) > 0 {
		netIP := net.ParseIP(splitIps[len(splitIps)-1])
		if netIP != nil {
			return netIP.String(), nil
		}
	}

	ip, _, err := net.SplitHostPort(request.RemoteAddr)
	if err != nil {
		return "", errMalformedAdr
	}

	if netIP := net.ParseIP(ip); netIP != nil {
		ip := netIP.String()
		if ip == "::1" {
			return "127.0.0.1", nil
		}

		return ip, nil
	}

	return "", errIPNotFound
}

func ipHandler(writer http.ResponseWriter, request *http.Request) {
	ipAddr, err := getIP(request)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)

		return
	}

	writer.WriteHeader(http.StatusOK)

	writtenLen, err := writer.Write([]byte(ipAddr))
	if err != nil || writtenLen != len(ipAddr) {
		writer.WriteHeader(http.StatusInternalServerError)

		return
	}
}

type Config struct {
	Addr              string
	ContextTimeOut    float64
	ReadHeaderTimeout float64
	ReadTimeout       float64
	WriteTimeout      float64
	IdleTimeout       float64
}

func getConfigValue(envVarName string) (float64, error) {
	defaultOptions := getDefaultOptions()

	paramEnv := os.Getenv(envVarName)
	if paramEnv != "" {
		param, err := strconv.ParseFloat(paramEnv, 64)
		if err != nil {
			return defaultOptions[envVarName], errBadFloatValue
		}

		return param, nil
	}

	return defaultOptions[envVarName], nil
}

func getConfig() (Config, error) {
	var config Config

	var err error

	appAddrEnv := os.Getenv("APP_ADDR")
	if appAddrEnv != "" {
		config.Addr = appAddrEnv
	} else {
		config.Addr = ":8080"
	}

	config.ContextTimeOut, err = getConfigValue("APP_CONTEXT_TIMEOUT")
	if err != nil {
		log.Println(err.Error())
	}

	config.ReadHeaderTimeout, err = getConfigValue("APP_READ_HEADER_TIMEOUT")
	if err != nil {
		log.Println(err.Error())
	}

	config.ReadTimeout, err = getConfigValue("APP_READ_TIMEOUT")
	if err != nil {
		log.Println(err.Error())
	}

	config.WriteTimeout, err = getConfigValue("APP_WRITE_TIMEOUT")
	if err != nil {
		log.Println(err.Error())
	}

	config.IdleTimeout, err = getConfigValue("APP_IDLE_TIMEOUT")
	if err != nil {
		log.Println(err.Error())
	}

	return config, err
}

func main() {
	log.SetOutput(os.Stdout)
	http.HandleFunc("/", ipHandler)

	config, err := getConfig()
	if err != nil {
		log.Fatal(errBadConfig)
	}

	server := http.Server{
		Addr:              config.Addr,
		ReadHeaderTimeout: time.Duration(config.ReadHeaderTimeout) * time.Second,
		ReadTimeout:       time.Duration(config.ReadTimeout) * time.Second,
		WriteTimeout:      time.Duration(config.WriteTimeout) * time.Second,
		IdleTimeout:       time.Duration(config.IdleTimeout) * time.Second,
		TLSNextProto:      nil,
		ErrorLog:          nil,
		Handler:           nil,
	}

	log.Fatal(server.ListenAndServe())
}
