package containers

import (
	"github.com/jmoiron/sqlx"
)

// DBProxy holds DB connection with proxy.
type DBProxy struct {
	DB       *sqlx.DB
	Upstream *Upstream
}

// Close close proxy related connections.
func (p *DBProxy) Close() {
	_ = p.DB.Close()
}
