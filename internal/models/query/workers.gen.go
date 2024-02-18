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

func newWorker(db *gorm.DB, opts ...gen.DOOption) worker {
	_worker := worker{}

	_worker.workerDo.UseDB(db, opts...)
	_worker.workerDo.UseModel(&models.Worker{})

	tableName := _worker.workerDo.TableName()
	_worker.ALL = field.NewAsterisk(tableName)
	_worker.ID = field.NewUint(tableName, "id")
	_worker.CreatedAt = field.NewTime(tableName, "created_at")
	_worker.UpdatedAt = field.NewTime(tableName, "updated_at")
	_worker.DeletedAt = field.NewField(tableName, "deleted_at")
	_worker.Name = field.NewString(tableName, "name")
	_worker.Slug = field.NewString(tableName, "slug")
	_worker.Status = field.NewString(tableName, "status")

	_worker.fillFieldMap()

	return _worker
}

type worker struct {
	workerDo workerDo

	ALL       field.Asterisk
	ID        field.Uint
	CreatedAt field.Time
	UpdatedAt field.Time
	DeletedAt field.Field
	Name      field.String
	Slug      field.String
	Status    field.String

	fieldMap map[string]field.Expr
}

func (w worker) Table(newTableName string) *worker {
	w.workerDo.UseTable(newTableName)
	return w.updateTableName(newTableName)
}

func (w worker) As(alias string) *worker {
	w.workerDo.DO = *(w.workerDo.As(alias).(*gen.DO))
	return w.updateTableName(alias)
}

func (w *worker) updateTableName(table string) *worker {
	w.ALL = field.NewAsterisk(table)
	w.ID = field.NewUint(table, "id")
	w.CreatedAt = field.NewTime(table, "created_at")
	w.UpdatedAt = field.NewTime(table, "updated_at")
	w.DeletedAt = field.NewField(table, "deleted_at")
	w.Name = field.NewString(table, "name")
	w.Slug = field.NewString(table, "slug")
	w.Status = field.NewString(table, "status")

	w.fillFieldMap()

	return w
}

func (w *worker) WithContext(ctx context.Context) IWorkerDo { return w.workerDo.WithContext(ctx) }

func (w worker) TableName() string { return w.workerDo.TableName() }

func (w worker) Alias() string { return w.workerDo.Alias() }

func (w worker) Columns(cols ...field.Expr) gen.Columns { return w.workerDo.Columns(cols...) }

func (w *worker) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := w.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (w *worker) fillFieldMap() {
	w.fieldMap = make(map[string]field.Expr, 7)
	w.fieldMap["id"] = w.ID
	w.fieldMap["created_at"] = w.CreatedAt
	w.fieldMap["updated_at"] = w.UpdatedAt
	w.fieldMap["deleted_at"] = w.DeletedAt
	w.fieldMap["name"] = w.Name
	w.fieldMap["slug"] = w.Slug
	w.fieldMap["status"] = w.Status
}

func (w worker) clone(db *gorm.DB) worker {
	w.workerDo.ReplaceConnPool(db.Statement.ConnPool)
	return w
}

func (w worker) replaceDB(db *gorm.DB) worker {
	w.workerDo.ReplaceDB(db)
	return w
}

type workerDo struct{ gen.DO }

