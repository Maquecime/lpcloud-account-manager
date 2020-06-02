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
		// Get a list of the every approvals.
		approvals, err := queryApprovals(ctx)
		if err != nil {
			msg := fmt.Sprintf("Could not get the approvals: %v", err)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}

		appJson, err := json.Marshal(Approvals{approvals})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write(appJson)

		return

	} else if r.Method == http.MethodPost {
		var approval Approval
		err := json.NewDecoder(r.Body).Decode(&approval)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		addApproval(w, ctx, &approval)

		appJson, err := json.Marshal(approval)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Println(string(appJson))

		return

	} else if r.Method == http.MethodDelete {
		keys, ok := r.URL.Query()["key"]

		if !ok || len(keys[0]) <1 {
			msg := fmt.Sprintf("Url Param 'key' is missing")
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}

		key := keys[0]

		err := deleteApproval(w, ctx, key)
		if err != nil {
			return
		}

		w.WriteHeader(200)

		fmt.Fprintf(w, "Approval supprimé avec la clé %s", key)

		return
	}
}


type Approval struct {
	Key *datastore.Key `datastore:"__key__"`
	Nom string `json:"nom"`
	ReponseManuelle string `json:"reponse_manuelle"`
}

type Approvals struct {
	Approvals []*Approval
}

func addApproval(w http.ResponseWriter,ctx context.Context, approval *Approval) {

	if approval.Key == nil {
		approval.Key = datastore.NameKey("Approval", approval.Nom, nil)
	}

	_, err := datastoreClient.Put(ctx, approval.Key, approval)

	if err != nil {
		msg := fmt.Sprintf("Could not add Approval: %v", err)
		http.Error(w, msg, http.StatusInternalServerError)
	}

	return
}

func deleteApproval(w http.ResponseWriter, ctx context.Context, key string) error {
	decodedKey, err := datastore.DecodeKey(key)

	if err != nil {
		msg := fmt.Sprintf("Invalid Key: %v", err)
		http.Error(w, msg, http.StatusInternalServerError)
		return err
	}

	err = datastoreClient.Delete(ctx, decodedKey)

	if err != nil {
		msg := fmt.Sprintf("Could not remove Approval: %v", err)
		http.Error(w, msg, http.StatusInternalServerError)

	}
	return err
}

func queryApprovals(ctx context.Context) ([]*Approval, error) {
	q := datastore.NewQuery("Approval")

	approvals := make([]*Approval, 0)
	_, err := datastoreClient.GetAll(ctx, q, &approvals)
	return approvals, err
}