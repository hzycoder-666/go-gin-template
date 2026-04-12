package model

import "time"

type AIGenerationTasks struct {
	ID           int64
	TaskId       string // 任务ID
	UserId       int64  // 关联用户ID
	Type         string // 任务类型
	Platform     string // 平台（midjourney）
	Prompt       string // 提示词
	Status       string // 任务状态
	RawResponse  string // 原始响应字符串
	ImageUrl     string
	ErrorMessage string
	CustomId     string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	FinishedAt   time.Time
}
