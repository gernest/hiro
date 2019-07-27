package access

import (
	"database/sql"

	"github.com/gernest/hiro/bus"
	"github.com/jmoiron/sqlx"

	// use postgres
	_ "github.com/lib/pq"
	"github.com/ory/ladon"
	manager "github.com/ory/ladon/manager/sql"
)

// New returns new instance of lodon.Lodon that uses postgres database for
// persistance.
func New(conn *sql.DB, eventbus *bus.Producer) (*ladon.Ladon, error) {
	db := sqlx.NewDb(conn, "postgres")
	mngr := manager.NewSQLManager(db, nil)
	if _, err := mngr.CreateSchemas("", ""); err != nil {
		return nil, err
	}
	w := &ladon.Ladon{
		Manager:     mngr,
		AuditLogger: NewAuditLogger(eventbus),
	}
	return w, nil
}

func NewAuditLogger(eventbus *bus.Producer) ladon.AuditLogger {
	if eventbus != nil {
		return &AuditLoggerNSQ{bus: eventbus}
	}
	return &ladon.AuditLoggerNoOp{}
}

// AuditLoggerNSQ implements ladon.AuditLogger interface. This sends the logged
// message to the nsq message bus.
type AuditLoggerNSQ struct {
	bus *bus.Producer
}

// LogRejectedAccessRequest sends an nsq message for failed deciders.
func (a *AuditLoggerNSQ) LogRejectedAccessRequest(r *ladon.Request, p ladon.Policies, d ladon.Policies) {
	msg := bus.AccessPolicy{}
	for _, v := range d {
		msg.Denied = append(msg.Denied, v.GetID())
	}
	a.bus.PublishAccessPolicy(msg)
}

// LogGrantedAccessRequest sends nsq message for succeeded deciders.
func (a *AuditLoggerNSQ) LogGrantedAccessRequest(r *ladon.Request, p ladon.Policies, d ladon.Policies) {
	msg := bus.AccessPolicy{}
	for _, v := range d {
		msg.Allowed = append(msg.Allowed, v.GetID())
	}
	a.bus.PublishAccessPolicy(msg)
}
