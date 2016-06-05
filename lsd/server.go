package main

import (
	"bufio"
	"bytes"
	"io"
)

func ServeClient(client Client, backend Backend) {
	defer client.Close()
	defer client.HandleErrors()

	log := client.Log()
	log.Debug("new connection")

	req := NewRequest(backend, client.Log())
	err := MapRequests(client,
		func(r string) bool {
			defer req.HandleErrors()
			rsp := req.ProcessRequest(r)
			rsp.SendResponse(client)
			return rsp.Keepalive()
		})

	if err != nil {
		log.Error("scanner error:", err)
	}

	log.Debug("end connection")
}

func MapRequests(reader io.Reader, callback func(string) bool) error {
	scanner := bufio.NewScanner(reader)
	scanner.Split(ScanRequest)

	for scanner.Scan() {
		if !callback(scanner.Text()) {
			break
		}
	}

	return scanner.Err()
}

func ScanRequest(data []byte, atEOF bool) (int, []byte, error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}

	i := -1
loop:
	// Look for "End of Request", signaled by an empty line
	for i < len(data) {
		switch n := bytes.IndexByte(data[i+1:], '\n'); {
		case i >= 0 && n == 0:
			return i + 2, data[0:i], nil
		case n >= 0:
			i += n + 1
		default:
			i++
			break loop
		}
	}

	// If we're at EOF, return what we've got
	if atEOF {
		if i == len(data) {
			// drop ending newline
			return len(data), data[0 : i-1], nil
		} else {
			return len(data), data, nil
		}
	}

	// Request more data.
	return 0, nil, nil
}
