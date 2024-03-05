// html.go

package newFolder

import (
	"html/template"
	"net/http"
)

// homeHandler handles the root endpoint
func HomeHandler(w http.ResponseWriter, r *http.Request) {
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
