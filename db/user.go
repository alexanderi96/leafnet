package db

import (
	"errors"
	"log"
	"github.com/alexanderi96/leafnet/types"
)

//CreateUser will create a new user. It takes as input a generic user type and optionally the UserId
func CreateUser(u types.User) (e error) {
	dCice := "insert into Utenti(NomeUtente, CognomeUtente, SessoUtente, DataNascitaUtente, EmailUtente, PasswordUtente) values(?,?,?,?,?,?)"
	dCiceId := "insert into Utenti(IdUtente, NomeUtente, CognomeUtente, SessoUtente, DataNascitaUtente, EmailUtente, PasswordUtente) values(?,?,?,?,?,?,?)"
	dGlobe := "insert into Utenti(NomeUtente, CognomeUtente, SessoUtente, DataNascitaUtente, EmailUtente, PasswordUtente) values(?,?,?,?,?,?)"
	dGlobeId := "insert into Utenti(IdUtente, NomeUtente, CognomeUtente, SessoUtente, DataNascitaUtente, EmailUtente, PasswordUtente) values(?,?,?,?,?,?,?)"

	switch u := u.(type) {
	default:
		return errors.New("Undefined user type")
	case types.Cicerone:
		if u.IdUtente == 0 {
			e = gQuery(dCice, u.Nome, u.Cognome, u.Sesso, u.DataNascita, u.Email, u.Password)
			if e == nil {
				u.IdUtente = GetUserID(u.Email)
			} else {
				return
			}
		} else if checkIdAviability(u.IdUtente) {
			e = gQuery(dCiceId, u.IdUtente, u.Nome, u.Cognome, u.Sesso, u.DataNascita, u.Email, u.Password)
		} else {
			return errors.New("Unavailable user id")
		}

		if e != nil {
			return
		} else {
			e = gQuery("insert into Ciceroni(IdCicerone, CodiceFiscaleCicerone, TelefonoCicerone, IbanCicerone) values (?,?,?,?)", u.IdUtente, u.Tel, u.Iban, u.CodFis)
		}
		
		
	case types.Globetrotter:
		if u.IdUtente == 0 {
			e = gQuery(dGlobe, u.Nome, u.Cognome, u.Sesso, u.DataNascita, u.Email, u.Password)
		} else if checkIdAviability(u.IdUtente) {
			e = gQuery(dGlobeId, u.IdUtente, u.Nome, u.Cognome, u.Sesso, u.DataNascita, u.Email, u.Password)
		} else {
			return errors.New("Unavailable user id")
		}
	}
	return
}

func checkIdAviability(uid int) (false bool) {
	var idCount int
	userSQL := "select count(IdUtente) from Utenti left join Ciceroni on Utenti.IdUtente = Ciceroni.IdCicerone where IdUtente = ?"
	log.Println("Checking if the id is already taken")
	rows := database.query(userSQL, uid)

	defer rows.Close()
	if rows.Next() {
		e := rows.Scan(&idCount)
		if e != nil {
			return
		} else if idCount != 0 {
			return
		}
	}
	return true
}

//ValidUser will check if the user exists in db and if exists if the username password
//combination is valid
func ValidUser(email, password string) bool {
	var passwordFromDB string
	//we don't want to validate a backed up user (id < 0)
	plainSQL := "select PasswordUtente from Utenti where (EmailUtente = ? and IdUtente > 0)"
	log.Print("validating user ", email)
	rows := database.query(plainSQL, email)

	defer rows.Close()
	if rows.Next() {
		err := rows.Scan(&passwordFromDB)
		if err != nil {
			return false
		}
		//If the password matches, return true
		if password == passwordFromDB {
			return true
		}
	}
	
	//by default return false
	return false
}

//GetUserID will get the user's ID from the database
func GetUserID(email string) (userID int) {
	userSQL := "select IdUtente from Utenti where EmailUtente = ?"
	rows := database.query(userSQL, email)

	defer rows.Close()
	if rows.Next() {
		if err := rows.Scan(&userID); err != nil {
			log.Println(err)
		}
	}
	return
}

