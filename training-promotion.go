package main

import (
	"html/template"
	"net/http"
	"strings"
)

type ClientInfo struct {
	Name, Email, Phone, TrainingDate string
	AlreadyClient                    bool
}

var confirmedTrainings = make([]*ClientInfo, 0)
var templates = make(map[string]*template.Template, 3)

func loadTemplates() {
	templateNames := [5]string{"main-page", "registration-form", "training-confirmed", "training-rejected", "participants"}
	for _, name := range templateNames {
		if parsedTemplate, err := template.ParseFiles("layout.html", name+".html"); err == nil {
			templates[name] = parsedTemplate
		} else {
			panic(err)
		}
	}
}

func mainPageHandler(writer http.ResponseWriter, request *http.Request) {
	templates["main-page"].Execute(writer, nil)
}

func participantsHandler(writer http.ResponseWriter, request *http.Request) {
	templates["participants"].Execute(writer, confirmedTrainings)
}

type ClientInfoFormData struct {
	*ClientInfo
	Errors []string
}

func registrationFormHandler(writer http.ResponseWriter, request *http.Request) {
	if request.Method == http.MethodGet {
		templates["registration-form"].Execute(writer, ClientInfoFormData{
			ClientInfo: &ClientInfo{}, Errors: []string{},
		})
	} else if request.Method == http.MethodPost {
		request.ParseForm()
		clientInfo := ClientInfo{
			Name:          cleanString(request.Form["name"][0]),
			Email:         cleanString(request.Form["email"][0]),
			Phone:         cleanString(request.Form["phone"][0]),
			TrainingDate:  request.Form["training"][0],
			AlreadyClient: request.Form["already-client"][0] == "true",
		}

		if errors := getValidationErrors(clientInfo); len(errors) > 0 {
			templates["registration-form"].Execute(writer, ClientInfoFormData{ClientInfo: &clientInfo, Errors: errors})
		} else {
			saveClientTraining(&clientInfo)
			templates["training-confirmed"].Execute(writer, clientInfo.Name)
		}
	}
}

func saveClientTraining(clientInfo *ClientInfo) {
	confirmedTrainings = append(confirmedTrainings, clientInfo)
}

func trainingRejectedFormHandler(writer http.ResponseWriter, request *http.Request) {
	templates["training-rejected"].Execute(writer, ClientInfoFormData{ClientInfo: &ClientInfo{}, Errors: []string{}})
}

func cleanString(str string) string {
	return strings.Trim(str, " \t\n\r")
}

func getValidationErrors(clientInfo ClientInfo) []string {
	errors := []string{}
	if clientInfo.Name == "" {
		errors = append(errors, "Необходимо указать имя")
	}
	if clientInfo.Email == "" {
		errors = append(errors, "Необходимо указать e-mail")
	}
	if clientInfo.Phone == "" {
		errors = append(errors, "Необходимо указать номер телефона")
	}
	return errors
}

func insertDummyClients() {
	saveClientTraining(&ClientInfo{"Анна", "anna@gmail.com", "+12 (34) 56-78-9", "7.08, Бачата", true})
	saveClientTraining(&ClientInfo{"Рома Ш.", "roman.sharikov@gmail.com", "+12 (34) 24-18-4", "9.08, Силовая тренировка", false})
	saveClientTraining(&ClientInfo{"Михаил Капитанов", "mimisha@gmail.com", "+12 (34) 66-44-6", "9.08, Силовая тренировка", true})
	saveClientTraining(&ClientInfo{"Маруся", "m.svekla@gmail.com", "+12 (34) 14-94-5", "13.08, Зумба", true})
	saveClientTraining(&ClientInfo{"Лилия Свекла", "l.svekla@gmail.com", "+12 (34) 99-21-1", "13.08, Зумба", false})
}

func main() {
	insertDummyClients()

	loadTemplates()

	http.HandleFunc("/", mainPageHandler)
	http.HandleFunc("/registration-form", registrationFormHandler)
	http.HandleFunc("/training-rejected", trainingRejectedFormHandler)
	http.HandleFunc("/participants", participantsHandler)

	if err := http.ListenAndServe(":8008", nil); err != nil {
		panic(err)
	}
}
