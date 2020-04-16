package main

import (
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"sync/atomic"
	"time"

	"github.com/juju/ratelimit"
	"github.com/spf13/pflag"
)

const (
	maxChunkSize      = 8192
	chunksPerInterval = 20
)

var (
	version     string
	byteCounter uint64
	reqsCounter uint64
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

func downloadLoop(
	url string,
	chunkSize int64,
	byteRate float64,
	reqsRateBucket *ratelimit.Bucket,
) {
	httpClient := http.Client{
		Transport: &http.Transport{
			MaxIdleConns:        1,
			MaxIdleConnsPerHost: 1,
			IdleConnTimeout:     60 * time.Second,
		},
	}
	buffer := make([]byte, chunkSize)
	var byteRateBucket *ratelimit.Bucket
	if byteRate > 0 {
		byteRateBucket = ratelimit.NewBucketWithRate(byteRate, chunkSize)
		byteRateBucket.TakeAvailable(chunkSize)
	}
	for {
		if reqsRateBucket != nil {
			reqsRateBucket.Wait(1)
		}
		resp, err := httpClient.Get(url)
		if err != nil {
			panic(err)
		}
		atomic.AddUint64(&reqsCounter, 1)
		rateLimitedBody := io.Reader(resp.Body)
		if byteRateBucket != nil {
			rateLimitedBody = ratelimit.Reader(resp.Body, byteRateBucket)
		}
		for {
			bytesCnt, err := rateLimitedBody.Read(buffer)
			if err == io.EOF {
				atomic.AddUint64(&byteCounter, uint64(bytesCnt))
				resp.Body.Close()
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
	pflag.CommandLine.SortFlags = false
	argHelp := pflag.BoolP("help", "h", false, "Just print help message and exit")
	argVersion := pflag.Bool("version", false, "Just print version and exit")
	argClients := pflag.Int64P(
		"clients-num", "n", 1,
		"Number of clients to make request in parallel",
	)
	argBitrate := pflag.Float64P(
		"client-bitrate", "b", 100e3,
		"Per-client download speed limit in bit/s (zero means no limit)",
	)
	argReqsPerSec := pflag.Float64P(
		"total-rps", "r", 1e3,
		"Total requests per second limit for all clients (zero means no limit)",
	)
	argInterval := pflag.Float64P(
		"interval", "i", 1,
		"Report interval in seconds",
	)
	argCount := pflag.Uint64P(
		"count", "c", 0,
		"Stop after given number of intervals (use zero to run non-stop)",
	)
	argURL := pflag.StringP(
		"url", "u", "", "HTTP URL to download (REQUIRED)",
	)
	pflag.Parse()
	if *argHelp {
		pflag.Usage()
		return
	}
	if *argVersion {
		fmt.Println("http-download-speed version " + version)
		return
	}
	clients := *argClients
	interval := *argInterval
	byteRate := *argBitrate / 8.0
	reqsPerSec := *argReqsPerSec
	maxCount := *argCount
	chunkSize := int64(byteRate * interval / chunksPerInterval)
	if byteRate == 0 {
		chunkSize = maxChunkSize
	}
	if chunkSize <= 10 {
		fmt.Printf("Number of bytes per interval is too low, bitrate may be wrong")
		chunkSize = 10
	}
	if chunkSize > maxChunkSize {
		chunkSize = maxChunkSize
	}
	url := *argURL
	if url == "" {
		fmt.Println("No HTTP URL specified")
		pflag.Usage()
		os.Exit(1)
	}
	var reqsRateBucket *ratelimit.Bucket
	if reqsPerSec > 0 {
		reqsRateBucket = ratelimit.NewBucketWithRate(reqsPerSec, clients)
		reqsRateBucket.TakeAvailable(clients)
	}
	for i := int64(0); i < clients; i++ {
		go downloadLoop(url, chunkSize, byteRate, reqsRateBucket)
	}
	fmt.Printf("Time    \tDownload speed (bit/s)\tRequests per second\n")
	var prevByteCnt, prevReqsCnt, counter uint64
	prevTime := time.Now()
	for {
		time.Sleep(time.Duration(interval*1e9) * time.Nanosecond)
		nowTime := time.Now()
		currByteCnt := atomic.LoadUint64(&byteCounter)
		currReqsCnt := atomic.LoadUint64(&reqsCounter)
		intervalNanosec := float64(nowTime.Sub(prevTime).Nanoseconds())
		bytesPerNanosec := float64(currByteCnt-prevByteCnt) / intervalNanosec
		reqsPerNanosec := float64(currReqsCnt-prevReqsCnt) / intervalNanosec
		fmt.Printf(
			"%s\t%s\t%s\n",
			nowTime.Format("15:04:05"),
			bignum2str(bytesPerNanosec*8e9),
			bignum2str(reqsPerNanosec*1e9),
		)
		prevTime = nowTime
		prevByteCnt = currByteCnt
		prevReqsCnt = currReqsCnt
		counter++
		if maxCount > 0 && counter >= maxCount {
			break
		}
	}
}
