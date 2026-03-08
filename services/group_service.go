package services

import (
	"time"

	"github.com/google/uuid"

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

	member := models.GroupMember{
		ID:       uuid.New(),
		GroupID:  group.ID,
		UserID:   req.CreatedBy,
		Role:     "admin",
		JoinedAt: time.Now(),
	}

	db.Create(&member)

	return &group, nil
}

func AddMember(groupID uuid.UUID, req dto.AddGroupMemberRequest) error {

	db := database.DB

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
