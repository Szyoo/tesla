// Package battery 提供按 VIN 估算电池容量的工具，供充电/行程能耗估算共用。
package battery

// capacityByVINPrefix 按 VIN 前三位（生产地/车型）粗略估算电池容量（kWh）。
// 注意：这是缺少精确车型配置时的兜底估算，并非区域特定值——同一前缀的车型
// 无论销往中国、日本还是其他市场，电池容量一致。
//
// VIN 前缀含义（前三位 = WMI，世界制造厂识别码）：
//   LRW = 上海工厂（Model 3 / Y，含出口日本的车型）
//   5YJ = 美国 Fremont 工厂（Model S / X / 3 / Y）
//   7SA = 美国 Fremont 工厂（较新批次）
//   XP7 = 德国柏林工厂（Model Y）
var capacityByVINPrefix = map[string]float64{
	"LRW": 60.0,
	"5YJ": 75.0,
	"7SA": 78.0,
	"XP7": 100.0,
}

// defaultCapacity 在 VIN 缺失或前缀未知时返回的兜底容量（kWh）。
const defaultCapacity = 60.0

// CapacityByVIN 根据 VIN 估算电池容量（kWh）。VIN 无效或前缀未知时返回兜底值。
func CapacityByVIN(vin string) float64 {
	if len(vin) < 3 {
		return defaultCapacity
	}
	if cap, ok := capacityByVINPrefix[vin[0:3]]; ok {
		return cap
	}
	return defaultCapacity
}
