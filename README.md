# OpenRedirect

Golang tool to look for open redirect vulnerabilities.

## Usage examples

Scan from a bunch of URLs in a given file:

```bash
openredirect -f urls.txt
```

It also supports pipeline mode:

```bash
cat urls.txt | openredirect
```

Send the request through a proxy:

```bash
openredirect -f urls.txt --proxy 'http://127.0.0.1:8080' --insecure
```

Do not check all the variables, just the most commons for this vulnerability:

```bash
openredirect -f urls.txt --smart
```

All options:

```bash
./openredirect --help
Usage of ./openredirect:
  -f string
    	File with the target URLs
  -insecure
    	Skip TLS verification
  -proxy string
    	Proxy with the format 'http://proxyhost:port'
  -smart
    	Do not check all the params
  -t int
    	Max number of URLs to check at the same time (default 150)
  -tpd int
    	Threads to execute per domain (default 3)
```