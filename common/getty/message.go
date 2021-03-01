package getty

type Data struct {
	Type     float64
	JsonData interface{}
}

type Message struct {
	Client *Client
	Data   *Data
}
