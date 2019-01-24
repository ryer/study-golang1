package broker

import (
	"strings"
	"study-golang1/data"
	"testing"
)

func TestBroker1_Invoke(t *testing.T) {
	target := data.Data{Items: []data.Item{
		{Url: "https://gege"},
	}}

	b := NewBroker1(target)

	b.Invoke(func(it data.Item) (data.Item, error) {
		it.Url = strings.Replace(it.Url, "/gege", "/Xyz", -1)
		return it, nil
	})

	it := <-b.Output()
	if "https://Xyz" != it.Url {
		t.Fatalf("Wrong output. expected=%q, got=%q", "https://Xyz", it.Url)
	}
}
