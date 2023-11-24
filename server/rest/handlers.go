package rest

import (
	"fangaoxs.com/go-chat/environment"
	"fangaoxs.com/go-chat/internal/auth"
	"fangaoxs.com/go-chat/internal/domain/group"
	"fangaoxs.com/go-chat/internal/domain/groupmember"
	"fangaoxs.com/go-chat/internal/domain/user"
	"fangaoxs.com/go-chat/internal/entity"
	"fangaoxs.com/go-chat/internal/infras/errors"
	"fangaoxs.com/go-chat/internal/infras/logger"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func NewHandlers(
	env environment.Env,
	logger logger.Logger,
	user user.User,
	group group.Group,
	groupMember groupmember.GroupMember,
) (Handlers, error) {
	return Handlers{
		logger:      logger,
		user:        user,
		group:       group,
		groupMember: groupMember,
	}, nil
}

type Handlers struct {
	logger logger.Logger

	user        user.User
	group       group.Group
	groupMember groupmember.GroupMember
}

func (h *Handlers) RegisterUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		// POST
		ctx := c.Request.Context()

		if c.PostForm("nickname") == "" {
			WrapGinError(c, errors.Newf(errors.InvalidArgument, nil, "invalid nickname"))
			return
		}
		if c.PostForm("username") == "" {
			WrapGinError(c, errors.Newf(errors.InvalidArgument, nil, "invalid username"))
			return
		}
		if c.PostForm("password") == "" {
			WrapGinError(c, errors.Newf(errors.InvalidArgument, nil, "invalid password"))
			return
		}
		if c.PostForm("phone") == "" {
			WrapGinError(c, errors.Newf(errors.InvalidArgument, nil, "invalid phone"))
			return
		}

		input := user.RegisterInput{
			Nickname: c.PostForm("nickname"),
			Username: c.PostForm("username"),
			Password: c.PostForm("password"),
			Phone:    c.PostForm("phone"),
		}
		subject, err := h.user.RegisterUser(ctx, input)
		if err != nil {
			WrapGinError(c, err)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"subject": subject,
		})
	}
}

func (h *Handlers) Me() gin.HandlerFunc {
	return func(c *gin.Context) {
		// GET
		ctx := c.Request.Context()
		ui := auth.FromContext(ctx)

		u, err := h.user.GetUserBySubject(ctx, ui.Subject)
		if err != nil {
			WrapGinError(c, err)
			return
		}

		c.JSON(http.StatusOK, u)
	}
}

func (h *Handlers) MyGroups() gin.HandlerFunc {
	return func(c *gin.Context) {
		// GET
		ctx := c.Request.Context()
		ui := auth.FromContext(ctx)

		groups, err := h.groupMember.ListGroupsOfUser(ctx, ui.Subject)
		if err != nil {
			WrapGinError(c, err)
			return
		}

		c.JSON(http.StatusOK, groups)
	}
}

func (h *Handlers) InsertGroup() gin.HandlerFunc {
	return func(c *gin.Context) {
		// POST
		ctx := c.Request.Context()
		ui := auth.FromContext(ctx)

		if c.PostForm("name") == "" {
			WrapGinError(c, errors.Newf(errors.InvalidArgument, nil, "invalid name"))
			return
		}
		groupType, ok := entity.GroupTypeFromString(c.PostForm("type"))
		if !ok {
			groupType = entity.DefaultGroupType
		}

		input := group.InsertGroupInput{
			Name:      c.PostForm("name"),
			Type:      groupType,
			CreatedBy: ui.Subject,
		}
		id, err := h.group.InsertGroup(ctx, input)
		if err != nil {
			WrapGinError(c, err)
			return
		}

		c.JSON(http.StatusOK, gin.H{"id": id})
	}
}

