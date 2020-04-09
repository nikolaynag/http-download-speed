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
build/bin$ ./http-download-speed --clients 4 --bitrate 10 --url http://google.com
Time    	Download speed (bit/s)	Requests per second
11:12:12	   39.907K	    3.997
11:12:13	   39.895K	    0.000
11:12:14	   39.900K	    0.000
11:12:15	   39.898K	    0.000
11:12:16	   39.895K	    0.000
^C
$
```
Take a look at possible arguments with `--help`:
```
Usage of ./http-download-speed:
      --help                            Just print help message and exit
      --version                         Just print version and exit
      --bitrate float                   Max download birate in kbit/s for single goroutine (default 100)
      --clients int                     Number of parallel download clients (default 1)
      --interval float                  Report interval in seconds (default 1)
      --min-chunks-per-interval float   Minimum number of download chunks per report interval (default 4)
      --url string                      HTTP URL to download (REQUIRED)
```
