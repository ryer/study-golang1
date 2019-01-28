package main

import (
	"flag"
	"fmt"
	"study-golang1/echosrv"
	"study-golang1/image_counter"
	"study-golang1/monkey"
	"study-golang1/socks_server"
)

var (
	Name     = ""
	Version  = ""
	Revision = ""
)

func main() {
	var (
		urlListJsonFile = flag.String("url-list", "", "画像カウント先URLのJSONデータファイル")
		socksPort      = flag.Int("socks-port", 0, "SOCKSサーバのポート番号")
		useMonkey       = flag.Bool("use-monkey", false, "モンキー言語を使う")
		echoPort        = flag.Int("echo-port", 0, "ECHOサーバのポート番号")
	)
	flag.Parse()

	if *urlListJsonFile != "" {
		image_counter.Main(*urlListJsonFile)
	} else if *socksPort != 0 {
		socks_server.Main(*socksPort)
	} else if *useMonkey {
		monkey.Main()
	} else if *echoPort != 0 {
		echosrv.Main(*echoPort)
	} else {
		fmt.Printf("%s %s (%s)\n", Name, Version, Revision)
		flag.Usage()
		return
	}
}
