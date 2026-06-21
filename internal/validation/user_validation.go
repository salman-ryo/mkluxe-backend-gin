package validation

import (
	"errors"
	"mkluxe-backend/internal/constants"
	"mkluxe-backend/internal/dto"
)

func ValidateAdminUser(req *dto.CreateUserRequest) error {
	switch req.Role {
	case constants.RoleSuperAdmin, constants.RoleAdmin, constants.RoleEditor, constants.RoleSupport:
		return nil
	}
	return errors.New("invalid role assignment")
}
