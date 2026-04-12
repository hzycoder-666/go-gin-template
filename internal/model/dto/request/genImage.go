package request

type GenerateImageRequest struct {
	Base64Array []string `json:"base64Array"`
	NotifyHook  string   `json:"notifyHook"`
	Prompt      string   `json:"prompt"`
	State       string   `json:"state"`
	BotType     string   `json:"botType"`
}

type UpdateGenerateImageActionRequest struct {
	ChooseSameChannel bool   `json:"chooseSameChannel" binding:"required"`
	CustomId          string `json:"customId" binding:"required"`
	TaskId            string `json:"taskId"`
	NotifyHook        string `json:"notifyHook"`
	State             string `json:"state"`
}

type UpdateGenerateImageModalRequest struct {
	MaskBase64 string `json:"maskBase64,omitempty"`
	Prompt     string `json:"prompt"`
	TaskId     string `json:"taskId"`
}
