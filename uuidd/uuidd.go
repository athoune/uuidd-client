package uuidd

/*
Spec are here : https://github.com/karelzak/util-linux/blob/master/misc-utils/uuidd.c
*/

import (
	"fmt"
	"io"
	"net"

	"encoding/binary"

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
	}
}

type Client struct {
	path string
}

func (c *Client) dial() (net.Conn, error) {
	return net.Dial("unix", c.path)
}

func (c *Client) TimeUUID() (uuid.UUID, error) {
	conn, err := c.dial()
	if err != nil {
		return uuid.Nil, err
	}
	return TimeUUID(conn)
}

func (c *Client) BulkTimeUUID(n int32, cb func(uuid.UUID) error) error {
	if n == 0 {
		return nil
	}
	if n < 0 {
		return fmt.Errorf("positive value only %d", n)
	}
	for n > 0 {
		conn, err := c.dial()
		if err != nil {
			return err
		}
		b, _, err := BulkTimeUUID(conn, n)
		if err != nil {
			return err
		}
		err = cb(b)
		if err != nil {
			return err
		}
		n--
	}
	return nil
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

func BulkTimeUUID(conn io.ReadWriter, n int32) (uuid.UUID, int32, error) {
	_, err := conn.Write([]byte{UUIDD_OP_BULK_TIME_UUID})
	if err != nil {
		return uuid.Nil, 0, err
	}
	err = binary.Write(conn, binary.LittleEndian, n)
	if err != nil {
		return uuid.Nil, 0, err
	}
	var len int32

	err = binary.Read(conn, binary.LittleEndian, &len)
	if err != nil {
		return uuid.Nil, 0, err
	}
	bulk := make([]byte, 16)
	_, err = io.ReadFull(conn, bulk)
	if err != nil {
		return uuid.Nil, 0, err
	}
	b, err := uuid.FromBytes(bulk)
	if err != nil {
		return uuid.Nil, 0, err
	}

	err = binary.Read(conn, binary.LittleEndian, &len)
	if err != nil {
		return uuid.Nil, 0, err
	}
	return b, len, nil
}
