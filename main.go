package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"study-golang1/broker"
	"study-golang1/data"
)

var (
	Name     = ""
	Version  = ""
	Revision = ""
)

func main() {
	var (
		jsonFile = flag.String("json", "", "JSONデータファイル")
		textFile = flag.String("text", "", "TEXTデータファイル")
	)
	flag.Parse()

	if *jsonFile == "" && *textFile == "" {
		fmt.Printf("%s %s (%s)\n", Name, Version, Revision)
		flag.Usage()
		return
	}

	fileData, err := ioutil.ReadFile(*jsonFile)
	if err != nil {
		panic(err)
	}

	var target data.Data
	err = json.Unmarshal(fileData, &target)
	if err != nil {
		panic(err)
	}

	b := broker.NewBroker2()
	b.Run()

	go func() {
		for _, it := range target.Items {
			b.Input <- it
		}
		close(b.Input)
	}()

	for {
		select {
		case output := <-b.Output:
			fmt.Println(output.Url)
		case <-b.Exit:
			goto LAST
		}
	}

LAST:
}
