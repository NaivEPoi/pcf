package notifyevent

import (
	"net/http"

	"github.com/free5gc/openapi"
	"github.com/free5gc/openapi/models"
	pcf_context "github.com/free5gc/pcf/internal/context"
	"github.com/free5gc/pcf/internal/logger"
	"github.com/free5gc/pcf/internal/util"
	"github.com/tim-ywliu/event"
)

const SendSMpolicyUpdateNotifyEventName event.Name = "SendSMpolicyUpdateNotify"

type SendSMpolicyUpdateNotifyEvent struct {
	uri     string
	request *models.SmPolicyNotification
}

func (e SendSMpolicyUpdateNotifyEvent) Handle() {
	logger.NotifyEventLog.Infof("Handle SendSMpolicyUpdateNotifyEvent\n")
	if e.uri == "" {
		logger.NotifyEventLog.Warnln("SM Policy Update Notification Error[URI is empty]")
		return
	}
	client := util.GetNpcfSMPolicyCallbackClient()
	logger.NotifyEventLog.Infof("Send SM Policy Update Notification to SMF")
	_, httpResponse, err :=
		client.DefaultCallbackApi.SmPolicyUpdateNotification(openapi.CreateContext(pcf_context.PCF_Self().OAuth, pcf_context.PCF_Self().NfId, pcf_context.PCF_Self().NrfUri, "PCF"), e.uri, *e.request)
	if err != nil {
		if httpResponse != nil {
			logger.NotifyEventLog.Warnf("SM Policy Update Notification Error[%s]", httpResponse.Status)
		} else {
			logger.NotifyEventLog.Warnf("SM Policy Update Notification Failed[%s]", err.Error())
		}
		return
	} else if httpResponse == nil {
		logger.NotifyEventLog.Warnln("SM Policy Update Notification Failed[HTTP Response is nil]")
		return
	}
	defer func() {
		if resCloseErr := httpResponse.Body.Close(); resCloseErr != nil {
			logger.NotifyEventLog.Errorf("NFInstancesStoreApi response body cannot close: %+v", resCloseErr)
		}
	}()
	if httpResponse.StatusCode != http.StatusOK && httpResponse.StatusCode != http.StatusNoContent {
		logger.NotifyEventLog.Warnf("SM Policy Update Notification Failed")
	} else {
		logger.NotifyEventLog.Tracef("SM Policy Update Notification Success")
	}
}
