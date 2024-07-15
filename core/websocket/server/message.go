package server

//	{
//		"scene": "test",
//		"sceneParams": {
//			"key1": 1234,
//			"key2": "value"
//		},
//		"action": "test",
//		"actionParams": {
//			"key1": 1234,
//			"key2": "value"
//		}
//	}
type ReceiveMessage struct {
	Scene        string                 `json:"scene"`
	SID          string                 `json:"sid"`
	SceneParams  map[string]interface{} `json:"sceneParams"`
	Action       string                 `json:"action"`
	ActionParams map[string]interface{} `json:"actionParams"`
}
