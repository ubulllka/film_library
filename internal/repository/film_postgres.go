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

type FilmPostgres struct {
	db *sql.DB
}

func NewFilmPostgres(db *sql.DB) *FilmPostgres {
	return &FilmPostgres{db: db}
}

func (r *FilmPostgres) getActorsInFilm(id int) ([]DTO.ActorInput, error) {
	var actorsID []int
	queryActors := fmt.Sprintf("SELECT actor_id FROM %s WHERE film_id=$1", db.FILMACTOR)
	rows, err := r.db.Query(queryActors, id)
	if err != nil {
		log.Panic(err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var actorId int
		rows.Scan(&actorId)
		actorsID = append(actorsID, actorId)
	}
	rows.Close()

	actors := make([]DTO.ActorInput, 0)
	for _, actorId := range actorsID {
		query := fmt.Sprintf("SELECT actor_name, sex, b_date FROM %s WHERE id=$1", db.ACTORS)
		rowsActors := r.db.QueryRow(query, actorId)
		if err != nil {
			log.Panic(err)
			return nil, err
		}

		var actorName, sex string
		var bDate time.Time

		if err := rowsActors.Scan(&actorName, &sex, &bDate); err != nil {
			log.Panic(err)
			return nil, err
		}

		actor := DTO.ActorInput{
			Name: actorName,
			Sex:  sex,
			Date: bDate.Format(db.PARSEDATE),
		}
		actors = append(actors, actor)
	}
	return actors, nil
}

func (r *FilmPostgres) GetAll(column, order string) ([]DTO.FilmDTO, error) {

	var query strings.Builder
	query.WriteString(fmt.Sprintf("SELECT id, film_name, description, release_date, rating FROM %s", db.FILMS))

	if column == "film_name" {
		query.WriteString(" ORDER BY film_name")
	} else if column == "release_date" {
		query.WriteString(" ORDER BY release_date")
	} else {
		query.WriteString(" ORDER BY rating")
	}

	if order == "ASC" {
		query.WriteString(" ASC")
	} else {
		query.WriteString(" DESC")
	}

	log.Println(query.String())
	rows, err := r.db.Query(query.String())
	if err != nil {
		log.Panic(err)
		return nil, err
	}
	defer rows.Close()

	filmsDTO := make([]DTO.FilmDTO, 0)

	for rows.Next() {

		var id int
		var filmName, description string
		var releaseDate time.Time
		var rating float64

		if err := rows.Scan(&id, &filmName, &description, &releaseDate, &rating); err != nil {
			log.Fatal(err)
			return nil, err
		}

		actors, err := r.getActorsInFilm(id)
		if err != nil {
			log.Panic(err)
			return nil, err
		}

		filmDTO := DTO.FilmDTO{
			Name:        filmName,
			Description: description,
			Date:        releaseDate.Format(db.PARSEDATE),
			Rating:      rating,
			Actors:      actors,
		}
		filmsDTO = append(filmsDTO, filmDTO)
	}
	return filmsDTO, nil
}

func (r *FilmPostgres) GetOne(id int) (DTO.FilmDTO, error) {
	query := fmt.Sprintf("SELECT film_name, description, release_date, rating FROM %s WHERE id=$1", db.FILMS)

	var filmName, description string
	var releaseDate time.Time
	var rating float64

	if err := r.db.QueryRow(query, id).Scan(&filmName, &description, &releaseDate, &rating); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return DTO.FilmDTO{}, errors.New("empty result")
		} else {
			log.Panic(err)
			return DTO.FilmDTO{}, err
		}
	}

	actors, err := r.getActorsInFilm(id)
	if err != nil {
		log.Panic(err)
		return DTO.FilmDTO{}, err
	}

	film := DTO.FilmDTO{
		Name:        filmName,
		Description: description,
		Date:        releaseDate.Format(db.PARSEDATE),
		Rating:      rating,
		Actors:      actors,
	}
	return film, nil
}

func (r *FilmPostgres) Create(film models.Film, actorsId []int) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		log.Fatal(err)
		return 0, err
	}

	var filmId int
	query := fmt.Sprintf("INSERT INTO %s (film_name, description, release_date, rating) values ($1, $2, $3, $4) RETURNING id", db.FILMS)
	row := tx.QueryRow(query, film.Name, film.Description, film.Date, film.Rating)
	if err := row.Scan(&filmId); err != nil {
		log.Panic(err)
		tx.Rollback()
		return 0, err
	}

	for _, aId := range actorsId {
		createQuery := fmt.Sprintf("INSERT INTO %s (film_id, actor_id) values ($1, $2)", db.FILMACTOR)
		if _, err := tx.Exec(createQuery, filmId, aId); err != nil {
			log.Panic(err)
			tx.Rollback()
			return 0, err
		}
	}

	return filmId, tx.Commit()
}

func (r *FilmPostgres) Update(id int, input DTO.FilmUpdate) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if input.Name != nil {
		setValues = append(setValues, fmt.Sprintf("film_name=$%d", argId))
		args = append(args, *input.Name)
		argId++
	}

	if input.Description != nil {
		setValues = append(setValues, fmt.Sprintf("description=$%d", argId))
		args = append(args, *input.Description)
		argId++
	}

	if input.Date != nil {
		setValues = append(setValues, fmt.Sprintf("release_date=$%d", argId))
		args = append(args, *input.Date)
		argId++
	}

	if input.Rating != nil {
		setValues = append(setValues, fmt.Sprintf("rating=$%d", argId))
		args = append(args, *input.Rating)
		argId++
	}

	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf("UPDATE %s SET %s WHERE id=$%d", db.FILMS, setQuery, argId)
	args = append(args, id)
	log.Println("query: " + query)

	tx, err := r.db.Begin()
	if err != nil {
		log.Fatal(err)
		return err
	}

	if input.Actors != nil {
		queryDelete := fmt.Sprintf("DELETE FROM %s WHERE film_id=$1", db.FILMACTOR)
		if _, err := tx.Exec(queryDelete, id); err != nil {
			log.Panic(err)
			tx.Rollback()
			return err
		}
		queryCreate := fmt.Sprintf("INSERT INTO %s (film_id, actor_id) values ($1, $2)", db.FILMACTOR)
		for _, actorId := range input.Actors {
			if _, err := tx.Exec(queryCreate, id, actorId); err != nil {
				log.Panic(err)
				tx.Rollback()
				return err
			}
		}
	}

	if argId > 1 {
		if _, err := tx.Exec(query, args...); err != nil {
			log.Panic(err)
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (r *FilmPostgres) Delete(id int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1", db.FILMS)

	if _, err := r.db.Exec(query, id); err != nil {
		log.Panic(err)
		return err
	}
	return nil
}
