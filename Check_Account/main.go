package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/", handle)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func handle(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodGet {

		ctx := context.Background()

		noms, ok := r.URL.Query()["nom"]

		if !ok ||len(noms[0]) <1 {
			msg := fmt.Sprintf("Url Param 'nom' is missing")
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}

		 nom := noms[0]

		 bankAccounts, err := queryBankAccounts(ctx, w)

		if err != nil {
			msg := fmt.Sprintf("Could not get the bank accounts: %v", err)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}

		 var foundedAcc *BankAccount

		for _, acc := range bankAccounts.BankAccounts {
			if nom == acc.Nom {
				foundedAcc = acc
			}
		}

		if foundedAcc != nil && (foundedAcc.Risk == "high" || foundedAcc.Risk == "low") {

			answ := answer{Risk: foundedAcc.Risk}
			answSerialized, err := json.Marshal(answ)

			if err != nil {
				msg := fmt.Sprintf("Could not create response: %v", err)
				http.Error(w, msg, http.StatusInternalServerError)
				return
			}

			w.Write(answSerialized)
		}

	}
}


type BankAccount struct {
	Nom string `json:"nom"`
	Prenom string `json:"prenom"`
	Montant float64 `json:"montant"`
	Risk string `json:"risk"`
}

type BankAccounts struct {
	BankAccounts []*BankAccount
}

type answer struct {
	Risk string `json:"risk"`
}

//Function that queries AccManager to get all the bank accounts
func queryBankAccounts(ctx context.Context, w http.ResponseWriter) (BankAccounts, error) {
	var bankaccounts BankAccounts

	resp, err := http.Get("https://lpcloud-project-1.ew.r.appspot.com/")

	if err != nil {
		return BankAccounts{}, err
	}

	defer resp.Body.Close()

	body, readErr := ioutil.ReadAll(resp.Body)

	if readErr != nil {
		return BankAccounts{}, readErr
	}

	jsonErr := json.Unmarshal(body, &bankaccounts)

	if jsonErr != nil {
		return BankAccounts{}, jsonErr
	}

	return bankaccounts, err

}