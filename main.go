package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

const letters = "abcdefghijkmnopqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ23456789"

func generateToken() string {
	rand.Seed(time.Now().UnixNano())

	bytes := make([]byte, 12)
	for i := range bytes {
		rn := rand.Int63()
		bytes[i] = letters[rn%int64(len(letters))]
	}
	return string(bytes)
}

func getLocalIPs() []string {
	localIPRegex, _ := regexp.Compile(`^(192|172|10)\.`)
	var addrs []string
	iAddrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Printf("Error getting interfaces: %s\n", err)
	}

	for i := range iAddrs {
		addr := strings.Split(iAddrs[i].String(), "/")[0]
		match := localIPRegex.MatchString(addr)
		if match {
			addrs = append(addrs, addr)
		}
	}
	return addrs
}

func getPublicIP() (string, error) {
	res, err := http.Get("https://checkip.amazonaws.com")
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	return string(body[:len(body)-1]), nil
}

func main() {
	dirPath := "./"
	port := "8080"
	args := os.Args[1:]

	if len(args) >= 1 {
		dirPath = args[0]
	}

	dirPath, _ = filepath.Abs(filepath.Dir(dirPath))

	if len(args) == 2 {
		port = args[1]
	}

	token := generateToken()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		status := http.StatusOK
		t := r.URL.Query().Get("t")

		if r.URL.Path == "/" && t != token {
			status = http.StatusForbidden
			http.Error(w, "Invalid token", status)
		} else {
			http.ServeFile(w, r, dirPath+r.URL.Path)
		}

		fmt.Printf("%s %s %d\n", r.Method, r.URL.String(), status)
	})

	fmt.Printf("Serving %s on http://0.0.0.0:%s?t=%s\n", dirPath, port, token)
	fmt.Println("\nAvailable local addresses:")

	addrs := getLocalIPs()
	for i := range addrs {
		fmt.Printf("http://%s:%s?t=%s\n", addrs[i], port, token)
	}

	addr, err := getPublicIP()
	if err != nil {
		fmt.Printf("Could not get public IP: %s\n", err)
	} else {
		fmt.Println("\nAvailable public addresses:")
		fmt.Printf("http://%s:%s?t=%s\n", addr, port, token)
	}

	http.ListenAndServe(":"+port, nil)
}
