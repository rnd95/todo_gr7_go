package database

import (
	"time"

	"github.com/BohdanBoriak/boilerplate-go-back/internal/domain"
	"github.com/upper/db/v4"
)

const TasksTableName = "tasks"

type task struct {
	Id          uint64            `db:"id,omitempty"`
	UserId      uint64            `db:"user_id"`
	Name        string            `db:"name"`
	Description string            `db:"description"`
	Deadline    time.Time         `db:"deadline"`
	Status      domain.TaskStatus `db:"status"`
	CreatedDate time.Time         `db:"created_date"`
	UpdatedDate time.Time         `db:"updated_date"`
	DeletedDate *time.Time        `db:"deleted_date"`
}

type TaskRepository interface {
	Save(t domain.Task) (domain.Task, error)
	FindById(id uint64) (domain.Task, error)
	Update(dev domain.Task) (domain.Task, error)
	Delete(id uint64) error
	FindForUser(id uint64) ([]domain.Task, error)
	MarkCompleted(id uint64) (domain.Task, error)
}

type taskRepository struct {
	coll db.Collection
	sess db.Session
}

func NewTaskRepository(session db.Session) TaskRepository {
	return taskRepository{
		coll: session.Collection(TasksTableName),
		sess: session,
	}
}

func (r taskRepository) Save(t domain.Task) (domain.Task, error) {
	tsk := r.mapDomainToModel(t)
	tsk.CreatedDate = time.Now()
	tsk.UpdatedDate = time.Now()
	err := r.coll.InsertReturning(&tsk)
	if err != nil {
		return domain.Task{}, err
	}
	t = r.mapModelToDomain(tsk)
	return t, nil
}

func (r taskRepository) FindById(id uint64) (domain.Task, error) {
	var t task
	err := r.coll.Find(db.Cond{"id": id, "deleted_date": nil}).One(&t)
	if err != nil {
		return domain.Task{}, err
	}

	return r.mapModelToDomain(t), nil
}

func (r taskRepository) Update(t domain.Task) (domain.Task, error) {
	m := r.mapDomainToModel(t)
	m.UpdatedDate = time.Now()
	err := r.coll.Find(db.Cond{"id": m.Id, "deleted_date": nil}).Update(&m)
	if err != nil {
		return domain.Task{}, err
	}
	return r.mapModelToDomain(m), nil
}

func (r taskRepository) MarkCompleted(id uint64) (domain.Task, error) {
	var m task
	err := r.coll.Find(db.Cond{"id": id, "deleted_date": nil}).One(&m)
	if err != nil {
		return domain.Task{}, err
	}
	m.UpdatedDate = time.Now()
	m.Status = domain.Done
	err = r.coll.Find(db.Cond{"id": m.Id, "deleted_date": nil}).Update(&m)
	if err != nil {
		return domain.Task{}, err
	}
	return r.mapModelToDomain(m), nil
}

func (r taskRepository) Delete(id uint64) error {
	return r.coll.Find(db.Cond{"id": id, "deleted_date": nil}).Update(map[string]interface{}{"deleted_date": time.Now()})
}

func (r taskRepository) FindForUser(id uint64) ([]domain.Task, error) {
	var tasks []task
	err := r.coll.Find(db.Cond{"user_id": id, "deleted_date": nil}).All(&tasks)
	if err != nil {
		return nil, err
	}
	res := r.mapModelToDomainCollection(tasks)
	return res, nil
}

func (r taskRepository) mapDomainToModel(t domain.Task) task {
	return task{
		Id:          t.Id,
		UserId:      t.UserId,
		Name:        t.Name,
		Description: t.Description,
		Deadline:    t.Deadline,
		Status:      t.Status,
		CreatedDate: t.CreatedDate,
		UpdatedDate: t.UpdatedDate,
		DeletedDate: t.DeletedDate,
	}
}

func (r taskRepository) mapModelToDomain(t task) domain.Task {
	return domain.Task{
		Id:          t.Id,
		UserId:      t.UserId,
		Name:        t.Name,
		Description: t.Description,
		Deadline:    t.Deadline,
		Status:      t.Status,
		CreatedDate: t.CreatedDate,
		UpdatedDate: t.UpdatedDate,
		DeletedDate: t.DeletedDate,
	}
}

func (r taskRepository) mapModelToDomainCollection(ts []task) []domain.Task {
	var tasks []domain.Task
	for _, t := range ts {
		tas := r.mapModelToDomain(t)
		tasks = append(tasks, tas)
	}
	return tasks
}
