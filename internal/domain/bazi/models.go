package bazi

import "encoding/json"

// BaziRequest 定义了八字排盘工具的输入参数结构。
type BaziRequest struct {
	Name     string `json:"name,omitempty" description:"姓名" default:"求测者"`
	Sex      int    `json:"sex" description:"性别 0男 1女" required:"true" enum:"0,1"`
	Type     int    `json:"type" description:"历类型 0农历 1公历" required:"true" enum:"0,1"`
	Year     int    `json:"year" description:"出生年 例: 1988" required:"true"`
	Month    int    `json:"month" description:"出生月 例: 8" required:"true"`
	Day      int    `json:"day" description:"出生日 例: 7" required:"true"`
	Hours    int    `json:"hours" description:"出生时 例: 12" required:"true"`
	Minute   int    `json:"minute,omitempty" description:"出生分 例: 30" default:"0"`
	Sect     int    `json:"sect,omitempty" description:"流派 1:晚子时日柱算明天 2:晚子时日柱算当天" default:"1"`
	Zhen     int    `json:"zhen,omitempty" description:"是否真太阳时 1:考虑真太阳时 2:不考虑真太阳时" default:"1"`
	Province string `json:"province,omitempty" description:"表示具体的省级行政区 最后面需要带上“省市区”等 例：北京市" x-enum:"data://provinces"`
	City     string `json:"city,omitempty" description:"表示具体的县市级行政区 最后面一般不带上“县市区”（除非带上后只有两个字） 例：北京" x-enum:"data://cities/{province}"`
	Lang     string `json:"lang,omitempty" description:"多语言:zh-cn、zh-tw" default:"zh-cn"`
}

// BaziResponse 定义了从外部 API 获取的八字排盘响应结构。
type BaziResponse struct {
	ErrCode int             `json:"errcode"`
	ErrMsg  string          `json:"errmsg"`
	Notice  string          `json:"notice"`
	Data    json.RawMessage `json:"data"`
}
