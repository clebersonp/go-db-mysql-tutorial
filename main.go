package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/go-sql-driver/mysql"
)

type Album struct {
	ID int64
	Title string
	Artist string
	Price float32
}

var db *sql.DB

func main() {
	// Capture connection properties.
	cfg := mysql.Config{
		User: os.Getenv("DBUSER"),
		Passwd: os.Getenv("DBPASS"),
		Net: "tcp",
		Addr: "localhost:3306",
		DBName: "recordings",
	}
	// Get a database handle.
	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Open might not immediately connect, depending on the driver.
	// You’re using Ping here to confirm that the database/sql package can connect when it needs to.
	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connected!")

	albums, err := albumsByArtist("John Coltrane")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Albums found: %v\n", albums)

	alb, err := albumByID(1)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Album found: %v\n", alb)

	albID, err := addAlbum(Album{
		Title:  "The Modern Sound of Betty Carter",
		Artist: "Betty Carter",
		Price:  49.99,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("ID of added album: %v\n", albID)

	albums, err = getAllAlbums()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("All Albums found: %v\n", albums)

	countAffected, err := updateAlbum(Album{
		Title:  "The Modern Sound of Bee",
		Artist: "Betty Cartner",
		Price:  45.99,
	}, albID)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%d rows affected by update by ID: %d\n", countAffected, albID)

	albums, err = getAllAlbums()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("All Albums found: %v\n", albums)

	countAffected, err = deleteAlbum(albID)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%d row affected by delete by ID: %d\n", countAffected, albID)

	albums, err = getAllAlbums()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("All Albums found: %v\n", albums)
}

// albumsByArtist queries for albums that have the specified artist name.
func albumsByArtist(name string) ([]Album, error) {
	// An albums slice to hold data from returned rows.
	var albums []Album

	rows, err := db.Query("SELECT * FROM album WHERE artist = ?", name)
	if err != nil {
		return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
	}
	defer rows.Close()
	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var alb Album
		if err := rows.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
			return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
		}
		albums = append(albums, alb)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
	}
	return albums, nil
}

// albums queries for all albums
func getAllAlbums() ([]Album, error) {
	// An albums slice to hod data from returned rows.
	var albums []Album

	rows, err := db.Query("SELECT * FROM album")
	if err != nil {
		return nil, fmt.Errorf("albums: %v", err)
	}
	defer rows.Close()
	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var alb Album
		if err := rows.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
			return nil, fmt.Errorf("albums: %v", err)
		}
		albums = append(albums, alb)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("albums: %v", err)
	}
	return albums, nil
}


// albumByID queries for the album with the specified ID.
func albumByID(id int64) (Album, error) {
	// An album to hold data from the returned row.
	var alb Album

	row := db.QueryRow("SELECT * FROM album WHERE id = ?", id)
	
	if err := row.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
		if err == sql.ErrNoRows {
			return alb, fmt.Errorf("albumByID %d: no such album", id)
		}
		return alb, fmt.Errorf("albumByID %d: %v", id, err)
	}
	return alb, nil
}

// addAlbum adds the specified album to the database,
// returning the album ID of the new entry
func addAlbum(alb Album) (int64, error) {
	result, err := db.Exec("INSERT INTO album (title, artist, price) VALUES (?, ?, ?)", alb.Title, alb.Artist, alb.Price)
	if err != nil {
		return 0, fmt.Errorf("addAlbum: %v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("addAlbum: %v", err)
	}
	return id, nil
}

// updateAlbum updates the specified album to the database,
// returning the number of affected rows by the update
func updateAlbum(alb Album, id int64) (int64, error) {
	result, err := db.Exec("UPDATE album SET title = ?, artist = ?, price = ? WHERE id = ?", alb.Title, alb.Artist, alb.Price, id)
	if err != nil {
		return 0, fmt.Errorf("updateAlbum id %d: %v", id, err)
	}
	count, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("updateAlbum id %d: %v", id, err)
	}
	return count, nil
}

// deleteAlbum deletes the specified album to the database,
// returning the number of affected rows by the delete
func deleteAlbum(id int64) (int64, error) {
	result, err := db.Exec("DELETE FROM album WHERE id = ?", id)
	if err != nil {
		return 0, fmt.Errorf("deleteAlbum id %d: %v", id, err)
	}
	count, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("deleteAlbum id %d: %v", id, err)
	}
	return count, nil
}

