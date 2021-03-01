package getty

type Data struct {
	Type    int
	Buffers []byte
}

type Message struct {
	Client *Client
	Data   *Data
}
