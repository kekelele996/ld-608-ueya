package constants

// UserRole backs the RBAC matrix enforced by rbacMiddleware and mirrored in
// frontend route guards and button visibility.
type UserRole string

const (
	RoleDispatcher UserRole = "DISPATCHER" // 地勤调度
	RoleCrew       UserRole = "CREW"       // 班组
	RoleResource   UserRole = "RESOURCE"   // 资源管理员
	RoleSupervisor UserRole = "SUPERVISOR" // 运行督导
)

var UserRoles = []UserRole{RoleDispatcher, RoleCrew, RoleResource, RoleSupervisor}

var UserRoleText = map[UserRole]string{
	RoleDispatcher: "地勤调度",
	RoleCrew:       "班组",
	RoleResource:   "资源管理员",
	RoleSupervisor: "运行督导",
}

// SeedAccounts lists demo accounts shipped with local seed data.
var SeedAccounts = []struct {
	Username string
	Password string
	Role     UserRole
	TeamID   string
	Name     string
}{
	{"dispatcher", "dispatch123", RoleDispatcher, "TEAM-DISPATCH", "调度员·林岚"},
	{"crew", "crew123", RoleCrew, "TEAM-CLEAN", "清洁班组长·周海"},
	{"resource", "resource123", RoleResource, "TEAM-RESOURCE", "资源管理员·高磊"},
	{"supervisor", "super123", RoleSupervisor, "TEAM-OPS", "运行督导·陈屿"},
}

func (r UserRole) Text() string { return UserRoleText[r] }
