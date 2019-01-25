package image_counter

type ItemWork func(Item) (Item, error)

type IBroker interface {
	Invoke(ItemWork)
	Output() chan Item
}
