package filestorage

import (
	"fmt"
  	"os"
	"math/rand"
	"log"
)
func TmpRefresh(){
	err := os.Mkdir("tmp", 0755)
	if err != nil && !os.IsExist(err){
	  log.Println("Ошибка при создании папки:", err)
	  return
	}
	val := rand.Intn(100)
	for i := 1; i <= val; i++ {
		fileName := fmt.Sprintf("tmp/somefile%d.txt", i)
		err = os.WriteFile(fileName, []byte(""), 0644)
		if err != nil {
			fmt.Println("Ошибка при создании файла:", err)
		}
	}

	flag := os.Getenv("FLAG")
	err = os.Mkdir("bills/admin", 0755)
	if err != nil && !os.IsExist(err){
		log.Println("Ошибка при создании папки admin:", err)
		return
	  }
	err = os.WriteFile("bills/admin/admin.txt", []byte(flag), 0644)
	if err != nil {
		log.Println("Ошибка при создании флага:", err)
	}
	
	log.Printf("Created admin.txt")
}
 
func Start() { 
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

	TmpRefresh()
	
}
