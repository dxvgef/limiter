package limiter

import (
	"context"
	"io"
	"net/http"
	"os"
	"time"

	"golang.org/x/time/rate"
)

// 服务端文件传输限速
func ServeFile(resp http.ResponseWriter, req *http.Request, filePath string, speed float64) error {
	// 当前连接的限速500KB/s
	speedLimiter := NewSpeedLimiter(speed)

	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}

	http.ServeContent(
		resp,
		req,
		filePath,
		fileInfo.ModTime(),
		NewReadSeeker(file, speedLimiter),
	)
	return nil
}

type Reader struct {
	reader  io.Reader
	limiter *rate.Limiter
}

func (reader *Reader) Read(buf []byte) (int, error) {
	n, err := reader.reader.Read(buf)
	if n <= 0 {
		return n, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = reader.limiter.WaitN(ctx, n)
	return n, err
}

type ReadSeeker struct {
	io.ReadSeeker
	reader io.Reader
}

func (rs ReadSeeker) Read(p []byte) (int, error) {
	return rs.reader.Read(p)
}

func NewSpeedLimiter(speed float64) *rate.Limiter {
	return rate.NewLimiter(rate.Limit(speed), int(speed))
}

func NewReader(reader io.Reader, limiter *rate.Limiter) io.Reader {
	return &Reader{
		reader:  reader,
		limiter: limiter,
	}
}

func NewReadSeeker(readSeeker io.ReadSeeker, limiter *rate.Limiter) io.ReadSeeker {
	return ReadSeeker{
		readSeeker,
		NewReader(readSeeker, limiter),
	}
}