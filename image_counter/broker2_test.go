package image_counter

import (
	"fmt"
	"strings"
	"testing"
)

func TestBroker2_Invoke(t *testing.T) {
	target := Data{Items: []Item{
		{Url: "https://gege"},
	}}

	input := make(chan Item)
	go func() {
		for _, it := range target.Items {
			input <- it
		}
		close(input)
	}()

	b := NewBroker2(input)

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
