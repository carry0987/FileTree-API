package handler

import (
	"net/http"

	"github.com/carry0987/FileTree-API/internal/utils"
	"github.com/carry0987/FileTree-API/pkg/api"
	jsoniter "github.com/json-iterator/go"
)

func HTTPMessage(w http.ResponseWriter, r *http.Request, message string) {
	utils.OutputMessage(w, utils.HTTPResponse, http.StatusOK, message)
}

func HTTPHandler(w http.ResponseWriter, r *http.Request) {
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	fileTreeResult, err := ProcessEncryptedPath(r)

	if err != nil {
		errorMsg := err.Error()
		apiErrorResponse := api.NewErrorResponse(errorMsg)
		w.WriteHeader(DetermineHTTPStatusCode(err))
		json.NewEncoder(w).Encode(apiErrorResponse)
		return
	}

	// Return the file tree
	response := api.NewSuccessResponse(fileTreeResult)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
