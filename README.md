# http-download-speed

A simple tool for simulating a load when many clients are downloading some big
file via HTTP with same fixed bitrate.

After parsing command line arguments, given number of parallel clients are
launched using goroutines. Each client sends requests to given URL and downloads
response in infinite loop. Download is made with specified bitrate limit.
During this process, shared counters of downloaded bytes and started requests
are incremented.

In the main goroutine shared counters are measured at regular intervals and
total download bitrate and requests rate are calculated and printed to stdout
in simple human readable format.

## Quick start

Ensure you have go version 1.13 or later:

```sh
go version
```

Clone repository and build binary  (`build/bin/http-download-speed`)

```sh
git clone git@github.com:nikolaynag/http-download-speed.git
cd http-download-speed/
make
make run
```
Usage example:
```
build/bin$ ./http-download-speed --clients-num 4 --client-bitrate 10e3 --count 4 --url http://google.com
Time    	Download speed (bit/s)	Requests per second
22:44:25	   25.787K	    3.999
22:44:26	   40.163K	    0.000
22:44:27	   40.168K	    0.000
22:44:28	   40.661K	    0.000
$
```

Take a look at possible arguments with `--help`:
```
Usage of ./build/bin/http-download-speed:
  -h, --help                   Just print help message and exit
      --version                Just print version and exit
  -n, --clients-num int        Number of clients to make request in parallel (default 1)
  -b, --client-bitrate float   Per-client download speed limit in bit/s (zero means no limit) (default 100)
  -r, --total-rps float        Total requests per second limit for all clients (zero means no limit) (default 10)
  -i, --interval float         Report interval in seconds (default 1)
  -c, --count uint             Stop after given number of intervals (use zero to run non-stop)
  -u, --url string             HTTP URL to download (REQUIRED)
```
