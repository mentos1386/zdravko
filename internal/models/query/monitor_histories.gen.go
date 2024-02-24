// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package query

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"

	"gorm.io/gen"
	"gorm.io/gen/field"

	"gorm.io/plugin/dbresolver"

	"code.tjo.space/mentos1386/zdravko/internal/models"
)

func newMonitorHistory(db *gorm.DB, opts ...gen.DOOption) monitorHistory {
	_monitorHistory := monitorHistory{}

	_monitorHistory.monitorHistoryDo.UseDB(db, opts...)
	_monitorHistory.monitorHistoryDo.UseModel(&models.MonitorHistory{})

	tableName := _monitorHistory.monitorHistoryDo.TableName()
	_monitorHistory.ALL = field.NewAsterisk(tableName)
	_monitorHistory.ID = field.NewUint(tableName, "id")
	_monitorHistory.CreatedAt = field.NewTime(tableName, "created_at")
	_monitorHistory.UpdatedAt = field.NewTime(tableName, "updated_at")
	_monitorHistory.DeletedAt = field.NewField(tableName, "deleted_at")
	_monitorHistory.Monitor = field.NewUint(tableName, "monitor")
	_monitorHistory.Status = field.NewString(tableName, "status")
	_monitorHistory.Note = field.NewString(tableName, "note")

	_monitorHistory.fillFieldMap()

	return _monitorHistory
}

type monitorHistory struct {
	monitorHistoryDo monitorHistoryDo

	ALL       field.Asterisk
	ID        field.Uint
	CreatedAt field.Time
	UpdatedAt field.Time
	DeletedAt field.Field
	Monitor   field.Uint
	Status    field.String
	Note      field.String

	fieldMap map[string]field.Expr
}

func (m monitorHistory) Table(newTableName string) *monitorHistory {
	m.monitorHistoryDo.UseTable(newTableName)
	return m.updateTableName(newTableName)
}

func (m monitorHistory) As(alias string) *monitorHistory {
	m.monitorHistoryDo.DO = *(m.monitorHistoryDo.As(alias).(*gen.DO))
	return m.updateTableName(alias)
}

func (m *monitorHistory) updateTableName(table string) *monitorHistory {
	m.ALL = field.NewAsterisk(table)
	m.ID = field.NewUint(table, "id")
	m.CreatedAt = field.NewTime(table, "created_at")
	m.UpdatedAt = field.NewTime(table, "updated_at")
	m.DeletedAt = field.NewField(table, "deleted_at")
	m.Monitor = field.NewUint(table, "monitor")
	m.Status = field.NewString(table, "status")
	m.Note = field.NewString(table, "note")

	m.fillFieldMap()

	return m
}

func (m *monitorHistory) WithContext(ctx context.Context) IMonitorHistoryDo {
	return m.monitorHistoryDo.WithContext(ctx)
}

func (m monitorHistory) TableName() string { return m.monitorHistoryDo.TableName() }

func (m monitorHistory) Alias() string { return m.monitorHistoryDo.Alias() }

func (m monitorHistory) Columns(cols ...field.Expr) gen.Columns {
	return m.monitorHistoryDo.Columns(cols...)
}

func (m *monitorHistory) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := m.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (m *monitorHistory) fillFieldMap() {
	m.fieldMap = make(map[string]field.Expr, 7)
	m.fieldMap["id"] = m.ID
	m.fieldMap["created_at"] = m.CreatedAt
	m.fieldMap["updated_at"] = m.UpdatedAt
	m.fieldMap["deleted_at"] = m.DeletedAt
	m.fieldMap["monitor"] = m.Monitor
	m.fieldMap["status"] = m.Status
	m.fieldMap["note"] = m.Note
}

func (m monitorHistory) clone(db *gorm.DB) monitorHistory {
	m.monitorHistoryDo.ReplaceConnPool(db.Statement.ConnPool)
	return m
}

func (m monitorHistory) replaceDB(db *gorm.DB) monitorHistory {
	m.monitorHistoryDo.ReplaceDB(db)
	return m
}

type monitorHistoryDo struct{ gen.DO }

