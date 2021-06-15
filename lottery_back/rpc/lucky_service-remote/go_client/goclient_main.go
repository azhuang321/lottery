// Autogenerated by Thrift Compiler (0.11.0)
// DO NOT EDIT UNLESS YOU ARE SURE THAT YOU KNOW WHAT YOU ARE DOING

package main

import (
	"context"
	"flag"
	"fmt"
	"git.apache.org/thrift.git/lib/go/thrift"
	"imooc.com/lottery/comm"
	"imooc.com/lottery/rpc"
	"net"
	"net/url"
	"os"
	"strings"
	"time"
)

func Usage() {
	fmt.Fprintln(os.Stderr, "Usage of ", os.Args[0], " [-h host:port] [-u url] [-f[ramed]] function [arg1 [arg2...]]:")
	flag.PrintDefaults()
	fmt.Fprintln(os.Stderr, "\nFunctions:")
	fmt.Fprintln(os.Stderr, "  DataResult DoLucky(i64 uid, string username, string ip, i64 now, string app, string sign)")
	fmt.Fprintln(os.Stderr, "   MyPrizeList(i64 uid, string username, string ip, i64 now, string app, string sign)")
	fmt.Fprintln(os.Stderr)
	os.Exit(0)
}

func main() {
	flag.Usage = Usage
	var host string = "localhost"
	var port int = 8080
	var protocol string = "json"
	var urlString string = "http://localhost:8080/rpc"
	var framed bool = true
	var useHttp bool = true
	var parsedUrl *url.URL
	var trans thrift.TTransport

	if len(urlString) > 0 {
		var err error
		parsedUrl, err = url.Parse(urlString)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error parsing URL: ", err)
			flag.Usage()
		}
		host = parsedUrl.Host
		useHttp = len(parsedUrl.Scheme) <= 0 || parsedUrl.Scheme == "http"
	} else if useHttp {
		_, err := url.Parse(fmt.Sprint("http://", host, ":", port))
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error parsing URL: ", err)
			flag.Usage()
		}
	}

	var err error
	if useHttp {
		trans, err = thrift.NewTHttpClient(parsedUrl.String())
	} else {
		portStr := fmt.Sprint(port)
		if strings.Contains(host, ":") {
			host, portStr, err = net.SplitHostPort(host)
			if err != nil {
				fmt.Fprintln(os.Stderr, "error with host:", err)
				os.Exit(1)
			}
		}
		trans, err = thrift.NewTSocket(net.JoinHostPort(host, portStr))
		if err != nil {
			fmt.Fprintln(os.Stderr, "error resolving address:", err)
			os.Exit(1)
		}
		if framed {
			trans = thrift.NewTFramedTransport(trans)
		}
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error creating transport", err)
		os.Exit(1)
	}
	defer trans.Close()
	var protocolFactory thrift.TProtocolFactory
	switch protocol {
	case "compact":
		protocolFactory = thrift.NewTCompactProtocolFactory()
		break
	case "simplejson":
		protocolFactory = thrift.NewTSimpleJSONProtocolFactory()
		break
	case "json":
		protocolFactory = thrift.NewTJSONProtocolFactory()
		break
	case "binary", "":
		protocolFactory = thrift.NewTBinaryProtocolFactoryDefault()
		break
	default:
		fmt.Fprintln(os.Stderr, "Invalid protocol specified: ", protocol)
		Usage()
		os.Exit(1)
	}
	iprot := protocolFactory.GetProtocol(trans)
	oprot := protocolFactory.GetProtocol(trans)
	// thrfit服务的客户端程序
	client := rpc.NewLuckyServiceClient(thrift.NewTStandardClient(iprot, oprot))
	if err := trans.Open(); err != nil {
		fmt.Fprintln(os.Stderr, "Error opening socket to ", host, ":", port, " ", err)
		os.Exit(1)
	}

	// 远程方法调用的参数
	var cmd string = "DoLucky"
	//var cmd string = "MyPrizeList"
	var uid int64 = 1
	username := "admin"
	ip := "127.0.0.1"
	now := time.Now().UnixNano()
	app := "web"
	str := fmt.Sprintf("uid=%d&username=%s&ip=%s&now=%d&app=%s",
		uid, username, ip, now, app)
	sign := comm.CreateSign(str)
	switch cmd {
	case "DoLucky":
		fmt.Println(cmd)
		fmt.Print(client.DoLucky(context.Background(), uid, username, ip, now, app, sign))
		fmt.Print("\n")
		break
	case "MyPrizeList":
		fmt.Println(cmd)
		fmt.Print(client.MyPrizeList(context.Background(), uid, username, ip, now, app, sign))
		fmt.Print("\n")
		break
	case "":
		Usage()
		break
	default:
		fmt.Fprintln(os.Stderr, "Invalid function ", cmd)
	}
}
