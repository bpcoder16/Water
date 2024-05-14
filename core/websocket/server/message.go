package server

//	{
//		"scene": "test",
//		"sceneParameters": {
//			"key1": 1234,
//			"key2": "value"
//		}
//	}
type ReceiveMessage struct {
	Scene        string                 `json:"scene"`
	SceneParams  map[string]interface{} `json:"sceneParams"`
	Action       string                 `json:"action"`
	ActionParams map[string]interface{} `json:"actionParams"`
}
