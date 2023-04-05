package stream

import (
	"io"
	"net/http"
)

// HTTPClient is an interface for http client.
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// Download downloads a file from a given url and send the data to the send function.
func Download(client HTTPClient, url string, bufferSize int64, send func([]byte)) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	buff := make([]byte, bufferSize)
	for {
		var bytesRead int
		bytesRead, err = res.Body.Read(buff)
		if err == io.EOF {
			err = nil
			// A Reader returning a non-zero number of bytes at the end of the input stream may
			// return either err == EOF or err == nil. The next Read should return 0, EOF.
			if bytesRead <= 0 {
				break
			}
		}
		if err != nil {
			return err
		}
		send(buff[:bytesRead])
	}
	return nil
}