type IMonitorHistoryDo interface {
	gen.SubQuery
	Debug() IMonitorHistoryDo
	WithContext(ctx context.Context) IMonitorHistoryDo
	WithResult(fc func(tx gen.Dao)) gen.ResultInfo
	ReplaceDB(db *gorm.DB)
	ReadDB() IMonitorHistoryDo
	WriteDB() IMonitorHistoryDo
	As(alias string) gen.Dao
	Session(config *gorm.Session) IMonitorHistoryDo
	Columns(cols ...field.Expr) gen.Columns
	Clauses(conds ...clause.Expression) IMonitorHistoryDo
	Not(conds ...gen.Condition) IMonitorHistoryDo
	Or(conds ...gen.Condition) IMonitorHistoryDo
	Select(conds ...field.Expr) IMonitorHistoryDo
	Where(conds ...gen.Condition) IMonitorHistoryDo
	Order(conds ...field.Expr) IMonitorHistoryDo
	Distinct(cols ...field.Expr) IMonitorHistoryDo
	Omit(cols ...field.Expr) IMonitorHistoryDo
	Join(table schema.Tabler, on ...field.Expr) IMonitorHistoryDo
	LeftJoin(table schema.Tabler, on ...field.Expr) IMonitorHistoryDo
	RightJoin(table schema.Tabler, on ...field.Expr) IMonitorHistoryDo
	Group(cols ...field.Expr) IMonitorHistoryDo
	Having(conds ...gen.Condition) IMonitorHistoryDo
	Limit(limit int) IMonitorHistoryDo
	Offset(offset int) IMonitorHistoryDo
	Count() (count int64, err error)
	Scopes(funcs ...func(gen.Dao) gen.Dao) IMonitorHistoryDo
	Unscoped() IMonitorHistoryDo
	Create(values ...*models.MonitorHistory) error
	CreateInBatches(values []*models.MonitorHistory, batchSize int) error
	Save(values ...*models.MonitorHistory) error
	First() (*models.MonitorHistory, error)
	Take() (*models.MonitorHistory, error)
	Last() (*models.MonitorHistory, error)
	Find() ([]*models.MonitorHistory, error)
	FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*models.MonitorHistory, err error)
	FindInBatches(result *[]*models.MonitorHistory, batchSize int, fc func(tx gen.Dao, batch int) error) error
	Pluck(column field.Expr, dest interface{}) error
	Delete(...*models.MonitorHistory) (info gen.ResultInfo, err error)
	Update(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	Updates(value interface{}) (info gen.ResultInfo, err error)
	UpdateColumn(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateColumnSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	UpdateColumns(value interface{}) (info gen.ResultInfo, err error)
	UpdateFrom(q gen.SubQuery) gen.Dao
	Attrs(attrs ...field.AssignExpr) IMonitorHistoryDo
	Assign(attrs ...field.AssignExpr) IMonitorHistoryDo
	Joins(fields ...field.RelationField) IMonitorHistoryDo
	Preload(fields ...field.RelationField) IMonitorHistoryDo
	FirstOrInit() (*models.MonitorHistory, error)
	FirstOrCreate() (*models.MonitorHistory, error)
	FindByPage(offset int, limit int) (result []*models.MonitorHistory, count int64, err error)
	ScanByPage(result interface{}, offset int, limit int) (count int64, err error)
	Scan(result interface{}) (err error)
	Returning(value interface{}, columns ...string) IMonitorHistoryDo
	UnderlyingDB() *gorm.DB
	schema.Tabler
}

func (m monitorHistoryDo) Debug() IMonitorHistoryDo {
	return m.withDO(m.DO.Debug())
}

func (m monitorHistoryDo) WithContext(ctx context.Context) IMonitorHistoryDo {
	return m.withDO(m.DO.WithContext(ctx))
}

func (m monitorHistoryDo) ReadDB() IMonitorHistoryDo {
	return m.Clauses(dbresolver.Read)
}

func (m monitorHistoryDo) WriteDB() IMonitorHistoryDo {
	return m.Clauses(dbresolver.Write)
}

func (m monitorHistoryDo) Session(config *gorm.Session) IMonitorHistoryDo {
	return m.withDO(m.DO.Session(config))
}

func (m monitorHistoryDo) Clauses(conds ...clause.Expression) IMonitorHistoryDo {
	return m.withDO(m.DO.Clauses(conds...))
}

func (m monitorHistoryDo) Returning(value interface{}, columns ...string) IMonitorHistoryDo {
	return m.withDO(m.DO.Returning(value, columns...))
}

func (m monitorHistoryDo) Not(conds ...gen.Condition) IMonitorHistoryDo {
	return m.withDO(m.DO.Not(conds...))
}

func (m monitorHistoryDo) Or(conds ...gen.Condition) IMonitorHistoryDo {
	return m.withDO(m.DO.Or(conds...))
}

func (m monitorHistoryDo) Select(conds ...field.Expr) IMonitorHistoryDo {
	return m.withDO(m.DO.Select(conds...))
}

func (m monitorHistoryDo) Where(conds ...gen.Condition) IMonitorHistoryDo {
	return m.withDO(m.DO.Where(conds...))
}

func (m monitorHistoryDo) Order(conds ...field.Expr) IMonitorHistoryDo {
	return m.withDO(m.DO.Order(conds...))
}

func (m monitorHistoryDo) Distinct(cols ...field.Expr) IMonitorHistoryDo {
	return m.withDO(m.DO.Distinct(cols...))
}

func (m monitorHistoryDo) Omit(cols ...field.Expr) IMonitorHistoryDo {
	return m.withDO(m.DO.Omit(cols...))
}

func (m monitorHistoryDo) Join(table schema.Tabler, on ...field.Expr) IMonitorHistoryDo {
	return m.withDO(m.DO.Join(table, on...))
}

func (m monitorHistoryDo) LeftJoin(table schema.Tabler, on ...field.Expr) IMonitorHistoryDo {
	return m.withDO(m.DO.LeftJoin(table, on...))
}

func (m monitorHistoryDo) RightJoin(table schema.Tabler, on ...field.Expr) IMonitorHistoryDo {
	return m.withDO(m.DO.RightJoin(table, on...))
}

func (m monitorHistoryDo) Group(cols ...field.Expr) IMonitorHistoryDo {
	return m.withDO(m.DO.Group(cols...))
}

func (m monitorHistoryDo) Having(conds ...gen.Condition) IMonitorHistoryDo {
	return m.withDO(m.DO.Having(conds...))
}

func (m monitorHistoryDo) Limit(limit int) IMonitorHistoryDo {
	return m.withDO(m.DO.Limit(limit))
}

func (m monitorHistoryDo) Offset(offset int) IMonitorHistoryDo {
	return m.withDO(m.DO.Offset(offset))
}

func (m monitorHistoryDo) Scopes(funcs ...func(gen.Dao) gen.Dao) IMonitorHistoryDo {
	return m.withDO(m.DO.Scopes(funcs...))
}

func (m monitorHistoryDo) Unscoped() IMonitorHistoryDo {
	return m.withDO(m.DO.Unscoped())
}

func (m monitorHistoryDo) Create(values ...*models.MonitorHistory) error {
	if len(values) == 0 {
		return nil
	}
	return m.DO.Create(values)
}

func (m monitorHistoryDo) CreateInBatches(values []*models.MonitorHistory, batchSize int) error {
	return m.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (m monitorHistoryDo) Save(values ...*models.MonitorHistory) error {
	if len(values) == 0 {
		return nil
	}
	return m.DO.Save(values)
}

func (m monitorHistoryDo) First() (*models.MonitorHistory, error) {
	if result, err := m.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*models.MonitorHistory), nil
	}
}

