package storage

import (
	"ucouturier.com/JOW_scrapper/models"
	"ucouturier.com/JOW_scrapper/db"
    "fmt"
	"database/sql"
    "errors"
)

func SaveUserToDB(user *models.User) error {

    exists, err := UserVerifyIfExistsInDB(user.Email)

    if err != nil {
        fmt.Println(err)
        return err
    }

    if exists {
        fmt.Println("[ERROR] Un utilisateur avec une email similaire existe déjà, utilisateur non créé")
        return errors.New("[ERROR]")
    }

    stmt, err := db.DB.Prepare(`
        INSERT INTO users (name, email, password) VALUES (?, ?, ?)
    `)
    if err != nil {
        return err
    }
    defer stmt.Close()

    res, err := stmt.Exec(user.Name, user.Email, user.Password)
    if err != nil {
        return err
    }

    _, err = res.LastInsertId()
    if err != nil {
        return err
    }

    return nil
}

func UserVerifyIfExistsInDB(email string) (bool, error) {
    row := db.DB.QueryRow(`SELECT id FROM users WHERE email = ?`, email)
    var id int
    err := row.Scan(&id)
    if err != nil {
        if err == sql.ErrNoRows {
            return false, nil
        }
        fmt.Println(err)
        return true, errors.New("[ERROR] Impossible de lire l'utilisateur")
    }
    return true, nil
}

func GetUserFromDBByEmail(email string) (*models.User, error) {
    row := db.DB.QueryRow(`SELECT id,name,email,password FROM users WHERE email = ?`, email)
    var user models.User
    err := row.Scan(&user.ID,&user.Name, &user.Email, &user.Password)
    if err != nil {
        return nil, err
    }
    return &user, nil
}