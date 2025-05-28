// @title People API
// @description API для управления данными людей
//
// @host localhost:8080
// @BasePath /
package rest

import (
	"EMtask/testtask/core"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
)

type Handlers struct {
	service *core.Service
	log     *slog.Logger
}

func NewHandlers(service *core.Service, log *slog.Logger) *Handlers {
	return &Handlers{
		service: service,
		log:     log,
	}
}

// CreatePersonHandler создает новую запись человека
// @Summary Создать человека
// @Description Создает запись с обогащением данных (возраст, пол, национальность)
// @Tags people
// @Accept json
// @Produce json
// @Param request body core.PersonRequest true "Данные для создания"
// @Success 201 {object} core.Person
// @Failure 400
// @Failure 500
// @Router /people [post]
func (h *Handlers) CreatePersonHandler(w http.ResponseWriter, r *http.Request) {
	var req core.PersonRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Error("failed to decode request", "error", err)
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	person, err := h.service.CreatePerson(r.Context(), req)
	if err != nil {
		if errors.Is(err, core.ErrAPIFailed) {
			http.Error(w, "failed to get data from API", http.StatusInternalServerError)
			return
		}
		http.Error(w, "failed to create person", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(person)
}

// GetPeopleHandler возвращает список людей
// @Summary Получить список людей
// @Description Возвращает список людей
// @Tags people
// @Produce json
// @Param age query string false "Фильтр по возрасту"
// @Param gender query string false "Фильтр по полу (male/female)"
// @Param nationality query string false "Фильтр по национальности"
// @Param limit query string false "Лимит записей (по умолчанию 10)"
// @Param offset query string false "Смещение (по умолчанию 0)"
// @Success 200 {array} core.Person
// @Failure 500
// @Router /people [get]
func (h *Handlers) GetPeopleHandler(w http.ResponseWriter, r *http.Request) {
	filters := core.PersonFilters{
		Age:         r.URL.Query().Get("age"),
		Gender:      r.URL.Query().Get("gender"),
		Nationality: r.URL.Query().Get("nationality"),
		Limit:       r.URL.Query().Get("limit"),
		Offset:      r.URL.Query().Get("offset"),
	}

	people, err := h.service.GetPeople(r.Context(), filters)
	if err != nil {
		http.Error(w, "failed to get people", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(people)
}

// GetPersonHandler возвращает данные человека
// @Summary Получить данные человека
// @Tags people
// @Produce json
// @Param id query string true "ID человека"
// @Success 200 {object} core.Person
// @Failure 400
// @Failure 404
// @Failure 500
// @Router /person [get]
func (h *Handlers) GetPersonHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	person, err := h.service.GetPerson(r.Context(), id)
	if err != nil {
		if errors.Is(err, core.ErrPersonNotFound) {
			http.Error(w, "person not found", http.StatusNotFound)
			return
		}
		http.Error(w, "failed to get person", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(person)
}

// UpdatePersonHandler обновляет данные человека
// @Summary Обновить данные человека
// @Tags people
// @Accept json
// @Produce json
// @Param request body core.Person true "Обновляемые данные"
// @Success 200
// @Failure 400
// @Failure 404
// @Failure 500
// @Router /person [put]
func (h *Handlers) UpdatePersonHandler(w http.ResponseWriter, r *http.Request) {
	var person core.Person
	if err := json.NewDecoder(r.Body).Decode(&person); err != nil {
		h.log.Error("failed to decode request", "error", err)
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if person.ID == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	if err := h.service.UpdatePerson(r.Context(), person); err != nil {
		if errors.Is(err, core.ErrPersonNotFound) {
			http.Error(w, "person not found", http.StatusNotFound)
			return
		}
		http.Error(w, "failed to update person", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "person updated"})
}

// DeletePersonHandler удаляет запись о человеке
// @Summary Удалить человека
// @Tags people
// @Produce json
// @Param id query string true "ID человека"
// @Success 200
// @Failure 400
// @Failure 404
// @Failure 500
// @Router /person [delete]
func (h *Handlers) DeletePersonHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	if err := h.service.DeletePerson(r.Context(), id); err != nil {
		if errors.Is(err, core.ErrPersonNotFound) {
			http.Error(w, "person not found", http.StatusNotFound)
			return
		}
		http.Error(w, "failed to delete person", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "person deleted"})
}
