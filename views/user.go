package views

import (
	"log"
	"net/http"

	"github.com/alexanderi96/leafnet/db"
	"github.com/alexanderi96/leafnet/types"
	"github.com/alexanderi96/leafnet/utils"
)

func UserPage(w http.ResponseWriter, r *http.Request) {
	prepareContext(w, r)

	if r.Method == "POST" {
		r.ParseForm()
		//TODO: implement user update
		http.Redirect(w, r, "/my-profile", http.StatusFound)
	} else if r.Method == "GET" {
		log.Println("Viewing: ", c.User.Email)
		if err := templates["user_profile"].Execute(w, c); err != nil {
			WriteError(w, err)
			return
		}
	}

}

func DeleteMyAccount(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/my-profile", http.StatusBadRequest)
		return
	}

	r.ParseForm()
	log.Println("Attempting to delete user: ", c.User)
	if err := db.DeleteSelectedUser(c.User.Email, c.User.Password); err != nil {
		WriteError(w, err)
		return
	}

	log.Println("Account deleted, redirecting...")
	http.Redirect(w, r, "/logout", http.StatusAccepted)
}

func parseUser(r *http.Request) (types.User, error) {
	r.ParseForm()

	hashedPwd, err := utils.EncryptStr(r.Form.Get("password"))
	if err != nil {
		return types.User{}, err
	}

	u := types.User{
		UserName: r.Form.Get("user_name"),
		Email:    r.Form.Get("email"),
		Password: hashedPwd,
		Person:   r.Form.Get("person"),
	}

	return u, nil
}

func checkIfUserExists(u types.User) (bool, error) {
	if email_user, err := db.GetUserInfoByEmail(u.Email); err != nil {
		return false, err
	} else if (email_user != types.User{}) {
		return true, nil
	}

	if name_user, err := db.GetUserInfoByUserName(u.UserName); err != nil {
		return false, err
	} else if (name_user != types.User{}) {
		return true, nil
	}

	return false, nil
}
