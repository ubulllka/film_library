package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"
	"vk/internal/db"
	"vk/internal/models"
	"vk/internal/models/DTO"
)

type ActorPostgres struct {
	db *sql.DB
}

func NewActorPostgres(db *sql.DB) *ActorPostgres {
	return &ActorPostgres{db: db}
}

func (r *ActorPostgres) getFilmNames(id int) ([]DTO.FilmForActor, error) {
	var filmsID []int
	queryFilm := fmt.Sprintf("SELECT film_id FROM %s WHERE actor_id=$1", db.FILMACTOR)
	rows, err := r.db.Query(queryFilm, id)
	if err != nil {
		log.Panic(err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var filmId int
		rows.Scan(&filmId)
		filmsID = append(filmsID, filmId)
	}
	rows.Close()

	films := make([]DTO.FilmForActor, 0)
	for _, filmId := range filmsID {
		queryName := fmt.Sprintf("SELECT film_name, description, release_date, rating FROM %s WHERE id=$1", db.FILMS)
		rowsName := r.db.QueryRow(queryName, filmId)

		var filmName, description string
		var releaseDate time.Time
		var rating float64
		if err := rowsName.Scan(&filmName, &description, &releaseDate, &rating); err != nil {
			log.Panic(err)
			return nil, err
		}

		film := DTO.FilmForActor{
			Name:        filmName,
			Description: description,
			Date:        releaseDate.Format(db.PARSEDATE),
			Rating:      rating,
		}

		films = append(films, film)
	}
	return films, nil
}

func (r *ActorPostgres) GetAll() ([]DTO.ActorDTO, error) {

	query := fmt.Sprintf("SELECT id, actor_name, sex, b_date FROM %s", db.ACTORS)
	rows, err := r.db.Query(query)
	if err != nil {
		log.Panic(err)
		return nil, err
	}
	defer rows.Close()

	actorsDTO := make([]DTO.ActorDTO, 0)
	for rows.Next() {
		var id int
		var actorName, sex string
		var bDate time.Time
		if err := rows.Scan(&id, &actorName, &sex, &bDate); err != nil {
			log.Fatal(err)
			return nil, err
		}

		films, err := r.getFilmNames(id)
		if err != nil {
			log.Panic(err)
			return nil, err
		}

		act := DTO.ActorDTO{
			Name:  actorName,
			Sex:   sex,
			Date:  bDate.Format(db.PARSEDATE),
			Films: films,
		}
		actorsDTO = append(actorsDTO, act)
	}
	return actorsDTO, nil
}

func (r *ActorPostgres) GetOne(id int) (DTO.ActorDTO, error) {
	query := fmt.Sprintf("SELECT actor_name, sex, b_date FROM %s WHERE id=$1", db.ACTORS)
	var actorName, sex string
	var bDate time.Time
	if err := r.db.QueryRow(query, id).Scan(&actorName, &sex, &bDate); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return DTO.ActorDTO{}, errors.New("empty result")
		} else {
			log.Panic(err)
			return DTO.ActorDTO{}, err
		}

	}
	filmsNames, err := r.getFilmNames(id)
	if err != nil {
		log.Panic(err)
		return DTO.ActorDTO{}, err
	}

	actor := DTO.ActorDTO{
		Name:  actorName,
		Sex:   sex,
		Date:  bDate.Format(db.PARSEDATE),
		Films: filmsNames,
	}
	return actor, nil
}

func (r *ActorPostgres) Create(actor models.Actor) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (actor_name, sex, b_date) values ($1, $2, $3) RETURNING id", db.ACTORS)
	row := r.db.QueryRow(query, actor.Name, actor.Sex, actor.Date)
	if err := row.Scan(&id); err != nil {
		log.Panic(err)
		return 0, err
	}

	return id, nil
}

func (r *ActorPostgres) Update(id int, input DTO.ActorUpdate) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if input.Name != nil {
		setValues = append(setValues, fmt.Sprintf("actor_name=$%d", argId))
		args = append(args, *input.Name)
		argId++
	}
	if input.Sex != nil {
		setValues = append(setValues, fmt.Sprintf("sex=$%d", argId))
		args = append(args, *input.Sex)
		argId++
	}

	if input.Date != nil {
		setValues = append(setValues, fmt.Sprintf("b_date=$%d", argId))
		args = append(args, *input.Date)
		argId++
	}

	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf("UPDATE %s SET %s WHERE id=$%d", db.ACTORS, setQuery, argId)
	args = append(args, id)
	log.Println("query: " + query)

	if _, err := r.db.Exec(query, args...); err != nil {
		log.Panic(err)
		return err
	}

	return nil
}

func (r *ActorPostgres) Delete(id int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1", db.ACTORS)

	if _, err := r.db.Exec(query, id); err != nil {
		log.Panic(err)
		return err
	}
	return nil
}
