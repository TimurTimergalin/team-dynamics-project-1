package services

import "team_dynamics/auth_sidecar/models"

type RoleService interface {
	ResolveRolesForUser(authorityMap models.AuthorityMap, userId models.UserId) []string
	ResolveRolesForServiceAccount(authorityMap models.AuthorityMap, serviceAccount string) []string
}

type roleServiceImpl struct{}

func NewRoleService() RoleService {
	return &roleServiceImpl{}
}

func collectRoles(authorityMap models.AuthorityMap, match func(models.Principal) bool) []string {
	roleSet := make(map[string]struct{})
	for key, roles := range authorityMap {
		if match(key) {
			for _, role := range roles {
				roleSet[role] = struct{}{}
			}
		}
	}
	result := make([]string, 0, len(roleSet))
	for role := range roleSet {
		result = append(result, role)
	}
	return result
}

func (s *roleServiceImpl) ResolveRolesForUser(authorityMap models.AuthorityMap, userId models.UserId) []string {
	return collectRoles(authorityMap, func(p models.Principal) bool {
		if p == (models.Principal{}) {
			return true // AnyUser
		}
		return p.UserId != nil && *p.UserId == userId
	})
}

func (s *roleServiceImpl) ResolveRolesForServiceAccount(authorityMap models.AuthorityMap, serviceAccount string) []string {
	return collectRoles(authorityMap, func(p models.Principal) bool {
		return p.ServiceAccount != nil && *p.ServiceAccount == serviceAccount
	})
}
