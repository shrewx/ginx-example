package global

import (
	"github.com/shrewx/ginx/pkg/conf"
	"github.com/shrewx/ginx/pkg/dbhelper"
	"github.com/shrewx/ginx/pkg/logx"
	"gorm.io/gorm"
)

// Configuration 全局配置
type Configuration struct {
	// 服务配置
	Server conf.Server `yaml:"server"`
	// 数据库配置
	DB conf.DB `yaml:"db"`
	// 日志配置
	Log conf.Log `yaml:"log"`
}

var (
	DB     *gorm.DB
	Config = &Configuration{}
)

func Load() {
	logx.Load(&Config.Log)
	dbLoad(Config.DB)
}

func dbLoad(conf conf.DB) {
	db, err := dbhelper.NewDB(conf)
	if err != nil {
		panic(err)
	}

	DB = db.DB
}
