package postgres

import (    
    "database/sql"
    _ "github.com/lib/pq"
    
    "os"
    "io"
    "bufio"
    "log"
)

type Database struct {
    // a database handle representing a pool of zero or more connections
    // safe for concurrent use
    dbHandle *sql.DB 
}

func ConnectDB(connStr string) (*Database, error) {  
    pg := &Database{}
    err := pg.OpenDB("postgres", connStr)
    return pg, err
}

func (database *Database) OpenDB(driverName string, connStr string) error {
	db, err := sql.Open(driverName, connStr)
    if err == nil {
        database.dbHandle = db
    }
	return err
}

func (database *Database) CreateTables(schemaFile string) error {
    // Read table schema file
    file, err := os.Open(schemaFile)
    if err != nil {
        log.Println(err)
        return err
    }
    defer file.Close()
    
    reader := bufio.NewReader(file)
    
    var line string
    for {
        lineBytes, _, err := reader.ReadLine()
        line = string(lineBytes)
        
        if err != nil {
            break
        } else {
            // Process line
            _, err = database.dbHandle.Exec("CREATE TABLE IF NOT EXISTS " + line)
            if err != nil {
                return err
            }
        }
    }
    
    if err != io.EOF {
        return err
    }
    
    return err
}

func (database *Database) Handle() *sql.DB {
    return database.dbHandle
}