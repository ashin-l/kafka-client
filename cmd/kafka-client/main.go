package main

import (
	"context"
	"fmt"
	"kafka-client/pkg/config"
	"kafka-client/pkg/kafka"
	"kafka-client/pkg/logger"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
)

func main() {
	// 加载配置
	cfg, err := config.Init("config.yaml")
	if err != nil {
		panic(fmt.Sprintf("无法加载配置: %v", err))
	}

	// 初始化日志记录器
	if err := logger.InitLogger(cfg); err != nil {
		panic(fmt.Sprintf("init logger failed, err: %v", err))
	}
	defer zap.L().Sync()

	// 初始化并启动Kafka消费者
	// 使用带取消功能的Context来管理消费者生命周期
	ctx, cancel := context.WithCancel(context.Background())
	kafkaHandler := kafka.NewConsumerGroupHandler()
	kafka.StartConsumerGroup(ctx, cfg, kafkaHandler)

	// 设置优雅关闭
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	// 等待中断信号
	<-interrupt
	cancel()
	logger.Info("收到中断信号，正在关闭...")
}
