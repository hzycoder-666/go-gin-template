package response

type GenerateImageResponse struct {
	Base64Array []string `json:"base64Array"`
	NotifyHook  string   `json:"notifyHook"`
	Prompt      string   `json:"prompt"`
	State       string   `json:"state"`
	BotType     string   `json:"botType"`
}

type QueryGeneratedImageResponse struct {
	Action      string                        `json:"action"`
	BotType     string                        `json:"botType"`
	Buttons     []QueryGeneratedImageButton   `json:"buttons"`
	CustomID    string                        `json:"customId"`
	Description string                        `json:"description"`
	FailReason  string                        `json:"failReason"`
	FinishTime  int64                         `json:"finishTime"` // 使用 int64 表示时间戳
	ID          string                        `json:"id"`
	ImageURL    string                        `json:"imageUrl"`
	MaskBase64  string                        `json:"maskBase64"`
	Progress    string                        `json:"progress"`
	Prompt      string                        `json:"prompt"`
	PromptEn    string                        `json:"promptEn"`
	Properties  QueryGeneratedImageProperties `json:"properties"`
	StartTime   int64                         `json:"startTime"`
	State       string                        `json:"state"`
	Status      string                        `json:"status"`
	SubmitTime  int64                         `json:"submitTime"`

	// 支持额外未知字段（对应 [property: string]: any）
	// 注意：如果确实需要动态字段，可保留 map[string]interface{}，
	// 但通常在严格 DTO 中不建议。此处按需保留。
	// 若不需要，可删除此字段。
	Extra map[string]interface{} `json:"-"`
}

type QueryGeneratedImageButton struct {
	CustomID string `json:"customId"`
	Emoji    string `json:"emoji"`
	Label    string `json:"label"`
	Style    int    `json:"style"`
	Type     int    `json:"type"`

	// 可选：支持额外字段
	Extra map[string]interface{} `json:"-"`
}

type QueryGeneratedImageProperties struct {
	FinalPrompt   string `json:"finalPrompt"`
	FinalZhPrompt string `json:"finalZhPrompt"`

	// 可选：支持额外字段
	Extra map[string]interface{} `json:"-"`
}

type TaskIdsResultWithPage[T any] struct {
	Items []T `json:"items"`
	PageBase
}
