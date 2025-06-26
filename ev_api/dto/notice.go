package dto

import (
	"errors"
	"time"
)

type NoticeLevel = string

type NoticeBtnJumpType = string

const (
	NoticeBtnJumpTypeInternal NoticeBtnJumpType = "internal"
	NoticeBtnJumpTypeRemote   NoticeBtnJumpType = "remote"
	NoticeBtnJumpTypeReload   NoticeBtnJumpType = "reload"
)

const (
	NoticeLevelPrimary NoticeLevel = "primary"
	NoticeLevelSuccess NoticeLevel = "success"
	NoticeLevelInfo    NoticeLevel = "info"
	NoticeLevelWarning NoticeLevel = "warning"
	NoticeLevelDanger  NoticeLevel = "danger"
)

type LiveBroadcastEvMsg2RolesReq struct {
	NoticeData *NoticeData `json:"notice_data"`
	RoleIds    []int       `json:"role_ids"`
}

type LiveBroadcastEvMsg2UsersReq struct {
	NoticeData *NoticeData `json:"notice_data"`
	UserIds    []int       `json:"user_ids"`
}

type LiveBroadcastEvMsg2AllReq struct {
	NoticeData *NoticeData `json:"notice_data"`
}

type NoticeData struct {
	Title         string         `json:"title"`           //标题
	Content       string         `json:"content"`         //内容
	Type          string         `json:"type"`            //通知类型（如 system, alert, announcement）
	Level         NoticeLevel    `json:"level"`           //消息严重程度 info, warn, danger,success,primary
	IsTask        bool           `json:"is_task"`         //是否定时任务
	FromUid       int            `json:"from_uid"`        //用户id
	PluginAlias   string         `json:"plugin_alias"`    //插件id
	Source        string         `json:"source"`          //来源
	NoticeJumpBtn *NoticeJumpBtn `json:"notice_jump_btn,omitempty"` //跳转按钮
	PublishTime   time.Time      `json:"publish_time"`    //发布时间
}

func (this *NoticeData) Validate() error {
	if this == nil {
		return errors.New("消息不能为空")
	}
	if this.Level == "" {
		return errors.New("消息等级不能为空")
	}
	if this.Title == "" {
		return errors.New("消息标题不能为空")
	}
	if this.Type == "" {
		return errors.New("消息类型不能为空")
	}
	return nil
}

type NoticeJumpBtn struct {
	Text     string            `json:"text"`      //按钮文案
	JumpUrl  string            `json:"jump_url"`  //跳转链接
	JumpType NoticeBtnJumpType `json:"jump_type"` // 跳转类型: "internal" 或 "external"
}
