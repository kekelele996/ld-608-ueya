package constants

// RBAC roles. Reaches: DB users table, auth/rbac middlewares, routes,
// frontend router guard, auth store, button visibility, log templates.
const (
	RoleDispatcher      = "DISPATCHER"       // 地勤调度
	RoleTeam            = "TEAM"             // 班组
	RoleResourceManager = "RESOURCE_MANAGER" // 资源管理员
	RoleSupervisor      = "SUPERVISOR"       // 运行督导
)

var Roles = []string{RoleDispatcher, RoleTeam, RoleResourceManager, RoleSupervisor}

var RoleText = map[string]string{
	RoleDispatcher:      "地勤调度",
	RoleTeam:            "保障班组",
	RoleResourceManager: "资源管理员",
	RoleSupervisor:      "运行督导",
}

func IsValidRole(value string) bool {
	for _, item := range Roles {
		if item == value {
			return true
		}
	}
	return false
}

// Delay types for DelayEvent.delay_type.
const (
	DelayLateArrival = "LATE_ARRIVAL"
	DelayWeather     = "WEATHER"
	DelayAirline     = "AIRLINE"
	DelaySecurity    = "SECURITY"
	DelayGround      = "GROUND"
)

var DelayType = []string{DelayLateArrival, DelayWeather, DelayAirline, DelaySecurity, DelayGround}

var DelayTypeText = map[string]string{
	DelayLateArrival: "晚到",
	DelayWeather:     "天气",
	DelayAirline:     "航司原因",
	DelaySecurity:    "安检",
	DelayGround:      "地勤原因",
}
