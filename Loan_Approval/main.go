package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

func main() {

	//Handling the route
	http.HandleFunc("/", handle)

	//Getting port based on local data
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func handle(w http.ResponseWriter, r *http.Request) {

	//Checking route pattern
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	//Initializing context
	ctx := context.Background()

	//Setting response content type
	w.Header().Set("Content-Type", "text/html")

	// GET
	if r.Method == http.MethodGet {

		//Fetching name
		noms, ok := r.URL.Query()["nom"]

		if !ok ||len(noms[0]) <1 {
			msg := fmt.Sprintf("Url Param 'nom' is missing")
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}

		nom := noms[0]

		//Fetching amount
		montants, ok := r.URL.Query()["montant"]

		if !ok ||len(montants[0]) <1 {
			msg := fmt.Sprintf("Url Param 'nom' is missing")
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}

		montant, err := strconv.ParseFloat(montants[0],64)

		if err != nil {
			http.Error(w, "Erreur parse int : " + err.Error(), http.StatusInternalServerError)
			return
		}

		//Getting the risk for the required user
		resp, err := http.Get("https://lpro-cloud-check-account-1.herokuapp.com/?nom=" + nom)

		if err != nil {
			msg := fmt.Sprintf("Aucun compte pour ce nom")
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}

		defer resp.Body.Close()

		body, readErr := ioutil.ReadAll(resp.Body)

		if readErr != nil {
			msg := fmt.Sprintf("Erreur de lecture :")
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}

		var answer RiskAnswer

		jsonErr := json.Unmarshal(body, &answer)

		if jsonErr != nil {
			msg := fmt.Sprintf("Aucun compte pour ce nom")
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}

		if answer.Risk == "high" {
			appManagerHandle(nom, "refused", w, ctx)

		} else if answer.Risk == "low" {
			//Check montant
			if montant > 10000 {
				appManagerHandle(nom, "refused", w, ctx)
			} else {
				appManagerHandle(nom, "accepted", w, ctx)
				accManagerIncr(nom, montant, w, ctx)
			}
		}

	}

	return
}

type RiskAnswer struct {
	Risk string `json:"risk"`
}

type Approval struct {
	Nom string `json:"nom"`
	ReponseManuelle string `json:"reponse_manuelle"`
}

type CreditOperation struct {
	Nom string `json:"nom"`
	Montant float64 `json:"montant"`
}

// Function that handles the approval post
func appManagerHandle(nom string, state string, w http.ResponseWriter, ctx context.Context) {
	approval := Approval{Nom: nom, ReponseManuelle: state}

	jsonStr, err := json.Marshal(approval)

	if err != nil {
		http.Error(w, "Json error : " + err.Error(), http.StatusInternalServerError)
		return
	}

	resp, err := http.Post("https://lp-cloud-app-manager-1.ew.r.appspot.com/", "application/json", bytes.NewBuffer(jsonStr))

	if err != nil {
		http.Error(w, "Request Post error" + err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println(resp.Body)
	fmt.Fprintf(w, "Loan approval state : %s\n", state)

	return
}

//Function that handles the account incrementation
func accManagerIncr(nom string, montant float64, w http.ResponseWriter,ctx context.Context)  {
	ope := CreditOperation{Nom: nom, Montant: montant}

	jsonStr, err := json.Marshal(ope)

	if err != nil {
		http.Error(w, "Json error : " + err.Error(), http.StatusInternalServerError)
		return
	}

	resp, err := http.Post("https://lpcloud-project-1.ew.r.appspot.com/add", "application/json", bytes.NewBuffer(jsonStr))

	if err != nil {
		http.Error(w, "Request Post error" + err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println(resp.Body)
	_, err = fmt.Fprintf(w, "Amount added : %f\n", montant )

	if err!=nil {
		http.Error(w, "Printing body error", http.StatusInternalServerError)
	}
}