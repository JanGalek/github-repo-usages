// api/getGoDependency.go
package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func GetGoDependency(w http.ResponseWriter, r *http.Request) {
	// Získání parametrů z URL
	repository := r.URL.Query().Get("repository")

	owner := "gouef"
	repo := "github-repo-usages"

	if repository == "" {
		http.Error(w, "Missing 'owner' or 'repo' query parameter", http.StatusBadRequest)
		return
	} else {
		// Rozdělení 'repository' na owner a repo
		parts := strings.Split(repository, "/")
		if len(parts) != 2 {
			http.Error(w, "Invalid 'repository' format. Expected 'owner/repo'", http.StatusBadRequest)
			return
		}
		owner = parts[0]
		repo = parts[1]
	}

	// URL pro GitHub API (získání počtu běhů akcí pro daný repozitář)
	url := fmt.Sprintf("https://api.github.com/search-code?q=path:**/go.mod %s/%s", owner, repo)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Nastavení hlavičky pro autorizaci GitHub API
	req.Header.Set("Authorization", "Bearer "+os.Getenv("GITHUB_TOKEN"))

	// Odeslání požadavku na GitHub API
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Failed to contact GitHub API", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Kontrola statusu odpovědi
	if resp.StatusCode != http.StatusOK {
		http.Error(w, "Error retrieving data from GitHub", http.StatusInternalServerError)
		return
	}

	// Načtení a parsování JSON odpovědi
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read response body", http.StatusInternalServerError)
		return
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		http.Error(w, "Failed to parse response body", http.StatusInternalServerError)
		return
	}

	actionsCount := result["total_count"].(float64)                  // GitHub API vrací číslo jako float64
	actionsCountStr := strconv.FormatFloat(actionsCount, 'f', 0, 64) // Převod na string

	// Vytvoření odpovědi ve formátu JSON
	response := map[string]interface{}{
		"schemaVersion": 1,
		"message":       actionsCountStr,
		"label":         "usages",
		"color":         "blue",
	}

	// Nastavení správných hlaviček pro JSON odpověď
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Odeslání odpovědi
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
