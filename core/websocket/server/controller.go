package server

type TextMessageController interface {
	Init(base TextMessageController)
	ParsePayload(c *Client, payload map[string]interface{}) error
	Process() error
}

type BaseTextMessageController struct {
	Client  *Client
	Payload map[string]interface{}
}

func (b *BaseTextMessageController) ParsePayload(client *Client, payload map[string]interface{}) (err error) {
	b.Client = client
	b.Payload = payload
	return
}

func (b *BaseTextMessageController) Init(_ TextMessageController) {}
