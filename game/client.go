package game

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
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

func (c *Client) Notifications(name string, hookurl string) bool {
	// http://10.14.12.68:8081
	cmd := fmt.Sprintf("SETHOOK %[1]s %[2]s/%[1]s NEARBY fleet MATCH %[1]s FENCE ROAM fleet * 1000", name, hookurl)
	fmt.Println(cmd)
	return c.post(cmd)
}

func (c *Client) post(cmd string) bool {
	// http://10.14.12.11:9851
	resp, err := c.c.Post(c.baseUrl, "", strings.NewReader(cmd))
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	val := map[string]interface{}{}
	json.NewDecoder(resp.Body).Decode(&val)
	if len(val) > 2 {
		fmt.Println(val)
		return false
	}
	fmt.Println(val["ok"])
	return val["ok"].(bool)
}
