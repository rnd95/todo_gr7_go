package app

import (
	"log"

	"github.com/BohdanBoriak/boilerplate-go-back/internal/domain"
	"github.com/BohdanBoriak/boilerplate-go-back/internal/infra/database"
)

type TaskService interface {
	Save(t domain.Task) (domain.Task, error)
	Find(id uint64) (interface{}, error)
	Update(d domain.Task) (domain.Task, error)
	Delete(id uint64) error
	FindForUser(id uint64) ([]domain.Task, error)
	MarkCompleted(id uint64) (domain.Task, error)
}

type taskService struct {
	taskRepo database.TaskRepository
}

func NewTaskService(tr database.TaskRepository) TaskService {
	return taskService{
		taskRepo: tr,
	}
}

func (s taskService) Find(id uint64) (interface{}, error) {
	rm, err := s.taskRepo.FindById(id)
	if err != nil {
		log.Printf("TaskService: %s", err)
		return domain.Task{}, err
	}
	return rm, err
}

func (s taskService) MarkCompleted(id uint64) (domain.Task, error) {
	rm, err := s.taskRepo.MarkCompleted(id)
	if err != nil {
		log.Printf("TaskService: %s", err)
		return domain.Task{}, err
	}
	return rm, err
}

func (s taskService) Save(t domain.Task) (domain.Task, error) {
	task, err := s.taskRepo.Save(t)
	if err != nil {
		log.Printf("TaskService: %s", err)
		return domain.Task{}, err
	}
	return task, nil
}

func (s taskService) Update(d domain.Task) (domain.Task, error) {
	d, err := s.taskRepo.Update(d)
	if err != nil {
		log.Printf("TaskService: %s", err)
		return domain.Task{}, err
	}

	return d, nil
}

func (s taskService) Delete(id uint64) error {
	err := s.taskRepo.Delete(id)
	if err != nil {
		log.Printf("TaskService: %s", err)
		return err
	}

	return nil
}

func (s taskService) FindForUser(id uint64) ([]domain.Task, error) {
	orgs, err := s.taskRepo.FindForUser(id)
	if err != nil {
		log.Printf("TaskService: %s", err)
		return nil, err
	}

	return orgs, nil
}
