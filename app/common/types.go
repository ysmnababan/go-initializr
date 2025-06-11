package common

type DBType string
type FrameWorkType string

const (
	DB_POSTGRES  DBType = "postgres"
	DB_MYSQL     DBType = "mysql"
	DB_SQLSERVER DBType = "sql_server"

	FRAMEWORK_ECHO FrameWorkType = "echo"
	FRAMEWORK_GIN  FrameWorkType = "gin"
)
