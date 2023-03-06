package store

import (
	"context"
	"fmt"
	"github.com/zhan3333/goutil"
	"gongfu/internal/model"
)

var _ Role = (*DBStore)(nil)

type Role interface {
	GetUserRoles(ctx context.Context, userID uint) ([]string, error)
	// SyncUserRoles 将用户的角色设置为传入的角色
	// 保持传入的角色与用户的角色一致
	SyncUserRoles(ctx context.Context, userId uint, roleNames []string) error
	GetRoleNames(ctx context.Context) ([]string, error)
}

func (s DBStore) GetRoleNames(ctx context.Context) ([]string, error) {
	roleNames := []string{}
	err := s.DB.WithContext(ctx).Model(&model.Role{}).Pluck("name", &roleNames).Error
	if err != nil {
		return nil, err
	}
	return roleNames, nil
}

func (s DBStore) GetUserRoles(ctx context.Context, userID uint) ([]string, error) {
	roleNames := []string{}
	err := s.DB.WithContext(ctx).
		Model(&model.UserHasRole{}).
		Where("user_id = ?", userID).
		Pluck("role_name", &roleNames).Error
	if err != nil {
		return nil, err
	}
	return roleNames, nil
}

func (s DBStore) SyncUserRoles(ctx context.Context, userId uint, roleNames []string) error {
	roles := []model.Role{}
	err := s.DB.WithContext(ctx).Model(&model.Role{}).Where("name in (?)", roleNames).Find(&roles).Error
	if err != nil {
		return err
	}
	roleIds := []uint{}
	roleIdToName := map[uint]string{}
	for _, role := range roles {
		roleIds = append(roleIds, role.ID)
		roleIdToName[role.ID] = role.Name
	}

	userHasRoles := []model.UserHasRole{}
	err = s.DB.WithContext(ctx).Model(&model.UserHasRole{}).Where("user_id = ?", userId).Find(&userHasRoles).Error
	if err != nil {
		return err
	}
	hasRoleIds := []uint{}
	for _, hasRole := range userHasRoles {
		hasRoleIds = append(hasRoleIds, hasRole.RoleID)
	}
	// 在 hasRoles 中，不在 roles 中
	waitRemove := util.Diff(hasRoleIds, roleIds)
	// 在 roles 中， 不在 hasRoles 中
	waitCreate := util.Diff(roleIds, hasRoleIds)
	if len(waitRemove) > 0 {
		err = s.DB.WithContext(ctx).Model(&model.UserHasRole{}).
			Where("user_id = ?", userId).
			Where("role_id in (?)", waitRemove).
			Delete(&model.UserHasRole{}).Error
		if err != nil {
			return fmt.Errorf("delete roles failed")
		}
	}
	if len(waitCreate) > 0 {
		hasRoles := []model.UserHasRole{}
		for _, roleId := range waitCreate {
			hasRoles = append(hasRoles, model.UserHasRole{
				UserID:   userId,
				RoleID:   roleId,
				RoleName: roleIdToName[roleId],
			})
		}
		err = s.DB.WithContext(ctx).Model(&model.UserHasRole{}).Create(&hasRoles).Error
		if err != nil {
			return fmt.Errorf("create roles failed")
		}
	}
	return nil
}
