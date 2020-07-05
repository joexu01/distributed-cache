package main

import (
	"flag"
	"fmt"
	"github.com/joexu01/distributed-cache/cache-benchmark/cacheClient"
)

var (
	server = flag.String("s", "localhost", "cache server address")
	operation = flag.String("c", "get", "command, could be get/set/del")
	key = flag.String("k", "", "key")
	value = flag.String("v", "", "value")
)

func init() {
	flag.Usage = func() {
		_, _ = fmt.Fprintf(flag.CommandLine.Output(), "More information at %s\n", `https://github.com/joexu01/distributed-cache`)
		_, _ = fmt.Fprintf(flag.CommandLine.Output(), "Usage of Cache Client:\n")
		flag.PrintDefaults()
	}
}

func main() {
	flag.Parse()
	client := cacheClient.New("tcp", *server)
	cmd := &cacheClient.Cmd{
		Name:  *operation,
		Key:   *key,
		Value: *value,
		Error: nil,
	}
	client.Run(cmd)
	if cmd.Error != nil {
		fmt.Println("error: ", cmd.Error)
	} else {
		fmt.Println(cmd.Value)
	}
}
