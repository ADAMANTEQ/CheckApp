package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net"
	"net/http"
	"os"
	"sync"
	"time"
)

type Point struct {
	Name         string    `json:"name"`
	IP           string    `json:"ip"`
	Address      string    `json:"address"`
	Status       bool      `json:"status"`
	StatusChange time.Time `json:"status_change"`
	OpenTime     time.Time `json:"open_time"`
}

var (
	points      []Point
	pointsMutex sync.Mutex
	whitelist   = map[string]bool{
		"46.147.129.30": true, "46.147.129.83": true, "46.147.129.73": true, "84.201.254.98": true, "46.147.129.29": true,
		"94.25.182.111": true, "84.201.245.231": true, "84.201.245.243": true, "84.201.244.155": true, "46.147.129.14": true,
		"84.201.249.43": true, "46.147.129.50": true, "84.201.246.151": true, "188.233.237.1": true, "46.147.129.49": true,
		"84.201.243.92": true, "46.147.129.21": true, "94.25.190.50": true, "84.201.246.79": true, "46.147.129.58": true,
		"46.147.129.72": true, "84.201.246.176": true, "46.147.129.79": true, "86.140.1.19": true, "46.147.129.55": true,
		"94.25.182.105": true, "84.201.242.38": true, "46.147.129.74": true, "46.147.129.25": true, "46.147.129.20": true,
		"78.85.35.221": true, "84.201.244.146": true, "84.201.254.70": true, "84.201.245.233": true, "185.93.254.68": true,
		"192.168.0.13": true, "92.55.4.50": true, "84.201.242.213": true}
	username = "admin"
	password = "666666"
	logger   *log.Logger // Глобальная переменная для логгера
)

func basicAuth(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()

		if !ok || user != username || pass != password {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		h(w, r)
	}
}

