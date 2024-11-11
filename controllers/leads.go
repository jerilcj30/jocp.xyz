package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/jerilcj30/jocp/models"
)

type LeadsHandler struct {
	LeadsService *models.LeadsService
}

func (l LeadsHandler) CreateLead(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"name":    r.FormValue("full_name"),
		"email":   strings.ToLower(r.FormValue("email")),
		"message": r.FormValue("message"),
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("Failed to marshal data: %v", err)
	}

	_, err = l.LeadsService.Create(jsonData)
	if err != nil {
		fmt.Println((err))
		http.Error(w, "Something went wrong.", http.StatusInternalServerError)
		return
	}

	cookie := http.Cookie{
		Name:     "successMessage",
		Value:    "Your lead has been successfully submitted!",
		Path:     "/",  // cookies can be accessed anywhere in the whole application
		HttpOnly: true, // cookie cannnot be accessed by javascript
	}

	http.SetCookie(w, &cookie)

	http.Redirect(w, r, "/", http.StatusFound)
}
