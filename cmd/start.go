/*
Copyright © 2019 yixy <youzhilane01@gmail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/yixy/gudong/log"

	"github.com/spf13/cobra"
)

var port *string
var readTimeout *int64
var writeTimeout *int64
var header *string
var body *string
var HeaderFile *string
var bodyFile *string
var noChunked *bool
var lineSeparate = "\n"
var logLevel *string

const WINDOWS = "windows"

func init() {
	if runtime.GOOS == WINDOWS {
		lineSeparate = "\r\n"
	}

	rootCmd.AddCommand(startCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// startCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	header = startCmd.Flags().StringP("header", "H", "", "specify response header by string "+
		"(multi-headers separate by ;)")
	body = startCmd.Flags().StringP("body", "B", "", "specify response body by string")
	HeaderFile = startCmd.Flags().String("header-file", "", "specify response header by file "+
		"(if --header-file is specified, -H will be ignored)")
	bodyFile = startCmd.Flags().String("body-file", "", "specify response body by file "+
		"(if --body-file is specified, -B will be ignored)")
	noChunked = startCmd.Flags().BoolP("no-chunked", "c", false, "specify body transfer encoding with noChunked when using --body-file")
	port = startCmd.Flags().StringP("port", "p", "7777", "specify http server port")
	readTimeout = startCmd.Flags().Int64P("read-timeout", "r", 3000, "specify http server read timeout (ms)")
	writeTimeout = startCmd.Flags().Int64P("write-timeout", "w", 3000, "specify http server write timeout (ms)")
	logLevel = startCmd.Flags().StringP("log-level", "l", "debug", "log level should be debug,error, default is debug")
}

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "start a echo http server",
	Long: `This tiny server print http request header and body to Standard output, 
		and return the response specified by -H or -B flag.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLogLevel(strings.ToUpper(*logLevel))
		mux := http.NewServeMux()
		mux.HandleFunc("/", mockHandler)
		server := http.Server{
			Addr:         fmt.Sprintf(":%s", *port),
			Handler:      mux,
			ReadTimeout:  time.Millisecond * time.Duration(*readTimeout),
			WriteTimeout: time.Millisecond * time.Duration(*writeTimeout),
		}
		err := server.ListenAndServe()
		if err != nil {
			panic(err)
		}
	},
}

func mockHandler(w http.ResponseWriter, r *http.Request) {
	log.Debug("%s %s %s\r\n", r.Method, r.RequestURI, r.Proto)
	// In Golang, https://github.com/golang/go/issues/7682
	// For incoming requests, the Host header is promoted to the
	// Request.Host field and removed from the Header map.
	log.Debug("HOST : %s\r\n", r.Host)
	for key, values := range r.Header {
		log.Debug("%s : %s\r\n", key, strings.Join(values, ","))
	}
	log.Debug("\r\n")
	if log.Level < log.ERROR {
		io.Copy(os.Stdout, r.Body)
	}

	log.Debug("\n========================================\n")
	log.Debug("#                gudong                #\n")
	defer func() {
		log.Debug("========================================\n")
	}()
	if *HeaderFile != "" {
		//read headers from head-file
		bytes, err := ioutil.ReadFile(*HeaderFile)
		if err != nil {
			log.Error("# Error: %s\n", err.Error())
			w.WriteHeader(500)
			return
		}
		lines := strings.Split(string(bytes), lineSeparate)
		for _, v := range lines {
			key, value := ParseKV(v)
			if key != "" {
				w.Header().Set(key, value)
			}
		}
	} else {
		//read headers from string specified by -H flag
		lines := strings.Split(*header, ";")
		for _, v := range lines {
			key, value := ParseKV(v)
			if key != "" {
				w.Header().Set(key, value)
			}
		}
	}
	if *bodyFile != "" {
		//read body from body-file
		if *noChunked {
			bytes, err := ioutil.ReadFile(*bodyFile)
			if err != nil {
				log.Error("# Error: %s\n", err.Error())
				w.WriteHeader(500)
				return
			}
			_, err = w.Write(bytes)
			if err != nil {
				log.Error("# Error: %s\n", err.Error())
				w.WriteHeader(500)
				return
			}
		} else {
			file, err := os.Open(*bodyFile)
			if err != nil {
				log.Error("# Error: %s\n", err.Error())
				w.WriteHeader(500)
				return
			}
			defer file.Close()
			_, err = io.Copy(w, file)
			if err != nil {
				log.Error("# Error: %s\n", err.Error())
				w.WriteHeader(500)
				return
			}
		}
	} else {
		//read body from string specified by -B flag
		_, err := w.Write([]byte(*body))
		if err != nil {
			log.Error("# Error: %s\n", err.Error())
			w.WriteHeader(500)
			return
		}
	}
}

func ParseKV(s string) (key, value string) {
	index := strings.Index(s, ":")
	if index != -1 {
		//For arrays or strings, the indices are in range if 0 <= low <= high <= len(a), otherwise they are out of range.
		//For slices, the upper index bound is the slice capacity cap(a) rather than the length.
		//
		//ref: https://golang.org/ref/spec#Slice_expressions
		//
		//The indices low and high select which elements of operand a appear in the result.
		//The result has indices starting at 0 and length equal to high - low
		//
		//For convenience, any of the indices may be omitted.
		//A missing low index defaults to zero; a missing high index defaults to the length of the sliced operand.
		//
		//the index x is in range if 0 <= x < len(a), otherwise it is out of range
		//
		//ref: https://golang.org/ref/spec#Index_expressions
		key = strings.TrimSpace(s[:index])
		value = strings.TrimSpace(s[index+1:])
	}
	return key, value
}
