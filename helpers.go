package main

import "strings"

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

func checkExt(file string) bool {
	if len(varBlacklistExts) == 0 {
		return true
	}
	for _, ext := range varBlacklistExts {
		if strings.HasSuffix(file, ext) {
			return false
		}
	}
	return true
}

func checkLen(in string) bool {
	if len(in) == 0 {
		return false
	}
	i := 0
	for _, v := range in {
		if v == ' ' {
			continue
		}
		i++
		if i > 2 {
			return true
		}
	}
	return false
}
