// point.go

package newFolder

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"time"
)

// Point represents a point data structure
type Point struct {
	Name         string    `json:"name"`
	IP           string    `json:"ip"`
	Address      string    `json:"address"`
	Status       bool      `json:"status"`
	StatusChange time.Time `json:"status_change"`
	OpenTime     time.Time `json:"open_time"`
}

func GetAllPointsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	pointsMutex.Lock()
	defer pointsMutex.Unlock()

	responseJSON, err := json.Marshal(points)
	if err != nil {
		http.Error(w, "Ошибка при кодировании JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseJSON)
}

func CheckInHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	// Получение IP-адреса клиента
	clientIP, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		http.Error(w, "Ошибка при получении IP-адреса клиента", http.StatusInternalServerError)
		return
	}

	// Проверка, есть ли IP-адрес в белом списке
	if _, ok := whitelist[clientIP]; ok {
		pointsMutex.Lock()
		defer pointsMutex.Unlock()

		for i := range points {
			if points[i].IP == clientIP {
				// Если смена уже открыта, возвращаем ошибку
				if points[i].Status {
					http.Error(w, "Смена уже открыта", http.StatusConflict)
					return
				}

				// Устанавливаем статус, время открытия и изменения статуса
				points[i].Status = true
				points[i].OpenTime = time.Now()
				points[i].StatusChange = time.Now()
				fmt.Printf("Точка %s: Смена открыта\n", points[i].Name)
				break
			}
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Статус успешно обновлен"))
	} else {
		http.Error(w, "IP-адрес не включен в белый список", http.StatusUnauthorized)
	}
}

func ResetStatus() {
	for {
		time.Sleep(1 * time.Minute) // Проверка каждую минуту

		pointsMutex.Lock()
		for i := range points {
			// Если смена открыта и прошло более 10 часов с момента открытия, закрываем смену
			if points[i].Status && time.Since(points[i].OpenTime) >= 10*time.Hour {
				points[i].Status = false
				points[i].StatusChange = time.Now()
				fmt.Printf("Точка %s: Смена закрыта (прошло более 10 часов)\n", points[i].Name)
			}
		}
		pointsMutex.Unlock()
	}
}
