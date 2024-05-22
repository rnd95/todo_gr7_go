package resources

import (
	"time"

	"github.com/BohdanBoriak/boilerplate-go-back/internal/domain"
)

type TaskDto struct {
	Id          uint64            `json:"id"`
	UserId      uint64            `json:"userId"`
	Name        string            `json:"name"`
	Description string            `json:"description,omitempty"`
	Deadline    time.Time         `json:"deadline"`
	Status      domain.TaskStatus `json:"status"`
}

type TasksDto struct {
	Tasks []TaskDto `json:"tasks"`
}

func (d TaskDto) DomainToDto(t domain.Task) TaskDto {
	return TaskDto{
		Id:          t.Id,
		UserId:      t.UserId,
		Name:        t.Name,
		Description: t.Description,
		Deadline:    t.Deadline,
		Status:      t.Status,
	}
}

func (d TasksDto) DomainToDto(orgs []domain.Task) TasksDto {
	var tasks []TaskDto
	for _, o := range orgs {
		var oDto TaskDto
		org := oDto.DomainToDto(o)
		tasks = append(tasks, org)
	}
	response := TasksDto{
		Tasks: tasks,
	}
	return response
}
