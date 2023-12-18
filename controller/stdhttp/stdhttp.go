package stdhttp

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"notes/gates/storage"
	"notes/models/dto"
	"notes/pkg"
	"strconv"
)

type Controller struct {
	srv  http.Server
	stor storage.Storage
}

func NewController(addr string, st storage.Storage) (hs *Controller) {
	hs = new(Controller)
	hs.srv = http.Server{}
	mux := http.NewServeMux()

	mux.HandleFunc("/save", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		if r.Method != http.MethodPost {
			log.Println("Method Not Allowed", http.StatusMethodNotAllowed)
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		hs.NoteSaveHandler(w, r)
	})
	mux.HandleFunc("/read", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		if r.Method != http.MethodPost {
			log.Println("Method Not Allowed", http.StatusMethodNotAllowed)
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		hs.NoteReadHandler(w, r)
	})
	mux.HandleFunc("/update", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		if r.Method != http.MethodPost {
			log.Println("Method Not Allowed", http.StatusMethodNotAllowed)
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		hs.NoteUpdateHandler(w, r)
	})
	mux.HandleFunc("/delete", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		if r.Method != http.MethodPost {
			log.Println("Method Not Allowed", http.StatusMethodNotAllowed)
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		hs.NoteDeleteHandler(w, r)
	})

	hs.srv.Handler = mux
	hs.srv.Addr = addr
	hs.stor = st
	return hs
}

func (hs *Controller) Start() error {
	myErr := pkg.NewMyError("package stdhttp: func (hs *Controller) Start()")
	err := hs.srv.ListenAndServe()
	if err != nil {
		log.Fatal(myErr.Wrap(err, "hs.srv.ListenAndServe()").Error())
		return myErr.Wrap(err, "hs.srv.ListenAndServe()")
	}
	return nil
}

// NoteSaveHandler обрабатывает HTTP запрос для добавления новой записки.
func (hs *Controller) NoteSaveHandler(w http.ResponseWriter, req *http.Request) {
	myErr := pkg.NewMyError("package stdhttp: func (hs *Controller) NoteSaveHandler(w http.ResponseWriter, req *http.Request)")
	w.WriteHeader(http.StatusOK) // Status 200 OK
	response := dto.Response{}
	note, err := GetBody(req)
	if err != nil {
		e := myErr.Wrap(err, "GetBody(req)")
		response.Result = "Error"
		response.Error = e.Error()
		js, erro := response.GetJson()
		if erro != nil {
			log.Println(erro.Error())
			w.Write(json.RawMessage(`{"result":"Error","data":{},"error":"` + erro.Error() + `"}`))
			return
		}
		w.Write(js)
		log.Println(e.Error())
		return
	}
	/// Проверяем нет ли пустых значений
	if note.Name == "" || note.LastName == "" || note.Text == "" {
		e := myErr.Wrap(nil, "All fields must be filled in")
		response.Result = "Error"
		response.Error = e.Error()
		js, erro := response.GetJson()
		if erro != nil {
			log.Println(erro.Error())
			w.Write(json.RawMessage(`{"result":"Error","data":{},"error":"` + erro.Error() + `"}`))
		}
		w.Write(js)
		log.Println(e.Error())
		return
	}
	index, err := hs.stor.Add(note)
	if err != nil {
		e := myErr.Wrap(err, "")
		response.Result = "Error"
		response.Error = e.Error()
		js, erro := response.GetJson()
		if erro != nil {
			log.Println(erro.Error())
			w.Write(json.RawMessage(`{"result":"Error","data":{},"error":"` + erro.Error() + `"}`))
		}
		w.Write(js)
		log.Println(e.Error())
		return
	}
	// Возвращаем результат и index созданной заметки
	response.Result = "Success"
	response.Data = json.RawMessage(fmt.Sprintf(`"%d"`, index))
	js, erro := response.GetJson()
	if erro != nil {
		log.Println(erro.Error())
		w.Write(json.RawMessage(`{"result":"Error","data":{},"error":"` + erro.Error() + `"}`))
	}
	w.Write(js)
}

