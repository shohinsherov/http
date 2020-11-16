package app

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/shohinsherov/http/pkg/banners"
)

// Server предостовляет собой логический сервер нашего приложения
type Server struct {
	mux        *http.ServeMux
	bannersSvc *banners.Service
}

// NewServer - функция-конструктор для создания сервера.
func NewServer(mux *http.ServeMux, bannersSvc *banners.Service) *Server {
	return &Server{mux: mux, bannersSvc: bannersSvc}
}

func (s *Server) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	s.mux.ServeHTTP(writer, request)
}

// Init инициализирует сервер (регистрирует все Handler-ы)
func (s *Server) Init() {
	s.mux.HandleFunc("/banners.getAll", s.handleGetAllBanners)
	s.mux.HandleFunc("/banners.getById", s.handleGetPostByID)
	s.mux.HandleFunc("/banners.save", s.handleSaveBanner)
	s.mux.HandleFunc("/banners.removeById", s.handleremoveByID)
	s.mux.HandleFunc("/process", s.process)
}

//  get all
func (s *Server) handleGetAllBanners(writer http.ResponseWriter, request *http.Request) {
	b, err := s.bannersSvc.All(request.Context())
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
		return
	}

	data, err := json.Marshal(b)
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	log.Print("ready")
	_, err = writer.Write([]byte(data))
	if err != nil {
		log.Print(err)
	}

}

// get by id
func (s *Server) handleGetPostByID(writer http.ResponseWriter, request *http.Request) {
	idParam := request.URL.Query().Get("id")

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	item, err := s.bannersSvc.ByID(request.Context(), id)
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(item)
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Contetn-Type", "applicatrion/json")
	_, err = writer.Write(data)
	if err != nil {
		log.Print(err)
	}

}

// add or update
func (s *Server) handleSaveBanner(writer http.ResponseWriter, request *http.Request) {

	idParam := request.FormValue("id")
	title := request.FormValue("title")
	content := request.FormValue("content")
	button := request.FormValue("button")
	link := request.FormValue("link")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	imageFile, imageHead, imageErr := request.FormFile("image")
	if imageErr != nil {
		banner := banners.Banner{
			ID:      id,
			Title:   title,
			Content: content,
			Button:  button,
			Link:    link,
			Image:   "",
		}
		bannerRes, err := s.bannersSvc.Save(request.Context(), &banner)
		if err != nil {
			log.Print(err)
			http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		data, err := json.Marshal(bannerRes)
		if err != nil {
			log.Print(err)
			http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		writer.Header().Set("Contetn-Type", "applicatrion/json")
		_, err = writer.Write(data)
		if err != nil {
			log.Print(err)
		}
		return
	}

	fileName := imageHead.Filename
	file, err := ioutil.ReadAll(imageFile)
	if err != nil {
		log.Print(err)
	}

	banner := banners.Banner{
		ID:      id,
		Title:   title,
		Content: content,
		Button:  button,
		Link:    link,
		Image:   fileName,
	}
	bannerRes, err := s.bannersSvc.Save(request.Context(), &banner)
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	ioutil.WriteFile("web/banners/"+bannerRes.Image, file, 0666)
	data, err := json.Marshal(bannerRes)
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	writer.Header().Set("Contetn-Type", "applicatrion/json")
	_, err = writer.Write(data)
	if err != nil {
		log.Print(err)
	}

}

// delete banner byID
func (s *Server) handleremoveByID(writer http.ResponseWriter, request *http.Request) {
	idParam := request.URL.Query().Get("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	dBanner, err := s.bannersSvc.RemoveByID(request.Context(), id)
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(dBanner)
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	writer.Header().Set("Contetn-Type", "applicatrion/json")
	_, err = writer.Write(data)
	if err != nil {
		log.Print(err)
	}
}

func (s *Server) process(writer http.ResponseWriter, request *http.Request) {
	log.Print(request.RequestURI) // полный урл
	log.Print(request.Method)     // метод
	/*	log.Print(request.Header)                     // все заголовки
		log.Print(request.Header.Get("Content-Type")) // конкретный заголовок

		log.Print(request.FormValue("tags"))     // только первое значение Query + POST
		log.Print(request.PostFormValue("tags")) // только первое значение POST*/

	body, err := ioutil.ReadAll(request.Body) // теле запроса
	if err != nil {
		log.Print(err)
	}
	log.Printf("%s", body)

	/*err = request.ParseMultipartForm(10 * 1024 * 1024)  // 10MB
	if err != nil {
		log.Print(err)
	}

	// доступно только после ParseForm (либо FormValue, PostFormValue)
	log.Print(request.Form)     // все значения формы
	log.Print(request.PostForm) // все значения формы

	// доступно только после ParseMultipart (FormValue, PostFromValue автоматически вызывают ParseMultipartForm)
	log.Print(request.FormFile("image"))
	// request.MultipartForm.Value - только "обычные поля"
	// request.MultipartForm.File - только файлы*/

}
