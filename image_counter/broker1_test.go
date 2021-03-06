package image_counter

import (
	"fmt"
	"strings"
	"testing"
)

func TestBroker1_Invoke(t *testing.T) {
	target := Data{Items: []Item{
		{Url: "https://gege"},
	}}

	b := NewBroker1(target)

	b.Invoke(func(it Item) (Item, error) {
		it.Url = strings.Replace(it.Url, "/gege", "/Xyz", -1)
		return it, nil
	})

	it := <-b.Output()
	if "https://Xyz" != it.Url {
		t.Fatalf("Wrong output. expected=%q, got=%q", "https://Xyz", it.Url)
	}

	it, ok := <-b.Output()
	if ok {
		t.Fatalf("Illegal fetch. expected=%s, got=%s", fmt.Sprint(false), fmt.Sprint(ok))
	}
}
