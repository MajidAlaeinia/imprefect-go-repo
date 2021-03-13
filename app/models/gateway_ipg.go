package models

import "github.com/jinzhu/gorm"

type GatewayIpg struct {
	gorm.Model
	UserId     uint64 `json:"user_id"`
	BusinessId uint64 `json:"business_id"`
	//Url string `json:"url"`
	Urls   string `json:"urls" mysql:"type:text"`
	Ips    string `json:"ips" mysql:"type:text"`
	ApiKey string `json:"api_key"`
	Status int    `json:"status"`
	Review bool   `json:"review"`
	//Ip []byte `json:"ip"`
	WageType bool `json:"wage_type"`
	//IsVip bool `json:"is_vip"`
	//HasPayedForAdditionalBusiness bool `json:"has_payed_for_additional_business"`
	MerchantId       string `json:"merchant_id"`
	MerchantPassword string `json:"-"`
	//WagePercent string `json:"wage_percent"`
	//WageMin string `json:"wage_min"`
	//WageMax string `json:"wage_max"`
}
