package handlers

import (
	"encoding/json"
	"net/http"
	"slate-rmm/database"

	"github.com/gorilla/mux"
)

// GetAllGroups handles the GET /api/groups route
func GetAllGroups(w http.ResponseWriter, r *http.Request) {
	groups, err := database.GetAllGroups()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(groups)
}

// GetGroup handles the GET /api/groups/{group_id} route
func GetGroup(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	groupID := vars["group_id"]

	group, err := database.GetGroup(groupID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(group)
}

// CreateGroup handles the POST /api/groups route
func CreateGroup(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		GroupName string `json:"group_name"`
	}

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = database.CreateGroup(payload.GroupName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Group created"))
}

// UpdateGroup handles the PUT /api/groups/{group_id} route
func UpdateGroup(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	groupID := vars["group_id"]
	groupName := r.FormValue("group_name")

	err := database.UpdateGroup(groupID, groupName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Group " + groupID + " updated"))
}

// DeleteGroup handles the DELETE /api/groups/{group_id} route
func DeleteGroup(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	groupID := vars["group_id"]

	err := database.DeleteGroup(groupID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Group " + groupID + " deleted"))
}