type IWorkerDo interface {
	gen.SubQuery
	Debug() IWorkerDo
	WithContext(ctx context.Context) IWorkerDo
	WithResult(fc func(tx gen.Dao)) gen.ResultInfo
	ReplaceDB(db *gorm.DB)
	ReadDB() IWorkerDo
	WriteDB() IWorkerDo
	As(alias string) gen.Dao
	Session(config *gorm.Session) IWorkerDo
	Columns(cols ...field.Expr) gen.Columns
	Clauses(conds ...clause.Expression) IWorkerDo
	Not(conds ...gen.Condition) IWorkerDo
	Or(conds ...gen.Condition) IWorkerDo
	Select(conds ...field.Expr) IWorkerDo
	Where(conds ...gen.Condition) IWorkerDo
	Order(conds ...field.Expr) IWorkerDo
	Distinct(cols ...field.Expr) IWorkerDo
	Omit(cols ...field.Expr) IWorkerDo
	Join(table schema.Tabler, on ...field.Expr) IWorkerDo
	LeftJoin(table schema.Tabler, on ...field.Expr) IWorkerDo
	RightJoin(table schema.Tabler, on ...field.Expr) IWorkerDo
	Group(cols ...field.Expr) IWorkerDo
	Having(conds ...gen.Condition) IWorkerDo
	Limit(limit int) IWorkerDo
	Offset(offset int) IWorkerDo
	Count() (count int64, err error)
	Scopes(funcs ...func(gen.Dao) gen.Dao) IWorkerDo
	Unscoped() IWorkerDo
	Create(values ...*models.Worker) error
	CreateInBatches(values []*models.Worker, batchSize int) error
	Save(values ...*models.Worker) error
	First() (*models.Worker, error)
	Take() (*models.Worker, error)
	Last() (*models.Worker, error)
	Find() ([]*models.Worker, error)
	FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*models.Worker, err error)
	FindInBatches(result *[]*models.Worker, batchSize int, fc func(tx gen.Dao, batch int) error) error
	Pluck(column field.Expr, dest interface{}) error
	Delete(...*models.Worker) (info gen.ResultInfo, err error)
	Update(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	Updates(value interface{}) (info gen.ResultInfo, err error)
	UpdateColumn(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateColumnSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	UpdateColumns(value interface{}) (info gen.ResultInfo, err error)
	UpdateFrom(q gen.SubQuery) gen.Dao
	Attrs(attrs ...field.AssignExpr) IWorkerDo
	Assign(attrs ...field.AssignExpr) IWorkerDo
	Joins(fields ...field.RelationField) IWorkerDo
	Preload(fields ...field.RelationField) IWorkerDo
	FirstOrInit() (*models.Worker, error)
	FirstOrCreate() (*models.Worker, error)
	FindByPage(offset int, limit int) (result []*models.Worker, count int64, err error)
	ScanByPage(result interface{}, offset int, limit int) (count int64, err error)
	Scan(result interface{}) (err error)
	Returning(value interface{}, columns ...string) IWorkerDo
	UnderlyingDB() *gorm.DB
	schema.Tabler
}

func (w workerDo) Debug() IWorkerDo {
	return w.withDO(w.DO.Debug())
}

func (w workerDo) WithContext(ctx context.Context) IWorkerDo {
	return w.withDO(w.DO.WithContext(ctx))
}

func (w workerDo) ReadDB() IWorkerDo {
	return w.Clauses(dbresolver.Read)
}

func (w workerDo) WriteDB() IWorkerDo {
	return w.Clauses(dbresolver.Write)
}

func (w workerDo) Session(config *gorm.Session) IWorkerDo {
	return w.withDO(w.DO.Session(config))
}

func (w workerDo) Clauses(conds ...clause.Expression) IWorkerDo {
	return w.withDO(w.DO.Clauses(conds...))
}

func (w workerDo) Returning(value interface{}, columns ...string) IWorkerDo {
	return w.withDO(w.DO.Returning(value, columns...))
}

func (w workerDo) Not(conds ...gen.Condition) IWorkerDo {
	return w.withDO(w.DO.Not(conds...))
}

func (w workerDo) Or(conds ...gen.Condition) IWorkerDo {
	return w.withDO(w.DO.Or(conds...))
}

func (w workerDo) Select(conds ...field.Expr) IWorkerDo {
	return w.withDO(w.DO.Select(conds...))
}

func (w workerDo) Where(conds ...gen.Condition) IWorkerDo {
	return w.withDO(w.DO.Where(conds...))
}

func (w workerDo) Order(conds ...field.Expr) IWorkerDo {
	return w.withDO(w.DO.Order(conds...))
}

func (w workerDo) Distinct(cols ...field.Expr) IWorkerDo {
	return w.withDO(w.DO.Distinct(cols...))
}

func (w workerDo) Omit(cols ...field.Expr) IWorkerDo {
	return w.withDO(w.DO.Omit(cols...))
}

func (w workerDo) Join(table schema.Tabler, on ...field.Expr) IWorkerDo {
	return w.withDO(w.DO.Join(table, on...))
}

func (w workerDo) LeftJoin(table schema.Tabler, on ...field.Expr) IWorkerDo {
	return w.withDO(w.DO.LeftJoin(table, on...))
}

func (w workerDo) RightJoin(table schema.Tabler, on ...field.Expr) IWorkerDo {
	return w.withDO(w.DO.RightJoin(table, on...))
}

func (w workerDo) Group(cols ...field.Expr) IWorkerDo {
	return w.withDO(w.DO.Group(cols...))
}

func (w workerDo) Having(conds ...gen.Condition) IWorkerDo {
	return w.withDO(w.DO.Having(conds...))
}

func (w workerDo) Limit(limit int) IWorkerDo {
	return w.withDO(w.DO.Limit(limit))
}

func (w workerDo) Offset(offset int) IWorkerDo {
	return w.withDO(w.DO.Offset(offset))
}

func (w workerDo) Scopes(funcs ...func(gen.Dao) gen.Dao) IWorkerDo {
	return w.withDO(w.DO.Scopes(funcs...))
}

func (w workerDo) Unscoped() IWorkerDo {
	return w.withDO(w.DO.Unscoped())
}

func (w workerDo) Create(values ...*models.Worker) error {
	if len(values) == 0 {
		return nil
	}
	return w.DO.Create(values)
}

func (w workerDo) CreateInBatches(values []*models.Worker, batchSize int) error {
	return w.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (w workerDo) Save(values ...*models.Worker) error {
	if len(values) == 0 {
		return nil
	}
	return w.DO.Save(values)
}

func (w workerDo) First() (*models.Worker, error) {
	if result, err := w.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*models.Worker), nil
	}
}

