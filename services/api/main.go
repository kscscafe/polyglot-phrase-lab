package main

import (
	"encoding/json"
	"errors"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type Phrase struct {
	ID     string   `json:"id"`
	Level  string   `json:"level"`
	ZH     string   `json:"zh"`
	Pinyin string   `json:"pinyin"`
	JA     string   `json:"ja"`
	Tokens []string `json:"tokens"`
}

type CheckRequest struct {
	PhraseID string   `json:"phraseId"`
	Tokens   []string `json:"tokens"`
}

type App struct {
	phrases []Phrase
	mu      sync.RWMutex
	rng     *rand.Rand
}

func loadPhrases(path string) ([]Phrase, error) {
	data, err := os.ReadFile(path)
	if err != nil { return nil, err }
	var phrases []Phrase
	if err := json.Unmarshal(data, &phrases); err != nil { return nil, err }
	if len(phrases) == 0 { return nil, errors.New("phrase dataset is empty") }
	for _, p := range phrases {
		if p.ID == "" || p.ZH == "" || len(p.Tokens) == 0 { return nil, errors.New("phrase dataset contains an incomplete record") }
	}
	return phrases, nil
}

func (a *App) routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/health", func(w http.ResponseWriter, r *http.Request) { writeJSON(w, http.StatusOK, map[string]string{"status":"ok"}) })
	mux.HandleFunc("GET /api/phrases", a.listPhrases)
	mux.HandleFunc("GET /api/practice/new", a.newPractice)
	mux.HandleFunc("POST /api/practice/check", a.checkPractice)
	return withCORS(mux)
}

func (a *App) listPhrases(w http.ResponseWriter, r *http.Request) {
	level := strings.TrimSpace(r.URL.Query().Get("level"))
	result := make([]Phrase, 0, len(a.phrases))
	for _, p := range a.phrases { if level == "" || strings.EqualFold(level, p.Level) { result = append(result, p) } }
	writeJSON(w, http.StatusOK, result)
}

func (a *App) newPractice(w http.ResponseWriter, r *http.Request) {
	level := strings.TrimSpace(r.URL.Query().Get("level"))
	choices := make([]Phrase, 0, len(a.phrases))
	for _, p := range a.phrases { if level == "" || strings.EqualFold(level, p.Level) { choices = append(choices, p) } }
	if len(choices) == 0 { writeJSON(w, http.StatusNotFound, map[string]string{"error":"該当レベルの例文がありません"}); return }
	a.mu.Lock()
	p := choices[a.rng.Intn(len(choices))]
	p.Tokens = append([]string(nil), p.Tokens...)
	originalTokens := append([]string(nil), p.Tokens...)
	a.rng.Shuffle(len(p.Tokens), func(i, j int) { p.Tokens[i], p.Tokens[j] = p.Tokens[j], p.Tokens[i] })
	if len(p.Tokens) > 1 && equalTokens(p.Tokens, originalTokens) {
		// Avoid a trivially solved exercise when a random shuffle is unchanged.
		p.Tokens[0], p.Tokens[1] = p.Tokens[1], p.Tokens[0]
	}
	a.mu.Unlock()
	writeJSON(w, http.StatusOK, map[string]any{"phraseId":p.ID, "level":p.Level, "tokens":p.Tokens, "prompt":"語句を正しい順番に並べてください。"})
}

func (a *App) checkPractice(w http.ResponseWriter, r *http.Request) {
	var req CheckRequest
	decoder := json.NewDecoder(http.MaxBytesReader(w, r.Body, 16*1024))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil { writeJSON(w, http.StatusBadRequest, map[string]string{"error":"リクエストJSONを確認してください"}); return }
	if req.PhraseID == "" || len(req.Tokens) == 0 || len(req.Tokens) > 30 { writeJSON(w, http.StatusBadRequest, map[string]string{"error":"phraseIdと1〜30個のtokensが必要です"}); return }
	for _, p := range a.phrases {
		if p.ID != req.PhraseID { continue }
		correct, err := checkOrder(p.Tokens, req.Tokens)
		if err != nil { writeJSON(w, http.StatusServiceUnavailable, map[string]string{"error":"語順エンジンを実行できません"}); return }
		writeJSON(w, http.StatusOK, map[string]any{"correct":correct, "zh":p.ZH, "pinyin":p.Pinyin, "ja":p.JA})
		return
	}
	writeJSON(w, http.StatusNotFound, map[string]string{"error":"例文が見つかりません"})
}

func equalTokens(a, b []string) bool {
	if len(a) != len(b) { return false }
	for i := range a { if a[i] != b[i] { return false } }
	return true
}

// PHRASE_ORDER_BIN enables delegation to the Rust scorer. Arguments are passed
// directly (never through a shell); unit separator is reserved as a delimiter.
func checkOrder(expected, submitted []string) (bool, error) {
	bin := os.Getenv("PHRASE_ORDER_BIN")
	if bin == "" { return equalTokens(expected, submitted), nil }
	cmd := exec.Command(bin, strings.Join(expected, "\x1f"), strings.Join(submitted, "\x1f"))
	output, err := cmd.Output()
	if err != nil { return false, err }
	return strings.TrimSpace(string(output)) == "correct", nil
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:4173")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		if r.Method == http.MethodOptions { w.WriteHeader(http.StatusNoContent); return }
		next.ServeHTTP(w, r)
	})
}

func main() {
	dataPath := os.Getenv("PHRASES_FILE")
	if dataPath == "" { dataPath = filepath.Join("..", "..", "data", "phrases.json") }
	phrases, err := loadPhrases(dataPath)
	if err != nil { log.Fatalf("load phrases: %v", err) }
	app := &App{phrases:phrases, rng:rand.New(rand.NewSource(time.Now().UnixNano()))}
	addr := os.Getenv("API_ADDR")
	if addr == "" { addr = ":8080" }
	log.Printf("Phrase Lab API listening on %s (%d phrases)", addr, len(phrases))
	log.Fatal(http.ListenAndServe(addr, app.routes()))
}
