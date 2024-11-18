package main

import (
	"fmt"
	"io/ioutil"
	"math/rand/v2"
	"net/http"
	"regexp"
	"text/template"
)

var Texte string

var LosePoint int
var ListeWord []string
var LosePointMax int
var Inconnu []string
var Begin bool
var LettreError []string

type Return struct {
	NbError int
	Live    int
	Mot     []string
	Error   bool
}

type Display struct {
	NbError int
	Live    int
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

type EndGame struct {
	Finish bool
	Mot    string
}

var Word string
var Essai []string

var StockageValueWord StockageValue = StockageValue{"", false}
var StockageValueLetter StockageValue = StockageValue{"", false}
var ErrorInputWord ErrorInput = ErrorInput{""}
var ErrorInputLetter ErrorInput = ErrorInput{""}

var Result Return
var Affichage Display

func Hangman() {

	//===============================================HANGMAN=============================================

	http.HandleFunc("/hangman_web/game", func(w http.ResponseWriter, r *http.Request) {
		temp, err := template.ParseFiles("./template/Hangman.html")
		if err != nil {
			fmt.Print("Il y a un problème de type :", err)
			return
		}
		Essai = []string{}
		LosePoint = 0

		Data, err := ioutil.ReadFile(Texte)
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
			Result = Return{LosePoint, 11 - LosePoint, ChooseWord(), true}
		} else {
			Result = Return{LosePoint, 11 - LosePoint, ChooseWord(), false}
		}
		temp.Execute(w, Result)
	})
	//===================================================================================================

	//===============================================TREATMENT=============================================

	http.HandleFunc("/hangman_web/treatment", func(w http.ResponseWriter, r *http.Request) {

		Affichage.Error = ""
		if r.Method != http.MethodPost {
			fmt.Println("Il y a un problème de traitement !")
			return
		}
		fmt.Println(Word)
		fmt.Println(r.FormValue("word"))
		fmt.Println(r.FormValue("letter"))

		checkValueWord, _ := regexp.MatchString("[a-zA-Z-]{1,64}$", r.FormValue("word"))
		checkValueLetter, _ := regexp.MatchString("[a-zA-Z-]{1,64}$", r.FormValue("letter"))

		if r.FormValue("word") != "" {
			if !checkValueWord {
				StockageValueWord = StockageValue{"", false}
				fmt.Println("Il y a un problème avec l'input mot")
				http.Redirect(w, r, "/hangman_web/display", http.StatusSeeOther)
				return
			}
		}
		if len(r.FormValue("word")) > 1 {
			if len(r.FormValue("word")) > len(Word) { //Vérifie si le mot est trop long..............................
				Affichage.Error = "Le mot est trop long"
				http.Redirect(w, r, "/hangman_web/display", http.StatusSeeOther)
				return
			} else {
				if r.FormValue("word") == Word { //Vérifie si l'input mot correspond au mot inconnu...................
					StockageValueWord = StockageValue{r.FormValue("word"), true}
					http.Redirect(w, r, "/hangman_web/win", http.StatusSeeOther)
					return
				} else if r.FormValue("word") != Word { //Si l'input mot ne l'est pas.................................
					Error := 0
					for _, k := range Essai {
						if string(k) == r.FormValue("word") { //Vérifie si on a pas déjà rentrer le mot...............
							Affichage.Error = "C'est pas le bon mot et tu l'as déjà rentrer"
							http.Redirect(w, r, "/hangman_web/display", http.StatusSeeOther)
							return
						}
						if r.FormValue("word") != string(k) {
							Error++
						}
					}
					if Error == len(Essai) { //Enregistre le mot rentrer dans la liste des erreurs....................
						if LosePoint >= 11 {
							http.Redirect(w, r, "/hangman_web/lose", http.StatusSeeOther)
							return
						} else {
							Essai = append(Essai, r.FormValue("word"))
							LosePoint += 2
							Affichage.Error = "C'est pas le bon mot"
							http.Redirect(w, r, "/hangman_web/display", http.StatusSeeOther)
							return
						}
					}
				}
			}
		}
		if r.FormValue("letter") != "" {
			if !checkValueLetter {
				StockageValueLetter = StockageValue{"", false}
				fmt.Println("Il y a un problème avec l'input lettre")
				http.Redirect(w, r, "/hangman_web/display", http.StatusSeeOther)
				return
			}
		}
		if len(r.FormValue("letter")) == 1 { //Vérifie si il n'y a qu'un seul paramètre.......................
			Error := 0
			for index, k := range Word {
				if string(k) == r.FormValue("letter") { //Vérifie si la lettre n'est pas dans le mot inconnue.
					StockageValueLetter = StockageValue{r.FormValue("letter"), true}
					Inconnu[index] = r.FormValue("letter")
				} else {
					Error++
				}
			}
			var mot string
			for _, Lettre := range Inconnu {
				mot += string(Lettre)
			}
			if mot == Word {
				http.Redirect(w, r, "/hangman_web/win", http.StatusSeeOther)
				return
			}
			if Error == len(Word) { //Si la lettre n'est pas dans le mot.......................................
				Error := 0
				for _, L := range Essai {
					if r.FormValue("letter") == string(L) { //Vérifie si la lettre a déjà été rentrer..........
						Affichage.Error = "C'est pas la bonne lettre et tu l'as déjà rentrer"
						http.Redirect(w, r, "/hangman_web/display", http.StatusSeeOther)
						return
					}
					if r.FormValue("letter") != string(L) {
						Error++
					}
				}
				if Error == len(Essai) {
					if LosePoint >= 11 {
						http.Redirect(w, r, "/hangman_web/lose", http.StatusSeeOther)
						return
					} else {
						Affichage.Error = "La lettre n'est pas dans le mot"
						LosePoint++
						Essai = append(Essai, r.FormValue("letter"))
					}
				}
			}
		} else if len(r.FormValue("letter")) > 1 {
			ErrorInputLetter = ErrorInput{"Trop de lettre !"}
		}
		http.Redirect(w, r, "/hangman_web/display", http.StatusSeeOther)
	})
	//===================================================================================================

	//===============================================DISPLAY=============================================

	http.HandleFunc("/hangman_web/display", func(w http.ResponseWriter, r *http.Request) {
		temp, err := template.ParseFiles("./template/HangmanDisplay.html")
		if err != nil {
			fmt.Print("Il y a un problème de type :", err)
			return
		}
		Affichage.IfError = false //Renvoie qu'il n'y a pas d'erreur...........................................
		if Affichage.Error != "" {
			Affichage.IfError = true
		}

		Affichage.Mot = Inconnu
		Affichage.Live = 11 - LosePoint
		Affichage.NbError = LosePoint
		Affichage.Try = Essai

		temp.Execute(w, Affichage)
	})

	//===================================================================================================

	//=================================================WIN===============================================

	http.HandleFunc("/hangman_web/win", func(w http.ResponseWriter, r *http.Request) {
		temp, err := template.ParseFiles("./template/Finish.html")
		if err != nil {
			fmt.Print("Il y a un problème de type : ", err)
			return
		}

		Win := EndGame{
			Finish: true,
			Mot:    Word,
		}

		Connect.NbGame++
		Connect.Score++

		temp.Execute(w, Win)
	})

	//===================================================================================================

	//===============================================GAME OVER===========================================

	http.HandleFunc("/hangman_web/lose", func(w http.ResponseWriter, r *http.Request) {
		temp, err := template.ParseFiles("./template/Finish.html")
		if err != nil {
			fmt.Print("Il y a un problème de type : ", err)
			return
		}
		GameOver := EndGame{
			Finish: false,
			Mot:    Word,
		}

		Connect.NbGame++

		temp.Execute(w, GameOver)
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
		if string(Word[i]) == "-" {
			Inconnu = append(Inconnu, "-")
		} else {
			Inconnu = append(Inconnu, "_")
		}
	}

	return Inconnu
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
