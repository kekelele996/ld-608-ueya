package constants

// RBACMatrix maps actions to the roles allowed to perform them. rbacMiddleware
// consults this table; frontend permission helpers mirror it for button hiding.
var RBACMatrix = map[string][]UserRole{
	ActionTurnaroundCreate:  {RoleDispatcher, RoleSupervisor},
	ActionTurnaroundArrive:  {RoleDispatcher, RoleSupervisor},
	ActionTurnaroundRelease: {RoleDispatcher, RoleSupervisor},

	ActionPlanGenerate: {RoleDispatcher, RoleSupervisor},
	ActionTaskDispatch: {RoleDispatcher, RoleSupervisor},
	ActionTaskSign:     {RoleCrew, RoleDispatcher, RoleSupervisor},
	ActionTaskFinish:   {RoleCrew, RoleDispatcher, RoleSupervisor},
	ActionTaskBlock:    {RoleCrew, RoleDispatcher, RoleSupervisor},

	ActionResourceCreate:      {RoleResource, RoleSupervisor},
	ActionResourceUpdate:      {RoleResource, RoleSupervisor},
	ActionResourceMaintenance: {RoleResource, RoleSupervisor},

	ActionBookingAdjust:  {RoleResource, RoleDispatcher, RoleSupervisor},
	ActionBookingConfirm: {RoleResource, RoleDispatcher, RoleSupervisor},
	ActionBookingRelease: {RoleResource, RoleDispatcher, RoleSupervisor},

	ActionDelayRegister: {RoleDispatcher, RoleSupervisor, RoleCrew},
	ActionDelayReplan:   {RoleDispatcher, RoleSupervisor},
	ActionDelayResolve:  {RoleSupervisor, RoleDispatcher},
}

func Can(role UserRole, action string) bool {
	allowed, ok := RBACMatrix[action]
	if !ok {
		return false
	}
	for _, candidate := range allowed {
		if candidate == role {
			return true
		}
	}
	return false
}
