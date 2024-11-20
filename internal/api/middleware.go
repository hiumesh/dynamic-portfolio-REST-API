package api

import (
	"bytes"
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type timeoutResponseWriter struct {
	gin.ResponseWriter
	sync.Mutex

	header      http.Header
	wroteHeader bool
	snapHeader  http.Header // snapshot of the header at the time WriteHeader was called
	statusCode  int
	buf         bytes.Buffer
}

func (t *timeoutResponseWriter) Header() http.Header {
	t.Lock()
	defer t.Unlock()

	return t.header
}

func (t *timeoutResponseWriter) Write(bytes []byte) (int, error) {
	t.Lock()
	defer t.Unlock()

	if !t.wroteHeader {
		t.writeHeaderLocked(http.StatusOK)
	}

	return t.buf.Write(bytes)
}

func (t *timeoutResponseWriter) WriteHeader(statusCode int) {
	t.Lock()
	defer t.Unlock()

	t.writeHeaderLocked(statusCode)
}

func (t *timeoutResponseWriter) writeHeaderLocked(statusCode int) {
	if t.wroteHeader {
		// ignore multiple calls to WriteHeader
		// once WriteHeader has been called once, a snapshot of the header map is taken
		// and saved in snapHeader to be used in finallyWrite
		return
	}

	t.statusCode = statusCode
	t.wroteHeader = true
	t.snapHeader = t.header.Clone()
}

func (t *timeoutResponseWriter) finallyWrite(w gin.ResponseWriter) {
	t.Lock()
	defer t.Unlock()

	dst := w.Header()
	for k, vv := range t.snapHeader {
		dst[k] = vv
	}

	if !t.wroteHeader {
		t.statusCode = http.StatusOK
	}

	w.WriteHeader(t.statusCode)
	if _, err := w.Write(t.buf.Bytes()); err != nil {
		logrus.WithError(err).Warn("Write failed")
	}
}

// TODO
func timeoutMiddleware(timeout time.Duration) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()

		w := c.Writer

		timeoutWriter := &timeoutResponseWriter{
			header: make(http.Header),
		}

		panicChan := make(chan any, 1)
		serverDone := make(chan struct{})
		go func() {
			defer func() {
				if p := recover(); p != nil {
					panicChan <- p
				}
			}()

			c.Request = c.Request.WithContext(ctx)
			c.Writer = timeoutWriter

			c.Next()

			close(serverDone)
		}()

		select {
		case p := <-panicChan:
			panic(p)

		case <-serverDone:
			timeoutWriter.finallyWrite(w)

		case <-ctx.Done():
			err := ctx.Err()

			if err == context.DeadlineExceeded {
				httpError := &HTTPError{
					HTTPStatus: http.StatusGatewayTimeout,
					ErrorCode:  ErrorCodeRequestTimeout,
					Message:    "Processing this request timed out, please retry after a moment.",
				}

				httpError = httpError.WithInternalError(err)
				c.Writer = w

				HandleResponseError(c, httpError)
			} else {
				// unrecognized context error, so we should wait for the server to finish
				// and write out the response
				<-serverDone

				timeoutWriter.finallyWrite(w)
			}
		}
	})

}
