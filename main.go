package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"time"
)

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Max-Age", "86400")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

type Cassette struct {
	Number       int  `json:"number"`
	Denomination int  `json:"denomination"`
	Count        int  `json:"count"`
	IsWorking    bool `json:"is_working"`
}

type Request struct {
	Amount    int        `json:"amount"`
	Cassettes []Cassette `json:"cassettes"`
}

type Response struct {
	Success bool             `json:"success"`
	Notes   []CassetteOutput `json:"notes"`
	TimeMs  string           `json:"time_ms"`
	Message string           `json:"message"`
}

type CassetteOutput struct {
	Number       int `json:"number"`
	Denomination int `json:"denomination"`
	Count        int `json:"count"`
}

func calculateCash(amount int, cassettes []Cassette) ([]CassetteOutput, bool) {
	calcFactor := 100
	var result []CassetteOutput
	target := amount / calcFactor
	if amount%calcFactor != 0 {
		return result, false
	}
	totalSum := 0
	for _, i := range cassettes {
		if i.IsWorking {
			totalSum += (i.Denomination / calcFactor) * i.Count
		}
	}
	if totalSum < target {
		return result, false
	}

	minNotes := make([]int, target+1)
	for i := 1; i <= target; i++ {
		minNotes[i] = -1
	}
	minNotes[0] = 0

	usedNotes := make([]map[int]int, target+1)
	for i := range usedNotes {
		usedNotes[i] = make(map[int]int)
	}

	denomByNumber := make(map[int]int)
	for _, cas := range cassettes {
		denomByNumber[cas.Number] = cas.Denomination / calcFactor
	}

	denomToCassettes := make(map[int][]Cassette)
	for _, cas := range cassettes {
		if cas.IsWorking && cas.Count > 0 {
			denomToCassettes[cas.Denomination/calcFactor] = append(denomToCassettes[cas.Denomination/calcFactor], cas)
		}
	}

	for sum := 0; sum < target; sum++ {
		if minNotes[sum] == -1 {
			continue
		}

		for denom, cassetteList := range denomToCassettes {
			newSum := sum + denom
			if newSum > target {
				continue
			}

			canAddNote := false
			var selectedCassette Cassette

			for _, cas := range cassetteList {
				if usedNotes[sum][cas.Number] < cas.Count {
					canAddNote = true
					selectedCassette = cas
					break
				}
			}

			if canAddNote {
				if minNotes[newSum] == -1 || minNotes[newSum] > minNotes[sum]+1 {
					minNotes[newSum] = minNotes[sum] + 1
					usedNotes[newSum] = make(map[int]int)
					for key, val := range usedNotes[sum] {
						usedNotes[newSum][key] = val
					}
					usedNotes[newSum][selectedCassette.Number]++
				}
			}
		}
	}

	if minNotes[target] == -1 {
		return nil, false
	}

	for num, count := range usedNotes[target] {
		if count > 0 {
			result = append(result, CassetteOutput{
				Number:       num,
				Denomination: denomByNumber[num] * calcFactor,
				Count:        count,
			})
		}
	}
	sort.Slice(result, func(i, j int) bool {
		if result[i].Denomination != result[j].Denomination {
			return result[i].Denomination > result[j].Denomination
		}
		return result[i].Number > result[j].Number
	})
	return result, true
}

func handler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	var request Request
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if len(request.Cassettes) < 1 || len(request.Cassettes) > 8 {
		http.Error(w, "Количество кассет должно быть от 1 до 8", http.StatusBadRequest)
		return
	}

	notes, success := calculateCash(request.Amount, request.Cassettes)
	duration := time.Since(start).Seconds() * 1000

	response := Response{
		Success: success,
		Notes:   notes,
		TimeMs:  fmt.Sprintf("%.3f", duration),
	}

	if success {
		response.Message = "Выдача возможна"
	} else {
		response.Message = "Невозможно выдать запрашиваемую сумму"
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/calculate", handler)

	log.Println("Сервер запущен на :8080")
	http.ListenAndServe(":8080", corsMiddleware(mux))
}
