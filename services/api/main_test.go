package main

import (
	"bytes"
	"encoding/json"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"
)

func testApp() *App {
	return &App{phrases: []Phrase{{ID:"x1", Level:"HSK1", ZH:"我喜欢中文。", Pinyin:"Wǒ xǐhuan Zhōngwén.", JA:"中国語が好きです。", Tokens:[]string{"我", "喜欢", "中文。"}}}}
}

func TestCheckPracticeCorrectAndIncorrect(t *testing.T) {
	app := testApp()
	for _, tc := range []struct{ tokens []string; want bool }{
		{[]string{"我", "喜欢", "中文。"}, true},
		{[]string{"喜欢", "我", "中文。"}, false},
	} {
		body, _ := json.Marshal(CheckRequest{PhraseID:"x1", Tokens:tc.tokens})
		req := httptest.NewRequest(http.MethodPost, "/api/practice/check", bytes.NewReader(body))
		res := httptest.NewRecorder()
		app.routes().ServeHTTP(res, req)
		if res.Code != http.StatusOK { t.Fatalf("status = %d", res.Code) }
		var got struct{ Correct bool `json:"correct"` }
		if err := json.Unmarshal(res.Body.Bytes(), &got); err != nil { t.Fatal(err) }
		if got.Correct != tc.want { t.Fatalf("correct = %v, want %v", got.Correct, tc.want) }
	}
}

func TestRejectsUnknownPhrase(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/practice/check", bytes.NewBufferString(`{"phraseId":"missing","tokens":["我"]}`))
	res := httptest.NewRecorder()
	testApp().routes().ServeHTTP(res, req)
	if res.Code != http.StatusNotFound { t.Fatalf("status = %d, want 404", res.Code) }
}

func TestNewPracticeDoesNotMutateAnswerKey(t *testing.T) {
	app := testApp()
	app.rng = rand.New(rand.NewSource(1))
	req := httptest.NewRequest(http.MethodGet, "/api/practice/new?level=HSK1", nil)
	res := httptest.NewRecorder()
	app.routes().ServeHTTP(res, req)
	if res.Code != http.StatusOK { t.Fatalf("status = %d", res.Code) }
	if !equalTokens(app.phrases[0].Tokens, []string{"我", "喜欢", "中文。"}) { t.Fatal("practice generation mutated answer key") }
}
