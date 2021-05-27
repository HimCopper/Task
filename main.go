package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"net/smtp"

				
	_ "github.com/go-sql-driver/mysql"
)

type user struct {
	Name     string
	Dob      string
	Email    string
	Phone_no string
}

var tpl *template.Template
var tpl1 *template.Template
var t *template.Template

var form *template.Template



func init() {

	tpl = template.Must(template.ParseFiles("index.gohtml"))

	t = template.Must(template.ParseFiles("registration.gohtml"))
	form = template.Must(template.ParseFiles("submitform.gohtml"))

}

func main() {
	handleRequest()
}

func handleRequest() {

	http.HandleFunc("/", index)
	http.HandleFunc("/r", registration)
	http.HandleFunc("/form", submitform)
	http.ListenAndServe(":8080", nil)
}

func index(w http.ResponseWriter, r *http.Request) {

	tpl.ExecuteTemplate(w, "index.gohtml", nil)
	fmt.Println("index Page")
}

func confiq() (db *sql.DB) {
	db, err := sql.Open("mysql", "root:@tcp(localhost:3306)/test")
	if err != nil {
		fmt.Println(err)
	}

     	return db
}

func registration(w http.ResponseWriter, r *http.Request) {

	db := confiq()

	var u user
	if r.Method == http.MethodPost {

		Name := r.FormValue("name")
		fmt.Println(Name)
		Dob := r.FormValue("dob")
		fmt.Println(Dob)
		Email := r.FormValue("email")
		fmt.Println(Email)
		Phone_no := r.FormValue("phone")
		fmt.Println(Phone_no)
		u = user{Name, Dob, Email, Phone_no}

		insert, err := db.Prepare("insert into form(Name , Dob, Email ,Phone_no) values(?,?,?,?)")
				
		fmt.Println(insert)

		if err != nil {
			panic(err.Error())
		}
		insert.Exec(Name, Dob, Email, Phone_no)

		defer db.Close()
	}

	t.ExecuteTemplate(w, "registration.gohtml", u)
	// 	Sender Data
	from := "kapoorhimanshu.176@gmail.com"
	password := "1762595e9"

	// Receiver Data
	to := []string{
		u.Email,
	}

	// smtp server configuration.
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	// Message.
	message := []byte("Your Form Is Submiited Successfully.")

	// Authentication.
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Sending email.
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Email Sent Successfully!")

	http.Redirect(w, r, "/form", 404)

}

func submitform(w http.ResponseWriter, r *http.Request) {

	db := confiq()
	sel, err := db.Query("select * from form")

	if err != nil {
		fmt.Println("error in selection")
		fmt.Println(err)
	}

	for sel.Next() {

		var u user
		err = sel.Scan(&u.Name, &u.Dob, &u.Email, &u.Phone_no)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(u.Name, u.Dob, u.Email, u.Phone_no)
	
	form.ExecuteTemplate(w,"submitform.gohtml",u)
	

	}

}
