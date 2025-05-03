package bazi

// BaseInfo 表示八字排盘的基本信息领域模型
type BaseInfo struct {
    Zhen     *ZhenInfo `json:"zhen,omitempty"` // 真太阳时信息
    Sex      string    `json:"sex"`           // 性别
    Name     string    `json:"name"`          // 姓名
    Gongli   string    `json:"gongli"`        // 公历年
    Nongli   string    `json:"nongli"`        // 农历年
    Qiyun    string    `json:"qiyun"`         // 起运
    Jiaoyun  string    `json:"jiaoyun"`       // 交运
    Zhengge  string    `json:"zhengge"`       // 八字正格
}

// ZhenInfo 表示真太阳时信息值对象
type ZhenInfo struct {
    Province string `json:"province"` // 省份
    City     string `json:"city"`      // 城市
    Jingdu   string `json:"jingdu"`    // 经度
    Weidu    string `json:"weidu"`     // 纬度
    Shicha   string `json:"shicha"`    // 时差
}

// BaziInfo 表示八字信息领域模型
type BaziInfo struct {
    Kw       string   `json:"kw"`        // 空亡
    TgCgGod  []string `json:"tg_cg_god"` // 天干十神（按年月日时排序）
    Bazi     []string `json:"bazi"`      // 八字四柱（按年月日时排序）
    DzCg     []string `json:"dz_cg"`     // 地支藏干（按年月日时排序）
    DzCgGod  []string `json:"dz_cg_god"` // 地支藏干十神（按年月日时排序，藏干中依次为本气中气余气）
    DayCs    []string `json:"day_cs"`    // 八字长生衰旺（按年月日时排序）
    NaYin    []string `json:"na_yin"`    // 五行纳音（按年月日时排序）
}

// DayunInfo 表示大运信息领域模型
type DayunInfo struct {
    BigGod              []string       `json:"big_god"`                 // 大运神（按时间排序）
    Big                 []string       `json:"big"`                     // 大运天干地支（按时间排序）
    BigCs               []string       `json:"big_cs"`                  // 大运长生衰旺（按时间排序）
    XuSui               []int          `json:"xu_sui"`                 // 虚岁（按时间排序）
    BigStartYear        []int          `json:"big_start_year"`          // 大运始于年份（按时间排序）
    BigStartYearLiuNian string         `json:"big_start_year_liu_nian"` // 大运流年
    BigEndYear          []int          `json:"big_end_year"`           // 大运止于年份（按时间排序）
    YearsInfo           [10][]YearChar `json:"-"`                      // 流年信息(使用数组替代原YearsInfo0-9)
}

// YearChar 表示流年年柱值对象
type YearChar struct {
    YearChar string `json:"year_char"` // 年柱
}

// StartInfo 表示起运信息领域模型
type StartInfo struct {
    Jishen []string `json:"jishen"` // 吉神凶煞（按年月日时排序）
    Xz     string   `json:"xz"`     // 星座
    Sx     string   `json:"sx"`     // 生肖
}

// DetailInfo 表示详细信息领域模型
type DetailInfo struct {
    Zhuxing     ZhuxingInfo   `json:"zhuxing"`     // 天干透出十神
    Sizhu       SizhuInfo     `json:"sizhu"`       // 四柱天干地支
    Canggan     CangganInfo   `json:"canggan"`     // 藏干信息
    Fuxing      FuxingInfo    `json:"fuxing"`      // 藏干十神
    Xingyun     XingyunInfo   `json:"xingyun"`     // 星运信息
    Zizuo       ZizuoInfo     `json:"zizuo"`       // 自坐信息
    Kongwang    KongwangInfo  `json:"kongwang"`    // 空亡信息
    Nayin       NayinInfo     `json:"nayin"`       // 纳音信息
    Shensha     ShenshaInfo   `json:"shensha"`     // 神煞信息
    Dayunshensha []DayunShensha `json:"dayunshensha"` // 大运神煞（按时间排序）
}

// Result 表示API返回结果
type Result struct {
    Errcode int    `json:"errcode"` // 错误码
    Errmsg  string `json:"errmsg"`  // 错误信息
    Notice  string `json:"notice"`  // 提示信息
    Data    struct {
        BaseInfo  BaseInfo  `json:"base_info"`  // 基本信息
        BaziInfo  BaziInfo  `json:"bazi_info"`  // 八字信息
        DayunInfo DayunInfo `json:"dayun_info"` // 大运信息
        StartInfo StartInfo `json:"start_info"` // 起运信息
        DetailInfo DetailInfo `json:"detail_info"` // 详细信息
    } `json:"data"` // 返回数据
}

