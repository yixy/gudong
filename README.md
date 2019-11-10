## How to use ##

Start a tiny http server using gudong.

```
$ ./gudong start -B=hello,world!
```

Call gudong server and you will get the specified response.

```
$ curl http://localhost:7777/mock.do -d mytest
hello,world! 
```

You can check up the request headers and body in standard output of gudong.

```
$ ./target/gudong start -B=hello,world! 
POST /mock.do HTTP/1.1
HOST : localhost:7777
Content-Type : application/x-www-form-urlencoded
User-Agent : curl/7.54.0
Accept : */*
Content-Length : 6

mytest
========================================
#                gudong                #
========================================
```

Using files to specify response.

```
$ cat > headers  <<EOF                  
Custom-Header : This is a custom header
EOF

$ cat > body  <<EOF                  
Hello,world!
EOF

$ ./gudong start --header-file=headers --body-file=body
```

Using `gudong start -h` to see more.

```
$ ./gudong start -h                      
This tiny server print http request header and body to Standard output, 
                and return the response specified by -H or -B flag.

Usage:
  gudong start [flags]

Flags:
  -B, --body string          specify response body by string
      --body-file string     specify response body by file (if --body-file is specified, -B will be ignored)
  -H, --header string        specify response header by string (multi-headers separate by ;)
      --header-file string   specify response header by file (if --header-file is specified, -H will be ignored)
  -h, --help                 help for start
  -c, --no-chunked           specify body transfer encoding with noChunked when using --body-file
  -p, --port string          specify http server port (default "7777")
  -r, --read-timeout int     specify http server read timeout (ms) (default 3000)
  -w, --write-timeout int    specify http server write timeout (ms) (default 3000)
```

## About gudong ##

[gudong's comming](https://zh.wikipedia.org/wiki/%E5%92%95%E5%92%9A%E6%9D%A5%E4%BA%86) is a Chinese fable.

## About cobra ##

[Cobra](https://github.com/spf13/cobra) is both a library for creating powerful modern CLI applications as well as a program to generate applications and command files.