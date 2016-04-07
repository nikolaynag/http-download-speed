# http-download-speed

A simple tool for creating and measuring high-bandwidth multi-connection HTTP
download. After parsing command line arguments, it launches requested number of
goroutines which in infinite loop download body of requested URL. During this
download process, a shared counter is incremented and total bitrate calculated
and printed to stdout in simple human readable table format.

Code does not have any external dependencies, so just clone repo and run `go
build` to get binary.

```
Usage of http-download-speed:
  -buffsize uint
        Size of download buffers in bytes (default 50000)
  -interval float
        Report interval in seconds (default 1)
  -routines int
        Number of parallel download routines (default 1)
  -url string
        HTTP URL to download (REQUIRED)
```

Example:
```
$ http-download-speed -routines 32 -interval 1 -url http://google.com
Time        Download speed (bit/s)
19:30:36       12.192M
19:30:37       15.197M
19:30:38       14.724M
19:30:39       13.312M
19:30:40       12.193M
19:30:41       14.419M
19:30:42       13.154M
^C
$
```
