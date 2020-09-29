package service

import (
	"ginProject/model"
	"strconv"
	"ginProject/dao"
)

type ShopService struct {
}

// 查询商铺列表数据
func (shopService *ShopService) ShopList(long, lat string) []model.Shop {
	longitude, err := strconv.ParseFloat(long, 10)
	if err != nil {
		return nil
	}
	latitude, err := strconv.ParseFloat(lat, 10)
	if err != nil {
		return nil
	}

	shopDao := dao.NewShopDao()
	return shopDao.QueryShops(longitude, latitude, "")
}

// 根据关键词查询对应的商家信息
 func (shopService *ShopService) SearchShops(long, lat, keyword string) []model.Shop {
	shopDao := dao.NewShopDao()
	longitude, err := strconv.ParseFloat(long, 10)
	if err != nil {
		return nil
	}
	latitude, err := strconv.ParseFloat(lat, 10)
	if err != nil {
		return nil
	}

	return shopDao.QueryShops(longitude, latitude, keyword)
}

// 根据商户id获取对应的服务
 func (shopService *ShopService) GetService(shopId int64) []model.Service {
	shopDao := dao.NewShopDao()
	return shopDao.QueryServiceByShopId(shopId)
}
