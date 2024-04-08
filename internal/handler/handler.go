package handler

import (
	"net/http"

	"github.com/carry0987/FileTree-API/internal/security"
	"github.com/carry0987/FileTree-API/internal/service"
	"github.com/carry0987/FileTree-API/internal/utils"
	"github.com/carry0987/FileTree-API/pkg/api"
	jsoniter "github.com/json-iterator/go"

	"github.com/gorilla/mux"
)

func DefaultHandler(w http.ResponseWriter, r *http.Request) {
	utils.OutputMessage(w, utils.HTTPResponse, http.StatusOK, "FileTree API")
}

// FileTreeHandler handles the file tree API request
func FileTreeHandler(w http.ResponseWriter, r *http.Request) {
	json := jsoniter.ConfigCompatibleWithStandardLibrary

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
	realPath, organize := utils.CheckOrganize(decryptedPath)
	fileTreeResult, err := service.GenerateFileTree(realPath, organize)
	if err != nil {
		errorMsg := "Error generating file tree"
		api.InternalServerError(errorMsg)
		response := api.NewErrorResponse(errorMsg)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Return the file tree
	response := api.NewSuccessResponse(fileTreeResult)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
