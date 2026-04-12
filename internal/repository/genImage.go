package repository

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"hzycoder.com/go-gin-template/internal/database"
	"hzycoder.com/go-gin-template/internal/model"
	"hzycoder.com/go-gin-template/internal/model/dto/request"
	resDto "hzycoder.com/go-gin-template/internal/model/dto/response"
)

func SaveGenImageTask(ctx context.Context, aigcTask model.AIGenerationTasks) error {
	now := time.Now()
	aigcTask.CreatedAt = now
	aigcTask.UpdatedAt = now

	query := `
		INSERT INTO ai_generation_tasks (task_id, user_id, type, platform, prompt, status, custom_id, raw_response, created_at, updated_at)
		VALUES (?,?,?,?,?,?,?,?,?,?);
	`

	_, err := database.DB.ExecContext(ctx, query,
		aigcTask.TaskId,
		aigcTask.UserId,
		aigcTask.Type,
		aigcTask.Platform,
		aigcTask.Prompt,
		aigcTask.Status,
		aigcTask.CustomId,
		aigcTask.RawResponse,
		aigcTask.CreatedAt,
		aigcTask.UpdatedAt,
	)

	if err != nil {
		slog.Error("insert ai_generation_tasks failed", "error", err)
		return err
	}

	return nil
}

func QueryGeneratedTaskIds(ctx context.Context, userId int64, page request.PageQuery) (*resDto.TaskIdsResultWithPage[string], error) {
	offset := (page.Page - 1) * page.PageSize
	limit := page.PageSize

	// 1. 查询总记录数
	var total int
	countQuery := `
		SELECT COUNT(*) FROM ai_generation_tasks WHERE user_id = ?
	`

	if err := database.DB.QueryRowContext(ctx, countQuery, userId).Scan(&total); err != nil {
		return nil, fmt.Errorf("failed to get total count: %w", err)
	}

	// 2. 查询当前页数据
	query := `
		SELECT task_id
		FROM ai_generation_tasks
		WHERE user_id = ?
		ORDER BY created_at DESC, id DESC
		LIMIT ? OFFSET ?
	`

	rows, err := database.DB.QueryContext(ctx, query, userId, limit, offset)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var taskIds []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("failed to scan task ID: %w", err)
		}
		taskIds = append(taskIds, id)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration: %w", err)
	}

	return &resDto.TaskIdsResultWithPage[string]{
		Items: taskIds,
		PageBase: resDto.PageBase{
			Page:     page.Page,
			PageSize: limit,
			Total:    total,
		},
	}, nil
}