func (m monitorHistoryDo) Take() (*models.MonitorHistory, error) {
	if result, err := m.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*models.MonitorHistory), nil
	}
}

func (m monitorHistoryDo) Last() (*models.MonitorHistory, error) {
	if result, err := m.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*models.MonitorHistory), nil
	}
}

func (m monitorHistoryDo) Find() ([]*models.MonitorHistory, error) {
	result, err := m.DO.Find()
	return result.([]*models.MonitorHistory), err
}

func (m monitorHistoryDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*models.MonitorHistory, err error) {
	buf := make([]*models.MonitorHistory, 0, batchSize)
	err = m.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (m monitorHistoryDo) FindInBatches(result *[]*models.MonitorHistory, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return m.DO.FindInBatches(result, batchSize, fc)
}

func (m monitorHistoryDo) Attrs(attrs ...field.AssignExpr) IMonitorHistoryDo {
	return m.withDO(m.DO.Attrs(attrs...))
}

func (m monitorHistoryDo) Assign(attrs ...field.AssignExpr) IMonitorHistoryDo {
	return m.withDO(m.DO.Assign(attrs...))
}

func (m monitorHistoryDo) Joins(fields ...field.RelationField) IMonitorHistoryDo {
	for _, _f := range fields {
		m = *m.withDO(m.DO.Joins(_f))
	}
	return &m
}

func (m monitorHistoryDo) Preload(fields ...field.RelationField) IMonitorHistoryDo {
	for _, _f := range fields {
		m = *m.withDO(m.DO.Preload(_f))
	}
	return &m
}

func (m monitorHistoryDo) FirstOrInit() (*models.MonitorHistory, error) {
	if result, err := m.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*models.MonitorHistory), nil
	}
}

func (m monitorHistoryDo) FirstOrCreate() (*models.MonitorHistory, error) {
	if result, err := m.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*models.MonitorHistory), nil
	}
}

func (m monitorHistoryDo) FindByPage(offset int, limit int) (result []*models.MonitorHistory, count int64, err error) {
	result, err = m.Offset(offset).Limit(limit).Find()
	if err != nil {
		return
	}

	if size := len(result); 0 < limit && 0 < size && size < limit {
		count = int64(size + offset)
		return
	}

	count, err = m.Offset(-1).Limit(-1).Count()
	return
}

func (m monitorHistoryDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = m.Count()
	if err != nil {
		return
	}

	err = m.Offset(offset).Limit(limit).Scan(result)
	return
}

func (m monitorHistoryDo) Scan(result interface{}) (err error) {
	return m.DO.Scan(result)
}

func (m monitorHistoryDo) Delete(models ...*models.MonitorHistory) (result gen.ResultInfo, err error) {
	return m.DO.Delete(models)
}

func (m *monitorHistoryDo) withDO(do gen.Dao) *monitorHistoryDo {
	m.DO = *do.(*gen.DO)
	return m
}
