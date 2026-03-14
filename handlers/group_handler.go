package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/Vedu3635/PRISM.git/dto"
	"github.com/Vedu3635/PRISM.git/services"
)

// CreateGroup godoc
//
//	@Summary		Create a group
//	@Description	Creates a new expense-splitting group. The caller is automatically added as admin.
//	@Tags			groups
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			body	body		dto.CreateGroupRequest	true	"Group payload"
//	@Success		201		{object}	models.Group			"created group"
//	@Failure		400		{object}	map[string]string		"validation error"
//	@Failure		500		{object}	map[string]string		"internal server error"
//	@Router			/groups [post]
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

// GetGroups godoc
//
//	@Summary		List all groups
//	@Description	Returns all active groups.
//	@Tags			groups
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{array}		models.Group		"list of groups"
//	@Failure		500	{object}	map[string]string	"internal server error"
//	@Router			/groups [get]
func GetGroups(c *gin.Context) {
	groups, err := services.GetGroups()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, groups)
}

// GetGroupsByID godoc
//
//	@Summary		Get group by ID
//	@Description	Returns a single group by UUID.
//	@Tags			groups
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		string				true	"Group UUID"
//	@Success		200	{object}	models.Group		"group"
//	@Failure		400	{object}	map[string]string	"invalid id"
//	@Failure		404	{object}	map[string]string	"group not found"
//	@Router			/groups/{id} [get]
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

// GetGroupsByUserID godoc
//
//	@Summary		Get groups for current user
//	@Description	Returns all groups the authenticated user is a member of.
//	@Tags			groups
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{array}		models.Group		"groups"
//	@Failure		401	{object}	map[string]string	"unauthorized"
//	@Failure		500	{object}	map[string]string	"internal server error"
//	@Router			/users/groups [get]
func GetGroupsByUserID(c *gin.Context) {
	callerID, ok := extractCallerID(c)
	if !ok {
		return
	}

	groups, err := services.GetGroupsByUserID(callerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, groups)
}

