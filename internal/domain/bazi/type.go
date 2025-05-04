package bazi

type Data struct {
	BaseInfo   BaseInfo   `json:"base_info"`   // 基本信息
	BaziInfo   BaziInfo   `json:"bazi_info"`   // 八字信息
	DayunInfo  DayunInfo  `json:"dayun_info"`  // 大运信息
	StartInfo  StartInfo  `json:"start_info"`  // 起运信息
	DetailInfo DetailInfo `json:"detail_info"` // 详细信息
}

// BaseInfo 表示八字排盘的基本信息领域模型
type BaseInfo struct {
	Zhen    *ZhenInfo `json:"zhen,omitempty"` // 真太阳时信息
	Sex     string    `json:"sex"`            // 性别
	Name    string    `json:"name"`           // 姓名
	Gongli  string    `json:"gongli"`         // 公历年
	Nongli  string    `json:"nongli"`         // 农历年
	Qiyun   string    `json:"qiyun"`          // 起运
	Jiaoyun string    `json:"jiaoyun"`        // 交运
	Zhengge string    `json:"zhengge"`        // 八字正格
}

// ZhenInfo 表示真太阳时信息值对象
type ZhenInfo struct {
	Province string `json:"province"` // 省份
	City     string `json:"city"`     // 城市
	Jingdu   string `json:"jingdu"`   // 经度
	Weidu    string `json:"weidu"`    // 纬度
	Shicha   string `json:"shicha"`   // 时差
}

// BaziInfo 表示八字信息领域模型
type BaziInfo struct {
	Kw      string   `json:"kw"`        // 空亡
	TgCgGod []string `json:"tg_cg_god"` // 天干十神（按年月日时排序）
	Bazi    []string `json:"bazi"`      // 八字四柱（按年月日时排序）
	DzCg    []string `json:"dz_cg"`     // 地支藏干（按年月日时排序）
	DzCgGod []string `json:"dz_cg_god"` // 地支藏干十神（按年月日时排序，藏干中依次为本气中气余气）
	DayCs   []string `json:"day_cs"`    // 八字长生衰旺（按年月日时排序）
	NaYin   []string `json:"na_yin"`    // 五行纳音（按年月日时排序）
}

// DayunInfo 表示大运信息领域模型
type DayunInfo struct {
	BigGod              []string   `json:"big_god"`                 // 大运神（按时间排序）
	Big                 []string   `json:"big"`                     // 大运天干地支（按时间排序）
	BigCs               []string   `json:"big_cs"`                  // 大运长生衰旺（按时间排序）
	XuSui               []int      `json:"xu_sui"`                  // 虚岁（按时间排序）
	BigStartYear        []int      `json:"big_start_year"`          // 大运始于年份（按时间排序）
	BigStartYearLiuNian string     `json:"big_start_year_liu_nian"` // 大运流年
	BigEndYear          []int      `json:"big_end_year"`            // 大运止于年份（按时间排序）
	YearsInfo0          []YearChar `json:"years_info0"`             // 大运第一个流年信息
	YearsInfo1          []YearChar `json:"years_info1"`             // 大运第二个流年信息
	YearsInfo2          []YearChar `json:"years_info2"`             // 大运第三个流年信息
	YearsInfo3          []YearChar `json:"years_info3"`             // 大运第四个流年信息
	YearsInfo4          []YearChar `json:"years_info4"`             // 大运第五个流年信息
	YearsInfo5          []YearChar `json:"years_info5"`             // 大运第六个流年信息
	YearsInfo6          []YearChar `json:"years_info6"`             // 大运第七个流年信息
	YearsInfo7          []YearChar `json:"years_info7"`             // 大运第八个流年信息
	YearsInfo8          []YearChar `json:"years_info8"`             // 大运第九个流年信息
	YearsInfo9          []YearChar `json:"years_info9"`             // 大运第十个流年信息
}

// StartInfo 表示起运信息领域模型
type StartInfo struct {
	Jishen []string `json:"jishen"` // 吉神凶煞（按年月日时排序）
	Xz     string   `json:"xz"`     // 星座
	Sx     string   `json:"sx"`     // 生肖
}

// DetailInfo 表示详细信息领域模型
type DetailInfo struct {
	Zhuxing      ZhuxingInfo    `json:"zhuxing"`      // 天干透出十神
	Sizhu        SizhuInfo      `json:"sizhu"`        // 四柱天干地支
	Canggan      CangganInfo    `json:"canggan"`      // 藏干信息
	Fuxing       FuxingInfo     `json:"fuxing"`       // 藏干十神
	Xingyun      XingyunInfo    `json:"xingyun"`      // 星运信息
	Zizuo        ZizuoInfo      `json:"zizuo"`        // 自坐信息
	Kongwang     KongwangInfo   `json:"kongwang"`     // 空亡信息
	Nayin        NayinInfo      `json:"nayin"`        // 纳音信息
	Shensha      ShenshaInfo    `json:"shensha"`      // 神煞信息
	Dayunshensha []DayunShensha `json:"dayunshensha"` // 大运神煞（按时间排序）
}

