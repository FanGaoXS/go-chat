package rest

import (
	"net/http"
	"strconv"

	"fangaoxs.com/go-chat/environment"
	"fangaoxs.com/go-chat/internal/auth"
	"fangaoxs.com/go-chat/internal/domain/group"
	"fangaoxs.com/go-chat/internal/domain/hub"
	"fangaoxs.com/go-chat/internal/domain/record"
	"fangaoxs.com/go-chat/internal/domain/user"
	"fangaoxs.com/go-chat/internal/entity"
	"fangaoxs.com/go-chat/internal/infras/errors"
	"fangaoxs.com/go-chat/internal/infras/logger"

	"github.com/gin-gonic/gin"
)

func newHandlers(
	env environment.Env,
	logger logger.Logger,
	user user.User,
	group group.Group,
	hub hub.Hub,
	record record.Record,
) (handlers, error) {
	return handlers{
		logger: logger,
		user:   user,
		group:  group,
		hub:    hub,
		record: record,
	}, nil
}

type handlers struct {
	logger logger.Logger

	user   user.User
	group  group.Group
	hub    hub.Hub
	record record.Record
}

func (h *handlers) RegisterUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		// POST
		ctx := c.Request.Context()

		nickname, ok := c.GetPostForm("nickname")
		if !ok {
			WrapGinError(c, errors.Newf(errors.InvalidArgument, nil, "invalid nickname"))
			return
		}
		username, ok := c.GetPostForm("username")
		if !ok {
			WrapGinError(c, errors.Newf(errors.InvalidArgument, nil, "invalid username"))
			return
		}
		password, ok := c.GetPostForm("password")
		if !ok {
			WrapGinError(c, errors.Newf(errors.InvalidArgument, nil, "invalid password"))
			return
		}
		phone, ok := c.GetPostForm("phone")
		if !ok {
			WrapGinError(c, errors.Newf(errors.InvalidArgument, nil, "invalid phone"))
			return
		}

		input := user.RegisterInput{
			Nickname: nickname,
			Username: username,
			Password: password,
			Phone:    phone,
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

func (h *handlers) Me() gin.HandlerFunc {
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

func (h *handlers) MyFriends() gin.HandlerFunc {
	return func(c *gin.Context) {
		// GET
		ctx := c.Request.Context()
		ui := auth.FromContext(ctx)

		friends, err := h.user.ListFriendsOfUser(ctx, ui.Subject)
		if err != nil {
			WrapGinError(c, err)
			return
		}

		c.JSON(http.StatusOK, friends)
	}
}

func (h *handlers) AssignFriends() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		ui := auth.FromContext(ctx)

		friendSubjects, ok := c.GetPostFormArray("friend_subject")
		if !ok {
			WrapGinError(c, errors.New(errors.InvalidArgument, nil, "empty friend subjects"))
			return
		}

		if err := h.user.AssignFriendsToUser(ctx, ui.Subject, friendSubjects...); err != nil {
			WrapGinError(c, err)
			return
		}

		c.Status(http.StatusOK)
	}
}

func (h *handlers) RemoveFriends() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		ui := auth.FromContext(ctx)

		friendSubjects, ok := c.GetPostFormArray("friend_subject")
		if !ok {
			WrapGinError(c, errors.New(errors.InvalidArgument, nil, "empty friend subjects"))
			return
		}

		if err := h.user.RemoveFriendsFromUser(ctx, ui.Subject, friendSubjects...); err != nil {
			WrapGinError(c, err)
			return
		}

		c.Status(http.StatusOK)
	}
}

// group

func (h *handlers) MyGroups() gin.HandlerFunc {
	return func(c *gin.Context) {
		// GET
		ctx := c.Request.Context()
		ui := auth.FromContext(ctx)

		groups, err := h.group.ListGroupsOfUser(ctx, ui.Subject)
		if len(groups) == 0 {
			WrapGinError(c, errors.New(errors.NotFound, nil, "没有群组"))
			return
		}
		if err != nil {
			WrapGinError(c, err)
			return
		}

		c.JSON(http.StatusOK, groups)
	}
}

func (h *handlers) InsertGroup() gin.HandlerFunc {
	return func(c *gin.Context) {
		// POST
		ctx := c.Request.Context()
		ui := auth.FromContext(ctx)

		name, ok := c.GetPostForm("name")
		if !ok {
			WrapGinError(c, errors.Newf(errors.InvalidArgument, nil, "invalid name"))
			return
		}

		groupType, ok := entity.GroupTypeFromString(c.PostForm("type"))
		if !ok {
			groupType = entity.DefaultGroupType
		}

		input := group.CreateGroupInput{
			Name:      name,
			Type:      groupType,
			CreatedBy: ui.Subject,
		}
		id, err := h.group.CreateGroup(ctx, input)
		if err != nil {
			WrapGinError(c, err)
			return
		}

		c.JSON(http.StatusOK, gin.H{"id": id})
	}
}

