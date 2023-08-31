package database

import (
	"context"
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"platzi.com/go/api/rest-ws/models"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(url string) (*PostgresRepository, error) {

	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}

	return &PostgresRepository{db}, nil
}

func (repo *PostgresRepository) InsertUser(ctx context.Context, user *models.User) error {
	_, err := repo.db.ExecContext(ctx, "INSERT INTO user (email, password) VALUES ($1, $2)", user.Email, user.Password)
	return err
}

func (repo *PostgresRepository) GetUserById(ctx context.Context, id string) (*models.User, error) {
	rows, err := repo.db.QueryContext(ctx, "SELECT id, email FROM users WHERE id = $1", id)
	if err != nil {
		log.Fatal(err.Error())
		return nil, err
	}
	defer rows.Close()

	var user = models.User{}
	for rows.Next() {
		if err := rows.Scan(&user.Id, &user.Email); err != nil {
			return &user, nil
		}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &user, nil

	//Nota:  esta parte puedo mejorar, tomar como ejemplo la forma en lo que lo hacen Javier en vedana

}

func (repo *PostgresRepository) Close() error {
	return repo.db.Close()
}
