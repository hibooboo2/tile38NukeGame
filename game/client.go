package game

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

type Client struct {
	c       *http.Client
	baseUrl string
}

func NewClient(baseUrl string) *Client {
	return &Client{&http.Client{Timeout: time.Second * 2}, baseUrl}
}

var boundAddr string

func (c *Client) Notifications(name string) bool {
	hookurl := "http://10.14.12.68:8081"
	if boundAddr != "" {
		hookurl = boundAddr
	}
	cmd := fmt.Sprintf("SETHOOK %[1]s %[2]s/%[1]s NEARBY fleet FENCE ROAM fleet %[1]s 1000", name, hookurl)
	fmt.Println(cmd)
	return c.post(cmd)
}

func (c *Client) post(cmd string) bool {
	// http://10.14.12.11:9851
	resp, err := c.c.Post(c.baseUrl, "", strings.NewReader(cmd))
	if err != nil {
		fmt.Println(err)
		return false
	}
	val := map[string]interface{}{}
	json.NewDecoder(resp.Body).Decode(&val)
	if len(val) > 2 {
		fmt.Println(val)
		return false
	}
	fmt.Println(val["ok"], val["elapsed"])
	return val["ok"].(bool)
}

var (
	chars    = map[string]*Character{}
	charLock = sync.Mutex{}
)

func init() {
	log.SetFlags(log.Lshortfile)
	http.HandleFunc("/", func(resp http.ResponseWriter, req *http.Request) {
		thing := Thing{}
		err := json.NewDecoder(req.Body).Decode(&thing)
		if err != nil {
			fmt.Println(err)
		}

		charLock.Lock()
		c, ok := chars[thing.Nearby.ID]
		charLock.Unlock()
		if ok {
			c.Things <- thing.KeyedPoint
		} else {
			log.Printf("%s %s %s", req.URL.String(), thing.ID, thing.Nearby.ID)
		}
		req.Body.Close()
		resp.WriteHeader(http.StatusOK)
	})
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		panic(err)
	}

	boundAddr = fmt.Sprintf("http://%s:%d", GetOutboundIP().String(), listener.Addr().(*net.TCPAddr).Port)
	go http.Serve(listener, nil)
}

// Get preferred outbound ip of this machine
func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}
