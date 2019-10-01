package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	httpip    = "https://localhost:8081" // 服务地址
	cacrt     = "./ve/ca.crt"            // 服务证书校验
	servercrt = "./ve/server.crt"        // 服务端证书
	serverkey = "./ve/server.key"        // 服务端密钥
	clientcrt = "./ve/client.crt"        // 客户端证书
	clientkey = "./ve/client.key"        // 客户端密钥
)

type myhandler struct {
}

var (
	check = flag.String("c", "https", "Start the HTTPS data service!")
)

func main() {
	flag.Parse()
	checkstring := *check
	switch checkstring {
	case "https":
		fmt.Println("Start the server-side validation service!")
		https_service()
	case "client":
		fmt.Println("Start the client validation service!")
		client_service()
	default:

	}
}

func https_service() {
	capem, err := ioutil.ReadFile(cacrt)
	if err != nil {
		fmt.Println("read ca.crt failed, err = %v", err)
	}

	ca := x509.NewCertPool()
	ok := ca.AppendCertsFromPEM(capem)
	if !ok {
		fmt.Println("ca append failed")
	}

	tlsConfig := &tls.Config{}
	tlsConfig.ClientCAs = ca
	tlsConfig.ClientAuth = tls.RequireAndVerifyClientCert

	s := &http.Server{}
	s.Addr = "localhost:8081"
	s.Handler = &myhandler{}
	s.TLSConfig = tlsConfig

	err = s.ListenAndServeTLS(servercrt, serverkey)
	if err != nil {
		fmt.Println("ListenAndServeTLS err = ", err)
	}
}

func client_service() {
	capem, err := ioutil.ReadFile(cacrt)
	if err != nil {
		fmt.Println("read ca.cer failed, err = ", err)
	}
	ca := x509.NewCertPool()
	ok := ca.AppendCertsFromPEM(capem)
	if !ok {
		fmt.Println("ca append failed")
	}

	tlsConfig := &tls.Config{}
	tlsConfig.RootCAs = ca
	tlsConfig.Certificates = make([]tls.Certificate, 1)
	tlsConfig.Certificates[0], err = tls.LoadX509KeyPair(clientcrt, clientkey)
	if err != nil {
		fmt.Println("load x509 failed, err = ", err)
	}

	tr := &http.Transport{}
	tr.TLSClientConfig = tlsConfig

	client := &http.Client{Transport: tr}
	resp, err := client.Get(httpip)
	if err != nil {
		fmt.Println("client get failed, err = ", err)
	}

	content, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(strings.TrimSpace(string(content)))
}

func (h *myhandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprintf(w, "service verified successfully!")
	if err != nil {
		fmt.Println(err)
	}
}