func (h *handlers) GetGroupByID() gin.HandlerFunc {
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

func (h *handlers) DeleteGroup() gin.HandlerFunc {
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

func (h *handlers) PrivateGroup() gin.HandlerFunc {
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

func (h *handlers) PublicGroup() gin.HandlerFunc {
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

func (h *handlers) AssignMembersToGroup() gin.HandlerFunc {
	return func(c *gin.Context) {
		// PUT
		ctx := c.Request.Context()
		ui := auth.FromContext(ctx)

		groupID, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			WrapGinError(c, err)
			return
		}

		subjects, ok := c.GetPostFormArray("user_subject")
		if !ok {
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

		if err = h.group.AssignMembersToGroup(ctx, g.ID, subjects...); err != nil {
			WrapGinError(c, err)
			return
		}

		c.Status(http.StatusOK)
	}
}

func (h *handlers) MembersOfGroup() gin.HandlerFunc {
	return func(c *gin.Context) {
		// GEt
		ctx := c.Request.Context()
		ui := auth.FromContext(ctx)

		groupID, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			WrapGinError(c, err)
			return
		}

		g, err := h.group.GetGroupByID(ctx, groupID)
		if err != nil {
			WrapGinError(c, err)
			return
		}

		members, err := h.group.ListMembersOfGroup(ctx, g.ID)
		if len(members) == 0 {
			WrapGinError(c, errors.New(errors.NotFound, nil, "没有群成员"))
			return
		}
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

// record

func (h *handlers) BroadcastMessage() gin.HandlerFunc {
	return func(c *gin.Context) {
		// POST
		message, ok := c.GetPostForm("message")
		if !ok {
			WrapGinError(c, errors.New(errors.InvalidArgument, nil, "empty message"))
			return
		}

		ctx := c.Request.Context()
		ui := auth.FromContext(ctx)

		err := h.hub.SendBroadcastMessage(ctx, ui.Subject, message)
		if err != nil {
			WrapGinError(c, err)
			return
		}
	}
}

func (h *handlers) GroupMessage() gin.HandlerFunc {
	return func(c *gin.Context) {
		// POST
		message, ok := c.GetPostForm("message")
		if !ok {
			WrapGinError(c, errors.New(errors.InvalidArgument, nil, "empty message"))
			return
		}
		groupID, err := strconv.ParseInt(c.Param("group_id"), 10, 64)
		if err != nil {
			WrapGinError(c, err)
			return
		}

		ctx := c.Request.Context()
		ui := auth.FromContext(ctx)

		err = h.hub.SendGroupMessage(ctx, ui.Subject, message, groupID)
		if err != nil {
			WrapGinError(c, err)
			return
		}
	}
}

func (h *handlers) PrivateMessage() gin.HandlerFunc {
	return func(c *gin.Context) {
		// POST
		message, ok := c.GetPostForm("message")
		if !ok {
			WrapGinError(c, errors.New(errors.InvalidArgument, nil, "empty message"))
			return
		}
		receiver, ok := c.GetPostForm("receiver")
		if !ok {
			WrapGinError(c, errors.New(errors.InvalidArgument, nil, "empty to"))
			return
		}

		ctx := c.Request.Context()
		ui := auth.FromContext(ctx)

		err := h.hub.SendPrivateMessage(ctx, ui.Subject, message, receiver)
		if err != nil {
			WrapGinError(c, err)
			return
		}
	}
}

func (h *handlers) RecordBroadcast() gin.HandlerFunc {
	return func(c *gin.Context) {
		// GET
		sender := c.Query("sender") // sender如果为空则查询所有

		ctx := c.Request.Context()
		res, err := h.record.ListRecordBroadcasts(ctx, sender)
		if err != nil {
			WrapGinError(c, err)
			return
		}

		c.JSON(http.StatusOK, res)
	}
}

func (h *handlers) RecordGroup() gin.HandlerFunc {
	return func(c *gin.Context) {
		// GET
		groupID, err := strconv.ParseInt(c.Param("group_id"), 10, 64)
		if err != nil {
			WrapGinError(c, errors.New(errors.InvalidArgument, err, "invalid group_id"))
			return
		}

		ctx := c.Request.Context()
		ui := auth.FromContext(ctx)

		ok, err := h.group.IsMemberOfGroup(ctx, groupID, ui.Subject)
		if err != nil {
			WrapGinError(c, err)
			return
		}
		if !ok {
			WrapGinError(c, errors.New(errors.PermissionDenied, nil, "你无法查看该群组"))
			return
		}

		res, err := h.record.ListRecordGroups(ctx, groupID)
		if err != nil {
			WrapGinError(c, err)
			return
		}

		c.JSON(http.StatusOK, res)
	}
}

func (h *handlers) RecordPrivate() gin.HandlerFunc {
	return func(c *gin.Context) {
		// GET
		ctx := c.Request.Context()
		ui := auth.FromContext(ctx)

		receiver := c.Param("receiver")
		ok, err := h.user.IsFriendOfUser(ctx, ui.Subject, receiver)
		if err != nil {
			WrapGinError(c, err)
			return
		}
		if !ok {
			WrapGinError(c, errors.Newf(errors.PermissionDenied, nil, "[%s]不是[%s]的朋友", receiver, ui.Subject))
			return
		}

		res, err := h.record.ListRecordPrivate(ctx, ui.Subject, receiver)
		if err != nil {
			WrapGinError(c, err)
			return
		}

		c.JSON(http.StatusOK, res)
	}
}