// UpdateGroup godoc
//
//	@Summary		Update a group
//	@Description	Updates mutable fields on a group. Only group admins can perform this action.
//	@Tags			groups
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		string					true	"Group UUID"
//	@Param			body	body		dto.UpdateGroupRequest	true	"Fields to update"
//	@Success		200		{object}	models.Group			"updated group"
//	@Failure		400		{object}	map[string]string		"invalid id or payload"
//	@Failure		403		{object}	map[string]string		"forbidden — admins only"
//	@Failure		500		{object}	map[string]string		"internal server error"
//	@Router			/groups/{id} [put]
func UpdateGroup(c *gin.Context) {
	id, err := parseGroupID(c)
	if err != nil {
		return
	}

	callerID, ok := extractCallerID(c)
	if !ok {
		return
	}

	// ── Admin check ───────────────────────────────────────────────────────────
	if !requireGroupAdmin(c, id, callerID) {
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

// DeleteGroup godoc
//
//	@Summary		Delete a group
//	@Description	Soft-deletes a group (is_active = false). Only group admins can perform this action.
//	@Tags			groups
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		string				true	"Group UUID"
//	@Success		200	{object}	map[string]string	"message"
//	@Failure		400	{object}	map[string]string	"invalid id"
//	@Failure		403	{object}	map[string]string	"forbidden — admins only"
//	@Failure		500	{object}	map[string]string	"internal server error"
//	@Router			/groups/{id} [delete]
func DeleteGroup(c *gin.Context) {
	id, err := parseGroupID(c)
	if err != nil {
		return
	}

	callerID, ok := extractCallerID(c)
	if !ok {
		return
	}

	// ── Admin check ───────────────────────────────────────────────────────────
	if !requireGroupAdmin(c, id, callerID) {
		return
	}

	if err := services.DeleteGroup(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "group deleted"})
}

// AddMember godoc
//
//	@Summary		Add a member to a group
//	@Description	Adds a user to the specified group with a given role (admin or member). Only group admins can add members.
//	@Tags			group-members
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		string						true	"Group UUID"
//	@Param			body	body		dto.AddGroupMemberRequest	true	"Member payload"
//	@Success		201		{object}	map[string]string			"message"
//	@Failure		400		{object}	map[string]string			"invalid id or payload"
//	@Failure		403		{object}	map[string]string			"forbidden — admins only"
//	@Failure		500		{object}	map[string]string			"internal server error"
//	@Router			/groups/{id}/members [post]
func AddMember(c *gin.Context) {
	id, err := parseGroupID(c)
	if err != nil {
		return
	}

	callerID, ok := extractCallerID(c)
	if !ok {
		return
	}

	// ── Admin check ───────────────────────────────────────────────────────────
	if !requireGroupAdmin(c, id, callerID) {
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

// GetGroupMembers godoc
//
//	@Summary		List group members
//	@Description	Returns all current members of the specified group.
//	@Tags			group-members
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		string				true	"Group UUID"
//	@Success		200	{array}		models.GroupMember	"members"
//	@Failure		400	{object}	map[string]string	"invalid id"
//	@Failure		500	{object}	map[string]string	"internal server error"
//	@Router			/groups/{id}/members [get]
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

// RemoveMember godoc
//
//	@Summary		Remove a member from a group
//	@Description	Removes the specified member from the group. Only group admins can remove members.
//	@Tags			group-members
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id			path		string				true	"Group UUID"
//	@Param			memberID	path		string				true	"Member UUID"
//	@Success		200			{object}	map[string]string	"message"
//	@Failure		400			{object}	map[string]string	"invalid id"
//	@Failure		403			{object}	map[string]string	"forbidden — admins only"
//	@Failure		500			{object}	map[string]string	"internal server error"
//	@Router			/groups/{id}/members/{memberID} [delete]
func RemoveMember(c *gin.Context) {
	groupID, err := parseGroupID(c)
	if err != nil {
		return
	}

	callerID, ok := extractCallerID(c)
	if !ok {
		return
	}

	// ── Admin check ───────────────────────────────────────────────────────────
	if !requireGroupAdmin(c, groupID, callerID) {
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

// LeaveGroup godoc
//
//	@Summary		Leave a group
//	@Description	Removes the authenticated user from the group. Admins must transfer ownership before leaving.
//	@Tags			groups
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		string				true	"Group UUID"
//	@Success		200	{object}	map[string]string	"message"
//	@Failure		400	{object}	map[string]string	"invalid id"
//	@Failure		401	{object}	map[string]string	"unauthorized"
//	@Failure		422	{object}	map[string]string	"admin cannot leave without transferring ownership"
//	@Failure		500	{object}	map[string]string	"internal server error"
//	@Router			/groups/{id}/leave [post]
func LeaveGroup(c *gin.Context) {
	groupID, err := parseGroupID(c)
	if err != nil {
		return
	}

	// ── Extract caller from token — do NOT trust request body ─────────────────
	callerID, ok := extractCallerID(c)
	if !ok {
		return
	}

	if err := services.LeaveGroup(groupID, callerID); err != nil {
		if err.Error() == "admin cannot leave the group without transferring ownership first" {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "left group successfully"})
}

// GetGroupBalances godoc
//
//	@Summary		Get group balances
//	@Description	Returns the net balance for each member in the group.
//	@Tags			groups
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		string					true	"Group UUID"
//	@Success		200	{object}	map[string]interface{}	"balances per member"
//	@Failure		400	{object}	map[string]string		"invalid id"
//	@Failure		500	{object}	map[string]string		"internal server error"
//	@Router			/groups/{id}/balances [get]
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

func parseGroupID(c *gin.Context) (uuid.UUID, error) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid group id"})
	}
	return id, err
}

// requireGroupAdmin checks that the caller is an admin of the group.
// Writes 403 and returns false if not.
func requireGroupAdmin(c *gin.Context, groupID uuid.UUID, callerID uuid.UUID) bool {
	isAdmin, err := services.IsGroupAdmin(groupID, callerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return false
	}
	if !isAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "only group admins can perform this action"})
		return false
	}
	return true
}
