package validation

import (
	"errors"
	"mkluxe-backend/internal/dto"
)

func ValidateChangePassword(req *dto.ChangePasswordRequest) error {
	if req.OldPassword == req.NewPassword {
		return errors.New("new password cannot be the identical to the old password")
	}
	return nil
}