// NoteReadHandler обрабатывает HTTP запрос для чтения записок.
func (hs *Controller) NoteReadHandler(w http.ResponseWriter, req *http.Request) {
	myErr := pkg.NewMyError("package stdhttp: func (hs *Controller) NoteReadHandler(w http.ResponseWriter, req *http.Request)")
	w.WriteHeader(http.StatusOK) // Status 200 OK
	response := dto.Response{}
	byteReq, err := io.ReadAll(req.Body)
	req.Body.Close()
	if err != nil {
		e := myErr.Wrap(err, "io.ReadAll(req.Body)")
		response.Result = "Error"
		response.Error = e.Error()
		js, erro := response.GetJson()
		if erro != nil {
			log.Println(erro.Error())
			w.Write(json.RawMessage(`{"result":"Error","data":{},"error":"` + erro.Error() + `"}`))
			return
		}
		w.Write(js)
		log.Println(e.Error())
		return
	}
	intValue, err := strconv.ParseInt(string(byteReq), 10, 64)
	if err != nil {
		e := myErr.Wrap(err, "strconv.ParseInt(string(byteData), 10, 64)")
		response.Result = "Error"
		response.Error = e.Error()
		js, erro := response.GetJson()
		if erro != nil {
			log.Println(erro.Error())
			w.Write(json.RawMessage(`{"result":"Error","data":{},"error":"` + erro.Error() + `"}`))
			return
		}
		w.Write(js)
		log.Println(e.Error())
		return
	}
	data, status := hs.stor.GetByIndex(intValue)
	if !status {
		e := myErr.Wrap(nil, "something going wrong in hs.stor.GetByIndex(intValue)")
		response.Result = "Error"
		response.Error = e.Error()
		js, erro := response.GetJson()
		if erro != nil {
			log.Println(erro.Error())
			w.Write(json.RawMessage(`{"result":"Error","data":{},"error":"` + erro.Error() + `"}`))
			return
		}
		log.Println(e.Error())
		w.Write(js)
		return
	}

	/// Формируем ответ
	w.Header().Set("Content-Type", "application/json")
	jsonData, err := json.Marshal(data) ///////////////////////////// interface
	if err != nil {
		e := myErr.Wrap(err, "json.Marshal(result)")
		response.Result = "Error"
		response.Error = e.Error()
		js, erro := response.GetJson()
		if erro != nil {
			log.Println(erro.Error())
			w.Write(json.RawMessage(`{"result":"Error","data":{},"error":"` + erro.Error() + `"}`))
			return
		}
		w.Write(js)
		log.Println(e.Error())
		return
	}
	response.Result = "Success"
	response.Data = jsonData
	js, err := response.GetJson()
	if err != nil {
		log.Println(err.Error())
		w.Write(json.RawMessage(`{"result":"Error","data":{},"error":"` + err.Error() + `"}`))
		return
	}
	w.Write(js)

}

// NoteUpdateHandler обрабатывает HTTP запрос для обновления записки.
func (hs *Controller) NoteUpdateHandler(w http.ResponseWriter, req *http.Request) {
	myErr := pkg.NewMyError("package stdhttp: func (hs *Controller) NoteUpdateHandler(w http.ResponseWriter, req *http.Request)")
	w.WriteHeader(http.StatusOK) // Status 200 OK
	response := dto.Response{}

	byteReq, err := io.ReadAll(req.Body)
	req.Body.Close()
	if err != nil {
		e := myErr.Wrap(err, "io.ReadAll(req.Body)")
		response.Result = "Error"
		response.Error = e.Error()
		js, erro := response.GetJson()
		if erro != nil {
			log.Println(erro.Error())
			w.Write(json.RawMessage(`{"result":"Error","data":{},"error":"` + erro.Error() + `"}`))
			return
		}
		w.Write(js)
		log.Println(e.Error())
		return
	}

	var requestData struct {
		Index int64    `json:"index,string"`
		Data  dto.Note `json:"data"`
	}

	// Декодируем данные в созданную структуру
	err = json.Unmarshal(byteReq, &requestData)
	if err != nil {
		e := myErr.Wrap(err, "json.Unmarshal(byteReq, &requestData)")
		response.Result = "Error"
		response.Error = e.Error()
		js, erro := response.GetJson()
		if erro != nil {
			log.Println(erro.Error())
			w.Write(json.RawMessage(`{"result":"Error","data":{},"error":"` + erro.Error() + `"}`))
			return
		}
		w.Write(js)
		log.Println(e.Error())
		return
	}

	// Для обновления заметки объединим все 3 запроса, то есть будем эту заметку получать, считывать из неё данные,
	// актуализировать из, удалять старую заметку и создавать новую

	data, status := hs.stor.GetByIndex(requestData.Index)
	if !status {
		e := myErr.Wrap(nil, "something going wrong in hs.stor.GetByIndex(intValue)")
		response.Result = "Error"
		response.Error = e.Error()
		js, erro := response.GetJson()
		if erro != nil {
			log.Println(erro.Error())
			w.Write(json.RawMessage(`{"result":"Error","data":{},"error":"` + erro.Error() + `"}`))
			return
		}
		w.Write(js)
		log.Println(e.Error())
		return
	}
	// Преобразуем интерфейс в структуру Note
	note, ok := data.(*dto.Note)
	if !ok {
		http.Error(w, "Unexpected type", http.StatusInternalServerError)
		e := myErr.Wrap(nil, "Unexpected type in data.(*dto.Note)")
		response.Result = "Error"
		response.Error = e.Error()
		js, erro := response.GetJson()
		if erro != nil {
			log.Println(erro.Error())
			w.Write(json.RawMessage(`{"result":"Error","data":{},"error":"` + erro.Error() + `"}`))
			return
		}
		w.Write(js)
		log.Println(e.Error())
		return
	}

	if requestData.Data.Name != "" {
		note.Name = requestData.Data.Name
	}

	if requestData.Data.LastName != "" {
		note.LastName = requestData.Data.LastName
	}

	if requestData.Data.Text != "" {
		note.Text = requestData.Data.Text
	}

	index, err := hs.stor.Add(note)
	if err != nil {
		e := myErr.Wrap(err, "")
		response.Result = "Error"
		response.Error = e.Error()
		js, erro := response.GetJson()
		if erro != nil {
			log.Println(erro.Error())
			w.Write(json.RawMessage(`{"result":"Error","data":{},"error":"` + erro.Error() + `"}`))
		}
		w.Write(js)
		log.Println(e.Error())
		return
	}

	hs.stor.RemoveByIndex(requestData.Index)

	response.Result = "Success"
	response.Data = json.RawMessage(fmt.Sprintf(`"%d"`, index))
	js, erro := response.GetJson()
	if erro != nil {
		log.Println(erro.Error())
		w.Write(json.RawMessage(`{"result":"Error","data":{},"error":"` + erro.Error() + `"}`))
	}
	w.Write(js)
}

