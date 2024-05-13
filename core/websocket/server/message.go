package server

//	{
//		"path": "test",
//		"payload": {
//			"key1": 1234,
//			"key2": "value"
//		}
//	}
type ReceiveMessage struct {
	Path    string                 `json:"path"`
	Payload map[string]interface{} `json:"payload"`
}
