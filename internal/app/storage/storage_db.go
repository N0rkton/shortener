package storage

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"log"
)

type links struct {
	link  string
	short string
}
type DBStorage struct {
	db   *pgx.Conn
	path string
}

func Ping(path string) error {
	ctx := context.Background()
	conf, _ := pgx.ParseConfig(path)
	db, err := pgx.ConnectConfig(ctx, conf)
	if err != nil {
		return errors.New("unable to connect")
	}
	defer db.Close(ctx)
	err = db.Ping(ctx)
	if err != nil {
		return errors.New("unable to ping")
	}
	return nil
}
func NewDBStorage(path string) (Storage, error) {
	ctx := context.Background()
	conf, _ := pgx.ParseConfig(path)
	db, err := pgx.ConnectConfig(ctx, conf)
	if err != nil {
		return nil, errors.New("unable to connect")
	}
	defer db.Close(ctx)
	query := `CREATE TABLE IF NOT EXISTS links(id text, link text,  
    short text primary key)`
	_, err = db.Exec(ctx, query)
	if err != nil {
		log.Printf("Error %s when creating product table", err)
		return nil, errors.New("unable to create")
	}
	return &DBStorage{db: db, path: path}, nil
}
func (dbs *DBStorage) AddURL(id string, code string, url string) error {
	ctx := context.Background()
	conf, _ := pgx.ParseConfig(dbs.path)
	db, err := pgx.ConnectConfig(ctx, conf)
	if err != nil {
		return errors.New("unable to connect")
	}
	defer db.Close(ctx)
	dbs.db.Exec(ctx, "insert into links (id, link, short) values ($1, $2)", id, url, code)
	return nil
}
func (dbs *DBStorage) GetURL(url string) (string, error) {
	ctx := context.Background()
	conf, _ := pgx.ParseConfig(dbs.path)
	db, err := pgx.ConnectConfig(ctx, conf)
	if err != nil {
		return "", errors.New("unable to connect")
	}
	defer db.Close(ctx)
	rows := dbs.db.QueryRow(ctx, "select link from links where short=$1 limit 1", url)
	var link string
	rows.Scan(&link)
	return link, nil
}

func (dbs *DBStorage) GetURLByID(id string) (map[string]string, error) {
	ctx := context.Background()
	conf, _ := pgx.ParseConfig(dbs.path)
	db, err := pgx.ConnectConfig(ctx, conf)
	if err != nil {
		return nil, errors.New("unable to connect")
	}
	defer db.Close(ctx)
	resp := make(map[string]string)
	rows, err := dbs.db.Query(ctx, "SELECT link, short from links where id=$1", id)
	if err != nil {
		return nil, errors.New("not found")
	}
	defer rows.Close()
	for rows.Next() {
		var v links
		err = rows.Scan(&v.link, &v.short)
		if err != nil {
			return nil, err
		}
		resp[v.short] = v.link
	}
	return resp, nil
}
