package image_counter

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

func Main(urlListJsonFile string) {
	target := loadJson(urlListJsonFile)

	var b IBroker
	if false {
		b = broker1(target)
	} else {
		b = broker2(target)
	}

	b.Invoke(work)

	for it := range b.Output() {
		fmt.Println(it.Url)
	}
}

func work(it Item) (Item, error) {
	cl := NewClient(nil)
	cnt, err := cl.CountImages(it.Url)
	if nil != err {
		return Item{}, err
	}
	it.Url = fmt.Sprintf("(%d)%s", cnt, it.Url)
	return it, nil
}

func broker1(target Data) IBroker {
	return NewBroker1(target)
}

func broker2(target Data) IBroker {
	input := make(chan Item)
	go func() {
		for _, it := range target.Items {
			input <- it
		}
		close(input)
	}()

	return NewBroker2(input)
}

func loadJson(jsonFile string) Data {
	fileData, err := ioutil.ReadFile(jsonFile)
	if err != nil {
		panic(err)
	}

	var target Data
	err = json.Unmarshal(fileData, &target)
	if err != nil {
		panic(err)
	}

	return target
}
