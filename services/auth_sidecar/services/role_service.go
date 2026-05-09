package services

import "team_dynamics/auth_sidecar/models"

type RoleService interface {
	ResolveRoles(authorityMap models.AuthorityMap, userId models.UserId) []string
}

type roleServiceImpl struct{}

func NewRoleService() RoleService {
	return &roleServiceImpl{}
}

func (s *roleServiceImpl) ResolveRoles(authorityMap models.AuthorityMap, userId models.UserId) []string {
	roleSet := make(map[string]struct{})
	for key, roles := range authorityMap {
		if key == (models.UserId{}) || key == userId {
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
