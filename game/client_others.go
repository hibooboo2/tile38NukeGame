// +build !js,!wasm

package game

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/hibooboo2/tile38NukeGame/tile38Client"

	"github.com/pkg/errors"
)

type Client struct {
	c       *http.Client
	baseUrl string
	tc      *tile38Client.Tile38RedisClient
}

func NewClient(baseUrl string, name string) *Client {
	c, err := tile38Client.New(":9851")
	if err != nil {
		panic(err)
	}
	return &Client{&http.Client{Timeout: time.Second * 2}, baseUrl, c}
}

func (c *Client) post(cmd string) error {
	// Tile38ServerURL
	// s := time.Now()
	// defer func() {
	// 	log.Println(time.Since(s))
	// }()
	// timePost := time.Now()

	req, err := http.NewRequest(http.MethodPost, c.baseUrl, strings.NewReader(cmd))
	if err != nil {
		return errors.Wrap(err, "failed to create request")
	}
	resp, err := c.c.Do(req)
	if err != nil {
		return errors.Wrap(err, "failed to do request")
	}
	// log.Println(time.Since(timePost).Nanoseconds() / int64(time.Millisecond))
	val := map[string]interface{}{}
	// timePost = time.Now()
	json.NewDecoder(resp.Body).Decode(&val)
	if len(val) != 2 {
		log.Println(val)
		return errors.Errorf("decoded val incorrect!")
	}
	// log.Println(time.Since(timePost).Nanoseconds() / int64(time.Millisecond))
	if !val["ok"].(bool) {
		return errors.Errorf("Request failed")
	}
	return nil
}

func (c *Client) Notifications(name string) error {
	// hookurl := "http://10.14.12.68:8081"
	// hookurl = "http://localhost:8081"

	// cmd := fmt.Sprintf("SETHOOK %[1]s %[2]s/%[1]s NEARBY fleet FENCE ROAM fleet %[1]s 1000", name, hookurl)
	// log.Println("Noticfications made!", cmd)
	// return c.post(cmd)
	return nil
}

func ClearNotifications() {
	c := NewClient(Tile38ServerURL, "")
	charLock.Lock()
	defer charLock.Unlock()
	for char := range chars {
		cmd := fmt.Sprintf("DELHOOK %s", char)
		c.post(cmd)
	}
}
