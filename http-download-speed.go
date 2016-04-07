package main

import (
	"flag"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"sync/atomic"
	"time"
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

func download_loop(url string, counter *uint64, buffsize uint64) {
	for {
		var size uint64 = 0
		res, err := http.Get(url)
		if err != nil {
			panic(err)
		}
		data := make([]byte, buffsize)
		for {
			n, err := res.Body.Read(data)
			if err != nil {
				if err == io.EOF {
					break
				}
				panic(err)
			}
			size += uint64(n)
			if size >= buffsize {
				atomic.AddUint64(counter, uint64(size))
				size = 0
			}
		}
		atomic.AddUint64(counter, uint64(size))
	}
}

func main() {
	routines_flag := flag.Int("routines", 1, "Number of parallel download routines")
	interval_flag := flag.Float64("interval", 1, "Report interval in seconds")
	buffsize_flag := flag.Uint64("buffsize", 50*1000, "Size of download buffers in bytes")
	url_flag := flag.String("url", "", "HTTP URL to download (REQUIRED)")
	flag.Parse()
	routines := *routines_flag
	interval := *interval_flag
	buffsize := *buffsize_flag
	url := *url_flag
	if url == "" {
		fmt.Println("No HTTP URL specified\n")
		flag.Usage()
		os.Exit(1)
	}
	var counter uint64 = 0
	var prev uint64 = 0
	for i := 0; i < routines; i++ {
		go download_loop(url, &counter, buffsize)
	}
	prev_time := time.Now()
	fmt.Printf("Time    \tDownload speed (bit/s)\n")
	for {
		time.Sleep(time.Duration(interval*1e9) * time.Nanosecond)
		now := time.Now()
		cnt := counter
		fmt.Printf("%s\t%s\n", now.Format("15:04:05"), bignum2str(8e9*float64(cnt-prev)/float64(now.Sub(prev_time))))
		prev_time = now
		prev = cnt
	}

}
