package main

import (
	"fmt"
	"io/ioutil"
	"math/rand/v2"
	"net/http"
	"regexp"
	"text/template"
)

var LosePoint int
var ListeWord []string
var LosePointMax int
var Inconnu []string
var Begin bool
var LettreError []string

type Return struct {
	NbError int
	Mot     []string
	Error   bool
}

type Display struct {
	NbError int
	Mot     []string
	Try     []string
	Error   string
	IfError bool
}

type StockageValue struct {
	Value      string
	CheckValue bool
}

type ErrorInput struct {
	Error string
}

var Word string
var Essai []string

var StockageValueWord StockageValue = StockageValue{"", false}
var StockageValueLetter StockageValue = StockageValue{"", false}
var ErrorInputWord ErrorInput = ErrorInput{""}
var ErrorInputLetter ErrorInput = ErrorInput{""}

var Result Return

func Hangman() {

	//===============================================HANGMAN=============================================

	http.HandleFunc("/game", func(w http.ResponseWriter, r *http.Request) {
		temp, err := template.ParseFiles("./template/Hangman.html")
		if err != nil {
			fmt.Print("Il y a un problème de type :", err)
			return
		}
		Essai = []string{}

		Data, err := ioutil.ReadFile("./txt/Global.txt") // lire le fichier text.txt
		if err != nil {
			fmt.Println(err)
		}
		ListeWord = []string{}
		ligneTexte := string(Data)

		for _, k := range ligneTexte {
			if k != '\n' && k != '\r' { // Utiliser && au lieu de ||
				Word += string(k)
			} else if Word != "" { // Ajouter le mot si non vide
				ListeWord = append(ListeWord, Word)
				Word = ""
			}
		}

		if LosePoint > 0 && LosePoint <= 11 {
			Result = Return{LosePoint, ChooseWord(), true}
		} else {
			Result = Return{LosePoint, ChooseWord(), false}
		}
		temp.Execute(w, Result)
	})
	//===================================================================================================

	//===============================================TREATMENT=============================================

	http.HandleFunc("/hangman/treatment", func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodPost {
			fmt.Println("toi t'a fais de la merde")
			return
		}

		checkValueWord, _ := regexp.MatchString("[a-zA-Z-]{1,64}$", r.FormValue("nom"))
		checkValueLetter, _ := regexp.MatchString("[a-zA-Z-]{1,64}$", r.FormValue("prenom"))

		if !checkValueWord {
			StockageValueWord = StockageValue{"", false}
			http.Redirect(w, r, "/hangman/display", http.StatusSeeOther)
			return
		}
		if !checkValueLetter {
			StockageValueLetter = StockageValue{"", false}
			http.Redirect(w, r, "/hangman/display", http.StatusSeeOther)
			return
		}

		if r.FormValue("word") == Word { //Vérifie si l'input mot correspond au mot inconnu....................................
			StockageValueWord = StockageValue{r.FormValue("word"), true}
			http.Redirect(w, r, "/hangman/win", http.StatusSeeOther)
			return
		} else if r.FormValue("word") != Word { //Si l'input mot ne l'est pas..................................................
			Error := 0
			for _, k := range Essai {
				if string(k) == r.FormValue("word") { //Vérifie si on a pas déjà rentrer le mot................................
					ErrorInputWord = ErrorInput{"C'est pas le bon mot et tu l'as déjà rentrer"}
					http.Redirect(w, r, "/hangman/display", http.StatusSeeOther)
					return
				}
				Error++
			}
			if Error == len(Essai) { //Enregistre le mot rentrer dans la liste des erreurs.....................................
				Essai = append(Essai, r.FormValue("word"))
				LosePoint += 2
				ErrorInputWord = ErrorInput{"C'est pas le bon mot"}
				http.Redirect(w, r, "/hangman/display", http.StatusSeeOther)
				return
			}
		}
		if len(r.FormValue("letter")) == 1 { //Vérifie si il n'y a qu'un seul paramètre........................................
			Error := 0
			for index, k := range Word {
				if string(k) == r.FormValue("letter") { //Vérifie si la lettre n'est pas dans le mot inconnue..................
					StockageValueLetter = StockageValue{r.FormValue("letter"), true}
					http.Redirect(w, r, "/hangman/display", http.StatusSeeOther)
					Inconnu[index] = r.FormValue("letter")
					return
				}
				Error++
			}
			if Error == len(Word) {
				Essai = append(Essai, r.FormValue("letter"))
				
			}
		} else if len(r.FormValue("letter")) > 1 {
			ErrorInputLetter = ErrorInput{"Trop de lettre !"}
		}

	})
	//===================================================================================================

	//===============================================DISPLAY=============================================

	http.HandleFunc("/hangman/display", func(w http.ResponseWriter, r *http.Request) {
		temp, err := template.ParseFiles("./template/HangmanDisplay.html")
		if err != nil {
			fmt.Print("Il y a un problème de type :", err)
			return
		}
		Return := Display{}    //Structure de return............................................................
		Return.IfError = false //Renvooie qu'il n'y a pas d'erreur...........................................

		fmt.Println(ErrorInputLetter)
		fmt.Println(ErrorInputWord)

		if ErrorInputLetter.Error != "" { //vérifie si il y a une erreur en lien avec l'input de lettre......
			Return.Error = ErrorInputLetter.Error
			Return.IfError = true
		}
		if ErrorInputWord.Error != "" { //vérifie si il y a une erreur en lien avec l'input de mot...........
			Return.Error = ErrorInputWord.Error
			Return.IfError = true
		}

		Return.Mot = Inconnu
		Return.NbError = LosePoint
		Return.Try = Essai
		fmt.Println(Return.Mot)
		fmt.Println(Return.Try)
		fmt.Println(Return.NbError)
		temp.Execute(w, Return)
	})

	//===================================================================================================
}

func ChooseWord() []string {
	//var SideWord string
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
