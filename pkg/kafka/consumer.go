package kafka

import (
	"context"
	"kafka-client/pkg/config"
	"kafka-client/pkg/logger"

	"github.com/IBM/sarama"
	"github.com/tidwall/gjson"
	"go.uber.org/zap"
)

// ConsumerGroupHandler 实现了sarama的消费者组处理器接口
type ConsumerGroupHandler struct {
}

func NewConsumerGroupHandler() *ConsumerGroupHandler {
	return &ConsumerGroupHandler{}
}

func (h *ConsumerGroupHandler) Setup(sarama.ConsumerGroupSession) error   { return nil }
func (h *ConsumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error { return nil }
func (h *ConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		logger.Debug("从Kafka接收到消息", zap.String("topic", message.Topic), zap.Int32("partition", message.Partition), zap.Int64("offset", message.Offset), zap.Int("value", len(message.Value)))
		session.MarkMessage(message, "")
		deviceId := gjson.GetBytes(message.Value, "deviceId").String()
		deviceName := gjson.GetBytes(message.Value, "deviceName").String()
		uuid := gjson.GetBytes(message.Value, "uuid").String()
		oid := gjson.GetBytes(message.Value, "oid").String()
		icaoAddress := gjson.GetBytes(message.Value, "icaoAddress").String()
		logger.Info("========= 收到消息", zap.String("topic", message.Topic), zap.String("deviceId", deviceId), zap.String("deviceName", deviceName), zap.String("uuid", uuid), zap.String("oid", oid), zap.String("icaoAddress", icaoAddress))
	}
	return nil
}

// StartConsumerGroup 启动Kafka消费者组
func StartConsumerGroup(ctx context.Context, cfg *config.AppConfig, handler *ConsumerGroupHandler) {
	saramaCfg := sarama.NewConfig()
	saramaCfg.Version = sarama.V2_8_1_0 // 选择一个合适的Kafka版本
	saramaCfg.Consumer.Offsets.Initial = sarama.OffsetOldest
	saramaCfg.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategyRoundRobin()

	consumerGroup, err := sarama.NewConsumerGroup(cfg.KafkaConfig.Brokers, cfg.KafkaConfig.GroupID, saramaCfg)
	if err != nil {
		logger.Fatal("创建Kafka消费者组失败", zap.Error(err))
	}

	go func() {
		defer consumerGroup.Close()
		for {
			if err := consumerGroup.Consume(ctx, cfg.KafkaConfig.Topics, handler); err != nil {
				logger.Error("消费者组消费失败", zap.Error(err))
			}
			if ctx.Err() != nil {
				return
			}
		}
	}()
	logger.Info("Kafka消费者组已启动")
}
