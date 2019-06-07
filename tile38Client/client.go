package tile38Client

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/hibooboo2/tile38NukeGame/game/model"
)

type Tile38RedisClient struct {
	p   *redis.Pool
	loc string
}

func New(loc string) (*Tile38RedisClient, error) {
	p := redis.NewPool(func() (redis.Conn, error) {
		log.Println("Made new redis client!")
		return redis.Dial("tcp", loc)
	}, 20)
	p.Wait = true
	p.IdleTimeout = time.Second * 30
	p.TestOnBorrow = func(c redis.Conn, t time.Time) error {
		_, err := c.Do("PING")
		if err != nil {
			log.Println(err)
			c.Close()
		}
		return err
	}
	return &Tile38RedisClient{p, loc}, nil
}

func CmdToArgs(cmd string) (string, []interface{}) {
	args := strings.Split(cmd, " ")
	cmdArgs := []interface{}{}
	for i, arg := range args {
		if i == 0 {
			continue
		}
		cmdArgs = append(cmdArgs, arg)
	}
	return args[0], cmdArgs
}

func (client *Tile38RedisClient) post(cmd string) error {
	arg, args := CmdToArgs(cmd)
	con := client.p.Get()
	defer con.Close()
	resp, err := con.Do(arg, args...)
	if err != nil {
		log.Println(err)
		log.Println(cmd)
	}

	if fmt.Sprintf("%v", resp) != "OK" {
		log.Printf("%v", resp)
	}
	return err
}

func (client *Tile38RedisClient) sub(cmdStr string) (chan interface{}, error) {
	vals := make(chan interface{})
	c, err := redis.Dial("tcp", client.loc)
	if err != nil {
		return nil, err
	}
	cmd, args := CmdToArgs(cmdStr)
	err = c.Send(cmd, args...)
	if err != nil {
		return nil, err
	}
	err = c.Flush()
	if err != nil {
		return nil, err
	}
	go func() {
		defer c.Close()
		for {
			r, err := c.Receive()
			if err != nil {
				if strings.Contains(err.Error(), "closed") {
					log.Println("connection closed!")
					os.Exit(0)
				}
				log.Println(err)
				continue
			}
			vals <- r
		}
	}()
	return vals, nil
}

func (client *Tile38RedisClient) Set(collection, key string, lat, lon float64) error {
	return client.post(fmt.Sprintf("SET %s %s POINT %f %f", collection, key, lat, lon))
}

func (client *Tile38RedisClient) Nearby(collectionToWatch string, collectionYouAreIn string, id string, distance float64) (chan model.Thing, error) {
	vals, err := client.sub(fmt.Sprintf("NEARBY %s FENCE ROAM %s %s %f", collectionToWatch, collectionYouAreIn, id, distance))
	if err != nil {
		return nil, err
	}
	things := make(chan model.Thing)
	go func() {
		for val := range vals {
			d, ok := val.([]uint8)
			if !ok {
				log.Printf("Ignored: %s", val)
				continue
			}
			t := model.Thing{}
			err := json.Unmarshal(d, &t)
			if err == nil {
				things <- t
			}
		}
	}()
	return things, nil
}
