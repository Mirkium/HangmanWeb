package main

import (
	"fmt"
	"io/ioutil"
	"math/rand/v2"
	"net/http"
	"text/template"
)

var LosePoint int = 6
var ListeWord []string
var LosePointMax int
var Inconnu []string
var Begin bool
var LettreError []string

func Hangman() {

	type Return struct {
		NbError int
		Mot     []string
		Error   bool
	}

	http.HandleFunc("/game", func(w http.ResponseWriter, r *http.Request) {
		temp, err := template.ParseFiles("./template/Hangman.html")
		if err != nil {
			fmt.Print("Il y a un problème de type :", err)
			return
		}

		Data, err := ioutil.ReadFile("./txt/Global.txt") // lire le fichier text.txt
		if err != nil {
			fmt.Println(err)
		}
		ListeWord = []string{}
		ligneTexte := string(Data)
		var word string
		for _, k := range ligneTexte {
			if k != '\n' && k != '\r' { // Utiliser && au lieu de ||
				word += string(k)
			} else if word != "" { // Ajouter le mot si non vide
				ListeWord = append(ListeWord, word)
				word = ""
			}
		}
		var Result Return
		if LosePoint > 0 && LosePoint <= 11 {
			Result = Return{LosePoint, ChooseWord(), true}
		} else {
			Result = Return{LosePoint, ChooseWord(), false}
		}
		fmt.Print(Result.Error)
		temp.Execute(w, Result)
	})

	http.HandleFunc("/game/mot", func(w http.ResponseWriter, r *http.Request) {
		tmp, err := template.ParseFiles("./template/HangmanMot.html")
		if err != nil {
			fmt.Print("Il y a un problème de type :", err)
			return
		}

		tmp.Execute(w, nil)
	})

	http.HandleFunc("/game/lettre", func(w http.ResponseWriter, r *http.Request) {
		tmp, err := template.ParseFiles("./template/HangmanLettre.html")
		if err != nil {
			fmt.Print("Il y a un problème de type :", err)
			return
		}

		tmp.Execute(w, nil)
	})
}

func ChooseWord() []string {
	//var SideWord string
	var Word string = ""
	Inconnu = []string{}
	LettreError = []string{}
	Word = ListeWord[rand.IntN(len(ListeWord))]
	Word = ToLower(Word)
	for i := 0; i <= len(Word)-1; i++ {
		Inconnu = append(Inconnu, "_")
	}

	return Inconnu
}

func RandomHelp(mot string) {
	help := mot[rand.IntN(len(mot))]
	for i, k := range mot {
		if string(help) == string(k) {
			Inconnu[i] = string(help)
		}
	}
	Begin = false
}

func ToLower(s string) string {
	result := ""
	for _, char := range s {
		if char >= 'A' && char <= 'Z' {
			result += string(char + 32)
		} else {
			result += string(char)
		}
	}
	return result
}