func init() {
	initialPoints := []Point{
		{Name: "9 Января", IP: "46.147.129.30", Address: "г. Ижевск, ул. 9 Января, 219А", Status: false, StatusChange: time.Now()},
		{Name: "Автовокзал", IP: "46.147.129.83", Address: "г. Ижевск, ул. Красноармейская, 127", Status: false, StatusChange: time.Now()},
		{Name: "Аврора", IP: "185.93.254.68", Address: "г. Ижевск, ул. Удмуртская, 304", Status: false, StatusChange: time.Now()},
		{Name: "Аксион", IP: "192.168.0.137", Address: "г. Ижевск, ул. Карла Маркса, 191", Status: false, StatusChange: time.Now()},
		{Name: "Гранат", IP: "46.147.129.73", Address: "г. Ижевск, ул. Ворошилова, 44а", Status: false, StatusChange: time.Now()},
		{Name: "Европа", IP: "84.201.254.98", Address: "г. Ижевск, ул. имени Вадима Сивкова, 150", Status: false, StatusChange: time.Now()},
		{Name: "ЖД Вокзал", IP: "46.147.129.29", Address: "г. Ижевск, ул. Гагарина, 29", Status: false, StatusChange: time.Now()},
		{Name: "Завьялово", IP: "94.25.182.111", Address: "УР, с. Завьялово, ул. Калинина, 5", Status: false, StatusChange: time.Now()},
		{Name: "Италмас", IP: "84.201.245.231", Address: "г. Ижевск, ул. Татьяны Барамзиной, 42", Status: false, StatusChange: time.Now()},
		{Name: "Кембридж", IP: "84.201.245.243", Address: "г. Ижевск, ул. 7-ая Подлесная, 96", Status: false, StatusChange: time.Now()},
		{Name: "Кит", IP: "84.201.244.155", Address: "г. Ижевск, ул. Ленина, 138", Status: false, StatusChange: time.Now()},
		{Name: "Клубная", IP: "192.168.0.100", Address: "г. Ижевск, ул. Клубная, 51", Status: false, StatusChange: time.Now()},
		{Name: "Кольцо", IP: "46.147.129.14", Address: "г. Ижевск, ул. 50 лет ВЛКСМ, 49", Status: false, StatusChange: time.Now()},
		{Name: "Кунгурцева", IP: "84.201.249.43", Address: "г. Ижевск, ул. Воткинское Шоссе, 39", Status: false, StatusChange: time.Now()},
		{Name: "Ленина 26", IP: "46.147.129.50", Address: "г. Ижевск, ул. Ленина, 26", Status: false, StatusChange: time.Now()},
		{Name: "Магнит", IP: "84.201.246.151", Address: "г. Ижевск, ул. Молодежная, 107а", Status: false, StatusChange: time.Now()},
		{Name: "Малина", IP: "188.233.237.1", Address: "г. Ижевск, Крылова, 20", Status: false, StatusChange: time.Now()},
		{Name: "Матрица", IP: "46.147.129.49", Address: "г. Ижевск, ул. Баранова, 87", Status: false, StatusChange: time.Now()},
		{Name: "Медведь", IP: "84.201.243.92", Address: "г. Ижевск, ул. Пушкинская, 161а", Status: false, StatusChange: time.Now()},
		{Name: "Меркурий", IP: "46.147.129.21", Address: "г. Ижевск, ул. Клубная, 23", Status: false, StatusChange: time.Now()},
		{Name: "Молодежная", IP: "94.25.190.50", Address: "г. Ижевск, ул. Молодежная, 10", Status: false, StatusChange: time.Now()},
		{Name: "Океан", IP: "84.201.246.79", Address: "г. Ижевск, ул. Кирова, 109", Status: false, StatusChange: time.Now()},
		{Name: "Острова", IP: "46.147.129.58", Address: "г. Ижевск, 50 лет ВЛКСМ, 6", Status: false, StatusChange: time.Now()},
		{Name: "Пушкинская 114", IP: "46.147.129.72", Address: "г. Ижевск, ул. Пушкинская, 114", Status: false, StatusChange: time.Now()},
		{Name: "П-2", IP: "84.201.246.176", Address: "г. Ижевск, ул. Петрова, 27а", Status: false, StatusChange: time.Now()},
		{Name: "Парашют", IP: "46.147.129.79", Address: "г. Ижевск, ул. Азина, 220", Status: false, StatusChange: time.Now()},
		{Name: "Парма", IP: "84.201.242.213", Address: "г. Ижевск, ул. Клубная, д. 67А", Status: false, StatusChange: time.Now()},
		{Name: "Петрова KFC", IP: "86.140.1.19", Address: "г. Ижевск, ул. Петрова, 1", Status: false, StatusChange: time.Now()},
		{Name: "Престиж", IP: "46.147.129.55", Address: "г. Ижевск, ул. Ленина, 112", Status: false, StatusChange: time.Now()},
		{Name: "Радиотехника", IP: "94.25.182.105", Address: "г. Ижевск, ул. Ленина, 140", Status: false, StatusChange: time.Now()},
		{Name: "Северный", IP: "84.201.242.38", Address: "г. Ижевск, ул. Буммашевская, 92в", Status: false, StatusChange: time.Now()},
		{Name: "Солнечный", IP: "46.147.129.74", Address: "г. Ижевск, ул. 5-ая Подлесная, 34", Status: false, StatusChange: time.Now()},
		{Name: "Техно", IP: "46.147.129.25", Address: "г. Ижевск, ул. Подлесная 7-Я, 34", Status: false, StatusChange: time.Now()},
		{Name: "Флагман", IP: "46.147.129.20", Address: "г. Ижевск, ул. Удмуртская, 255И", Status: false, StatusChange: time.Now()},
		{Name: "Фурманова", IP: "78.85.35.221", Address: "г. Ижевск, ул. Азина, 340 ", Status: false, StatusChange: time.Now()},
		{Name: "Холмы", IP: "192.168.0.8", Address: "г. Ижевск, ул. Пушкинская, 290", Status: false, StatusChange: time.Now()},
		{Name: "ЦУМ", IP: "84.201.244.146", Address: "г. Ижевск, ул. Карла Маркса, 244", Status: false, StatusChange: time.Now()},
		{Name: "Эльгрин", IP: "84.201.254.70", Address: "г. Ижевск, ул. 10 лет Октября, 53", Status: false, StatusChange: time.Now()},
		{Name: "Южный", IP: "84.201.245.233", Address: "г. Ижевск, ул. Маяковского, 47", Status: false, StatusChange: time.Now()},
		{Name: "Test", IP: "92.55.4.50", Address: "Test", Status: false, StatusChange: time.Now()},
	}

	points = append(points, initialPoints...)
}

