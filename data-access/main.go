package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/go-sql-driver/mysql"
)

var db *sql.DB

//Album Struct
type Album struct {
	ID     int64
	Title  string
	Artist string
	Price  float32
}

func main() {
	// CAPTURE CONNECTION PROP

	cfg := mysql.Config{
		User:   os.Getenv("DBUSER"),
		Passwd: os.Getenv("DBPASS"),
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
		DBName: "recordings",
	}

	//GET A DATABASE HANDLE
	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}

	fmt.Println("Connected !")

	albums, err := albumByArtist("John Coltrane")

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Album found: %v\n", albums)

	//GET ALBUM BY ID
	alb, err := albumByID(2)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Album found: %v\n", alb)

	//ADD ALBUM

	albID, err := addAlbum(Album{
		Title:  "Gaint",
		Artist: "Burna Boy",
		Price:  50.99,
	})

	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("ID of addes album: %v\n", albID)
}

// FUNC TO QUERY ALBUM
func albumByArtist(name string) ([]Album, error) {
	var albums []Album

	rows, err := db.Query(" SELECT * FROM album WHERE artist =?", name)
	if err != nil {
		return nil, fmt.Errorf(" album by Artist %q: %v", name, err)
	}

	// DEFER DELAY THE EXECUTION OF FUNC
	defer rows.Close()

	for rows.Next() {
		var alb Album
		if err := rows.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
			return nil, fmt.Errorf("albumByArtist %q: %v", name, err)
		}
		albums = append(albums, alb)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("albumByArtist %q: %v", name, err)
	}
	return albums, nil
}

// QUERY ALBUM BY ID
func albumByID(id int64) (Album, error) {
	var alb Album

	row := db.QueryRow(" SELECT * FROM album WHERE id = ?", id)
	if err := row.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
		if err == sql.ErrNoRows {
			return alb, fmt.Errorf("album by ID %d: no such album", id)
		}
		return alb, fmt.Errorf("album by id %d: %v", id, err)
	}
	return alb, nil
}

//ADD ALBUM INTO DB
func addAlbum(alb Album) (int64, error) {
	result, err := db.Exec("INSERT INTO ALBUM (title, artist, price) VALUES (?,?,?)", alb.Title, alb.Title, alb.Price)
	if err != nil {
		return 0, fmt.Errorf("addAlbum: %v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("addAlbum: %v", err)
	}

	return id, nil
}
