package service

import (
	"comfyui_endpoint/database"
	"comfyui_endpoint/dto/request"
	"comfyui_endpoint/dto/response"
	"comfyui_endpoint/model"
	"fmt"
	"strings"
	"time"
)

type EndpointService interface {
	Create(p request.EndpointCreateRequest) error
	Update(p request.EndpointUpdateRequest) error
	Index(p request.EndpointIndexRequest) (data response.EndpointIndexResponse, err error)
}

type endpointService struct{}

func NewEndpointService() EndpointService {
	return &endpointService{}
}

func (s *endpointService) Create(p request.EndpointCreateRequest) error {
	if p.Path == "" {
		return fmt.Errorf("请填写路径")
	}

	p.Path = strings.TrimPrefix(p.Path, "/")

	syncPath := strings.Join([]string{p.Path, "sync"}, "/")
	var old model.Endpoint
	SqliteDb().Model(&model.Endpoint{}).Where("path = ?", p.Path).First(&old)

	if old.Id > 0 {
		return fmt.Errorf("路径已存在")
	}

	err := SqliteDb().Debug().Model(&model.Endpoint{}).Create(&model.Endpoint{
		Path:        p.Path,
		SyncPath:    syncPath,
		Description: p.Description,
		ApiJson:     p.ApiJson,
		Workflow:    p.Workflow,
		CallbackUrl: p.CallbackUrl,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}).Error

	if err != nil {
		return err

	}
	// 注册同步处理方法
	RegisterSyncHandle(GetRouter(), syncPath)
	RegisterAsyncHandle(GetRouter(), p.Path)
	return nil
}

func (s *endpointService) Update(p request.EndpointUpdateRequest) error {
	var old model.Endpoint
	SqliteDb().Model(&model.Endpoint{}).Where("id = ?", p.Id).First(&old)

	if old.Id == 0 {
		return fmt.Errorf("endpoint不存在")
	}

	SqliteDb().Model(&model.Endpoint{}).Where("id = ?", p.Id).Updates(map[string]interface{}{
		"api_json":     p.ApiJson,
		"description":  p.Description,
		"workflow":     p.Workflow,
		"callback_url": p.CallbackUrl,
		"updated_at":   time.Now(),
	})

	return nil
}

func (s *endpointService) Index(p request.EndpointIndexRequest) (data response.EndpointIndexResponse, err error) {
	db := SqliteDb().Model(&model.Endpoint{})

	if p.Path != "" {
		db = db.Where("path LIKE ?", "%"+p.Path+"%")
	}

	if p.Description != "" {
		db = db.Where("description LIKE ?", "%"+p.Description+"%")
	}

	if p.CallbackUrl != "" {
		db = db.Where("callback_url LIKE ?", "%"+p.CallbackUrl+"%")
	}

	db.Count(&data.Total)
	err = db.Order("created_at desc").Scopes(database.Paginate(p.Page, p.PageSize)).Find(&data.Data).Error

	return
}
