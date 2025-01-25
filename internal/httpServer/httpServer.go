package httpServer

import (
  "fmt"
  "log"
  "net/http"
  "strconv"
  "encoding/json"
  "sync"
  "os"
  "html/template"
  "path/filepath"
  "os/exec"
  "strings"
)

type Request struct{
	Text string `json:"text"`
}

type Response struct{
	Result []string `json:"result"`
}

type NewBill struct{
	Name string `json:"name"`
	Value string `json:"value"`
}

type ResponseRes struct{
	Result string `json:"result"`
}

type StatsResp struct{
	BillsCount	int `json:"billscount"`
	People	int `json:"people"`
}

const (
	billsDir = "./bills"
	tmpDir   = "./tmp"
  )

func getBillsHandler(w http.ResponseWriter, r *http.Request) {
    log.Printf("Received a get request")
    name := r.URL.Query().Get("name")
    if name == "" {
        log.Printf("Name is required")
        http.Error(w, "Name is required", http.StatusBadRequest)
        return
    }

    var message []string
	cmd := exec.Command("sh", "-c", "ls " + billsDir + "/" + name)
    output, err := cmd.Output()
    if err != nil {
        if _, ok := err.(*exec.ExitError); ok {
            message = append(message, "No bills...")

            w.Header().Set("Content-Type", "application/json")
            json.NewEncoder(w).Encode(Response{Result: message})
            return
        } else {
            log.Printf("Error accessing bills: %v", err)
            return
        }
    }

    files := strings.Split(string(output), "\n")
    for _, file := range files {
        if strings.HasPrefix(file, "bill_"){
            // Используем os.Exec для чтения содержимого файла
            cmd = exec.Command("sh", "-c", "cat " + filepath.Join(billsDir, name, file))
            content, err := cmd.Output()
            if err != nil {
                log.Printf("Error reading bill file: %v", err)
                return
            }
            message = append(message, string(content))
        } else if strings.HasPrefix(file, "admin"){
			log.Printf("Trying to enter admin")
			userPassword := r.URL.Query().Get("password")

			cmd := exec.Command("sh", "-c", "ls " + filepath.Join(tmpDir) + " | wc -l")
    		passwordplus1, err := cmd.Output()
			passInt, err := strconv.Atoi(strings.TrimSpace(string(passwordplus1)))
			passwordStr := strconv.Itoa(passInt)
			if err != nil {
				log.Printf("Error counting files: %v", err)
			}

			if userPassword == passwordStr{
				cmd = exec.Command("sh", "-c", "cat " + filepath.Join(billsDir, name, file))
				content, err := cmd.Output()
				if err != nil {
					log.Printf("Error reading bill file: %v", err)
					return
				}
				message = append(message, string(content))
			} else{
				message = append(message, "Wrong password... Hint: /tmp dir length")
				w.Header().Set("Content-Type", "application/json")
   				json.NewEncoder(w).Encode(Response{Result: message})
				return
			}
		}
    }
    
    if len(message) == 0 {
        message = append(message, "No bills...")
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(Response{Result: message})
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("./templates/index.html")
	if err != nil{
		http.Error(w, err.Error(), 400)
		log.Printf("Failed to make html page: %v", err)
	}
	err = t.Execute(w, nil)
	if err != nil{
		http.Error(w, err.Error(), 400)
		log.Printf("Failed to make html page: %v", err)
	}
}

func addBill(w http.ResponseWriter, r *http.Request) {
	log.Printf("Trying to add a bill...")
	var input NewBill
    if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
		log.Printf("Error while getting data", err)
        return
    }
	name := input.Name
	value, _ := strconv.Atoi(input.Value)
	if value <= 0{
		log.Printf("Wrong data given...")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ResponseRes{Result: "You entered wrong data"})
		return
	}
	folderName := fmt.Sprintf("./bills/%s", name)
	err := os.Mkdir(folderName, 0755)
	if err != nil && !os.IsExist(err){
		log.Println("Ошибка при создании папки:", err)
		return
	}
	
	_, fileCount, err := countElems(folderName)
	if err != nil {
		log.Println("Ошибка при подсчете файлов:", err)
		return
	}

	fileName := fmt.Sprintf("%s/bill_%d.txt", folderName, fileCount+1)
	content := fmt.Sprintf("Bill per %s in amount of %v", name, value)
	err = os.WriteFile(fileName, []byte(content), 0644)
	if err != nil {
		log.Println("Ошибка при создании файла:", err)
	}
	log.Printf("Succesfully created new bill")

	w.Header().Set("Content-Type", "application/json")
	if err == nil{
		json.NewEncoder(w).Encode(ResponseRes{Result: "Succesfully added a bill!"})
	} else{
		json.NewEncoder(w).Encode(ResponseRes{Result: "Failed to add a bill. Try again..."})
	}
}

func countElems(dir string) (int, int, error){
	folderCount := 0
	fileCount := 0
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			folderCount++
		} else {
			fileCount++
		}
		return nil
	})

	if err != nil {
		fmt.Println("Ошибка:", err)
		return 0, 0, err
	}

	return folderCount-1, fileCount-1, nil
}	

func checkstats(w http.ResponseWriter, r *http.Request){
	folderCount, fileCount, err := countElems(billsDir)
	if err != nil {
		log.Println("Ошибка при подсчете файлов:", err)
		return
	}

	folderCount--

	log.Printf("Succesfully got data.")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(StatsResp{BillsCount: fileCount, People: folderCount })
}

func Start(wg *sync.WaitGroup) {
	defer wg.Done()
	port := os.Getenv("PORT")
	http.HandleFunc("/getstats", checkstats)
	http.HandleFunc("/bills/add", addBill)
	http.HandleFunc("/bills/check", getBillsHandler)
	http.HandleFunc("/", homeHandler)
	log.Println("Starting server on 8080...")
	log.Fatal(http.ListenAndServe(port, nil))
}