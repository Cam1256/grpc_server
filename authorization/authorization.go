package authorization

import (
	"log"

	"github.com/SalviCF/authorization-server/database"
	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
)

type Enforcer struct {
	Enforcer *casbin.Enforcer
}

var db = database.InitDatabase()

// Init Casbin Enforcer loading RBAC with tenants model
func InitEnforcer() *Enforcer {

	a, err := gormadapter.NewAdapterByDB(db.Database)

	if err != nil {
		log.Fatalf("error: adapter: %s", err)
		return nil
	}

	enforcer, err := casbin.NewEnforcer("model/rbac_with_domains_model.conf", a)

	if err != nil {
		log.Fatalf("error: enforcer: %s", err)
		return nil
	}

	return &Enforcer{enforcer}
}
