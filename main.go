package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type TranslationRequest struct {
	Text     string `json:"text"`
	FromLang string `json:"from_lang"`
	ToLang   string `json:"to_lang"`
}

type TranslationResponse struct {
	TranslatedText string `json:"translated_text"`
}

func translateText(request TranslationRequest) (TranslationResponse, error) {

	apiURL := "http://translationapi.com/translate"

	requestBody, err := json.Marshal(request)
	if err != nil {
		return TranslationResponse{}, err
	}

	response, err := http.Post(apiURL, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return TranslationResponse{}, err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return TranslationResponse{}, err
	}

	var translationResponse TranslationResponse
	if err := json.Unmarshal(body, &translationResponse); err != nil {
		return TranslationResponse{}, err
	}

	return translationResponse, nil
}

func translateHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Метод не підтримується", http.StatusMethodNotAllowed)
		return
	}

	var translationRequest TranslationRequest
	if err := json.NewDecoder(r.Body).Decode(&translationRequest); err != nil {
		http.Error(w, "Неправильний формат даних запиту", http.StatusBadRequest)
		return
	}

	translationResult, err := translateText(translationRequest)
	if err != nil {
		http.Error(w, "Помилка перекладу тексту", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(translationResult)
}

func main() {

	http.HandleFunc("/translate", translateHandler)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Помилка при запуску сервера:", err)
		return
	}
}
