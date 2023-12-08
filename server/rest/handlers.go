package rest

import (
	"net/http"
	"strconv"
	"strings"

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

		nickname := strings.TrimSpace(c.PostForm("nickname"))
		if nickname == "" {
			WrapGinError(c, errors.Newf(errors.InvalidArgument, nil, "invalid username"))
			return
		}
		username := strings.TrimSpace(c.PostForm("username"))
		if username == "" {
			WrapGinError(c, errors.Newf(errors.InvalidArgument, nil, "invalid username"))
			return
		}
		password := strings.TrimSpace(c.PostForm("password"))
		if password == "" {
			WrapGinError(c, errors.Newf(errors.InvalidArgument, nil, "invalid password"))
			return
		}
		phone := strings.TrimSpace(c.PostForm("phone"))
		if phone == "" {
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

func (h *handlers) MyGroups() gin.HandlerFunc {
	return func(c *gin.Context) {
		// GET
		ctx := c.Request.Context()
		ui := auth.FromContext(ctx)

		groups, err := h.group.ListGroupsOfUser(ctx, ui.Subject)
		if err != nil {
			WrapGinError(c, err)
			return
		}

		c.JSON(http.StatusOK, groups)
	}
}

func (h *handlers) AssignFriends() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		ui := auth.FromContext(ctx)

		subjects := c.PostFormArray("friend_subject")
		if len(subjects) == 0 {
			WrapGinError(c, errors.New(errors.InvalidArgument, nil, "empty friend subjects"))
			return
		}
		for _, subject := range subjects {
			if subject = strings.TrimSpace(subject); subject == "" {
				WrapGinError(c, errors.New(errors.InvalidArgument, nil, "empty friend subject"))
				return
			}
		}

		if err := h.user.AssignFriendsToUser(ctx, ui.Subject, subjects...); err != nil {
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

		subjects := c.PostFormArray("friend_subject")
		if len(subjects) == 0 {
			WrapGinError(c, errors.New(errors.InvalidArgument, nil, "empty friend subjects"))
			return
		}
		for _, subject := range subjects {
			if subject = strings.TrimSpace(subject); subject == "" {
				WrapGinError(c, errors.New(errors.InvalidArgument, nil, "empty friend subject"))
				return
			}
		}

		if err := h.user.RemoveFriendsFromUser(ctx, ui.Subject, subjects...); err != nil {
			WrapGinError(c, err)
			return
		}

		c.Status(http.StatusOK)
	}
}

// group

func (h *handlers) CreateGroup() gin.HandlerFunc {
	return func(c *gin.Context) {
		// POST
		ctx := c.Request.Context()
		ui := auth.FromContext(ctx)

		name := strings.TrimSpace(c.PostForm("name"))
		if name == "" {
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
		// 当群为公开或者访问者是群成员的时候才可以查询

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
		isMember, err := h.group.IsMemberOfGroup(ctx, id, ui.Subject)
		if err != nil {
			WrapGinError(c, err)
			return
		}

		if !g.IsPublic && !isMember {
			WrapGinError(c, errors.New(errors.PermissionDenied, nil, "你不可以查看该群"))
			return
		}

		c.JSON(http.StatusOK, g)
	}
}

func (h *handlers) DeleteGroup() gin.HandlerFunc {
	return func(c *gin.Context) {
		// DELETE
		// 只有群创建者能删除群

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

func (h *handlers) MakeGroupPublic() gin.HandlerFunc {
	return func(c *gin.Context) {
		// PUT
		// 群管理员可以将群公开

		ctx := c.Request.Context()
		ui := auth.FromContext(ctx)

		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			WrapGinError(c, err)
			return
		}

		ok, err := h.group.IsAdminOfGroup(ctx, id, ui.Subject)
		if err != nil {
			WrapGinError(c, err)
			return
		}
		if !ok {
			WrapGinError(c, errors.New(errors.PermissionDenied, nil, "你不可以操作该群"))
			return
		}

		err = h.group.MakeGroupPublic(ctx, id)
		if err != nil {
			WrapGinError(c, err)
			return
		}

		c.Status(http.StatusOK)
	}
}

func (h *handlers) MakeGroupPrivate() gin.HandlerFunc {
	return func(c *gin.Context) {
		// PUT
		// 群管理员可以将群私有

		ctx := c.Request.Context()
		ui := auth.FromContext(ctx)

		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			WrapGinError(c, err)
			return
		}

		ok, err := h.group.IsAdminOfGroup(ctx, id, ui.Subject)
		if err != nil {
			WrapGinError(c, err)
			return
		}
		if !ok {
			WrapGinError(c, errors.New(errors.PermissionDenied, nil, "你不可以操作该群"))
			return
		}

		err = h.group.MakeGroupPrivate(ctx, id)
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
		// 群管理员可以添加成员

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
		for _, subject := range subjects {
			if subject = strings.TrimSpace(subject); subject == "" {
				WrapGinError(c, errors.Newf(errors.InvalidArgument, nil, "empty user subject"))
				return
			}
			if subject == ui.Subject {
				WrapGinError(c, errors.Newf(errors.InvalidArgument, nil, "无法操作自己"))
				return
			}
		}

		ok, err := h.group.IsAdminOfGroup(ctx, groupID, ui.Subject)
		if err != nil {
			WrapGinError(c, err)
			return
		}
		if !ok {
			WrapGinError(c, errors.New(errors.PermissionDenied, nil, "你不可以操作该群"))
			return
		}

		if err = h.group.AssignMembersToGroup(ctx, groupID, subjects...); err != nil {
			WrapGinError(c, err)
			return
		}

		c.Status(http.StatusOK)
	}
}

func (h *handlers) RemoveMembersFromGroup() gin.HandlerFunc {
	return func(c *gin.Context) {
		// PUT
		// 群管理员可以移除成员

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
		for _, subject := range subjects {
			if subject = strings.TrimSpace(subject); subject == "" {
				WrapGinError(c, errors.Newf(errors.InvalidArgument, nil, "empty user subject"))
				return
			}
			if subject == ui.Subject {
				WrapGinError(c, errors.Newf(errors.InvalidArgument, nil, "无法操作自己"))
				return
			}
		}

		ok, err := h.group.IsAdminOfGroup(ctx, groupID, ui.Subject)
		if err != nil {
			WrapGinError(c, err)
			return
		}
		if !ok {
			WrapGinError(c, errors.New(errors.PermissionDenied, nil, "你不可以操作该群"))
			return
		}

		if err = h.group.RemoveMembersFromGroup(ctx, groupID, subjects...); err != nil {
			WrapGinError(c, err)
			return
		}

		c.Status(http.StatusOK)
	}
}

func (h *handlers) MembersOfGroup() gin.HandlerFunc {
	return func(c *gin.Context) {
		// GET
		// 只有群成员可以查看该群

		ctx := c.Request.Context()
		ui := auth.FromContext(ctx)

		groupID, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			WrapGinError(c, err)
			return
		}

		ok, err := h.group.IsMemberOfGroup(ctx, groupID, ui.Subject)
		if err != nil {
			WrapGinError(c, err)
			return
		}
		if !ok {
			WrapGinError(c, errors.New(errors.PermissionDenied, nil, "你不可以查看该群"))
			return
		}

		members, err := h.group.ListMembersOfGroup(ctx, groupID)
		if err != nil {
			WrapGinError(c, err)
			return
		}
		if len(members) == 0 {
			WrapGinError(c, errors.New(errors.NotFound, nil, "没有群成员"))
			return
		}

		c.JSON(http.StatusOK, members)
	}
}

func (h *handlers) AssignAdminsToGroup() gin.HandlerFunc {
	return func(c *gin.Context) {
		// PUT
		// 只有群创建者可以添加管理员

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
		for _, subject := range subjects {
			if subject = strings.TrimSpace(subject); subject == "" {
				WrapGinError(c, errors.Newf(errors.InvalidArgument, nil, "empty user subject"))
				return
			}
			if subject == ui.Subject {
				WrapGinError(c, errors.Newf(errors.InvalidArgument, nil, "无法操作自己"))
				return
			}
		}

		g, err := h.group.GetGroupByID(ctx, groupID)
		if err != nil {
			WrapGinError(c, err)
			return
		}
		if g.CreatedBy != ui.Subject {
			WrapGinError(c, errors.New(errors.PermissionDenied, nil, "你不是群组创建者"))
			return
		}

		err = h.group.AssignAdminsToGroup(ctx, groupID, subjects...)
		if err != nil {
			WrapGinError(c, err)
			return
		}

		c.Status(http.StatusOK)
	}
}

func (h *handlers) RemoveAdminsFromGroup() gin.HandlerFunc {
	return func(c *gin.Context) {
		// PUT
		// 只有群创建者可以移除管理员

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
		for _, subject := range subjects {
			if subject = strings.TrimSpace(subject); subject == "" {
				WrapGinError(c, errors.Newf(errors.InvalidArgument, nil, "empty user subject"))
				return
			}
			if subject == ui.Subject {
				WrapGinError(c, errors.Newf(errors.InvalidArgument, nil, "无法操作自己"))
				return
			}
		}

		g, err := h.group.GetGroupByID(ctx, groupID)
		if err != nil {
			WrapGinError(c, err)
			return
		}
		if g.CreatedBy != ui.Subject {
			WrapGinError(c, errors.New(errors.PermissionDenied, nil, "你不是群组创建者"))
			return
		}

		err = h.group.RemoveAdminsFromGroup(ctx, groupID, subjects...)
		if err != nil {
			WrapGinError(c, err)
			return
		}

		c.Status(http.StatusOK)
	}
}

func (h *handlers) AdminsOfGroup() gin.HandlerFunc {
	return func(c *gin.Context) {
		// GET
		// 只有群成员可以访问群管理员

		ctx := c.Request.Context()
		ui := auth.FromContext(ctx)

		groupID, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			WrapGinError(c, err)
			return
		}

		ok, err := h.group.IsMemberOfGroup(ctx, groupID, ui.Subject)
		if err != nil {
			WrapGinError(c, err)
			return
		}
		if !ok {
			WrapGinError(c, errors.New(errors.PermissionDenied, nil, "你不可以查看该群"))
			return
		}

		admins, err := h.group.ListAdminsOfGroup(ctx, groupID)
		if err != nil {
			WrapGinError(c, err)
			return
		}
		if len(admins) == 0 {
			WrapGinError(c, errors.New(errors.NotFound, nil, "没有群管理员"))
			return
		}

		c.JSON(http.StatusOK, admins)
	}
}

