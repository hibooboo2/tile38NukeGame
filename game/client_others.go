// +build !js,!wasm

package game

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	c       *http.Client
	baseUrl string
}

func NewClient(baseUrl string) *Client {
	return &Client{&http.Client{Timeout: time.Second * 2}, baseUrl}
}

func (c *Client) post(cmd string) bool {
	// http://10.14.12.11:9851
	// s := time.Now()
	// defer func() {
	// 	log.Println(time.Since(s))
	// }()
	// timePost := time.Now()

	req, err := http.NewRequest(http.MethodPost, c.baseUrl, strings.NewReader(cmd))
	if err != nil {
		log.Println(err)
		return false
	}
	resp, err := c.c.Do(req)
	if err != nil {
		log.Println(err)
		return false
	}
	// log.Println(time.Since(timePost).Nanoseconds() / int64(time.Millisecond))
	val := map[string]interface{}{}
	// timePost = time.Now()
	json.NewDecoder(resp.Body).Decode(&val)
	if len(val) != 2 {
		log.Println(val)
		return false
	}
	// log.Println(time.Since(timePost).Nanoseconds() / int64(time.Millisecond))
	// log.Println(val["ok"], val["elapsed"])
	return val["ok"].(bool)
}

func (c *Client) Notifications(name string) bool {
	hookurl := "http://10.14.12.68:8081"
	if boundAddr != "" {
		hookurl = boundAddr
	}
	cmd := fmt.Sprintf("SETHOOK %[1]s %[2]s/%[1]s NEARBY fleet FENCE ROAM fleet %[1]s 1000", name, hookurl)
	fmt.Println(cmd)
	return c.post(cmd)
}

func ClearNotifications() {
	c := NewClient("http://10.14.12.11:9851")
	charLock.Lock()
	defer charLock.Unlock()
	for char := range chars {
		cmd := fmt.Sprintf("DELHOOK %s", char)
		c.post(cmd)
	}
}

var boundAddr string

func init() {
	log.SetFlags(log.Lshortfile)
	http.HandleFunc("/", func(resp http.ResponseWriter, req *http.Request) {
		thing := Thing{}
		err := json.NewDecoder(req.Body).Decode(&thing)
		if err != nil {
			fmt.Println(err)
		}
		events <- thing
		req.Body.Close()
		resp.WriteHeader(http.StatusOK)
	})
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		panic(err)
	}

	boundAddr = fmt.Sprintf("http://%s:%d", GetBoundAddr(), listener.Addr().(*net.TCPAddr).Port)
	go http.Serve(listener, nil)
}

// Get preferred outbound ip of this machine
func GetBoundAddr() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP.String()
}
