package repository

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
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

func (r *ActorPostgres) getFilmNames(id string) ([]string, error) {

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

	filmsNames := make([]string, 0)
	for _, filmId := range filmsID {
		queryName := fmt.Sprintf("SELECT film_name FROM %s WHERE id=$1", db.FILMS)
		rowsName, err := r.db.Query(queryName, filmId)
		if err != nil {
			log.Panic(err)
			return nil, err
		}
		defer rowsName.Close()
		var fName string
		rowsName.Scan(&fName)
		filmsNames = append(filmsNames, fName)
	}
	return filmsNames, nil
}

func (r *ActorPostgres) GetAllActors() ([]DTO.ActorDTO, error) {
	var actorsDTO []DTO.ActorDTO
	query := fmt.Sprintf("SELECT id, actor_name, sex, b_date FROM %s", db.ACTORS)
	rows, err := r.db.Query(query)
	if err != nil {
		log.Panic(err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var id, actorName, sex string
		var bDate time.Time
		if err := rows.Scan(&id, &actorName, &sex, &bDate); err != nil {
			log.Fatal(err)
			return nil, err
		}

		filmsNames, err := r.getFilmNames(id)
		if err != nil {
			log.Panic(err)
			return nil, err
		}

		act := DTO.ActorDTO{
			Name:  actorName,
			Sex:   sex,
			Date:  bDate.Format(db.PARSEDATE),
			Films: filmsNames,
		}
		actorsDTO = append(actorsDTO, act)
	}
	return actorsDTO, nil
}

func (r *ActorPostgres) GetActor(id string) (DTO.ActorDTO, error) {
	query := fmt.Sprintf("SELECT actor_name, sex, b_date FROM %s WHERE id=$1", db.ACTORS)
	var actorName, sex string
	var bDate time.Time
	idn, _ := strconv.Atoi(id)
	if err := r.db.QueryRow(query, idn).Scan(&actorName, &sex, &bDate); err != nil {
		log.Panic(err)
		return DTO.ActorDTO{}, err
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

func (r *ActorPostgres) CreateActor(actor models.Actor) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (actor_name, sex, b_date) values ($1, $2, $3) RETURNING id", db.ACTORS)
	row := r.db.QueryRow(query, actor.Name, actor.Sex, actor.Date)
	if err := row.Scan(&id); err != nil {
		log.Panic(err)
		return 0, err
	}

	return id, nil
}