func getAllPointsHandler(w http.ResponseWriter, r *http.Request) {
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

func checkInHandler(w http.ResponseWriter, r *http.Request) {
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

func resetStatus() {
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

func getNetworkIP() (string, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "", err
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP.String(), nil
}

// getClientIP возвращает IP-сеть клиента из запроса
func getClientIP(r *http.Request) (string, error) {
	// Извлекаем IP-адрес из заголовка X-Real-IP или X-Forwarded-For
	ip := r.Header.Get("X-Real-IP")
	if ip == "" {
		ip = r.Header.Get("X-Forwarded-For")
	}

	// Если IP-адрес не найден в заголовках, используем RemoteAddr
	if ip == "" {
		// Убедитесь, что RemoteAddr не является пустым
		if addr, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
			ip = addr
		}
	}

	// Преобразуем строку в объект IP
	ipAddress := net.ParseIP(ip)
	if ipAddress == nil {
		return "", fmt.Errorf("неверный IP-адрес: %s", ip)
	}

	// Преобразуем IP-адрес в строку
	ipString := ipAddress.String()

	return ipString, nil
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	pointsMutex.Lock()
	defer pointsMutex.Unlock()

	openCount := 0
	closedCount := 0

	for _, point := range points {
		if point.Status {
			openCount++
		} else {
			closedCount++
		}
	}

	tmpl, err := template.New("index").Parse(`
	<!DOCTYPE html>
	<html lang="en">

	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>Список Смен</title>
		<link href="https://fonts.googleapis.com/css?family=Roboto:300,400,500,700&display=swap" rel="stylesheet">
		<link href="https://cdnjs.cloudflare.com/ajax/libs/animate.css/4.1.1/animate.min.css" rel="stylesheet">
		<link href="https://cdnjs.cloudflare.com/ajax/libs/hover.css/2.3.1/css/hover-min.css" rel="stylesheet">
		<link href="https://cdnjs.cloudflare.com/ajax/libs/twitter-bootstrap/5.3.0/css/bootstrap.min.css" rel="stylesheet">
		<style>
			body {
				font-family: 'Roboto', sans-serif;
				margin: 0;
				padding: 0;
				display: flex;
				flex-direction: column;
				align-items: center;
				min-height: 100vh;
				position: relative;
				background-color: #f8f9fa;
			}

			#particles-js {
				position: fixed;
				width: 100%;
				height: 100%;
				top: 0;
				left: 0;
				z-index: -1;
			}

			header {
				background-color: transparent;
				color: dark;
				text-align: center;
				padding: 1em;
				width: 100%;
				box-sizing: border-box;
				border-bottom: 2px solid #343a40;
				margin-bottom: 20px;
				animation: fadeInDown 1s ease-in-out;
			}

			.container {
				width: 90%;
				overflow-x: auto;
				max-width: 100%;
				border-radius: 10px;
				box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
				background-color: rgba(255, 255, 255, 0.9);
				padding: 20px;
			}

			table {
				width: 100%;
				border-collapse: collapse;
				margin-top: 20px;
				background-color: #ffffff;
				border-radius: 10px;
			}

			th,
			td {
				padding: 12px;
				text-align: left;
				border-bottom: 1px solid #dee2e6;
				color: #333333;
			}

			th {
				background-color: #343a40;
				color: white;
				font-size: 1.2em;
				text-transform: uppercase;
				letter-spacing: 1px;
				border-top-left-radius: 10px;
				border-top-right-radius: 10px;
			}

			tbody tr:hover {
				background-color: #f0f0f0;
			}

			.status-button {
				padding: 8px 0;
				border-radius: 5px;
				font-size: 0.9em;
				font-weight: bold;
				text-align: center;
				cursor: pointer;
				transition: background-color 0.3s ease;
				display: block;
				width: 100%;
			}

			.status-button:hover {
				transform: scale(1.1);
				background-color: #6c757d;
			}

			.status-button-true {
				background-color: #20c997;
				color: white;
			}

			.status-button-false {
				background-color: #f56565;
				color: white;
			}

			.counter-container {
				display: flex;
				justify-content: space-around;
				align-items: center;
				width: 100%;
				margin-bottom: 20px;
			}

			.counter {
				padding: 8px 16px;
				background-color: #343a40;
				color: white;
				border-radius: 5px;
			}

			.counter span {
				font-size: 1.2em;
			}

			footer {
				text-align: center;
				font-size: 0.8em;
				color: #6c757d;
				margin-top: 20px;
			}
		</style>
	</head>

	<body>
		<div id="particles-js"></div>
		<header class="animate__animated animate__fadeInDown">
			<h1 class="text-3xl">Список Смен</h1>
		</header>
		<div class="container">
			<div class="counter-container">
				<div class="counter">
					Смен открыто: <span class="open-counter">{{.OpenCount}}</span>
				</div>
				<div class="counter">
					Смен закрыто: <span class="closed-counter">{{.ClosedCount}}</span>
				</div>
			</div>
			<table class="table table-striped">
				<thead>
					<tr>
						<th class="px-4 py-2">Название</th>
						<th class="px-4 py-2">Адрес</th>
						<th class="px-4 py-2">Статус</th>
						<th class="px-4 py-2">Время изменения статуса</th>
						<th class="px-4 py-2">Время открытия</th>
					</tr>
				</thead>
				<tbody>
					{{range .Points}}
					<tr>
						<td class="border px-4 py-2">{{.Name}}</td>
						<td class="border px-4 py-2">{{.Address}}</td>
						<td class="border px-4 py-2">
							<div class="status-button {{if .Status}}status-button-true{{else}}status-button-false{{end}}">
								{{if .Status}}Открыта{{else}}Закрыта{{end}}
							</div>
						</td>
						<td class="border px-4 py-2">{{.StatusChange.Format "02.01.2006 15:04:05"}}</td>
						<td class="border px-4 py-2">{{.OpenTime.Format "02.01.2006 15:04:05"}}</td>
					</tr>
					{{end}}
				</tbody>
			</table>
		</div>
		<footer>
			Developed by ADMNT for A-VAPE 2024 ©
		</footer>

		<!-- Подключение библиотек и скриптов -->
		<script src="https://cdn.jsdelivr.net/particles.js/2.0.0/particles.min.js"></script>
		<script>
			particlesJS("particles-js", {
				"particles": {
					"number": {
						"value": 80,
						"density": {
							"enable": true,
							"value_area": 800
						}
					},
					"color": {
						"value": "#000000"
					},
					"shape": {
						"type": "circle",
						"stroke": {
							"width": 0,
							"color": "#000000"
						},
						"polygon": {
							"nb_sides": 5
						},
						"image": {
							"src": "img/github.svg",
							"width": 100,
							"height": 100
						}
					},
					"opacity": {
						"value": 0.5,
						"random": false,
						"anim": {
							"enable": false,
							"speed": 1,
							"opacity_min": 0.1,
							"sync": false
						}
					},
					"size": {
						"value": 5,
						"random": true,
						"anim": {
							"enable": false,
							"speed": 40,
							"size_min": 0.1,
							"sync": false
						}
					},
					"line_linked": {
						"enable": true,
						"distance": 150,
						"color": "#000000",
						"opacity": 0.4,
						"width": 1
					},
					"move": {
						"enable": true,
						"speed": 6,
						"direction": "none",
						"random": false,
						"straight": false,
						"out_mode": "out",
						"bounce": false,
						"attract": {
							"enable": false,
							"rotateX": 600,
							"rotateY": 1200
						}
					}
				},
				"interactivity": {
					"detect_on": "canvas",
					"events": {
						"onhover": {
							"enable": true,
							"mode": "grab"
						},
						"onclick": {
							"enable": true,
							"mode": "push"
						},
						"resize": true
					},
					"modes": {
						"grab": {
							"distance": 140,
							"line_linked": {
								"opacity": 1
							}
						},
						"bubble": {
							"distance": 400,
							"size": 40,
							"duration": 2,
							"opacity": 8,
							"speed": 3
						},
						"repulse": {
							"distance": 200,
							"duration": 0.4
						},
						"push": {
							"particles_nb": 4
						},
						"remove": {
							"particles_nb": 2
						}
					}
				},
				"retina_detect": true
			});
		</script>
	</body>

	</html>
	`)
	if err != nil {
		http.Error(w, "Ошибка при создании HTML-шаблона: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Передаем данные в шаблон
	data := struct {
		Points      []Point
		OpenCount   int
		ClosedCount int
	}{
		Points:      points,
		OpenCount:   openCount,
		ClosedCount: closedCount,
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Ошибка при выполнении HTML-шаблона: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func main() {
	// Открыть файл для записи лога
	logFile, err := os.OpenFile("server.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("Не удалось открыть файл лога:", err)
		return
	}
	defer logFile.Close()

	// Настроить логгер на использование файла лога
	logger = log.New(logFile, "SERVER: ", log.Ldate|log.Ltime|log.Lshortfile)

	go resetStatus()

	http.HandleFunc("/getallpoints", getAllPointsHandler)
	http.HandleFunc("/checkin", checkInHandler)
	http.HandleFunc("/", basicAuth(homeHandler))

	port := 8080
	fmt.Println("Developed by ADMNT for A-VAPE 2024 ©")
	fmt.Println("")
	fmt.Printf("Check-In Server запущен на порту: %d\n", port)
	logger.Printf("Check-In Server запущен на порту: %d\n", port) // Записать в лог информацию о запуске сервера

	err = http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		fmt.Println("Ошибка при запуске сервера:", err)
		logger.Println("Ошибка при запуске сервера:", err) // Записать в лог информацию об ошибке при запуске сервера
	}
}
