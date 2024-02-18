package stream_test

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/rayspock/mastering-go-examples/stream"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"io"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"testing"
	"time"
)

const (
	// TestServerPort to use the port that is not in use on the local machine.
	TestServerPort = 50432
)

func TestHttpClientSuite(t *testing.T) {
	suite.Run(t, new(HttpClientSuite))
}

type HttpClientSuite struct {
	suite.Suite
}

func (s *HttpClientSuite) SetupSuite() {
	go HandleRequests()
}

func (s *HttpClientSuite) TearDownSuite() {
	p, _ := os.FindProcess(os.Getpid())
	p.Signal(syscall.SIGINT)
}

func (s *HttpClientSuite) TestReliableDownloadWithHttpServer() {
	// We run it several times to make sure the result is what we expect each time as we experience an arbitrary result
	// if we didn't check if there was a last chunk of data returned by the http body reader, even though the read
	// return err == EOF
	for i := 0; i < 10; i++ {
		s.T().Run(strconv.Itoa(i), func(t *testing.T) {
			t.Parallel()
			// For demonstration purposes, we sleep for 2 seconds to make sure asciinema can capture the entire output.
			time.Sleep(2 * time.Second)
			expectedBytes := []byte("some dummy data to test download")
			var dataRcv bytes.Buffer
			w := bufio.NewWriter(&dataRcv)
			bufferSize := int64(5)
			send := func(b []byte) {
				_, err := w.Write(b)
				assert.NoError(t, err)
			}
			err := stream.ReliableDownload(&http.Client{}, fmt.Sprintf("http://localhost:%d/file", TestServerPort), bufferSize, send)
			assert.NoError(t, err)
			err = w.Flush()
			assert.NoError(t, err)
			assert.Equal(t, expectedBytes, dataRcv.Bytes())
		})
	}
}

func (s *HttpClientSuite) TestReliableDownloadWithMockClient() {
	// It is idempotent if we do a mock client test.
	for i := 0; i < 10; i++ {
		s.T().Run(strconv.Itoa(i), func(t *testing.T) {
			t.Parallel()
			// For demonstration purposes, we sleep for 2 seconds to make sure asciinema can capture the entire output.
			time.Sleep(2 * time.Second)
			expectedBytes := []byte("some dummy data")
			mockClient := &mockHTTPClient{
				File: expectedBytes,
			}
			var dataRcv bytes.Buffer
			w := bufio.NewWriter(&dataRcv)
			bufferSize := int64(5)
			send := func(b []byte) {
				_, err := w.Write(b)
				assert.NoError(t, err)
			}
			err := stream.ReliableDownload(mockClient, "", bufferSize, send)
			assert.NoError(t, err)
			err = w.Flush()
			assert.NoError(t, err)
			assert.Equal(t, expectedBytes, dataRcv.Bytes())

		})
	}
}

func (s *HttpClientSuite) TestNaughtyDownloadWithMockClient() {
	// It is idempotent if we do a mock client test.
	for i := 0; i < 10; i++ {
		s.T().Run(strconv.Itoa(i), func(t *testing.T) {
			t.Parallel()
			// For demonstration purposes, we sleep for 2 seconds to make sure asciinema can capture the entire output.
			time.Sleep(2 * time.Second)
			expectedBytes := []byte("some dummy data")
			mockClient := &mockHTTPClient{
				File: expectedBytes,
			}
			var dataRcv bytes.Buffer
			w := bufio.NewWriter(&dataRcv)
			bufferSize := int64(5)
			send := func(b []byte) {
				_, err := w.Write(b)
				assert.NoError(t, err)
			}
			err := stream.NaughtyDownload(mockClient, "", bufferSize, send)
			assert.NoError(t, err)
			err = w.Flush()
			assert.NoError(t, err)
			assert.Equal(t, expectedBytes, dataRcv.Bytes())

		})
	}
}

func (s *HttpClientSuite) TestNaughtyDownloadWithHttpServer() {
	// We run it several times to make sure the result is what we expect each time as we experience an arbitrary result
	// if we didn't check if there was a last chunk of data returned by the http body reader, even though the read
	// return err == EOF
	for i := 0; i < 10; i++ {
		s.T().Run(strconv.Itoa(i), func(t *testing.T) {
			t.Parallel()
			// For demonstration purposes, we sleep for 2 seconds to make sure asciinema can capture the entire output.
			time.Sleep(2 * time.Second)
			expectedBytes := []byte("some dummy data to test download")
			var dataRcv bytes.Buffer
			w := bufio.NewWriter(&dataRcv)
			bufferSize := int64(5)
			send := func(b []byte) {
				_, err := w.Write(b)
				assert.NoError(t, err)
			}
			err := stream.NaughtyDownload(&http.Client{}, fmt.Sprintf("http://localhost:%d/file", TestServerPort), bufferSize, send)
			assert.NoError(t, err)
			err = w.Flush()
			assert.NoError(t, err)
			assert.Equal(t, expectedBytes, dataRcv.Bytes())
		})
	}
}

// mockHTTPClient is a mock http client.
type mockHTTPClient struct {
	File []byte
}

// Do is a mock http client Do method.
func (m *mockHTTPClient) Do(_ *http.Request) (*http.Response, error) {
	return &http.Response{
		Body: io.NopCloser(bytes.NewReader(m.File)),
	}, nil
}

func HandleRequests() {
	http.HandleFunc("/file", fileHandler)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", TestServerPort))
	if err != nil {
		panic(err)
	}

	go func() {
		err = http.Serve(listener, nil)
		if err != nil {
			panic(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	fmt.Println("quiting http server")
}

func fileHandler(writer http.ResponseWriter, req *http.Request) {
	// Set up the HTTP response header
	writer.Header().Set("content-type", "application/pdf")
	writer.Header().Set("transfer-encoding", "chunked")

	// Write the header to the response
	writer.WriteHeader(http.StatusOK)
	writer.(http.Flusher).Flush()

	// Write the chunk to the response
	b := []byte("some dummy data to test download")
	_, err := writer.Write(b)
	if err != nil {
		panic(err)
	}
	writer.(http.Flusher).Flush()
}
