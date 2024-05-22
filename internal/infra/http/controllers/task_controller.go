package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/BohdanBoriak/boilerplate-go-back/internal/app"
	"github.com/BohdanBoriak/boilerplate-go-back/internal/domain"
	"github.com/BohdanBoriak/boilerplate-go-back/internal/infra/http/requests"
	"github.com/BohdanBoriak/boilerplate-go-back/internal/infra/http/resources"
)

type TaskController struct {
	taskService app.TaskService
}

func NewTaskController(ts app.TaskService) TaskController {
	return TaskController{
		taskService: ts,
	}
}

func (c TaskController) Save() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		task, err := requests.Bind(r, requests.TaskRequest{}, domain.Task{})
		if err != nil {
			log.Printf("TaskController: %s", err)
			BadRequest(w, err)
			return
		}

		user := r.Context().Value(UserKey).(domain.User)
		task.UserId = user.Id
		task.Status = domain.New
		task, err = c.taskService.Save(task)
		if err != nil {
			log.Printf("TaskController: %s", err)
			InternalServerError(w, err)
			return
		}

		var tDto resources.TaskDto
		tDto = tDto.DomainToDto(task)
		Created(w, tDto)
	}
}

func (c TaskController) FindById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(UserKey).(domain.User)
		task := r.Context().Value(TaskKey).(domain.Task)

		if task.UserId != user.Id {
			err := fmt.Errorf("access denied")
			Forbidden(w, err)
			return
		}

		var taskDto resources.TaskDto
		Success(w, taskDto.DomainToDto(task))
	}
}

func (c TaskController) MarkCompleted() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(UserKey).(domain.User)
		task := r.Context().Value(TaskKey).(domain.Task)
		if task.UserId != user.Id {
			err := fmt.Errorf("access denied")
			Forbidden(w, err)
			return
		}

		t, err := c.taskService.MarkCompleted(task.Id)
		if err != nil {
			log.Printf("TaskController: %s", err)
			InternalServerError(w, err)
			return
		}
		var taskDto resources.TaskDto
		Success(w, taskDto.DomainToDto(t))
	}
}

func (c TaskController) FindForUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(UserKey).(domain.User)
		tasks, err := c.taskService.FindForUser(user.Id)
		if err != nil {
			log.Printf("TaskController: %s", err)
			InternalServerError(w, err)
			return
		}

		var tasksDto resources.TasksDto
		response := tasksDto.DomainToDto(tasks)
		Success(w, response)
	}
}

func (c TaskController) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		task := r.Context().Value(TaskKey).(domain.Task)
		err := c.taskService.Delete(task.Id)
		if err != nil {
			log.Printf("TaskController: %s", err)
			InternalServerError(w, err)
			return
		}
		Ok(w)
	}
}

func (c TaskController) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		task, err := requests.Bind(r, requests.TaskRequest{}, domain.Task{})
		if err != nil {
			log.Printf("TaskController: %s", err)
			BadRequest(w, err)
			return
		}

		t := r.Context().Value(TaskKey).(domain.Task)
		u := r.Context().Value(UserKey).(domain.User)
		if t.UserId != u.Id {
			err := fmt.Errorf("access denied")
			Forbidden(w, err)
			return
		}
		t.Name = task.Name
		t.Description = task.Description
		t.Deadline = task.Deadline
		t.Status = task.Status
		task, err = c.taskService.Update(t)
		if err != nil {
			log.Printf("TaskController: %s", err)
			InternalServerError(w, err)
			return
		}

		var taskDto resources.TaskDto
		Success(w, taskDto.DomainToDto(task))
	}
}
