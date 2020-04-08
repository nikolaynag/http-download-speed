# http-download-speed

A simple tool for simulating high-bandwidth multiple client HTTP requests load.
After parsing command line arguments, it launches requested number of
goroutines which in infinite loop request given URL and download response body.
During this download process, a shared counter is incremented and total bitrate
calculated and printed to stdout in simple human readable table format.

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
$ ./http-download-speed --clients 4 --bitrate 10 --url http://google.com
Time    	Download speed (bit/s)
23:16:03	   39.920K
23:16:04	   39.883K
23:16:05	   39.909K
23:16:06	   39.902K
23:16:07	   39.902K
^C
$
```
Take a look at possible arguments with `--help`:
```
Usage of http-download-speed:
      --help                        Just print help message and exit
      --version                     Just print version and exit
      --bitrate float               Max download birate in kbit/s for single goroutine (default 100)
      --clients int                 Number of parallel download clients (default 1)
      --interval float              Report interval in seconds (default 1)
      --chunks-per-interval float   Number of download chunks per report interval (default 4)
      --url string                  HTTP URL to download (REQUIRED)
```
