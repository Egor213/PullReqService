package service

import e "app/internal/entity"

func CompareMembers(a, b []e.TeamMember) bool {
	if len(a) != len(b) {
		return false
	}

	amap := make(map[string]struct{}, len(a))
	for _, m := range a {
		amap[m.UserID] = struct{}{}
	}

	for _, m := range b {
		if _, ok := amap[m.UserID]; !ok {
			return false
		}
	}
	return true
}
