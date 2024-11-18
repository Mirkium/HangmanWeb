package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"text/template"
)

type Profil struct {
	Name     string `json:"name"`
	Score    int    `json:"score"`
	NbGame   int    `json:"nb_game"`
	Password string `json:"password"`
}

type User struct {
	Name    string
	Score   int
	NbGame  int
	Connect bool
}

type Error struct {
	IfError bool
	Error   string
}

var ValueError Error

var Connect User

var FirstConnect bool

func Menu() {

	//======================================================MENU================================================================

	http.HandleFunc("/hangman_web/menu", func(w http.ResponseWriter, r *http.Request) {
		temp, err := template.ParseFiles("./template/Menu.html")
		if err != nil {
			fmt.Print("Tu as une erreur de type : ", err)
			return
		}
		if FirstConnect {
			temp.Execute(w, nil)
		} else {
			temp.Execute(w, Connect)
		}
	})

	//==========================================================================================================================

	//======================================================CHOOSE TXT===========================================================

	http.HandleFunc("/hangman_web/treatment/txt", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Redirect(w, r, "/error", http.StatusSeeOther)
			return
		}
		if FirstConnect {
			http.Redirect(w, r, "/hangman_web/connect", http.StatusSeeOther)
			return
		} else {
			Texte = r.FormValue("name")
			http.Redirect(w, r, "/hangman_web/game", http.StatusSeeOther)
			return
		}
	})

	//=========================================================PROFIL===========================================================

	http.HandleFunc("/hangman_web/profil", func(w http.ResponseWriter, r *http.Request) {
		temp, err := template.ParseFiles("./template/Profil.html")
		if err != nil {
			fmt.Print("Tu as une erreur de type : ", err)
			return
		}
		if FirstConnect {
			http.Redirect(w, r, "/hangman_web/connect", http.StatusSeeOther)
			return
		} else {
			temp.Execute(w, Connect)
			return
		}
	})

	//==========================================================================================================================

	//======================================================CONNECT==============================================================

	http.HandleFunc("/hangman_web/connect", func(w http.ResponseWriter, r *http.Request) {
		temp, err := template.ParseFiles("./template/Connect.html")
		if err != nil {
			fmt.Print("Tu as une erreur de type : ", err)
		}

		temp.Execute(w, nil)
	})

	//=========================================================================================================================

	//=====================================================VERIFI CONNECTION===================================================

	http.HandleFunc("/hangman_web/verif_account", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {

			return
		}

		CheckValueName, _ := regexp.MatchString("[a-zA-Z-]{1,64}$", r.FormValue("nom"))
		CheckValuePassword, _ := regexp.MatchString("[a-zA-Z-]{1,64}$", r.FormValue("password"))

		if (!CheckValueName && !CheckValuePassword) || !CheckValueName || !CheckValuePassword {
			ValueError.IfError = true
			ValueError.Error = "Your name or password are misspelled!"
			http.Redirect(w, r, "/hangman_web/connect", http.StatusSeeOther)
			return
		}

		if VerifUser(r.FormValue("nom"), r.FormValue("password")) {
			http.Redirect(w, r, "/hangman_web/profil", http.StatusSeeOther)
			return
		}
	})

	//=========================================================================================================================

	//======================================================REGISTER===========================================================

	http.HandleFunc("/hangman_web/register", func(w http.ResponseWriter, r *http.Request) {
		temp, err := template.ParseFiles("./template/Register.html")
		if err != nil {
			fmt.Print("Tu as un problème de type : ", err)
			return
		}

		temp.Execute(w, ValueError)
	})

	//=========================================================================================================================

	//======================================================ADD ACCOUNT========================================================

	http.HandleFunc("/hangman_web/add", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {

			return
		}

		CheckValueName, _ := regexp.MatchString("[a-zA-Z-]{1,64}$", r.FormValue("nom"))

		if !CheckValueName {
			ValueError.IfError = true
			ValueError.Error = "Your name are misspelled!"
			http.Redirect(w, r, "/hangman_web/register", http.StatusSeeOther)
			return
		} else {
			ValueError.IfError = false

			AddUser(r.FormValue("nom"), 0, 0, r.FormValue("password"))
			http.Redirect(w, r, "/hangman_web/connect", http.StatusSeeOther)
			return
		}
	})

	//==========================================================================================================================

	//=====================================================SCOREBOARD===========================================================

	http.HandleFunc("/hangman_web/scoreboard", func(w http.ResponseWriter, r *http.Request) {
		temp, err := template.ParseFiles("./template/ScoreBoard.html")
		if err != nil {
			fmt.Print("Tu as une rerreur de type : ", err)
		}

		temp.Execute(w, nil)
	})

	//==========================================================================================================================
}

//========================================================ADD JSON==========================================================

func AddUser(name string, score int, nbGame int, password string) error {

	newUser := Profil{
		Name:     name,
		Score:    score,
		NbGame:   nbGame,
		Password: password,
	}

	// Lire les utilisateurs existants à partir du fichier User.json
	var users []Profil
	file, err := os.OpenFile("User.json", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("erreur lors de l'ouverture du fichier: %w", err)
	}
	defer file.Close()

	fileContent, err := ioutil.ReadAll(file)
	if err != nil {
		return fmt.Errorf("erreur lors de la lecture du fichier: %w", err)
	}

	if len(fileContent) > 0 {
		err = json.Unmarshal(fileContent, &users)
		if err != nil {
			return fmt.Errorf("erreur lors de la désérialisation: %w", err)
		}
	}

	// Ajouter le nouvel utilisateur à la liste
	users = append(users, newUser)

	// Sérialiser les utilisateurs en JSON
	updatedContent, err := json.MarshalIndent(users, "", "  ")
	if err != nil {
		return fmt.Errorf("erreur lors de la sérialisation: %w", err)
	}

	// Réécrire le fichier avec les nouvelles données
	err = ioutil.WriteFile("User.json", updatedContent, 0644)
	if err != nil {
		return fmt.Errorf("erreur lors de l'écriture dans le fichier: %w", err)
	}

	return nil
}

//=========================================================================================================================

func VerifUser(name, password string) bool {
	// Lire les utilisateurs à partir du fichier User.json
	var users []Profil
	file, err := os.Open("User.json")
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false
		}
		return false
	}
	defer file.Close()

	fileContent, err := ioutil.ReadAll(file)
	if err != nil {
		return false
	}

	err = json.Unmarshal(fileContent, &users)
	if err != nil {
		return false
	}

	// Parcourir la liste des utilisateurs
	for _, user := range users {
		if user.Name == name && user.Password == password {

			FirstConnect = false
			Connect.Name = name
			Connect.NbGame = user.NbGame
			Connect.Score = user.Score
			Connect.Connect = true

			return true
		}
	}
	return false
}