// BaziService 八字排盘领域服务接口
type BaziService interface {
    Calculate(baseInfo BaseInfo) (*Result, error)
    Analyze(baziInfo BaziInfo) (string, error)
    Predict(dayunInfo DayunInfo) (string, error)
}


// ZhuxingInfo 表示天干透出十神信息
type ZhuxingInfo struct {
    NianGan string `json:"nian_gan"` // 年干透出十神
    YueGan  string `json:"yue_gan"`  // 月干透出十神
    RiGan   string `json:"ri_gan"`    // 日干透出十神
    ShiGan  string `json:"shi_gan"`   // 时干透出十神
}

// SizhuInfo 表示四柱天干地支信息
type SizhuInfo struct {
    NianZhu string `json:"nian_zhu"` // 年柱
    YueZhu  string `json:"yue_zhu"`  // 月柱
    RiZhu   string `json:"ri_zhu"`   // 日柱
    ShiZhu  string `json:"shi_zhu"`  // 时柱
}

// CangganInfo 表示藏干信息
type CangganInfo struct {
    NianZhi []string `json:"nian_zhi"` // 年支藏干（按本气中气余气排序）
    YueZhi  []string `json:"yue_zhi"`  // 月支藏干（按本气中气余气排序）
    RiZhi   []string `json:"ri_zhi"`   // 日支藏干（按本气中气余气排序）
    ShiZhi  []string `json:"shi_zhi"`  // 时支藏干（按本气中气余气排序）
}

// FuxingInfo 表示藏干十神信息
type FuxingInfo struct {
    NianZhi []string `json:"nian_zhi"` // 年支藏干十神（按本气中气余气排序）
    YueZhi  []string `json:"yue_zhi"`  // 月支藏干十神（按本气中气余气排序）
    RiZhi   []string `json:"ri_zhi"`   // 日支藏干十神（按本气中气余气排序）
    ShiZhi  []string `json:"shi_zhi"`  // 时支藏干十神（按本气中气余气排序）
}

// XingyunInfo 表示星运信息
type XingyunInfo struct {
    Xing  string `json:"xing"`  // 星运
    Yun   string `json:"yun"`   // 运程
    Jieqi string `json:"jieqi"` // 节气
}

// ZizuoInfo 表示自坐信息
type ZizuoInfo struct {
    NianZhu string `json:"nian_zhu"` // 年柱自坐
    YueZhu  string `json:"yue_zhu"`  // 月柱自坐
    RiZhu   string `json:"ri_zhu"`   // 日柱自坐
    ShiZhu  string `json:"shi_zhu"`  // 时柱自坐
}

// KongwangInfo 表示空亡信息
type KongwangInfo struct {
    NianZhu string `json:"nian_zhu"` // 年柱空亡
    YueZhu  string `json:"yue_zhu"`  // 月柱空亡
    RiZhu   string `json:"ri_zhu"`   // 日柱空亡
    ShiZhu  string `json:"shi_zhu"`  // 时柱空亡
}

// NayinInfo 表示纳音信息
type NayinInfo struct {
    NianZhu string `json:"nian_zhu"` // 年柱纳音
    YueZhu  string `json:"yue_zhu"`  // 月柱纳音
    RiZhu   string `json:"ri_zhu"`   // 日柱纳音
    ShiZhu  string `json:"shi_zhu"`  // 时柱纳音
}

// ShenshaInfo 表示神煞信息
type ShenshaInfo struct {
    NianZhu []string `json:"nian_zhu"` // 年柱神煞（按出现顺序排序）
    YueZhu  []string `json:"yue_zhu"`  // 月柱神煞（按出现顺序排序）
    RiZhu   []string `json:"ri_zhu"`   // 日柱神煞（按出现顺序排序）
    ShiZhu  []string `json:"shi_zhu"`  // 时柱神煞（按出现顺序排序）
}

// DayunShensha 表示大运神煞信息
type DayunShensha struct {
    BigYun  string   `json:"big_yun"`  // 大运
    Shensha []string `json:"shensha"`  // 对应神煞（按出现顺序排序）
}
