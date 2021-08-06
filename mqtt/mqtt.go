// Package mqtt
// 参考 https://next.api.aliyun.com/api/OnsMqtt/2020-04-20/RegisterDeviceCredential?params={%22ClientId%22:%22ClientId%22,%22InstanceId%22:%22InstanceId%22}
package mqtt

import (
	"basic/cipher"
	"basic/color"
	"encoding/base64"
	"fmt"
	openapi "github.com/alibabacloud-go/darabonba-openapi/client"
	onsmqtt20200420 "github.com/alibabacloud-go/onsmqtt-20200420/client"
	"github.com/alibabacloud-go/tea/tea"
	"log"
	"strings"
	"time"
)

var Mqtt *server
var mqttServer *onsmqtt20200420.Client

type auth struct {
	ClientId   string `json:"client_id"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	ExpireTime int64  `json:"expire_time"`
}

type server struct {
	instanceId  *string
	accessKeyId string
}

type CommonDeviceCredential struct {
	RequestId             *string
	UpdateTime            *int64
	DeviceAccessKeyId     *string
	CreateTime            *int64
	InstanceId            *string
	DeviceAccessKeySecret *string
	ClientId              *string
}

// ClientsInfo 查询设备信息
func (s server) ClientsInfo(clientIdList ...string) {

	if len(clientIdList) < 1 {
		return
	}

	var _clientIdList = make([]*string, len(clientIdList))

	for i, clientId := range clientIdList {
		_clientIdList[i] = tea.String(clientId)
	}

	batchQuerySessionByClientIdsRequest := &onsmqtt20200420.BatchQuerySessionByClientIdsRequest{
		ClientIdList: _clientIdList,
		InstanceId:   s.instanceId,
	}
	// 复制代码运行请自行打印 API 的返回值
	result, err := mqttServer.BatchQuerySessionByClientIds(batchQuerySessionByClientIdsRequest)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(result)
}

// Send 发送消息
func (s server) Send(topic, payload string) {
	sendMessageRequest := &onsmqtt20200420.SendMessageRequest{
		MqttTopic:  tea.String(topic),
		InstanceId: s.instanceId,
		Payload:    tea.String(payload),
	} // 复制代码运行请自行打印 API 的返回值
	result, err := mqttServer.SendMessage(sendMessageRequest)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(result)
}

// Token 创建临时访问 Token
func (s server) Token(expireTime int64, resources, actions string) (body *onsmqtt20200420.ApplyTokenResponseBody, err error) {
	applyTokenRequest := &onsmqtt20200420.ApplyTokenRequest{
		InstanceId: s.instanceId,
		ExpireTime: tea.Int64(expireTime), //Token失效那天的毫秒时间戳多是30天
		Resources:  tea.String(resources), //MQTT的Topic，多个Topic以逗号（,），支持配符语法
		Actions:    tea.String(actions),   //"R","W","R,W"
	}

	result, err := mqttServer.ApplyToken(applyTokenRequest)
	if err != nil {
		log.Println(err)
		return
	}
	return result.Body, nil
}

// TokenQuery 校验 Token 的有效性
func (s server) TokenQuery(token string) {
	queryTokenRequest := &onsmqtt20200420.QueryTokenRequest{
		Token:      tea.String(token),
		InstanceId: s.instanceId,
	}
	result, err := mqttServer.QueryToken(queryTokenRequest)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(result)
}

// TokenRevoke 吊销 Token
func (s server) TokenRevoke(token string) {
	revokeTokenRequest := &onsmqtt20200420.RevokeTokenRequest{
		Token:      tea.String(token),
		InstanceId: s.instanceId,
	}

	result, err := mqttServer.RevokeToken(revokeTokenRequest)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(result)
}

//拼接ClientId
func ClientId(groupId, deviceId string) *string {
	return tea.String(fmt.Sprintf(`%s@@@%s`, groupId, deviceId))
}

// Register 为某个设备注册访问凭证，一般用于服务端
func (s server) Register(groupId, deviceId string) (commonDeviceCredential *CommonDeviceCredential, err error) {
	//查询是否存在
	commonDeviceCredential, err = s.GetDeviceCredential(groupId, deviceId)
	if err != nil {
		//没有查到，重新注册
		registerDeviceCredentialRequest := &onsmqtt20200420.RegisterDeviceCredentialRequest{
			ClientId:   ClientId(groupId, deviceId),
			InstanceId: s.instanceId,
		}
		result := new(onsmqtt20200420.RegisterDeviceCredentialResponse)
		result, err = mqttServer.RegisterDeviceCredential(registerDeviceCredentialRequest)
		if err != nil {
			log.Println(err)
			return
		}

		//重新组装结果
		commonDeviceCredential.RequestId = result.Body.RequestId
		commonDeviceCredential.UpdateTime = result.Body.DeviceCredential.UpdateTime
		commonDeviceCredential.DeviceAccessKeyId = result.Body.DeviceCredential.DeviceAccessKeyId
		commonDeviceCredential.CreateTime = result.Body.DeviceCredential.UpdateTime
		commonDeviceCredential.InstanceId = result.Body.DeviceCredential.InstanceId
		commonDeviceCredential.DeviceAccessKeySecret = result.Body.DeviceCredential.DeviceAccessKeySecret
		commonDeviceCredential.ClientId = result.Body.DeviceCredential.ClientId
		return
	}
	return
}

//UnRegister 注销设备的访问凭证
func (s server) UnRegister(groupId, deviceId string) (err error) {
	unRegisterDeviceCredentialRequest := &onsmqtt20200420.UnRegisterDeviceCredentialRequest{
		ClientId:   ClientId(groupId, deviceId),
		InstanceId: s.instanceId,
	}
	_, err = mqttServer.UnRegisterDeviceCredential(unRegisterDeviceCredentialRequest)
	if err != nil {
		log.Println(err)
		return
	}
	return
}

//RefreshDeviceCredential 更新设备的访问凭证
func (s server) RefreshDeviceCredential(groupId, deviceId string) (*CommonDeviceCredential, error) {
	refreshDeviceCredentialRequest := &onsmqtt20200420.RefreshDeviceCredentialRequest{
		ClientId:   ClientId(groupId, deviceId),
		InstanceId: s.instanceId,
	}
	result, err := mqttServer.RefreshDeviceCredential(refreshDeviceCredentialRequest)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &CommonDeviceCredential{
		RequestId:             result.Body.RequestId,
		UpdateTime:            result.Body.DeviceCredential.UpdateTime,
		DeviceAccessKeyId:     result.Body.DeviceCredential.DeviceAccessKeyId,
		CreateTime:            result.Body.DeviceCredential.CreateTime,
		InstanceId:            result.Body.DeviceCredential.InstanceId,
		DeviceAccessKeySecret: result.Body.DeviceCredential.DeviceAccessKeySecret,
		ClientId:              result.Body.DeviceCredential.ClientId,
	}, nil
}

//GetDeviceCredential 查询设备的访问凭证
func (s server) GetDeviceCredential(groupId, deviceId string) (*CommonDeviceCredential, error) {
	getDeviceCredentialRequest := &onsmqtt20200420.GetDeviceCredentialRequest{
		ClientId:   ClientId(groupId, deviceId),
		InstanceId: s.instanceId,
	}
	result, err := mqttServer.GetDeviceCredential(getDeviceCredentialRequest)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &CommonDeviceCredential{
		RequestId:             result.Body.RequestId,
		UpdateTime:            result.Body.DeviceCredential.UpdateTime,
		DeviceAccessKeyId:     result.Body.DeviceCredential.DeviceAccessKeyId,
		CreateTime:            result.Body.DeviceCredential.CreateTime,
		InstanceId:            result.Body.DeviceCredential.InstanceId,
		DeviceAccessKeySecret: result.Body.DeviceCredential.DeviceAccessKeySecret,
		ClientId:              result.Body.DeviceCredential.ClientId,
	}, nil
}

// AuthToken 参考 https://help.aliyun.com/document_detail/54226.html?spm=a2c4g.11186623.6.594.34331ac3VAk6PV
// 用于分发给客户端，有读写权限和时效
//func (s server) AuthToken(groupId, deviceId, resources, actions string) (_auth *auth, err error) {
func (s server) AuthToken(groupId, deviceId string, resourcesActions map[string]string) (_auth *auth, err error) {
	//resourcesActions
	//{resources:actions,resources:actions}
	ts := time.Now().AddDate(0, 1, 0).UnixNano() / 1e6

	var password = ""
	var token *onsmqtt20200420.ApplyTokenResponseBody

	for resources, actions := range resourcesActions {
		token, err = s.Token(ts, resources, actions)
		if err != nil {
			log.Println(err)
			return
		}
		password += fmt.Sprintf("%s|%s|", strings.Replace(actions, ",", "", -1), *token.Token)
	}

	//生成账号密码
	return &auth{
		ClientId:   *ClientId(groupId, deviceId),
		Username:   fmt.Sprintf("Token|%s|%s", s.accessKeyId, *s.instanceId),
		Password:   strings.TrimRight(password, "|"),
		ExpireTime: ts,
	}, nil
}

// AuthRegister 用于服务端，没有限制 https://help.aliyun.com/document_detail/189437.htm?spm=a2c4g.11186623.2.6.31d03bd5kflDAi#concept-1995234
func (s server) AuthRegister(groupId, deviceId string) (_auth *auth, err error) {
	register, err := s.Register(groupId, deviceId)
	if err != nil {
		return nil, nil
	}

	return &auth{
		ClientId: *register.ClientId,
		Username: fmt.Sprintf("DeviceCredential|%s|%s", *register.DeviceAccessKeyId, *s.instanceId),
		Password: base64.StdEncoding.EncodeToString(
			cipher.HmacSha1Byte([]byte(*register.ClientId), []byte(*register.DeviceAccessKeySecret)),
		),
		ExpireTime: 0,
	}, nil
}

type Server struct {
	Endpoint        string
	InstanceId      string
	AccessKeyId     string
	AccessKeySecret string
}

func (s Server) CreateServer() {
	if Mqtt != nil {
		return
	}
	var err error
	config := &openapi.Config{
		// 您的AccessKey ID
		AccessKeyId: tea.String(s.AccessKeyId),
		// 您的AccessKey Secret
		AccessKeySecret: tea.String(s.AccessKeySecret),
	}
	// 访问的域名
	config.Endpoint = tea.String(s.Endpoint)
	//mqttServer = &onsmqtt20200420.Client{}
	mqttServer, err = onsmqtt20200420.NewClient(config)
	if err != nil {
		log.Fatalln(color.Red, err, color.Reset)
	}

	Mqtt = &server{
		instanceId:  tea.String(s.InstanceId),
		accessKeyId: s.AccessKeyId,
	}
	color.Success("[mqtt] create client success")
}
