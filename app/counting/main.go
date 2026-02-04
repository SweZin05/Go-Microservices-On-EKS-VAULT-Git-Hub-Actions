package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

var (
	counter     int
	counterLock sync.Mutex
)

type CountData struct {
	Service   string    `json:"service"`
	Count     int       `json:"count"`
	Timestamp time.Time `json:"timestamp"`
	Message   string    `json:"message"`
}

const htmlTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Counting Service</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            min-height: 100vh;
            display: flex;
            align-items: center;
            justify-content: center;
            padding: 20px;
        }
        .container {
            background: white;
            border-radius: 20px;
            box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
            padding: 60px 40px;
            max-width: 600px;
            width: 100%;
            text-align: center;
        }
        h1 {
            color: #667eea;
            font-size: 2.5em;
            margin-bottom: 10px;
            font-weight: 700;
        }
        .subtitle {
            color: #666;
            font-size: 1.1em;
            margin-bottom: 40px;
        }
        .counter-display {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
            font-size: 5em;
            font-weight: bold;
            padding: 40px;
            border-radius: 15px;
            margin: 30px 0;
            box-shadow: 0 10px 30px rgba(102, 126, 234, 0.3);
            transition: transform 0.3s ease;
        }
        .counter-display:hover {
            transform: scale(1.05);
        }
        .info {
            background: #f8f9fa;
            padding: 20px;
            border-radius: 10px;
            margin: 20px 0;
            color: #495057;
        }
        .info strong {
            color: #667eea;
        }
        .buttons {
            margin-top: 30px;
            display: flex;
            gap: 15px;
            justify-content: center;
            flex-wrap: wrap;
        }
        button {
            background: #667eea;
            color: white;
            border: none;
            padding: 15px 30px;
            border-radius: 10px;
            font-size: 1.1em;
            cursor: pointer;
            transition: all 0.3s ease;
            font-weight: 600;
        }
        button:hover {
            background: #764ba2;
            transform: translateY(-2px);
            box-shadow: 0 5px 15px rgba(102, 126, 234, 0.4);
        }
        button.reset {
            background: #dc3545;
        }
        button.reset:hover {
            background: #c82333;
        }
        .timestamp {
            color: #999;
            font-size: 0.9em;
            margin-top: 20px;
        }
        .badge {
            display: inline-block;
            background: #28a745;
            color: white;
            padding: 5px 15px;
            border-radius: 20px;
            font-size: 0.9em;
            margin-bottom: 20px;
        }
    </style>
    <script>
        function refreshCount() {
            window.location.reload();
        }
        function resetCount() {
            fetch('/reset')
                .then(() => window.location.reload())
                .catch(err => console.error('Error resetting counter:', err));
        }
        // Auto-refresh every 5 seconds
        setTimeout(() => {
            window.location.reload();
        }, 5000);
    </script>
</head>
<body>
    <div class="container">
        <div class="badge">🚀 Service Active</div>
        <h1>📊 Counting Service</h1>
        <p class="subtitle">Microservice Counter Demo</p>
        
        <div class="counter-display">
            {{.Count}}
        </div>
        
        <div class="info">
            <p><strong>Service:</strong> {{.Service}}</p>
            <p><strong>Message:</strong> {{.Message}}</p>
        </div>
        
        <div class="buttons">
            <button onclick="refreshCount()">🔄 Refresh Count</button>
            <button class="reset" onclick="resetCount()">🔄 Reset Counter</button>
        </div>
        
        <div class="timestamp">
            Last updated: {{.Timestamp.Format "2006-01-02 15:04:05"}}
        </div>
    </div>
</body>
</html>`

func incrementAndGet() int {
	counterLock.Lock()
	defer counterLock.Unlock()
	counter++
	return counter
}

func resetCounter() int {
	counterLock.Lock()
	defer counterLock.Unlock()
	oldCount := counter
	counter = 0
	return oldCount
}

func getCounter() int {
	counterLock.Lock()
	defer counterLock.Unlock()
	return counter
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	currentCount := incrementAndGet()

	data := CountData{
		Service:   "Counting Service",
		Count:     currentCount,
		Timestamp: time.Now(),
		Message:   fmt.Sprintf("This service has been called %d time(s)", currentCount),
	}

	// Check Accept header for JSON vs HTML
	if r.Header.Get("Accept") == "application/json" {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(data)
		return
	}

	// Render HTML
	tmpl, err := template.New("home").Parse(htmlTemplate)
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		log.Printf("Template parse error: %v", err)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tmpl.Execute(w, data); err != nil {
		log.Printf("Template execution error: %v", err)
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "healthy",
		"service": "counting",
	})
}

func resetHandler(w http.ResponseWriter, r *http.Request) {
	oldCount := resetCounter()

	data := map[string]interface{}{
		"service":        "Counting Service",
		"message":        "Counter reset",
		"previous_count": oldCount,
		"current_count":  0,
		"timestamp":      time.Now(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "9003"
	}

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/reset", resetHandler)

	addr := ":" + port
	log.Printf("Counting Service starting on port %s", port)
	log.Printf("Access the service at http://localhost:%s", port)

	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal(err)
	}
}
