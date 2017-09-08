package client

import (
	"github.com/gorilla/websocket"
	"log"
	"net/url"
	"time"
)

//
const (
	_host       = "localhost:8080"
	_uri        = "/whecho"
	_retrytime  = 1e9
	_retrycount = 0
	retry       = "[client]retry conn"
)

//
var Lived bool = false
var Retry_ string

//
type Client struct {
	Host       string
	Uri        string
	Retrytime  time.Duration
	Retrycount int
	conn       *websocket.Conn
}

// Create new chat client.
func NewClient(host string, uri string, retrytime time.Duration, retrycount int) *Client {

	if host == "" {
		host = _host
	}
	if uri == "" {
		uri = _uri
	}

	if retrytime > 0 {
		retrytime = retrytime
	} else {
		retrytime = _retrytime
	}

	if retrycount > 0 {
		retrycount = retrycount
	} else {
		retrycount = _retrycount
	}

	return &Client{host, uri, retrytime, retrycount, nil}
}

//
var url_ string

//
func (c *Client) InitWSConn() {

	u := url.URL{Scheme: "ws", Host: c.Host, Path: c.Uri}
	url_ = u.String()
	log.Printf("connecting to %s", url_)

	c.initConn()
	c.listener()
}

//
func (c *Client) reOpen() {
	//c.conn.Close()
	//log.Println("###################",c.Retrytime)
	time.Sleep(c.Retrytime)

	str := retry
	if "" != Retry_ {
		str = Retry_
	}
	c.initConn(str)

}

//
func (c *Client) initConn(msg ...string) {

	conn, _, err := websocket.DefaultDialer.Dial(url_, nil)
	if err != nil {
		//log.Fatal("<dial>:", err)
		log.Println("<dial>:", err)

		Lived = false
		conn = nil
		c.reOpen()

		return
	}

	Lived = true
	c.conn = conn

	//
	if len(msg) > 0 {
		c.SendMsg(msg[0])
	}
}

//
func (c *Client) listener() {
	go c.read()
}

//
func (c *Client) read() {

	for {
		conn := c.conn
		mt, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("read err:", err)

			c.conn.Close()
			conn = nil
			c.conn = nil
			c.reOpen()
			// return
		}

		messageHandle(mt, message)
	}
}

//
func messageHandle(messageType int, message []byte) {

	switch messageType {
	case websocket.TextMessage:

		log.Printf("recv: %s", message, "[len]:", len(message))

		if nil != readHandle && len(message) > 0 {
			readHandle(message)
		}

	case websocket.BinaryMessage:
		log.Println("Not support binary protocol.")
	}
}

//
func (c *Client) SendMsg(msg string) {

	conn := c.conn
	err := conn.WriteMessage(websocket.TextMessage, []byte(msg))
	if err != nil {
		log.Println("write err:", err)
		//return
	}
}

//
var readHandle func(msg []byte)

func (c *Client) SetReadHandle(fn func(msg []byte)) {
	readHandle = fn
}
