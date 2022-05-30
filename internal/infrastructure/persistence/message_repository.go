package persistence

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"git.nova.net.cn/nova/misc/wx-public/proxy/internal/pkg/httputil"

	"git.nova.net.cn/nova/misc/wx-public/proxy/internal/consts"
	errors3 "git.nova.net.cn/nova/misc/wx-public/proxy/internal/pkg/errorx"

	"gorm.io/gorm"

	"git.nova.net.cn/nova/misc/wx-public/proxy/internal/config"
	"git.nova.net.cn/nova/misc/wx-public/proxy/internal/domain/entity"
	"git.nova.net.cn/nova/misc/wx-public/proxy/internal/utils"

	errors2 "errors"

	log "github.com/sirupsen/logrus"
)

type MessageRepo struct {
	DB *gorm.DB
}

var defaultMessageRepo *MessageRepo

func NewMessageRepo() {
	if defaultMessageRepo == nil {
		defaultMessageRepo = &MessageRepo{
			DB: CommonRepositories.DB,
		}
	}
}

func DefaultMessageRepo() *MessageRepo {
	return defaultMessageRepo
}

func (a *MessageRepo) SendTmplMsgFromRequest(ctx context.Context, param entity.SendTmplMsgRemoteReq) (entity.SendTmplMsgRemoteResp, error) {
	traceID := utils.ShouldGetTraceID(ctx)
	log.Debugf("SendTmplMsgFromRequest traceID:%s", traceID)
	// 请求wx msg send
	bs, err := json.Marshal(param)
	if err != nil {
		log.Errorf("SendTmplMsgFromRequest json marshal send msg req failed,traceID:%s,err:%+v", traceID, err)
		return entity.SendTmplMsgRemoteResp{}, err
	}
	requestProperty := httputil.GetRequestProperty(http.MethodPost, config.WXMsgTmplSendURL+fmt.Sprintf("?access_token=%s", param.AccessToken),
		bs, make(map[string]string))
	statusCode, body, _, err := httputil.RequestWithContextAndRepeat(ctx, requestProperty, traceID)
	if err != nil {
		log.Errorf("SendTmplMsgFromRequest request wx msg send failed, traceID:%s, error:%+v", traceID, err)
		return entity.SendTmplMsgRemoteResp{}, err
	}
	if statusCode != http.StatusOK {
		log.Errorf("SendTmplMsgFromRequest request wx msg send failed, statusCode:%d,traceID:%s, error:%+v", statusCode, traceID, err)
		return entity.SendTmplMsgRemoteResp{}, err
	}
	var msgResp entity.SendTmplMsgRemoteResp
	err = json.Unmarshal(body, &msgResp)
	if err != nil {
		log.Errorf("SendTmplMsgFromRequest get wx msg send failed by unmarshal, resp:%s, traceID:%s, err:%+v", string(body), traceID, err)
		return entity.SendTmplMsgRemoteResp{}, err
	}
	// 获取失败
	if msgResp.ErrCode != errors3.CodeOK {
		log.Errorf("SendTmplMsgFromRequest get wx msg send failed,resp:%s,traceID:%s,errMsg:%s", string(body), traceID, msgResp.ErrMsg)
		return msgResp, fmt.Errorf("get wx msg send failed,errMsg:%s", msgResp.ErrMsg)
	}
	return msgResp, nil
}

func (a *MessageRepo) SaveMsgLog(ctx context.Context, param entity.MsgLog) error {
	traceID := utils.ShouldGetTraceID(ctx)
	log.Debugf("SaveMsgLog traceID:%s", traceID)
	if err := a.DB.Create(&param).Error; err != nil {
		log.Errorf("SaveMsgLog create failure msg log failed,traceID:%s,err:%+v", traceID, err)
		return err
	}
	return nil
}

func (a *MessageRepo) BatchSaveMsgLog(ctx context.Context, msgLogs []entity.MsgLog) error {
	traceID := utils.ShouldGetTraceID(ctx)
	log.Debugf("BatchSaveMsgLog traceID:%s", traceID)
	if err := a.DB.Create(&msgLogs).Error; err != nil {
		log.Errorf("BatchSaveMsgLog batch insert msg log failed,traceID:%s,err:%+v", traceID, err)
		return err
	}
	return nil
}

// UpdateMsgLog sup non-zero value
func (a *MessageRepo) UpdateMsgLog(ctx context.Context, msgLog entity.MsgLog) error {
	traceID := utils.ShouldGetTraceID(ctx)
	log.Debugf("UpdateMsgLog traceID:%s", traceID)
	err := a.DB.Model(&msgLog).Updates(msgLog).Error
	if err != nil {
		log.Errorf("UpdateMsgLog update failed,traceID:%s,err:%+v", traceID, err)
		return err
	}
	return nil
}

func (a *MessageRepo) UpdateMsgLogSendStatus(ctx context.Context, msgLog entity.MsgLog) error {
	traceID := utils.ShouldGetTraceID(ctx)
	log.Debugf("UpdateMsgLogStatus traceID:%s", traceID)
	err := a.DB.Model(&entity.MsgLog{}).Where("id = ?", msgLog.ID).Updates(map[string]interface{}{
		"cause":       msgLog.Cause,
		"status":      msgLog.Status,
		"count":       msgLog.Count,
		"update_time": msgLog.UpdateTime,
	}).Error
	if err != nil {
		log.Errorf("UpdateMsgLogStatus update failed,traceID:%s,err:%+v", traceID, err)
		return err
	}
	return nil
}

