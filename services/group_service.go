package services

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/Vedu3635/PRISM.git/database"
	"github.com/Vedu3635/PRISM.git/dto"
	"github.com/Vedu3635/PRISM.git/models"
)

func CreateGroup(req dto.CreateGroupRequest) (*models.Group, error) {
	db := database.DB

	group := models.Group{
		ID:          uuid.New(),
		CreatedBy:   req.CreatedBy,
		Name:        req.Name,
		Description: req.Description,
		Type:        req.Type,
		Currency:    req.Currency,
		InviteCode:  uuid.New().String()[:8],
		IsActive:    true,
		IsPersonal:  false,
		CreatedAt:   time.Now(),
	}

	if err := db.Create(&group).Error; err != nil {
		return nil, err
	}

	// Creator is automatically made admin
	member := models.GroupMember{
		ID:       uuid.New(),
		GroupID:  group.ID,
		UserID:   req.CreatedBy,
		Role:     "admin",
		JoinedAt: time.Now(),
	}

	if err := db.Create(&member).Error; err != nil {
		return nil, err
	}

	return &group, nil
}

func GetGroups() ([]models.Group, error) {
	db := database.DB

	var groups []models.Group
	err := db.Where("is_active = ?", true).Find(&groups).Error

	return groups, err
}

func GetGroupsByID(id uuid.UUID) (*models.Group, error) {
	db := database.DB

	var group models.Group
	if err := db.First(&group, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return &group, nil
}

// GetGroupsByUserID returns all active groups the user is a member of.
func GetGroupsByUserID(userID uuid.UUID) ([]models.Group, error) {
	db := database.DB

	var groups []models.Group
	err := db.Joins("JOIN group_members ON group_members.group_id = groups.id").
		Where("group_members.user_id = ? AND groups.is_active = ?", userID, true).
		Find(&groups).Error

	return groups, err
}

func UpdateGroup(id uuid.UUID, req dto.UpdateGroupRequest) (*models.Group, error) {
	db := database.DB

	var group models.Group
	if err := db.First(&group, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("group not found")
		}
		return nil, err
	}

	if req.Name != nil {
		group.Name = *req.Name
	}
	if req.CreatedBy != nil {
		group.CreatedBy = *req.CreatedBy
	}
	if req.Description != nil {
		group.Description = req.Description
	}
	if req.Type != nil {
		group.Type = *req.Type
	}
	if req.Currency != nil {
		group.Currency = *req.Currency
	}

	if err := db.Save(&group).Error; err != nil {
		return nil, err
	}

	return &group, nil
}

func DeleteGroup(id uuid.UUID) error {
	db := database.DB

	result := db.Model(&models.Group{}).Where("id = ?", id).Update("is_active", false)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("group not found")
	}

	return nil
}

func AddMember(groupID uuid.UUID, req dto.AddGroupMemberRequest) error {
	db := database.DB

	// Prevent duplicate membership
	var existing models.GroupMember
	err := db.Where("group_id = ? AND user_id = ?", groupID, req.UserID).First(&existing).Error
	if err == nil {
		return errors.New("user is already a member of this group")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	member := models.GroupMember{
		ID:       uuid.New(),
		GroupID:  groupID,
		UserID:   req.UserID,
		Role:     req.Role,
		JoinedAt: time.Now(),
	}

	return db.Create(&member).Error
}

func GetGroupMembers(groupID uuid.UUID) ([]models.GroupMember, error) {
	db := database.DB

	var members []models.GroupMember
	err := db.Where("group_id = ?", groupID).Find(&members).Error

	return members, err
}

func RemoveMember(groupID uuid.UUID, memberID uuid.UUID) error {
	db := database.DB

	result := db.Where("group_id = ? AND user_id = ?", groupID, memberID).Delete(&models.GroupMember{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("member not found in group")
	}

	return nil
}

// LeaveGroup now takes the callerID directly from the token — no DTO needed.
func LeaveGroup(groupID uuid.UUID, callerID uuid.UUID) error {
	db := database.DB

	var member models.GroupMember
	if err := db.Where("group_id = ? AND user_id = ?", groupID, callerID).First(&member).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("you are not a member of this group")
		}
		return err
	}

	if member.Role == "admin" {
		return errors.New("admin cannot leave the group without transferring ownership first")
	}

	result := db.Where("group_id = ? AND user_id = ?", groupID, callerID).Delete(&models.GroupMember{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("member not found in group")
	}

	return nil
}

// GetGroupBalances returns all balances for a given group.
func GetGroupBalances(groupID uuid.UUID) ([]models.Balance, error) {
	db := database.DB

	var balances []models.Balance
	if err := db.Where("group_id = ?", groupID).Find(&balances).Error; err != nil {
		return nil, err
	}

	return balances, nil
}

// ─── Internal helpers ─────────────────────────────────────────────────────────

// IsGroupAdmin checks whether the given user has the admin role in the group.
func IsGroupAdmin(groupID uuid.UUID, userID uuid.UUID) (bool, error) {
	db := database.DB

	var member models.GroupMember
	err := db.Where("group_id = ? AND user_id = ? AND role = ?", groupID, userID, "admin").
		First(&member).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return true, nil
}

// IsGroupMember checks whether the given user is a member of the group (any role).
func IsGroupMember(groupID uuid.UUID, userID uuid.UUID) (bool, error) {
	db := database.DB

	var member models.GroupMember
	err := db.Where("group_id = ? AND user_id = ?", groupID, userID).
		First(&member).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return true, nil
}
