package handler

import (
	"FileTree-API/internal/security"
	"FileTree-API/internal/service"
	"FileTree-API/internal/utils"
	"FileTree-API/pkg/api"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func DefaultHandler(w http.ResponseWriter, r *http.Request) {
	utils.OutputMessage(w, utils.HTTPResponse, http.StatusOK, "FileTree API")
}

// FileTreeHandler handles the file tree API request
func FileTreeHandler(w http.ResponseWriter, r *http.Request) {
	// Get the signature and encrypted parameters from the route or query parameters
	vars := mux.Vars(r)
	encryptedPath := vars["encrypted"]
	if encryptedPath == "" {
		errorMsg := "Missing encrypted parameter"
		api.UnauthorizedError(errorMsg)
		utils.OutputMessage(w, utils.HTTPResponse, http.StatusBadRequest, errorMsg)
		return
	}

	// Decrypt the path
	decryptedPath, err := security.Decrypt(encryptedPath)
	if err != nil {
		errorMsg := "Failed to decrypt the path"
		api.UnauthorizedError(errorMsg)
		response := api.NewErrorResponse(errorMsg)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Generate the file tree using the decrypted path
	fileTree, err := service.GenerateFileTree(decryptedPath)
	if err != nil {
		errorMsg := "Error generating file tree"
		api.InternalServerError(errorMsg)
		response := api.NewErrorResponse(errorMsg)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Return the file tree
	response := api.NewSuccessResponse(fileTree)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