func (a *MessageRepo) GetMsgLogByMsgID(ctx context.Context, msgID int64) (entity.MsgLog, error) {
	traceID := utils.ShouldGetTraceID(ctx)
	log.Debugf("GetMsgLogByMsgID traceID:%s", traceID)
	var m entity.MsgLog
	err := a.DB.Where("msg_id = ?", msgID).First(&m).Error
	if err != nil {
		log.Errorf("GetMsgLogByMsgID failed to get msg log by msg id,msgID:%d,traceID:%s,err:%+v", msgID, traceID, err)
		return entity.MsgLog{}, err
	}
	return m, err
}

func (a *MessageRepo) IsExistMsgLogFromDB(ctx context.Context, fromUserName string, createTime int64) (bool, error) {
	traceID := utils.ShouldGetTraceID(ctx)
	log.Debugf("IsExistMsgLogFromDB traceID:%s", traceID)
	var failureMsgLog entity.MsgLog
	err := a.DB.Where("to_user = ? AND create_time = ?", fromUserName, createTime).First(&failureMsgLog).Error
	if err != nil {
		// 不存在记录
		if errors2.Is(err, gorm.ErrRecordNotFound) {
			log.Errorf("IsExistMsgLogFromDB record is not found,traceID:%s,err:%+v", traceID, err)
			return false, nil
		}
		log.Errorf("IsExistMsgLogFromDB failed,traceID:%s,err:%+v", traceID, err)
		return false, err
	}
	return true, nil
}

func (a *MessageRepo) ListPendingMsgLogs(ctx context.Context) ([]entity.MsgLog, error) {
	traceID := utils.ShouldGetTraceID(ctx)
	log.Debugf("GetPendingMsgLog traceID:%s", traceID)
	var msgLogs []entity.MsgLog
	err := a.DB.Where("status = ? AND count < ?", consts.SendPending, consts.MaxRetryCount).Find(&msgLogs).Error
	if err != nil {
		log.Errorf("GetListPendingMsgLogs get list pending msg logs failed,traceID:%s,err:%+v", traceID, err)
		return nil, err
	}
	return msgLogs, err
}

func (a *MessageRepo) UpdateMaxRetryCntMsgLogsStatus(ctx context.Context) error {
	traceID := utils.ShouldGetTraceID(ctx)
	log.Debugf("UpdateMaxRetryCntMsgLogsStatus traceID:%s", traceID)
	err := a.DB.Model(&entity.MsgLog{}).Where("status = ? AND count >= ?", consts.SendPending, consts.MaxRetryCount).Update("status", consts.SendFailure).Error
	if err != nil {
		log.Errorf("UpdateMaxRetryCntMsgLogsStatus update max retry msg logs failed,traceID:%s,err:%+v", traceID, err)
		return err
	}
	return nil
}

func (a *MessageRepo) UpdateTimeoutMsgLogsStatus(ctx context.Context) error {
	traceID := utils.ShouldGetTraceID(ctx)
	log.Debugf("UpdateTimeoutMsgLogsStatus traceID:%s", traceID)
	err := a.DB.Model(&entity.MsgLog{}).Where("status = ? AND count >= ? AND create_time <= ?", consts.Sending, consts.MaxRetryCount, time.Now().Unix()-consts.MaxWXCallBackTime).
		Update("status", consts.SendFailure).Error
	if err != nil {
		log.Errorf("UpdateTimeoutMsgLogsStatus update time out msg log failed,traceID:%s,err:%+v", traceID, err)
		return err
	}
	return nil
}

func (a *MessageRepo) getListMsgLogsByReqIDSession(ctx context.Context, requestID string) *gorm.DB {
	traceID := utils.ShouldGetTraceID(ctx)
	log.Debugf("GetListMsgLogsByRequestIDSession traceID:%s", traceID)
	return a.DB.Where("request_id = ?", requestID)
}

func (a *MessageRepo) ListMsgLogsByReqIDCnt(ctx context.Context, requestID string) (int64, error) {
	traceID := utils.ShouldGetTraceID(ctx)
	log.Debugf("ListMsgLogsByReqIDCnt traceID:%s", traceID)
	db := a.getListMsgLogsByReqIDSession(ctx, requestID)
	var count int64
	err := db.Model(&entity.MsgLog{}).Count(&count).Error
	if err != nil {
		log.Errorf("ListMsgLogsByReqID find list msg logs by request id failed,traceID:%s,err:%+v", traceID, err)
		return 0, err
	}
	return count, nil
}

func (a *MessageRepo) ListMsgLogsByReqID(ctx context.Context, requestID string) ([]entity.MsgLog, error) {
	traceID := utils.ShouldGetTraceID(ctx)
	log.Debugf("ListMsgLogsByReqID traceID:%s", traceID)
	var msgLogs []entity.MsgLog
	db := a.getListMsgLogsByReqIDSession(ctx, requestID)
	err := db.Find(&msgLogs).Error
	if err != nil {
		log.Errorf("ListMsgLogsByReqID find list msg logs by request id failed,traceID:%s,err:%+v", traceID, err)
		return nil, err
	}
	return msgLogs, nil
}
