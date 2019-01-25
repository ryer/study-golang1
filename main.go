package main

import (
	"flag"
	"fmt"
	ic "study-golang1/image_counter"
	"study-golang1/monkey"
	"study-golang1/socks4server"
)

var (
	Name     = ""
	Version  = ""
	Revision = ""
)

func main() {
	var (
		urlListJsonFile = flag.String("url-list", "", "画像カウント先URLのJSONデータファイル")
		socks4port      = flag.Int("socks4-port", 0, "SOCKS4サーバのポート番号")
		useMonkey       = flag.Bool("use-monkey", false, "モンキー言語を使う")
	)
	flag.Parse()

	if *urlListJsonFile != "" {
		ic.Main(*urlListJsonFile)
	} else if *socks4port != 0 {
		socks4server.Main(*socks4port)
	} else if *useMonkey {
		monkey.Main()
	} else {
		fmt.Printf("%s %s (%s)\n", Name, Version, Revision)
		flag.Usage()
		return
	}
}
