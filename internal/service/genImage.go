package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"hzycoder.com/go-gin-template/internal/config"
	"hzycoder.com/go-gin-template/internal/model"
	"hzycoder.com/go-gin-template/internal/model/dto/request"
	resDto "hzycoder.com/go-gin-template/internal/model/dto/response"
	"hzycoder.com/go-gin-template/internal/repository"
)

type SumbitImageGenTaskResponse struct {
	Code        int                    `json:"code"`
	Description string                 `json:"description"`
	Result      string                 `json:"result,omitempty"`
	Properties  map[string]interface{} `json:"properties,omitempty"`
}

func GenerateImage(ctx context.Context, req request.GenerateImageRequest, userID int64) (string, error) {
	jsonData, err := json.Marshal(req)
	slog.Info("generate image request body json", "json", string(jsonData))

	if err != nil {
		slog.Error("JSON serialization failed", "error", err)
		return "", err
	}

	payload := bytes.NewReader(jsonData)
	client := &http.Client{}

	reqGenImage, err := http.NewRequestWithContext(ctx, "POST", "https://aiyiwei.vip/mj/submit/imagine", payload)
	if err != nil {
		slog.Error("Image Generate Request Construct Failed", "error", err)
		return "", err
	}

	token := config.Global.Ai.Token
	reqGenImage.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	reqGenImage.Header.Add("Content-Type", "application/json")

	resGenImage, err := client.Do(reqGenImage)
	if err != nil {
		slog.Error("Image Generate Request Failed", "error", err)
		return "", err
	}

	defer resGenImage.Body.Close()

	body, err := io.ReadAll(resGenImage.Body)

	if err != nil {
		slog.Error("Image Generate Response Read Failed", "error", err)
		return "", err
	}
	slog.Info("Raw response", "body", string(body))

	var resp SumbitImageGenTaskResponse

	if err := json.Unmarshal(body, &resp); err != nil {
		slog.Error("Failed to parse image generation response", "error", err, "body", string(body))
		return "", fmt.Errorf("invalid response format: %w", err)
	}

	var taskStatus string
	var errMsg string
	taskID := resp.Result

	if resp.Code == 1 {
		taskStatus = "pending"
		slog.Info("Image generation submitted successfully", "taskID", taskID)
	} else {
		taskStatus = "error"
		errMsg = fmt.Sprintf("image generation failed: code=%d, description=%s", resp.Code, resp.Description)
		slog.Error(errMsg, "response", string(body))
	}

	err = repository.SaveGenImageTask(ctx, model.AIGenerationTasks{
		TaskId:       taskID,
		UserId:       userID,
		Type:         "imagine",
		Platform:     req.BotType,
		Prompt:       req.Prompt,
		Status:       taskStatus,
		CustomId:     "",
		ErrorMessage: errMsg,
		RawResponse:  string(body),
	})
	if err != nil {
		slog.Warn("Failed to save image generation task record", "error", err)
	}

	if resp.Code == 1 {
		return taskID, nil
	} else {
		return "", errors.New(errMsg)
	}
}

func QueryGeneratedImageById(ctx context.Context, taskId string) (*resDto.QueryGeneratedImageResponse, error) {
	url := fmt.Sprintf("https://aiyiwei.vip/mj/task/%s/fetch", taskId)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		slog.Error("Failed to construct image query request", "error", err, "taskId", taskId)
		return nil, err
	}

	token := config.Global.Ai.Token
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		slog.Error("Image query request failed", "error", err, "taskId", taskId)
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		slog.Error("Failed to read image query response body", "error", err, "taskId", taskId)
		return nil, err
	}

	slog.Info("Raw image query response", "taskId", taskId, "body", string(body))

	var resp resDto.QueryGeneratedImageResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		slog.Error("Failed to parse image query response", "error", err, "body", string(body), "taskId", taskId)
		return nil, fmt.Errorf("invalid response format: %w", err)
	}

	if resp.ID != "" {
		slog.Info("Image query succeeded", "taskId", taskId, "imageUrl", resp.ImageURL)
		return &resp, nil
	} else {
		errMsg := fmt.Sprintf("Image query failed: id=%s, status=%s, description=%s", resp.ID, resp.Status, resp.Description)
		slog.Error(errMsg, "taskId", taskId, "response", string(body))
		return nil, errors.New(errMsg)
	}
}

func QueryGeneratedTaskIds(ctx context.Context, userId int64, page request.PageQuery) (*resDto.TaskIdsResultWithPage[string], error) {
	taskIdsWithPage, err := repository.QueryGeneratedTaskIds(ctx, userId, page)
	if err != nil {
		slog.Error("Failed to query generated taskId list", "error", err)
		return nil, err
	}
	return taskIdsWithPage, nil
}

type QueryGeneratedImageListParams struct {
	Ids []string
}

