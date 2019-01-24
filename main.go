package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"strings"
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

	target := loadJson(*jsonFile)

	var b broker.IBroker
	if false {
		b = broker1(target)
	} else {
		b = broker2(target)
	}

	b.Invoke(func(it data.Item) data.Item {
		it.Url = strings.Replace(it.Url, "https://", "/", -1)
		return it
	})

	for it := range b.Output() {
		fmt.Println(it.Url)
	}
}

func broker1(target *data.Data) broker.IBroker {
	return broker.NewBroker1(target)
}

func broker2(target *data.Data) broker.IBroker {
	input := make(chan data.Item)
	go func() {
		for _, it := range target.Items {
			input <- it
		}
		close(input)
	}()

	return broker.NewBroker2(input)
}

func loadJson(jsonFile string) *data.Data {
	fileData, err := ioutil.ReadFile(jsonFile)
	if err != nil {
		panic(err)
	}

	var target data.Data
	err = json.Unmarshal(fileData, &target)
	if err != nil {
		panic(err)
	}

	return &target
}
