package postgres

import (
	"database/sql"
	"ecoswap/config"
	"fmt"
	_"github.com/lib/pq"
)


func ConnectDB()(*sql.DB, error){
	cfg := config.Load()
	connector := fmt.Sprintf("host = %s port = %s user = %s dbname = %s password = %s sslmode = disable", 
							cfg.DB_HOST, cfg.DB_PORT, cfg.DB_USER, cfg.DB_NAME, cfg.DB_PASSWORD)
	return sql.Open("postgres", connector)
}