package filestorage

import (
	"fmt"
  	"os"
	"math/rand"
	"log"
	"io/ioutil"
)
func TmpRefresh() int {
	tmpDir := "./tmp"
	files, err := ioutil.ReadDir(tmpDir)
	if err != nil && !os.IsExist(err){
		log.Println("Ошибка при удалении существующих файлов: %s", err)
	}

	for _, file := range files {
		err := os.Remove(tmpDir + "/" + file.Name())
		if err != nil{
			log.Println("Ошибка при удалении существующих файлов: %s", err)
		}
	}
	log.Printf("Files deleted")

	err = os.Mkdir("tmp", 0755)
	if err != nil && !os.IsExist(err){
		panic(err)
	}
	
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890_*/")
	for i := 1; i <= 10; i++ {
		fileName := fmt.Sprintf("tmp/somefile%d.txt", i)
		text := make([]rune, 18)
		for j := range(text){
			text[j] = letters[rand.Intn(len(letters))]
		}
		err = os.WriteFile(fileName, []byte(string(text)), 0644)
		if err != nil {
			fmt.Println("Ошибка при создании файла:", err)
		}
	}

	flag := os.Getenv("FLAG")
	err = os.Mkdir("bills/admin", 0755)
	if err != nil && !os.IsExist(err){
		panic(err)
	  }
	err = os.WriteFile("bills/admin/admin.txt", []byte(flag), 0644)
	if err != nil {
		log.Println("Ошибка при создании флага:", err)
	}
	
	log.Printf("Created admin.txt")

	val := rand.Intn(10)
	return val
}
 
func Start(){ 
	// Создаем папку bills
	err := os.Mkdir("bills", 0755)
	if err != nil && !os.IsExist(err){
	  log.Println("Ошибка при создании папки:", err)
	  return
	}
  
	// Создаем папки и файлы
	for i := 1; i <= 10; i++ {
		name := fmt.Sprintf("Name%d", i)
		folderName := fmt.Sprintf("bills/Name%d", i)
		err := os.Mkdir(folderName, 0755)
		if err != nil && !os.IsExist(err){
			log.Println("Ошибка при создании папки:", err)
			continue
		}
		

		val := rand.Intn(1500)
		fileName := fmt.Sprintf("%s/bill_1.txt", folderName)
		content := fmt.Sprintf("Bill per %s in amount of %v", name, val)
		err = os.WriteFile(fileName, []byte(content), 0644)
		if err != nil {
			log.Println("Ошибка при создании файла:", err)
		}
	}

	log.Printf("Created and filled bills directory.")
	
}
