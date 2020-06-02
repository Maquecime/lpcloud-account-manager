// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/datastore"
	"google.golang.org/appengine"
)

var datastoreClient *datastore.Client

func main() {
	ctx := context.Background()

	// Set this in app.yaml when running in production.
	projectID := os.Getenv("GCLOUD_DATASET_ID")

	var err error
	datastoreClient, err = datastore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", handle)
	http.HandleFunc("/add", handleAdd)
	appengine.Main()
}

func handle(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	ctx := context.Background()

	if r.Method == http.MethodGet {
		// Get a list of every bank accounts.
		bankAccounts, err := queryBankAccounts(ctx)
		if err != nil {
			msg := fmt.Sprintf("Could not get the approvals: %v", err)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}

		accJson, err := json.Marshal(BankAccounts{bankAccounts})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		_, err = w.Write(accJson)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		return

	} else if r.Method == http.MethodPut {
		var account BankAccount
		err := json.NewDecoder(r.Body).Decode(&account)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		addBankAccount(w, ctx, &account)

		accJson, err := json.Marshal(account)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		_, err = w.Write(accJson)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		return

	} else if r.Method == http.MethodDelete {
		keys, ok := r.URL.Query()["key"]

		if !ok || len(keys[0]) <1 {
			msg := fmt.Sprintf("Url Param 'key' is missing")
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}

		key := keys[0]

		err := deleteBankAccount(w, ctx, key)
		if err != nil {
			return
		}

		w.WriteHeader(200)

		_, err = fmt.Fprintf(w, "Bank account supprimé avec la clé %s", key)

		if err!= nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	return
}

func handleAdd(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/add" {
		http.NotFound(w, r)
		return
	}

	ctx := context.Background()

	w.Header().Set("Content-Type", "application/json")


	if r.Method == http.MethodPost {
		var op CreditOperation
		err := json.NewDecoder(r.Body).Decode(&op)
		if err != nil && op.Nom != "" {
			http.Error(w, "json decode error : " + err.Error(), http.StatusBadRequest)
			return
		}

		err = incrBankAccount(ctx, &op)

		if err != nil {
			http.Error(w, "incr error : " + err.Error(), http.StatusInternalServerError)
		}
		opSerialized, errJson := json.Marshal(op)

		if errJson != nil {
			http.Error(w, "serialized result error : "+errJson.Error(), http.StatusInternalServerError)
		}
		_, err = w.Write(opSerialized)

		if err != nil {
			http.Error(w, "write error :"+err.Error(), http.StatusInternalServerError)
		}
	}

	return
}



type BankAccount struct {
	Key *datastore.Key `datastore:"__key__"`
	Nom string `json:"nom"`
	Prenom string `json:"prenom"`
	Montant float64 `json:"montant"`
	Risk string `json:"risk"`
}

type BankAccounts struct {
	BankAccounts []*BankAccount
}

type CreditOperation struct {
	Nom string `json:"nom"`
	Montant float64 `json:"montant"`
}

func addBankAccount(w http.ResponseWriter,ctx context.Context, account *BankAccount) {

	if account.Key == nil {
		account.Key = datastore.NameKey("BankAccount", account.Nom, nil)
	}

	_, err := datastoreClient.Put(ctx, account.Key, account)

	if err != nil {
		msg := fmt.Sprintf("Could not add BankAccount: %v", err)
		http.Error(w, msg, http.StatusInternalServerError)
	}

	return
}

func deleteBankAccount(w http.ResponseWriter, ctx context.Context, key string) error {
	decodedKey, err := datastore.DecodeKey(key)

	if err != nil {
		msg := fmt.Sprintf("Invalid Key: %v", err)
		http.Error(w, msg, http.StatusInternalServerError)
		return err
	}

	err = datastoreClient.Delete(ctx, decodedKey)

	if err != nil {
		msg := fmt.Sprintf("Could not remove BankAccount: %v", err)
		http.Error(w, msg, http.StatusInternalServerError)

	}
	return err
}

func queryBankAccounts(ctx context.Context) ([]*BankAccount, error) {
	q := datastore.NewQuery("BankAccount")

	bankAccounts := make([]*BankAccount, 0)
	_, err := datastoreClient.GetAll(ctx, q, &bankAccounts)
	return bankAccounts, err
}

func incrBankAccount(ctx context.Context, op *CreditOperation) error {
	var bankAccounts []*BankAccount

	q := datastore.NewQuery("BankAccount").Filter("Nom =", op.Nom).Limit(1)
	if _, err := datastoreClient.GetAll(ctx, q, &bankAccounts); err != nil {
		return err
	}

	if bankAccounts == nil {
		return errors.New("null bank account")
	}
	bankAccount := bankAccounts[0]
	bankAccount.Montant += op.Montant

	_, err := datastoreClient.Put(ctx, bankAccount.Key, bankAccount)

	if err != nil {
		return err
	}

	return nil
}