package containers

import (
	toxiproxy "github.com/Shopify/toxiproxy/v2/client"
	"github.com/jmoiron/sqlx"
)

// DBProxy holds DB connection with proxy.
type DBProxy struct {
	DB    *sqlx.DB
	Proxy *toxiproxy.Proxy
}

// Close close proxy related connections.
func (p *DBProxy) Close() {
	_ = p.DB.Close()
}