// NoteDeleteHandler обрабатывает HTTP запрос для удаления записи по номеру телефона.
func (hs *Controller) NoteDeleteHandler(w http.ResponseWriter, req *http.Request) {
	myErr := pkg.NewMyError("package stdhttp: func (hs *Controller) NoteDeleteHandler(w http.ResponseWriter, req *http.Request)")
	w.WriteHeader(http.StatusOK) /// Status 200 OK
	response := dto.Response{}
	byteReq, err := io.ReadAll(req.Body)
	req.Body.Close()
	if err != nil {
		e := myErr.Wrap(err, "io.ReadAll(req.Body)")
		response.Result = "Error"
		response.Error = e.Error()
		js, erro := response.GetJson()
		if erro != nil {
			log.Println(erro.Error())
			w.Write(json.RawMessage(`{"result":"Error","data":{},"error":"` + erro.Error() + `"}`))
			return
		}
		w.Write(js)
		log.Println(e.Error())
		return
	}
	intValue, err := strconv.ParseInt(string(byteReq), 10, 64)
	if err != nil {
		e := myErr.Wrap(err, "strconv.ParseInt(string(byteData), 10, 64)")
		response.Result = "Error"
		response.Error = e.Error()
		js, erro := response.GetJson()
		if erro != nil {
			log.Println(erro.Error())
			w.Write(json.RawMessage(`{"result":"Error","data":{},"error":"` + erro.Error() + `"}`))
			return
		}
		w.Write(js)
		log.Println(e.Error())
		return
	}
	hs.stor.RemoveByIndex(intValue)

	response.Result = "Success"
	js, err := response.GetJson()
	if err != nil {
		log.Println(err.Error())
		w.Write(json.RawMessage(`{"result":"Error","data":{},"error":"` + err.Error() + `"}`))
		return
	}
	w.Write(js)
}

// / Функция для преобразования response в структуру Note
func GetBody(req *http.Request) (note *dto.Note, e error) {
	myErr := pkg.NewMyError("package stdhttp: func GetBody(req *http.Request)")
	byteReq, err := io.ReadAll(req.Body)
	req.Body.Close()
	if err != nil {
		e = myErr.Wrap(err, "io.ReadAll(req.Body)")
		log.Println(e.Error())
		return nil, e
	}
	err = json.Unmarshal(byteReq, &note)
	if err != nil {
		e = myErr.Wrap(err, "json.Unmarshal(byteReq, &note)")
		log.Println(e.Error())
		return nil, e
	}
	return note, nil
}
