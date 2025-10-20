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
	conf.ShowLog = true
	db, err := dbhelper.NewDB(conf)
	if err != nil {
		panic(err)
	}
	// 表数据自动迁移
	if err := dbhelper.Migrate(db); err != nil {
		panic(err)
	}

	DB = db.DB
	
	// 执行用户表迁移和种子数据
	if err := migrateUserTables(); err != nil {
		panic(err)
	}
}

// migrateUserTables 迁移用户相关表
func migrateUserTables() error {
	// 导入models包以避免循环依赖
	// 这里我们直接在这个函数中处理迁移
	return nil // 实际迁移将在main函数中处理
}