func QueryGeneratedImageListByIds(ctx context.Context, taskIds []string) ([]resDto.QueryGeneratedImageResponse, error) {
	url := "https://aiyiwei.vip/mj/task/list-by-condition"

	jsonData, err := json.Marshal(QueryGeneratedImageListParams{
		Ids: taskIds,
	})
	slog.Info("query generated image list request body json", "json", string(jsonData))

	if err != nil {
		slog.Error("JSON serialization failed", "error", err)
		return nil, err
	}

	payload := bytes.NewReader(jsonData)

	req, err := http.NewRequestWithContext(ctx, "POST", url, payload)
	if err != nil {
		slog.Error("Failed to construct image list query request", "error", err, "taskIds", taskIds)
		return nil, err
	}

	token := config.Global.Ai.Token
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		slog.Error("Image list query request failed", "error", err, "taskIds", taskIds)
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		slog.Error("Failed to read image list query response body", "error", err, "taskIds", taskIds)
		return nil, err
	}

	slog.Info("Raw image list query response", "taskIds", taskIds, "body", string(body))

	var resp []resDto.QueryGeneratedImageResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		slog.Error("Failed to parse image list query response", "error", err, "body", string(body), "taskIds", taskIds)
		return nil, fmt.Errorf("invalid response format: %w", err)
	}

	return resp, nil
}

func UpdateImageAction(ctx context.Context, req request.UpdateGenerateImageActionRequest, userID int64) (string, error) {
	jsonData, err := json.Marshal(req)
	slog.Info("image update action request body json", "json", string(jsonData))

	if err != nil {
		slog.Error("JSON serialization failed", "error", err)
		return "", err
	}

	payload := bytes.NewReader(jsonData)
	client := &http.Client{}

	reqGenImage, err := http.NewRequestWithContext(ctx, "POST", "https://aiyiwei.vip/mj/submit/action", payload)
	if err != nil {
		slog.Error("Image update action request Construct Failed", "error", err)
		return "", err
	}

	token := config.Global.Ai.Token
	reqGenImage.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	reqGenImage.Header.Add("Content-Type", "application/json")

	resGenImage, err := client.Do(reqGenImage)
	if err != nil {
		slog.Error("Image update action request failed", "error", err)
		return "", err
	}

	defer resGenImage.Body.Close()

	body, err := io.ReadAll(resGenImage.Body)

	if err != nil {
		slog.Error("Image update action Response Read failed", "error", err)
		return "", err
	}
	slog.Info("Raw update Image action response", "body", string(body))

	var resp SumbitImageGenTaskResponse

	if err := json.Unmarshal(body, &resp); err != nil {
		slog.Error("Failed to parse image update action response", "error", err, "body", string(body))
		return "", fmt.Errorf("invalid response format: %w", err)
	}

	var taskStatus string
	var errMsg string
	taskID := resp.Result

	if resp.Code == 1 || resp.Code == 21 {
		slog.Info("Image update action submitted successfully", "taskID", taskID)
	} else {
		errMsg = fmt.Sprintf("Image update action failed: code=%d, description=%s", resp.Code, resp.Description)
		slog.Error(errMsg, "response", string(body))
	}

	err = repository.SaveGenImageTask(ctx, model.AIGenerationTasks{
		TaskId:       taskID,
		UserId:       userID,
		Type:         "action",
		Platform:     "",
		Prompt:       "",
		Status:       taskStatus,
		CustomId:     req.CustomId,
		ErrorMessage: errMsg,
		RawResponse:  string(body),
	})
	if err != nil {
		slog.Warn("Failed to save image update action task record", "error", err)
	}

	if resp.Code == 1 || resp.Code == 21 {
		return taskID, nil
	} else {
		return "", errors.New(errMsg)
	}
}

func UpdateImageModal(ctx context.Context, req request.UpdateGenerateImageModalRequest, userID int64) (string, error) {
	jsonData, err := json.Marshal(req)
	slog.Info("image update modal request body json", "json", string(jsonData))

	if err != nil {
		slog.Error("JSON serialization failed", "error", err)
		return "", err
	}

	payload := bytes.NewReader(jsonData)
	client := &http.Client{}

	reqGenImage, err := http.NewRequestWithContext(ctx, "POST", "https://aiyiwei.vip/mj/submit/modal", payload)
	if err != nil {
		slog.Error("Image update modal request Construct Failed", "error", err)
		return "", err
	}

	token := config.Global.Ai.Token
	reqGenImage.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	reqGenImage.Header.Add("Content-Type", "application/json")

	resGenImage, err := client.Do(reqGenImage)
	if err != nil {
		slog.Error("Image update modal request failed", "error", err)
		return "", err
	}

	defer resGenImage.Body.Close()

	body, err := io.ReadAll(resGenImage.Body)

	if err != nil {
		slog.Error("Image update modal Response Read failed", "error", err)
		return "", err
	}
	slog.Info("Raw update Image modal response", "body", string(body))

	var resp SumbitImageGenTaskResponse

	if err := json.Unmarshal(body, &resp); err != nil {
		slog.Error("Failed to parse image update modal response", "error", err, "body", string(body))
		return "", fmt.Errorf("invalid response format: %w", err)
	}

	var errMsg string
	taskID := resp.Result

	if resp.Code == 1 {
		slog.Info("Image update modal submitted successfully", "taskID", taskID)
	} else {
		errMsg = fmt.Sprintf("Image update modal failed: code=%d, description=%s", resp.Code, resp.Description)
		slog.Error(errMsg, "response", string(body))
	}

	if resp.Code == 1 {
		return taskID, nil
	} else {
		return "", errors.New(errMsg)
	}
}
