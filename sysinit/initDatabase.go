package sysinit

import (
	_ "ticket/models"
	// "time"
	"fmt"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

// InitDatabase -
func InitDatabase() {
	// read the configuration file and set the  database parameters.
	// database type
	fmt.Println("Start InitDatabase")
	dbType := beego.AppConfig.String("db_type")
	// database alias
	dbAlias := beego.AppConfig.String(dbType + "::db_alias")
	// database name
	dbName := beego.AppConfig.String(dbType + "::db_name")
	// user
	dbUser := beego.AppConfig.String(dbType + "::db_user")
	// password
	dbPwd := beego.AppConfig.String(dbType + "::db_pwd")
	// address
	dbHost := beego.AppConfig.String(dbType + "::db_host")
	beego.Debug(dbHost)
	// port
	dbPort := beego.AppConfig.String(dbType + "::db_port")
	switch dbType {
	case "sqlite3":
		orm.RegisterDataBase(dbAlias, dbType, dbName)
	case "mysql":
		dbCharset := beego.AppConfig.String(dbType + "::db_charset")
		orm.RegisterDataBase(dbAlias, dbType, dbUser+":"+dbPwd+"@tcp("+dbHost+":"+
			dbPort+")/"+dbName+"?charset="+dbCharset+"&loc=Asia%2FShanghai", 30)
	}
	// If it is in development mode, display command information
	isDev := (beego.AppConfig.String("runmode") == "dev")
	// Automatic table creation
	//orm.RunSyncdb("default", false, isDev)
	if isDev {
		orm.Debug = isDev
	}
	// orm.DefaultTimeLoc = time.UTC
}