// record

func (h *handlers) BroadcastMessage() gin.HandlerFunc {
	return func(c *gin.Context) {
		// POST
		message := strings.TrimSpace(c.PostForm("message"))
		if message == "" {
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

		c.Status(http.StatusOK)
	}
}

func (h *handlers) GroupMessage() gin.HandlerFunc {
	return func(c *gin.Context) {
		// POST
		message := strings.TrimSpace(c.PostForm("message"))
		if message == "" {
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

		ok, err := h.group.IsMemberOfGroup(ctx, groupID, ui.Subject)
		if err != nil {
			WrapGinError(c, err)
			return
		}
		if !ok {
			WrapGinError(c, errors.New(errors.PermissionDenied, nil, "你不是该群成员"))
			return
		}

		err = h.hub.SendGroupMessage(ctx, ui.Subject, message, groupID)
		if err != nil {
			WrapGinError(c, err)
			return
		}

		c.Status(http.StatusOK)
	}
}

func (h *handlers) PrivateMessage() gin.HandlerFunc {
	return func(c *gin.Context) {
		// POST
		message := strings.TrimSpace(c.PostForm("message"))
		if message == "" {
			WrapGinError(c, errors.New(errors.InvalidArgument, nil, "empty message"))
			return
		}
		receiver := strings.TrimSpace(c.PostForm("receiver"))
		if receiver == "" {
			WrapGinError(c, errors.New(errors.InvalidArgument, nil, "empty receiver"))
			return
		}

		ctx := c.Request.Context()
		ui := auth.FromContext(ctx)

		ok, err := h.user.IsFriendOfUser(ctx, ui.Subject, receiver)
		if err != nil {
			WrapGinError(c, err)
			return
		}
		if !ok {
			WrapGinError(c, errors.New(errors.PermissionDenied, nil, "对方不是你的好友"))
			return
		}

		err = h.hub.SendPrivateMessage(ctx, ui.Subject, message, receiver)
		if err != nil {
			WrapGinError(c, err)
			return
		}

		c.Status(http.StatusOK)
	}
}

func (h *handlers) GetRecordBroadcast() gin.HandlerFunc {
	return func(c *gin.Context) {
		// GET
		sender := c.Query("sender")

		ctx := c.Request.Context()
		var res []*entity.RecordBroadcast
		var err error
		if sender = strings.TrimSpace(sender); sender == "" {
			res, err = h.record.ListAllRecordBroadcasts(ctx)
			if err != nil {
				WrapGinError(c, err)
				return
			}
		} else {
			res, err = h.record.ListRecordBroadcastsBySender(ctx, sender)
			if err != nil {
				WrapGinError(c, err)
				return
			}
		}

		c.JSON(http.StatusOK, res)
	}
}

func (h *handlers) GetRecordGroup() gin.HandlerFunc {
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

func (h *handlers) GetRecordPrivate() gin.HandlerFunc {
	return func(c *gin.Context) {
		// GET
		receiver := c.Param("receiver")

		ctx := c.Request.Context()
		ui := auth.FromContext(ctx)

		res, err := h.record.ListRecordPrivate(ctx, ui.Subject, receiver)
		if err != nil {
			WrapGinError(c, err)
			return
		}

		c.JSON(http.StatusOK, res)
	}
}
