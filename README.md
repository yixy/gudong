## how to use ##

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

## about cobra ##

[Cobra](https://github.com/spf13/cobra) is both a library for creating powerful modern CLI applications as well as a program to generate applications and command files.