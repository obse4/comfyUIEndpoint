package service

import (
	"comfyui_endpoint/dto/request"
	"comfyui_endpoint/model"
	"fmt"
	"time"
)

type EndpointParamService interface {
	Set(p request.EndpointParamSetRequest) error
	Find(p request.EndpointParamFindRequest) (data []model.EndpointParam, err error)
	FindOne(path, paramKey string) (data model.EndpointParam, err error)
}

type endpointParamService struct{}

func NewEndpointParamService() EndpointParamService {
	return &endpointParamService{}
}

func (s *endpointParamService) Set(p request.EndpointParamSetRequest) error {
	// 删除旧的参数
	// 创建新参数
	SqliteDb().Model(&model.EndpointParam{}).Debug().Where("endpoint_id = ?", p.EndpointId).Delete(&model.EndpointParam{})

	for _, v := range p.Items {
		SqliteDb().Model(&model.EndpointParam{}).Create(&model.EndpointParam{
			EndpointId:  p.EndpointId,
			ParamName:   v.ParamName,
			ParamKey:    v.ParamKey,
			Description: v.Description,
			ParamType:   v.ParamType,
			JsonKey:     v.JsonKey,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		})
	}
	return nil
}

func (s *endpointParamService) Find(p request.EndpointParamFindRequest) (data []model.EndpointParam, err error) {
	err = SqliteDb().Model(&model.EndpointParam{}).Where("endpoint_id = ?", p.EndpointId).Find(&data).Error
	return
}

func (s *endpointParamService) FindOne(path, paramKey string) (data model.EndpointParam, err error) {
	var endpoint model.Endpoint

	SqliteDb().Model(&model.Endpoint{}).Where("path = ?", path).First(&endpoint)

	if endpoint.Id == 0 {
		err = fmt.Errorf("找不到路径对应endpoint [%s]", path)
		return
	}
	SqliteDb().Model(&model.EndpointParam{}).Where("param_key = ?", paramKey).Where("endpoint_id = ?", endpoint.Id).First(&data)

	if data.Id == 0 {
		err = fmt.Errorf("找不到 路径 [%s] 参数 [%s]", path, paramKey)
		return
	}
	return
}
