// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package query

import (
	"context"
	"database/sql"

	"gorm.io/gorm"

	"gorm.io/gen"

	"gorm.io/plugin/dbresolver"
)

var (
	Q              = new(Query)
	Cronjob        *cronjob
	CronjobHistory *cronjobHistory
	Monitor        *monitor
	MonitorHistory *monitorHistory
	OAuth2State    *oAuth2State
	Worker         *worker
)

func SetDefault(db *gorm.DB, opts ...gen.DOOption) {
	*Q = *Use(db, opts...)
	Cronjob = &Q.Cronjob
	CronjobHistory = &Q.CronjobHistory
	Monitor = &Q.Monitor
	MonitorHistory = &Q.MonitorHistory
	OAuth2State = &Q.OAuth2State
	Worker = &Q.Worker
}

func Use(db *gorm.DB, opts ...gen.DOOption) *Query {
	return &Query{
		db:             db,
		Cronjob:        newCronjob(db, opts...),
		CronjobHistory: newCronjobHistory(db, opts...),
		Monitor:        newMonitor(db, opts...),
		MonitorHistory: newMonitorHistory(db, opts...),
		OAuth2State:    newOAuth2State(db, opts...),
		Worker:         newWorker(db, opts...),
	}
}

type Query struct {
	db *gorm.DB

	Cronjob        cronjob
	CronjobHistory cronjobHistory
	Monitor        monitor
	MonitorHistory monitorHistory
	OAuth2State    oAuth2State
	Worker         worker
}

func (q *Query) Available() bool { return q.db != nil }

func (q *Query) clone(db *gorm.DB) *Query {
	return &Query{
		db:             db,
		Cronjob:        q.Cronjob.clone(db),
		CronjobHistory: q.CronjobHistory.clone(db),
		Monitor:        q.Monitor.clone(db),
		MonitorHistory: q.MonitorHistory.clone(db),
		OAuth2State:    q.OAuth2State.clone(db),
		Worker:         q.Worker.clone(db),
	}
}

func (q *Query) ReadDB() *Query {
	return q.ReplaceDB(q.db.Clauses(dbresolver.Read))
}

func (q *Query) WriteDB() *Query {
	return q.ReplaceDB(q.db.Clauses(dbresolver.Write))
}

func (q *Query) ReplaceDB(db *gorm.DB) *Query {
	return &Query{
		db:             db,
		Cronjob:        q.Cronjob.replaceDB(db),
		CronjobHistory: q.CronjobHistory.replaceDB(db),
		Monitor:        q.Monitor.replaceDB(db),
		MonitorHistory: q.MonitorHistory.replaceDB(db),
		OAuth2State:    q.OAuth2State.replaceDB(db),
		Worker:         q.Worker.replaceDB(db),
	}
}

type queryCtx struct {
	Cronjob        ICronjobDo
	CronjobHistory ICronjobHistoryDo
	Monitor        IMonitorDo
	MonitorHistory IMonitorHistoryDo
	OAuth2State    IOAuth2StateDo
	Worker         IWorkerDo
}

func (q *Query) WithContext(ctx context.Context) *queryCtx {
	return &queryCtx{
		Cronjob:        q.Cronjob.WithContext(ctx),
		CronjobHistory: q.CronjobHistory.WithContext(ctx),
		Monitor:        q.Monitor.WithContext(ctx),
		MonitorHistory: q.MonitorHistory.WithContext(ctx),
		OAuth2State:    q.OAuth2State.WithContext(ctx),
		Worker:         q.Worker.WithContext(ctx),
	}
}

func (q *Query) Transaction(fc func(tx *Query) error, opts ...*sql.TxOptions) error {
	return q.db.Transaction(func(tx *gorm.DB) error { return fc(q.clone(tx)) }, opts...)
}

func (q *Query) Begin(opts ...*sql.TxOptions) *QueryTx {
	tx := q.db.Begin(opts...)
	return &QueryTx{Query: q.clone(tx), Error: tx.Error}
}

type QueryTx struct {
	*Query
	Error error
}

func (q *QueryTx) Commit() error {
	return q.db.Commit().Error
}

func (q *QueryTx) Rollback() error {
	return q.db.Rollback().Error
}

func (q *QueryTx) SavePoint(name string) error {
	return q.db.SavePoint(name).Error
}

func (q *QueryTx) RollbackTo(name string) error {
	return q.db.RollbackTo(name).Error
}
