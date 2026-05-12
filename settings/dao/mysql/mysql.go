package mysql

import (
	"bluebell/settings"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB

// Init 初始化MySQL连接
func Init(cfg *settings.MySQLConfig) (err error) {
	// "user:password@tcp(host:port)/dbname"
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&loc=Local", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DB)
	db, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		return
	}
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	return
}

// Close 关闭MySQL连接
func Close() {
	_ = db.Close()
}

//这个是用viper配置传参，上面更新后是用结构体传参
// func Init() (err error) {
// 	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True",
// 		viper.GetString("mysql.user"),
// 		viper.GetString("mysql.password"),
// 		viper.GetString("mysql.host"),
// 		viper.GetInt("mysql.port"),
// 		viper.GetString("mysql.dbname"),
// 	)
// 	db, err = sqlx.Connect("mysql", dsn)
// 	if err != nil {
// 		fmt.Printf("mysql sqlx connect fail err : %v", err)
// 		return
// 	}

// 	db.SetMaxOpenConns(viper.GetInt("max_open_connes"))
// 	db.SetMaxIdleConns(viper.GetInt("max_idle_connes"))

// 	return
// }
