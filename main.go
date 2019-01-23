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

	b := broker.NewBroker(&target)
	b.Run()

	fmt.Println("b.Result")
	for {
		select {
		case r := <-b.Result:
			fmt.Println(r.Url)
		case <-b.Wait:
			goto L
		}
	}

L:
	fmt.Println("target.Items")
	for _, v := range target.Items {
		fmt.Println(v.Url)
	}
}
