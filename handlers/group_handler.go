package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/Vedu3635/PRISM.git/dto"
	"github.com/Vedu3635/PRISM.git/services"
)

func CreateGroup(c *gin.Context) {
	var req dto.CreateGroupRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	group, err := services.CreateGroup(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, group)
}

func GetGroups(c *gin.Context) {
	groups, err := services.GetGroups()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, groups)
}

func GetGroupsByID(c *gin.Context) {
	id, err := parseGroupID(c)
	if err != nil {
		return
	}

	group, err := services.GetGroupsByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "group not found"})
		return
	}

	c.JSON(http.StatusOK, group)
}

// GetGroupsByUserID returns all groups the user is a member of.
func GetGroupsByUserID(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	groups, err := services.GetGroupsByUserID(userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, groups)
}

func UpdateGroup(c *gin.Context) {
	id, err := parseGroupID(c)
	if err != nil {
		return
	}

	var req dto.UpdateGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	group, err := services.UpdateGroup(id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, group)
}

func DeleteGroup(c *gin.Context) {
	id, err := parseGroupID(c)
	if err != nil {
		return
	}

	if err := services.DeleteGroup(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "group deleted"})
}

func AddMember(c *gin.Context) {
	id, err := parseGroupID(c)
	if err != nil {
		return
	}

	var req dto.AddGroupMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := services.AddMember(id, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "member added"})
}

func GetGroupMembers(c *gin.Context) {
	id, err := parseGroupID(c)
	if err != nil {
		return
	}

	members, err := services.GetGroupMembers(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, members)
}

func RemoveMember(c *gin.Context) {
	groupID, err := parseGroupID(c)
	if err != nil {
		return
	}

	memberIDParam := c.Param("memberID")
	memberID, err := uuid.Parse(memberIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid member id"})
		return
	}

	if err := services.RemoveMember(groupID, memberID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "member removed"})
}

func LeaveGroup(c *gin.Context) {
	groupID, err := parseGroupID(c)
	if err != nil {
		return
	}

	var req dto.LeaveGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := services.LeaveGroup(groupID, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "left group successfully"})
}

func GetGroupBalances(c *gin.Context) {
	id, err := parseGroupID(c)
	if err != nil {
		return
	}

	balances, err := services.GetGroupBalances(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, balances)
}

// ─── Helpers ──────────────────────────────────────────────────────────────────

// parseGroupID parses and validates the :id route param as a group UUID.
func parseGroupID(c *gin.Context) (uuid.UUID, error) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid group id"})
	}
	return id, err
}
