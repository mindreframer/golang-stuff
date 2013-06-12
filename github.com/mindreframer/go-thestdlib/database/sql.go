package main

import (
    "database/sql"
    "flag"
    _ "github.com/mattn/go-sqlite3"
    // An example if you want to use Postgres
    // _ "github.com/bmizerany/pq"
    "log"
)

var rollback = flag.Bool("rollback", false, "Rollback in the insert transaction")

type Show struct {
    Name, Country string
}

func openSqlite() (*sql.DB, err) {
    return sql.Open("sqlite3", "go-thestdlib.db")
}

func openPostgres() (*sql.DB, err) {
    return sql.Open("postgres", "user=bob password=secret host=1.2.3.4 port=5432 dbname=mydb sslmode=verify-full")
}

func openDB() *sql.DB {
    db, err := openSqlite()
    // db, err := openPostgres()
    if err != nil {
        log.Fatalf("failed opening database: %s", err)
    }
    return db
}

func removeTable(db *sql.DB) {
    _, err := db.Exec("DROP TABLE IF EXISTS shows")
    if err != nil {
        log.Fatalf("failed dropping table: %s", err)
    } else {
        log.Println("dropped table (if it existed) shows")
    }
}

func createTable(db *sql.DB) {
    _, err := db.Exec("CREATE TABLE shows (name TEXT, country TEXT)")
    if err != nil {
        log.Fatalf("failed creating table: %s", err)
    } else {
        log.Println("created table shows")
    }
}

func insertRow(db *sql.DB) {
    // For postgres we use $1 and $2 instead of ?
    res, err := db.Exec("INSERT INTO shows (name, country) VALUES (?, ?)", "Nöjesmaskinen", "SE")
    if err != nil {
        log.Fatalf("failed inserting Swedish show: %s", err)
    } else {
        log.Println("inserted 1 Swedish TV show")
    }

    if id, err := res.LastInsertId(); err != nil {
        log.Printf("failed retrieving LastInsertId: %s", err)
    } else {
        log.Printf("LastInsertId: %d", id)
    }

    if n, err := res.RowsAffected(); err != nil {
        log.Printf("failed retrieving RowsAffected: %s", err)
    } else {
        log.Printf("RowsAffected: %d", n)
    }
}

func insertRows(db *sql.DB) {
    tx, err := db.Begin()
    if err != nil {
        log.Fatalf("failed starting transaction: %s", err)
    }

    shows := []Show{
        Show{"Top Gear", "UK"},
        Show{"Wilfred", "AU"},
        Show{"Top Gear", "US"},
        Show{"Arctic Air", "CA"},
    }

    stmt, err := tx.Prepare("INSERT INTO shows (name, country) VALUES (?, ?)")
    if err != nil {
        log.Fatalf("failed preparing statement: %s", err)
    }

    for _, show := range shows {
        _, err := stmt.Exec(show.Name, show.Country)
        if err != nil {
            log.Fatalf("failed insert show %s (%s): %s", show.Name, show.Country, err)
        } else {
            log.Printf("inserted show %#v for country %#v", show.Name, show.Country)
        }
    }

    if *rollback {
        if err := tx.Rollback(); err != nil {
            log.Fatalf("failed rolling back transaction: %s", err)
        } else {
            log.Println("rolled back transaction, nothing inserted")
        }
    } else {
        if err := tx.Commit(); err != nil {
            log.Fatalf("failed committing transaction: %s", err)
        } else {
            log.Println("committed transaction, 4 new shows added")
        }
    }
}

func queryCount(db *sql.DB) {
    row := db.QueryRow("SELECT COUNT(*) FROM shows")
    var count int
    if err := row.Scan(&count); err != nil {
        log.Fatalf("failed getting count: %s", err)
    }
    log.Printf("there are %d TV shows in the database", count)
}

func queryRow(db *sql.DB) {
    row := db.QueryRow("SELECT * FROM shows WHERE country = ? LIMIT 1", "CA")
    show := Show{}
    if err := row.Scan(&show.Name, &show.Country); err != nil {
        log.Printf("failed scanning single row: %s", err)
    } else {
        log.Printf("Found 1 %s TV show: %s", show.Country, show.Name)
    }
}

func queryRows(db *sql.DB) {
    name := "Top Gear"
    rows, err := db.Query("SELECT * FROM shows WHERE name = ?", name)
    if err != nil {
        log.Fatalf("failed querying multiple rows: %s", err)
    }
    shows := make([]Show, 0)
    for rows.Next() {
        show := Show{}
        if err := rows.Scan(&show.Name, &show.Country); err != nil {
            log.Fatalf("failed scanning row: %s", err)
        }
        shows = append(shows, show)
    }
    log.Printf("found %d shows named %#v", len(shows), name)
    for _, show := range shows {
        log.Printf("\t...in country %s", show.Country)
    }
    if err := rows.Err(); err != nil {
        log.Fatalf("got unexpected error during iteration: %s", err)
    }
}

func deleteRows(db *sql.DB) {
    _, err := db.Exec("DELETE FROM shows")
    if err != nil {
        log.Fatalf("failed deleting rows: %s", err)
    }
}

func main() {
    flag.Parse()
    db := openDB()
    defer db.Close()

    removeTable(db)
    createTable(db)
    insertRow(db)
    insertRows(db)
    queryCount(db)
    queryRow(db)
    // Sleep here...
    queryRows(db)
    deleteRows(db)
}
