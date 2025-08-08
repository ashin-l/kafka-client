package mqtt

import (
	"kafka-client/pkg/config"
	"kafka-client/pkg/logger"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"go.uber.org/zap"
)

// NewClient 创建并返回一个新的MQTT客户端实例
func NewClient(cfg *config.AppConfig, username string) mqtt.Client {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(cfg.MqttConfig.Broker)
	// opts.SetClientID(cfg.MqttConfig.ClientID)
	opts.SetUsername(username)
	// opts.SetPassword(cfg.MqttConfig.Password)
	opts.SetConnectTimeout(5 * time.Second)

	// --- 实现MQTT重连逻辑 ---
	opts.SetAutoReconnect(true)                   // 开启自动重连功能
	opts.SetMaxReconnectInterval(3 * time.Minute) // 设置最大重连间隔
	opts.SetConnectRetry(true)                    // 如果初始连接失败，则尝试重连

	// 设置连接成功和连接丢失的回调函数
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler

	logger.Info("正在连接MQTT服务器...")
	client := mqtt.NewClient(opts)
	// 使用异步连接，这样即使首次连接失败，程序也不会退出，
	// 而是由库在后台进行重连尝试。
	token := client.Connect()
	if token.WaitTimeout(5*time.Second) && token.Error() != nil {
		logger.Fatal("首次MQTT连接失败。将在后台自动重试...", zap.Error(token.Error()))
	}
	logger.Info("MQTT服务器连接成功", zap.String("username", username))
	return client
}

// connectHandler 会在客户端成功连接或重连后被调用
var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	logger.Info("已连接到MQTT Broker")
	// 如果有需要持久化的订阅，可以在这里重新订阅主题
	// client.Subscribe("some/topic", 1, nil)
}

// connectLostHandler 会在客户端连接丢失时被调用
var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	logger.Error("与MQTT Broker的连接丢失。正在尝试重连...", zap.Any("error", err))
}
