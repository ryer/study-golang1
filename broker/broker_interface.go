package broker

import "study-golang1/data"

type ItemWork func(data.Item) data.Item

type IBroker interface {
	Invoke(ItemWork)
	Output() chan data.Item
}
