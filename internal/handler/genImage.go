package handler

import (
	"log/slog"
	"sort"

	"github.com/gin-gonic/gin"
	"hzycoder.com/go-gin-template/internal/model/dto/request"
	dto "hzycoder.com/go-gin-template/internal/model/dto/request"
	resDto "hzycoder.com/go-gin-template/internal/model/dto/response"
	"hzycoder.com/go-gin-template/internal/repository"
	"hzycoder.com/go-gin-template/internal/service"
	"hzycoder.com/go-gin-template/pkg/response"
)

func GenerateImage(c *gin.Context) {
	var req dto.GenerateImageRequest
	userIDAny, exists := c.Get("user_id")
	if !exists {
		response.FailWithCode(c, response.CodeParamInvalid)
		return
	}

	userID, ok := ExtractInt64(userIDAny)
	if !ok {
		response.FailWithCode(c, response.CodeParamInvalid)
		return
	}

	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		c.Set("validation_error", err)
		c.Abort()
		return
	}

	ctx := c.Request.Context()

	generateImageTaskID, err := service.GenerateImage(ctx, req, userID)
	if err != nil {
		slog.Error("generate image failed", "error", err)
		response.HandleError(c, err)
		return
	}

	response.Success(c, generateImageTaskID)
}

func QueryGeneratedImageById(c *gin.Context) {
	taskId := c.Query("taskId")

	if taskId == "" {
		slog.Error("query generated image query empty", "error", "参数为空")
		response.FailWithCode(c, response.CodeParamInvalid)
		return
	}

	ctx := c.Request.Context()

	generatedImageInfo, err := service.QueryGeneratedImageById(ctx, taskId)

	if err != nil {
		slog.Error("query generated image task failed", "error", err)
		response.HandleError(c, err)
		return
	}

	response.Success(c, generatedImageInfo)
}

func deduplicateByID(items []resDto.QueryGeneratedImageResponse) []resDto.QueryGeneratedImageResponse {
	seen := make(map[string]bool)
	var result []resDto.QueryGeneratedImageResponse

	for _, item := range items {
		if !seen[item.ID] {
			seen[item.ID] = true
			result = append(result, item)
		}
	}
	return result
}

func QueryGeneratedImageList(c *gin.Context) {
	userIDAny, exists := c.Get("user_id")
	if !exists {
		response.FailWithCode(c, response.CodeParamInvalid)
		return
	}

	userID, ok := ExtractInt64(userIDAny)
	if !ok {
		response.FailWithCode(c, response.CodeParamInvalid)
		return
	}

	var pq dto.PageQuery

	if err := c.ShouldBindQuery(&pq); err != nil {
		response.FailWithCode(c, response.CodeParamInvalid)
		return
	}

	if pq.Page == 0 {
		pq.Page = 1
	}
	if pq.PageSize == 0 {
		pq.PageSize = 10
	}

	ctx := c.Request.Context()

	taskIdsWithPage, err := repository.QueryGeneratedTaskIds(ctx, userID, pq)

	if err != nil {
		slog.Error("query generated taskId list failed", "error", err)
		response.HandleError(c, err)
		return
	}

	generatedImageList, err := service.QueryGeneratedImageListByIds(ctx, taskIdsWithPage.Items)

	if err != nil {
		slog.Error("query generated image List task failed", "error", err)
		response.HandleError(c, err)
		return
	}

	sort.Slice(generatedImageList, func(i, j int) bool {
		return generatedImageList[i].SubmitTime > generatedImageList[j].SubmitTime
	})

	generatedImageList = deduplicateByID(generatedImageList)

	response.Success(c, gin.H{
		"items": generatedImageList,
		"page":  taskIdsWithPage.Page,
		"total": taskIdsWithPage.Total,
	})
}

func UpdateImageAction(c *gin.Context) {
	var req request.UpdateGenerateImageActionRequest
	userIDAny, exists := c.Get("user_id")
	if !exists {
		response.FailWithCode(c, response.CodeParamInvalid)
		return
	}

	userID, ok := ExtractInt64(userIDAny)
	if !ok {
		response.FailWithCode(c, response.CodeParamInvalid)
		return
	}

	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		c.Set("validation_error", err)
		c.Abort()
		return
	}

	ctx := c.Request.Context()

	updateImageActionTaskID, err := service.UpdateImageAction(ctx, req, userID)

	if err != nil {
		slog.Error("update image action failed", "error", err)
		response.HandleError(c, err)
		return
	}

	response.Success(c, updateImageActionTaskID)
}

func UpdateImageModal(c *gin.Context) {
	var req request.UpdateGenerateImageModalRequest
	userIDAny, exists := c.Get("user_id")
	if !exists {
		response.FailWithCode(c, response.CodeParamInvalid)
		return
	}

	userID, ok := ExtractInt64(userIDAny)
	if !ok {
		response.FailWithCode(c, response.CodeParamInvalid)
		return
	}

	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		c.Set("validation_error", err)
		c.Abort()
		return
	}

	ctx := c.Request.Context()

	updateImageModalTaskID, err := service.UpdateImageModal(ctx, req, userID)

	if err != nil {
		slog.Error("update image modal failed", "error", err)
		response.HandleError(c, err)
		return
	}

	response.Success(c, updateImageModalTaskID)
}
