package global

import (
	"github.com/go-redis/redis/v8"
	"golang.org/x/sync/singleflight"

	"go.uber.org/zap"

	"project/config"

	"github.com/spf13/viper"
	"gorm.io/gorm"
)

var (
	GSD_DB                  *gorm.DB
	GSD_REDIS               *redis.Client
	GSD_CONFIG              config.Server
	GSD_VP                  *viper.Viper
	GSD_LOG                 *zap.Logger
	GSD_Concurrency_Control = &singleflight.Group{}
)
