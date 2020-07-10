#!/bin/bash
sudo ifconfig lo:1 10.29.1.1/16
sudo ifconfig lo:2 10.29.1.2/16
sudo ifconfig lo:3 10.29.1.3/16

go run main.go -node="10.29.1.1" &
go run main.go -node="10.29.1.2" -cluster="10.29.1.1" &
go run main.go -node="10.29.1.3" -cluster="10.29.1.2" &