func (h *Handlers) GetGroupByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		// GET
		ctx := c.Request.Context()
		ui := auth.FromContext(ctx)

		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			WrapGinError(c, err)
			return
		}

		g, err := h.group.GetGroupByID(ctx, id)
		if err != nil {
			WrapGinError(c, err)
			return
		}
		if !g.IsPublic && g.CreatedBy != ui.Subject {
			WrapGinError(c, errors.New(errors.PermissionDenied, nil, "你不可以查看该群"))
			return
		}

		c.JSON(http.StatusOK, g)
	}
}

func (h *Handlers) DeleteGroup() gin.HandlerFunc {
	return func(c *gin.Context) {
		// DELETE
		ctx := c.Request.Context()
		ui := auth.FromContext(ctx)

		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			WrapGinError(c, err)
			return
		}

		g, err := h.group.GetGroupByID(ctx, id)
		if err != nil {
			WrapGinError(c, err)
			return
		}
		if g.CreatedBy != ui.Subject {
			WrapGinError(c, errors.New(errors.PermissionDenied, nil, "你不可以删除该群"))
			return
		}

		err = h.group.DeleteGroup(ctx, id)
		if err != nil {
			WrapGinError(c, err)
			return
		}

		c.Status(http.StatusOK)
	}
}

func (h *Handlers) PrivateGroup() gin.HandlerFunc {
	return func(c *gin.Context) {
		// PUT
		ctx := c.Request.Context()
		ui := auth.FromContext(ctx)

		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			WrapGinError(c, err)
			return
		}

		g, err := h.group.GetGroupByID(ctx, id)
		if err != nil {
			WrapGinError(c, err)
			return
		}
		if g.CreatedBy != ui.Subject {
			WrapGinError(c, errors.New(errors.PermissionDenied, nil, "你不可以操作该群"))
			return
		}

		err = h.group.PrivateGroup(ctx, id)
		if err != nil {
			WrapGinError(c, err)
			return
		}

		c.Status(http.StatusOK)
	}
}

func (h *Handlers) PublicGroup() gin.HandlerFunc {
	return func(c *gin.Context) {
		// PUT
		ctx := c.Request.Context()
		ui := auth.FromContext(ctx)

		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			WrapGinError(c, err)
			return
		}

		g, err := h.group.GetGroupByID(ctx, id)
		if err != nil {
			WrapGinError(c, err)
			return
		}
		if g.CreatedBy != ui.Subject {
			WrapGinError(c, errors.New(errors.PermissionDenied, nil, "你不可以操作该群"))
			return
		}

		err = h.group.PublicGroup(ctx, id)
		if err != nil {
			WrapGinError(c, err)
			return
		}

		c.Status(http.StatusOK)
	}
}

func (h *Handlers) AssignUsersToGroup() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		ui := auth.FromContext(ctx)

		groupID, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			WrapGinError(c, err)
			return
		}

		subjects := c.PostFormArray("user_subject")
		if len(subjects) == 0 {
			WrapGinError(c, errors.New(errors.InvalidArgument, nil, "empty user subjects"))
			return
		}

		g, err := h.group.GetGroupByID(ctx, groupID)
		if err != nil {
			WrapGinError(c, err)
			return
		}
		if g.CreatedBy != ui.Subject {
			WrapGinError(c, errors.New(errors.PermissionDenied, nil, "你不可以操作该群"))
			return
		}

		for _, subject := range subjects {
			if err = h.groupMember.AddUserToGroup(ctx, subject, groupID); err != nil {
				WrapGinError(c, err)
				return
			}
		}

		c.Status(http.StatusOK)
	}
}

func (h *Handlers) GroupMembers() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		ui := auth.FromContext(ctx)

		groupID, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			WrapGinError(c, err)
			return
		}

		members, err := h.groupMember.ListUsersOfGroup(ctx, groupID)
		if err != nil {
			WrapGinError(c, err)
			return
		}

		// 如果当前用户存在该群组中
		for _, member := range members {
			if member.Subject == ui.Subject {
				c.JSON(http.StatusOK, members)
				return
			}
		}

		WrapGinError(c, errors.New(errors.PermissionDenied, nil, "你无法查看该群组成员"))
		return
	}
}
