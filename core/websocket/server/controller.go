package server

type TextMessageController interface {
	Init(base TextMessageController)
	ParsePayload(c *Client, message ReceiveMessage) error
	Process() error
}

type BaseTextMessageController struct {
	Client       *Client
	Action       string
	ActionParams map[string]interface{}
}

func (b *BaseTextMessageController) ParsePayload(client *Client, message ReceiveMessage) (err error) {
	b.Client = client
	if len(message.Scene) > 0 {
		b.Client.State.Scene = message.Scene
	}
	if len(message.SceneParams) > 0 {
		b.Client.State.SceneParams = message.SceneParams
	}
	if len(message.SID) > 0 {
		b.Client.State.SID = message.SID
	}
	b.Action = message.Action
	b.ActionParams = message.ActionParams
	return
}

func (b *BaseTextMessageController) Init(_ TextMessageController) {}