// ZhuxingInfo 表示天干透出十神信息
type ZhuxingInfo struct {
	Year  string `json:"year"`  // 年干透出十神
	Month string `json:"month"` // 月干透出十神
	Day   string `json:"day"`   // 日干透出十神
	Hour  string `json:"hour"`  // 时干透出十神
}

// SizhuInfo 表示四柱天干地支信息
type SizhuInfo struct {
	Year  YearInfo  `json:"year"`  // 年柱
	Month MonthInfo `json:"month"` // 月柱
	Day   DayInfo   `json:"day"`   // 日柱
	Hour  HourInfo  `json:"hour"`  // 时柱
}

// YearInfo 表示年柱信息
type YearInfo struct {
	Tg string `json:"tg"` // 年干
	Dz string `json:"dz"` // 年支
}

// MonthInfo 表示月柱信息
type MonthInfo struct {
	Tg string `json:"tg"` // 月干
	Dz string `json:"dz"` // 月支
}

// DayInfo 表示日柱信息
type DayInfo struct {
	Tg string `json:"tg"` // 日干
	Dz string `json:"dz"` // 日支
}

// HourInfo 表示时柱信息
type HourInfo struct {
	Tg string `json:"tg"` // 时干
	Dz string `json:"dz"` // 时支
}

// CangganInfo 表示藏干信息
type CangganInfo struct {
	Year  []string `json:"year"`  // 年支藏干（按本气中气余气排序）
	Month []string `json:"month"` // 月支藏干（按本气中气余气排序）
	Day   []string `json:"day"`   // 日支藏干（按本气中气余气排序）
	Hour  []string `json:"hour"`  // 时支藏干（按本气中气余气排序）
}

// FuxingInfo 表示藏干十神信息
type FuxingInfo struct {
	Year  []string `json:"year"`  // 年支藏干十神（按本气中气余气排序）
	Month []string `json:"month"` // 月支藏干十神（按本气中气余气排序）
	Day   []string `json:"day"`   // 日支藏干十神（按本气中气余气排序）
	Hour  []string `json:"hour"`  // 时支藏干十神（按本气中气余气排序）
}

// XingyunInfo 表示星运信息
type XingyunInfo struct {
	Year  string `json:"year"`  // 年柱星运
	Month string `json:"month"` // 月柱星运
	Day   string `json:"day"`   // 日柱星运
	Hour  string `json:"hour"`  // 时柱星运
}

// ZizuoInfo 表示自坐信息
type ZizuoInfo struct {
	Year  string `json:"year"`  // 年柱自坐
	Month string `json:"month"` // 月柱自坐
	Day   string `json:"day"`   // 日柱自坐
	Hour  string `json:"hour"`  // 时柱自坐
}

// KongwangInfo 表示空亡信息
type KongwangInfo struct {
	Year  string `json:"year"`  // 年柱空亡
	Month string `json:"month"` // 月柱空亡
	Day   string `json:"day"`   // 日柱空亡
	Hour  string `json:"hour"`  // 时柱空亡
}

// NayinInfo 表示纳音信息
type NayinInfo struct {
	Year  string `json:"year"`  // 年柱纳音
	Month string `json:"month"` // 月柱纳音
	Day   string `json:"day"`   // 日柱纳音
	Hour  string `json:"hour"`  // 时柱纳音
}

// ShenshaInfo 表示神煞信息
type ShenshaInfo struct {
	Year  string `json:"year"`  // 年柱神煞
	Month string `json:"month"` // 月柱神煞
	Day   string `json:"day"`   // 日柱神煞
	Hour  string `json:"hour"`  // 时柱神煞
}

// DayunshenshaInfo 表示大运神煞信息
type DayunshenshaInfo struct {
	Tgdz    string `json:"tgdz"`    // 大运天干地支
	Shensha string `json:"shensha"` // 对应神煞
}

// YearChar 表示流年年柱值对象
type YearChar struct {
	YearChar string `json:"year_char"` // 年柱
}

// DayunShensha 表示大运神煞信息
type DayunShensha struct {
	Tgdz    string `json:"tgdz"`    // 大运天干地支
	Shensha string `json:"shensha"` // 对应神煞（字符串形式）
}

// BaziService 八字排盘领域服务接口
type BaziService interface {
	Calculate(baseInfo BaseInfo) (*Data, error)
	Analyze(baziInfo BaziInfo) (string, error)
	Predict(dayunInfo DayunInfo) (string, error)
}
