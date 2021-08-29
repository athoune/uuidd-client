package uuidd

/*
Spec are here : https://github.com/karelzak/util-linux/blob/master/misc-utils/uuidd.c
*/

import (
	"io"
	"net"
	"sync"

	"github.com/google/uuid"
)

// See https://github.com/karelzak/util-linux/blob/master/libuuid/src/uuidd.h
const (
	UUIDD_OP_GETPID           uint8 = 0
	UUIDD_OP_GET_MAXOP        uint8 = 1
	UUIDD_OP_TIME_UUID        uint8 = 2
	UUIDD_OP_RANDOM_UUID      uint8 = 3
	UUIDD_OP_BULK_TIME_UUID   uint8 = 4
	UUIDD_OP_BULK_RANDOM_UUID uint8 = 5
)

func New(path string) *Client {
	return &Client{
		path: path,
		lock: &sync.Mutex{},
	}
}

type Client struct {
	path string
	lock *sync.Mutex
	conn net.Conn
}

func (c *Client) dial() (net.Conn, error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if c.conn != nil {
		return c.conn, nil
	}
	var err error
	c.conn, err = net.Dial("unix", c.path)
	if err != nil {
		return nil, err
	}
	return c.conn, nil
}

func (c *Client) TimeUUID() (uuid.UUID, error) {
	conn, err := c.dial()
	if err != nil {
		return uuid.Nil, err
	}
	return TimeUUID(conn)
}

func TimeUUID(conn io.ReadWriter) (uuid.UUID, error) {
	_, err := conn.Write([]byte{UUIDD_OP_TIME_UUID})
	if err != nil {
		return uuid.Nil, err
	}
	lenRaw := make([]byte, 4)
	_, err = io.ReadFull(conn, lenRaw)
	if err != nil {
		return uuid.Nil, err
	}
	// FIXME assert lenght
	u := make([]byte, 16)
	_, err = io.ReadFull(conn, u)
	if err != nil {
		return uuid.Nil, err
	}
	return uuid.FromBytes(u)
}
