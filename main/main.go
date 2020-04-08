package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"sync/atomic"
	"time"

	"github.com/juju/ratelimit"
	flag "github.com/spf13/pflag"
)

var (
	version     string
	byteCounter uint64
)

func bignum2str(num float64) string {
	for _, suffix := range []string{" ", "K", "M", "G", "T", "P", "E", "Z"} {
		if math.Abs(num) < 1000.0 {
			return fmt.Sprintf("%9.3f%s", num, suffix)
		}
		num /= 1000.0
	}
	return fmt.Sprintf("%9.3f%s", num, "Y")
}

func downloadLoop(url string, chunkSize int64, byteRate float64) {
	rateBucket := ratelimit.NewBucketWithRate(byteRate, chunkSize)
	rateBucket.TakeAvailable(chunkSize)
	for {
		resp, err := http.Get(url)
		if err != nil {
			panic(err)
		}
		rateLimitedBody := ratelimit.Reader(resp.Body, rateBucket)
		for {
			bytesCnt, err := io.CopyN(ioutil.Discard, rateLimitedBody, chunkSize)
			if err == io.EOF {
				resp.Body.Close()
				atomic.AddUint64(&byteCounter, uint64(bytesCnt))
				break
			}
			if err != nil {
				panic(err)
			}
			atomic.AddUint64(&byteCounter, uint64(bytesCnt))
		}
	}
}

func main() {
	flag.CommandLine.SortFlags = false
	argHelp := flag.Bool("help", false, "Just print help message and exit")
	argVersion := flag.Bool("version", false, "Just print version and exit")
	argBitrate := flag.Float64(
		"bitrate", 100, "Max download birate in kbit/s for single goroutine",
	)
	argClients := flag.Int(
		"clients", 1, "Number of parallel download clients",
	)
	argInterval := flag.Float64("interval", 1, "Report interval in seconds")
	argChunksPerInterval := flag.Float64(
		"chunks-per-interval", 4, "Number of download chunks per report interval",
	)
	argURL := flag.String("url", "", "HTTP URL to download (REQUIRED)")
	flag.Parse()
	if *argHelp {
		flag.Usage()
		return
	}
	if *argVersion {
		fmt.Println("http-download-speed version " + version)
		return
	}
	clients := *argClients
	interval := *argInterval
	byteRate := *argBitrate * 1e3 / 8.0
	chunkSize := int64(byteRate * interval / (*argChunksPerInterval))
	if chunkSize <= 10 {
		fmt.Printf("Number of bytes per interval is too low, bitrate may be wrong")
		chunkSize = 10
	}
	url := *argURL
	if url == "" {
		fmt.Println("No HTTP URL specified")
		flag.Usage()
		os.Exit(1)
	}
	for i := 0; i < clients; i++ {
		go downloadLoop(url, chunkSize, byteRate)
	}
	fmt.Printf("Time    \tDownload speed (bit/s)\n")
	var prevByteCnt uint64
	prevTime := time.Now()
	for {
		time.Sleep(time.Duration(interval*1e9) * time.Nanosecond)
		nowTime := time.Now()
		currByteCnt := byteCounter
		intervalNanosec := float64(nowTime.Sub(prevTime).Nanoseconds())
		bytesPerNanosec := float64(currByteCnt-prevByteCnt) / intervalNanosec
		fmt.Printf(
			"%s\t%s\n",
			nowTime.Format("15:04:05"),
			bignum2str(bytesPerNanosec*8e9),
		)
		prevTime = nowTime
		prevByteCnt = currByteCnt
	}
}
