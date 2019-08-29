package util

import (
	"encoding/hex"
	"errors"
	"fmt"
	"net"
	"time"
)

func EscapeStr(str string) ([]byte, error) {
	start := -1
	end := -1
	for i, v := range str {
		if v == '$' {
			if start != -1 {
				end = i
			} else {
				start = i
			}
		}
	}
	if start != -1 && end == -1 {
		return nil, errors.New("unmatched $, need pair of $")
	}

	if start == -1 {
		return hex.DecodeString(str)
	}
	escaped, err := hex.DecodeString(str[:start])
	if err != nil {
		return nil, err
	}
	escaped = append(escaped, []byte(str[start+1:end])...)
	arr, err := hex.DecodeString(str[end+1:])
	if err != nil {
		return nil, err
	}
	escaped = append(escaped, arr...)
	return escaped, nil
}

func PrintReceived(conn *net.TCPConn) {
	buf := make([]byte, 1024)
	for n, err := conn.Read(buf); err == nil; n, err = conn.Read(buf) {
		if n > 0 {
			fmt.Println()
			fmt.Println("--------------------receive msg start--------------------")
			fmt.Println(hex.Dump(buf[:n]))
			fmt.Println("---------------------receive msg end---------------------")
		} else {
			time.Sleep(500 * time.Millisecond)
		}
	}
}
