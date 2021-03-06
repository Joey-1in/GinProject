package controller

import (
	"github.com/gin-gonic/gin"
	"ginProject/service"
	"ginProject/tool"
	"fmt"
)

type ShopController struct {
}

// shop模块的路由解析
func (sc *ShopController) Router(app *gin.Engine) {
	app.GET("/api/shops", sc.GetShopList)
	app.GET("api/searchShops", sc.SearchShop)
}

// 获取商铺列表
func (sc *ShopController) GetShopList(context *gin.Context) {

	longitude := context.Query("longitude")
	latitude := context.Query("latitude")

	fmt.Println(longitude, latitude)
	if longitude == "" || longitude == "undefined" || latitude == "" || latitude == "undefined" {
		longitude = "116.34" //北京
		latitude = "40.34"
	}

	fmt.Println(longitude, latitude)

	shopService := service.ShopService{}
	shops := shopService.ShopList(longitude, latitude)
	if len(shops) == 0 {
		tool.Failed(context, "暂未获取到商户信息")
		return
	}
	
	for _, shop := range shops {
		shopServices := shopService.GetService(shop.Id)
		if len(shopServices) == 0 {
			shop.Supports = nil
		} else {
			shop.Supports = shopServices
		}
	}
	
	tool.Success(context, shops)
}

// 关键词搜索商铺信息
 func (sc *ShopController) SearchShop(context *gin.Context) {
	longitude := context.Query("longitude")
	latitude := context.Query("latitude")
	keyword := context.Query("keyword")

	if keyword == "" {
		tool.Failed(context, "重新输入商铺名称")
		return
	}

	if longitude == "" || longitude == "undefined" || latitude == "" || latitude == "undefined" {
		longitude = "116.34" //北京
		latitude = "40.34"
	}

	//执行真实的搜索逻辑
	shopService := service.ShopService{}
	shopService.SearchShops(longitude, latitude, keyword)

}
