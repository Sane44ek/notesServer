package main

import (
	"flag"
	"log"
	"notes/controller/stdhttp"
	"notes/gates/storage"
	"notes/gates/storage/list"
	"notes/gates/storage/mp"
)

func main() {
	var st storage.Storage
	// Определение флагов
	flagList := flag.Bool("l", false, "Use list storage")
	flagMap := flag.Bool("m", false, "Use map storage")

	// Разбор аргументов командной строки
	flag.Parse()

	// Проверка, что установлен один из флагов
	if !(*flagList || *flagMap) || (*flagList && *flagMap) {
		log.Println("Error: you must add flag -l to use list storage or flag -m to use map storage")
		flag.PrintDefaults()
		return
	}

	// Если есть другие аргументы, выдаем ошибку
	if len(flag.Args()) > 0 {
		log.Println("Error: bad arguments")
		flag.PrintDefaults()
		return
	}

	// Проверка значений флагов
	if *flagList {
		st = list.NewList()
	}
	if *flagMap {
		st = mp.NewMap()
	}
	hs := stdhttp.NewController(":4040", st)
	hs.Start()

}
