package main

func checkUserRole(userRoles []string) bool {
	if len(roles) == 0 {
		return true
	}

	if len(userRoles) == 0 {
		return false
	}

	for _, role := range userRoles {
		if _, ok := roles[role]; ok {
			return true
		}
	}

	return false
}