func (w workerDo) Take() (*models.Worker, error) {
	if result, err := w.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*models.Worker), nil
	}
}

func (w workerDo) Last() (*models.Worker, error) {
	if result, err := w.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*models.Worker), nil
	}
}

func (w workerDo) Find() ([]*models.Worker, error) {
	result, err := w.DO.Find()
	return result.([]*models.Worker), err
}

func (w workerDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*models.Worker, err error) {
	buf := make([]*models.Worker, 0, batchSize)
	err = w.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (w workerDo) FindInBatches(result *[]*models.Worker, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return w.DO.FindInBatches(result, batchSize, fc)
}

func (w workerDo) Attrs(attrs ...field.AssignExpr) IWorkerDo {
	return w.withDO(w.DO.Attrs(attrs...))
}

func (w workerDo) Assign(attrs ...field.AssignExpr) IWorkerDo {
	return w.withDO(w.DO.Assign(attrs...))
}

func (w workerDo) Joins(fields ...field.RelationField) IWorkerDo {
	for _, _f := range fields {
		w = *w.withDO(w.DO.Joins(_f))
	}
	return &w
}

func (w workerDo) Preload(fields ...field.RelationField) IWorkerDo {
	for _, _f := range fields {
		w = *w.withDO(w.DO.Preload(_f))
	}
	return &w
}

func (w workerDo) FirstOrInit() (*models.Worker, error) {
	if result, err := w.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*models.Worker), nil
	}
}

func (w workerDo) FirstOrCreate() (*models.Worker, error) {
	if result, err := w.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*models.Worker), nil
	}
}

func (w workerDo) FindByPage(offset int, limit int) (result []*models.Worker, count int64, err error) {
	result, err = w.Offset(offset).Limit(limit).Find()
	if err != nil {
		return
	}

	if size := len(result); 0 < limit && 0 < size && size < limit {
		count = int64(size + offset)
		return
	}

	count, err = w.Offset(-1).Limit(-1).Count()
	return
}

func (w workerDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = w.Count()
	if err != nil {
		return
	}

	err = w.Offset(offset).Limit(limit).Scan(result)
	return
}

func (w workerDo) Scan(result interface{}) (err error) {
	return w.DO.Scan(result)
}

func (w workerDo) Delete(models ...*models.Worker) (result gen.ResultInfo, err error) {
	return w.DO.Delete(models)
}

func (w *workerDo) withDO(do gen.Dao) *workerDo {
	w.DO = *do.(*gen.DO)
	return w
}
