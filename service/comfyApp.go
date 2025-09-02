package service

import (
	"comfyui_endpoint/database"
	"comfyui_endpoint/dto/request"
	"comfyui_endpoint/dto/response"
	"comfyui_endpoint/logger"
	"comfyui_endpoint/model"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type ComfyAppService interface {
	Create(p request.ComfyAppCreateRequest) error
	Remove(p request.ComfyAppRemoveRequest) error
	Update(p request.ComfyAppUpdateRequest) error
	Index(p request.ComfyAppIndexRequest) (data response.ComfyAppIndexResponse, err error)
	RestartWs(p request.ComfyAppRestartWsRequest) error
	InitWs() error
}

type comfyAppService struct{}

func NewComfyAppService() ComfyAppService {
	return &comfyAppService{}
}

func (s *comfyAppService) Create(p request.ComfyAppCreateRequest) error {
	var old model.ComfyApp
	SqliteDb().Model(&model.ComfyApp{}).Where("addr = ?", p.Addr).First(&old)

	if old.Id > 0 {
		return fmt.Errorf("地址已存在")
	}

	var new = model.ComfyApp{
		Uid:         uuid.NewString(),
		Addr:        p.Addr,
		Description: p.Description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	SqliteDb().Model(&model.ComfyApp{}).Create(&new)

	SqliteDb().Model(&model.ComfyAppInfo{}).Create(&model.ComfyAppInfo{
		Uid:        new.Uid,
		RunningNum: 0,
		PendingNum: 0,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	})

	err := NewWsClient(new.Uid, new.Addr).Start()

	if err != nil {
		return err
	}
	SqliteDb().Model(&model.ComfyAppInfo{}).Where("uid = ?", new.Uid).Update("ws_status", "connect")
	return nil
}

func (s *comfyAppService) Remove(p request.ComfyAppRemoveRequest) error {
	var old model.ComfyApp

	SqliteDb().Model(&model.ComfyApp{}).Where("id = ?", p.Id).First(&old)

	if old.Id == 0 {
		return fmt.Errorf("app不存在")
	}

	SqliteDb().Model(&model.ComfyApp{}).Delete(&old, p.Id)

	SqliteDb().Model(&model.ComfyAppInfo{}).Where("uid = ?", old.Uid).Delete(&model.ComfyAppInfo{})

	// ws 关闭连接
	if clientMap[old.Uid] != nil {
		clientMap[old.Uid].Close()
	}
	return nil
}

func (s *comfyAppService) Update(p request.ComfyAppUpdateRequest) error {
	var old model.ComfyApp

	SqliteDb().Model(&model.ComfyApp{}).Where("id = ?", p.Id).First(&old)

	if old.Id == 0 {
		return fmt.Errorf("app不存在")
	}

	SqliteDb().Model(&model.ComfyApp{}).Where("id = ?", p.Id).Updates(map[string]interface{}{
		"addr":        p.Addr,
		"description": p.Description,
	})

	if p.Addr != old.Addr {
		//  ws 连接
		if clientMap[old.Uid] != nil {
			err := clientMap[old.Uid].Close()
			if err != nil {
				return fmt.Errorf("关闭ws连接失败 %s", err.Error())
			}

			SqliteDb().Model(&model.ComfyAppInfo{}).Where("uid = ?", old.Uid).Update("ws_status", "close")
		}
		err := NewWsClient(old.Uid, p.Addr).Start()

		if err != nil {
			return err
		}

		SqliteDb().Model(&model.ComfyAppInfo{}).Where("uid = ?", old.Uid).Update("ws_status", "connect")
	}

	return nil
}

func (s *comfyAppService) Index(p request.ComfyAppIndexRequest) (data response.ComfyAppIndexResponse, err error) {
	db := SqliteDb().Table("comfy_apps ca").Select("ca.id, ca.uid, ca.addr, ca.description, ca.created_at, ca.updated_at, cai.running_num as running_num, cai.pending_num as pending_num, cai.ws_status as ws_status").Joins("LEFT JOIN comfy_app_infos cai ON ca.uid = cai.uid")

	if p.Addr != "" {
		db = db.Where("ca.addr LIKE ?", "%"+p.Addr+"%")
	}

	if p.Description != "" {
		db = db.Where("ca.description LIKE ?", "%"+p.Description+"%")
	}

	db.Count(&data.Total)

	err = db.Order("ca.created_at desc").Scopes(database.Paginate(p.Page, p.PageSize)).Find(&data.Data).Error

	if err != nil {
		return
	}

	return
}

func (s *comfyAppService) InitWs() error {
	var list []model.ComfyApp
	SqliteDb().Model(&model.ComfyApp{}).Find(&list)

	for _, app := range list {
		SqliteDb().Model(&model.ComfyAppInfo{}).Where("uid = ?", app.Uid).Update("ws_status", "close")
		if clientMap[app.Uid] != nil {
			clientMap[app.Uid].Close()
		}
		wrong := NewWsClient(app.Uid, app.Addr).Start()
		if wrong != nil {
			logger.Error("error %s", "ws连接失败")
		}
		SqliteDb().Model(&model.ComfyAppInfo{}).Where("uid = ?", app.Uid).Update("ws_status", "connect")
	}
	return nil
}

func (s *comfyAppService) RestartWs(p request.ComfyAppRestartWsRequest) error {
	var app model.ComfyApp

	SqliteDb().Model(&model.ComfyApp{}).Where("id = ?", p.Id).First(&app)

	if clientMap[app.Uid] != nil {
		clientMap[app.Uid].Close()
		SqliteDb().Model(&model.ComfyAppInfo{}).Where("uid = ?", app.Uid).Update("ws_status", "close")
	}

	err := NewWsClient(app.Uid, app.Addr).Start()
	if err != nil {
		SqliteDb().Model(&model.ComfyAppInfo{}).Where("uid = ?", app.Uid).Update("ws_status", "connect")
	}
	return nil
}
