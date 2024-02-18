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
	Q                      = new(Query)
	Cronjob                *cronjob
	CronjobHistory         *cronjobHistory
	HealthcheckHttp        *healthcheckHttp
	HealthcheckHttpHistory *healthcheckHttpHistory
	HealthcheckTcp         *healthcheckTcp
	HealthcheckTcpHistory  *healthcheckTcpHistory
	OAuth2State            *oAuth2State
	Worker                 *worker
)

func SetDefault(db *gorm.DB, opts ...gen.DOOption) {
	*Q = *Use(db, opts...)
	Cronjob = &Q.Cronjob
	CronjobHistory = &Q.CronjobHistory
	HealthcheckHttp = &Q.HealthcheckHttp
	HealthcheckHttpHistory = &Q.HealthcheckHttpHistory
	HealthcheckTcp = &Q.HealthcheckTcp
	HealthcheckTcpHistory = &Q.HealthcheckTcpHistory
	OAuth2State = &Q.OAuth2State
	Worker = &Q.Worker
}

func Use(db *gorm.DB, opts ...gen.DOOption) *Query {
	return &Query{
		db:                     db,
		Cronjob:                newCronjob(db, opts...),
		CronjobHistory:         newCronjobHistory(db, opts...),
		HealthcheckHttp:        newHealthcheckHttp(db, opts...),
		HealthcheckHttpHistory: newHealthcheckHttpHistory(db, opts...),
		HealthcheckTcp:         newHealthcheckTcp(db, opts...),
		HealthcheckTcpHistory:  newHealthcheckTcpHistory(db, opts...),
		OAuth2State:            newOAuth2State(db, opts...),
		Worker:                 newWorker(db, opts...),
	}
}

type Query struct {
	db *gorm.DB

	Cronjob                cronjob
	CronjobHistory         cronjobHistory
	HealthcheckHttp        healthcheckHttp
	HealthcheckHttpHistory healthcheckHttpHistory
	HealthcheckTcp         healthcheckTcp
	HealthcheckTcpHistory  healthcheckTcpHistory
	OAuth2State            oAuth2State
	Worker                 worker
}

func (q *Query) Available() bool { return q.db != nil }

func (q *Query) clone(db *gorm.DB) *Query {
	return &Query{
		db:                     db,
		Cronjob:                q.Cronjob.clone(db),
		CronjobHistory:         q.CronjobHistory.clone(db),
		HealthcheckHttp:        q.HealthcheckHttp.clone(db),
		HealthcheckHttpHistory: q.HealthcheckHttpHistory.clone(db),
		HealthcheckTcp:         q.HealthcheckTcp.clone(db),
		HealthcheckTcpHistory:  q.HealthcheckTcpHistory.clone(db),
		OAuth2State:            q.OAuth2State.clone(db),
		Worker:                 q.Worker.clone(db),
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
		db:                     db,
		Cronjob:                q.Cronjob.replaceDB(db),
		CronjobHistory:         q.CronjobHistory.replaceDB(db),
		HealthcheckHttp:        q.HealthcheckHttp.replaceDB(db),
		HealthcheckHttpHistory: q.HealthcheckHttpHistory.replaceDB(db),
		HealthcheckTcp:         q.HealthcheckTcp.replaceDB(db),
		HealthcheckTcpHistory:  q.HealthcheckTcpHistory.replaceDB(db),
		OAuth2State:            q.OAuth2State.replaceDB(db),
		Worker:                 q.Worker.replaceDB(db),
	}
}

type queryCtx struct {
	Cronjob                ICronjobDo
	CronjobHistory         ICronjobHistoryDo
	HealthcheckHttp        IHealthcheckHttpDo
	HealthcheckHttpHistory IHealthcheckHttpHistoryDo
	HealthcheckTcp         IHealthcheckTcpDo
	HealthcheckTcpHistory  IHealthcheckTcpHistoryDo
	OAuth2State            IOAuth2StateDo
	Worker                 IWorkerDo
}

func (q *Query) WithContext(ctx context.Context) *queryCtx {
	return &queryCtx{
		Cronjob:                q.Cronjob.WithContext(ctx),
		CronjobHistory:         q.CronjobHistory.WithContext(ctx),
		HealthcheckHttp:        q.HealthcheckHttp.WithContext(ctx),
		HealthcheckHttpHistory: q.HealthcheckHttpHistory.WithContext(ctx),
		HealthcheckTcp:         q.HealthcheckTcp.WithContext(ctx),
		HealthcheckTcpHistory:  q.HealthcheckTcpHistory.WithContext(ctx),
		OAuth2State:            q.OAuth2State.WithContext(ctx),
		Worker:                 q.Worker.WithContext(ctx),
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
