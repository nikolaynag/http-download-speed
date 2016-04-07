package main

import (
	"flag"
	"net/http"
	"os"
	"fmt"
	"io"
	"sync/atomic"
	"time"
	"math"
)
func bignum2str(num float64) string {
    for _, suffix := range []string{" ","K","M","G","T","P","E","Z"} {
        if math.Abs(num) < 1000.0 {
            return fmt.Sprintf("%9.3f%s", num, suffix)
		}
        num /= 1000.0
	}
    return fmt.Sprintf("%9.3f%s", num, "Y")
}

func download(url string, counter *uint64, chunk_size int) {
	for {
		res, err := http.Get(url)
		if err != nil {
			panic(err)
		}
		data := make([]byte, chunk_size)
		for {
			n, err := res.Body.Read(data)
			if err == io.EOF {
				break
			}
			if err != nil {
				panic(err)
			}
			atomic.AddUint64(counter, uint64(n))
		}
	}
}

func main() {
	routines_flag := flag.Int("routines", 1, "Number of parallel download routines")
	interval_flag := flag.Float64("interval", 1, "Report interval in seconds")
	buffsize_flag := flag.Int("buffsize", 500*1000, "Size of download buffers in bytes")
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
		go download(url, &counter, buffsize)
	}
	fmt.Printf("Time    \tDownload speed (bit/s)\n")
	for {
		time.Sleep(time.Duration(interval*1e9)*time.Nanosecond)
		fmt.Printf("%s\t%s\n", time.Now().Format("15:04:05"), bignum2str(8*float64(counter - prev)/interval))
		prev = counter
	}

}