func GetUserEmail(uid int) (email string) {
	userSQL := "select EmailUtente from Utenti where IdUtente = ?"
	rows := database.query(userSQL, uid)

	defer rows.Close()
	if rows.Next() {
		if err := rows.Scan(&email); err != nil {
			log.Println(err)
		}
	}
	return
}

func GetUserInfo(email string) (u types.User, e error){
	if IsCicerone(GetUserID(email)) {
		user := types.Cicerone{}
		userSQL := "select IdUtente, NomeUtente, CognomeUtente, SessoUtente, DataNascitaUtente, EmailUtente, CodiceFiscaleCicerone, TelefonoCicerone, IbanCicerone from Utenti join Ciceroni on Utenti.IdUtente = Ciceroni.IdCicerone where EmailUtente = ?"
		rows := database.query(userSQL, email)
		defer rows.Close() //must defer after every database interaction

		if rows.Next() {
			log.Println(rows)
			if e := rows.Scan(&user.IdUtente, &user.Nome, &user.Cognome, &user.Sesso, &user.DataNascita, &user.Email, &user.Tel, &user.Iban, &user.CodFis); e != nil {
				return nil, e
			}

			return user, nil
		}

	} else {
		user := types.Globetrotter{}
		userSQL := "select IdUtente, NomeUtente, CognomeUtente, SessoUtente, DataNascitaUtente, EmailUtente from Utenti where EmailUtente = ?"
		rows := database.query(userSQL, email)
		defer rows.Close() //must defer after every database interaction

		if rows.Next() {
			if e := rows.Scan(&user.IdUtente, &user.Nome, &user.Cognome, &user.Sesso, &user.DataNascita, &user.Email); e != nil {
				return nil, e
			}
		}
		return user, nil
	}
	return
}

func AddCicerone(uid, tel int, iban, fcode string) error {
	err := gQuery("insert into Ciceroni(IdCicerone, CodiceFiscaleCicerone, TelefonoCicerone, IbanCicerone) values (?,?,?,?)", uid, tel, iban, fcode)
	return err
}

func IsCicerone(uid int) (false bool) {
	userSQL := "select count(IdCicerone) from Ciceroni where IdCicerone = ?"
	rows := database.query(userSQL, uid)

	defer rows.Close()
	var nr int
	if rows.Next() {
		err := rows.Scan(&nr)
		if err != nil {
			log.Println(err)
			return
		} else if nr < 1 {
			return 
		}
		return true
	}
	return
}

func DeleteSelectedUser(email, password string) (e error) {
	if ValidUser(email, password) {
		uid := GetUserID(email)
		if IsCicerone(uid) {
			ciceSQL := "delete from Ciceroni where IdCicerone = ?"
			if e = gQuery(ciceSQL, uid); e != nil {
				return
			}
		}
		globeSQL := "delete from Utenti where EmailUtente = ?"
		e = gQuery(globeSQL, email)
	} else {
		e = errors.New("Invalid User")
	}
	return
}

func DeleteUserById(uid int) (e error) {
	if IsCicerone(uid) {
		ciceSQL := "delete from Ciceroni where IdCicerone = ?"
		if e = gQuery(ciceSQL, uid); e != nil {
			return
		}
	}
	globeSQL := "delete from Utenti where IdUtente = ?"
	e = gQuery(globeSQL, uid)
	return 
}

func InvertUserId(uid int) (e error) {
	//let's delete the old backup (if there is one)
	DeleteUserById(uid * -1)
	updateSQL := " update Utenti set IdUtente = ? where IdUtente = ?"
	if e = gQuery(updateSQL, uid * -1, uid); e != nil && IsCicerone(uid) {
		updateSQL := "update Ciceroni set IdCicerone = ? where IdCicerone = ?"
		e = gQuery(updateSQL, uid * -1, uid)
	}
	return
}