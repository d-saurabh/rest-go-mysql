package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

//template structure
type App struct {
	Router *mux.Router
	Db *sql.DB
}

//functions of App struct
func (a *App) Init(user,password,db string){
	var err error
	a.Db, err = sql.Open("mysql",fmt.Sprintf("%s:%s@tcp(localhost:3306)/%s",user,password,db))

	if err!=nil{
		log.Fatal(err)
	}
	//start the router
	a.Router = mux.NewRouter()

	//initialize the end points
	a.InitializeEndPoints()
}

func (a *App) Start(route string) {
	log.Fatal(http.ListenAndServe(":8000", a.Router))
}

func (a *App) InitializeEndPoints()  {
	a.Router.HandleFunc("/order",a.createOrder).Methods("POST")
	a.Router.HandleFunc("/orders",a.getOrders).Methods("GET")
	a.Router.HandleFunc("/order/{id:[0-9]+}",a.getOrder).Methods("GET")
	a.Router.HandleFunc("/order/{id:[0-9]+}", a.updateOrder).Methods("PUT")
	a.Router.HandleFunc("/order/{id:[0-9]+}", a.deleteOrder).Methods("DELETE")
}

//endpoints handlers
func (a *App) getOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])

	if err !=nil{
		logError(w,http.StatusBadRequest,err.Error())
		return
	}

	o := order{Id:id}

	obj,err:= o.getOrderById(a.Db)
	if  err != nil {
		logError(w, http.StatusBadRequest, err.Error())
		return
	}
	jsonResponse(w,http.StatusOK,obj)
}

func (a *App) getOrders(w http.ResponseWriter,r *http.Request) {

	orders, err := getOrders(a.Db)
	if err!=nil{
		logError(w,http.StatusBadRequest,"Some error")
	}

	jsonResponse(w,http.StatusOK,orders)
}

func (a *App) createOrder(w http.ResponseWriter, r *http.Request){
	var o order
	decoder := json.NewDecoder(r.Body)
	if err:= decoder.Decode(&o); err!= nil{
		logError(w, http.StatusBadRequest,"invalid json")
		return
	}

	defer r.Body.Close()

	if err := o.createOrder(a.Db); err!= nil{
		logError(w, http.StatusInternalServerError,err.Error())
		return
	}

	jsonResponse(w, http.StatusCreated, o)

}

func (a *App) deleteOrder(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		logError(w, http.StatusBadRequest, "Invalid id")
		return
	}

	o := order{Id:id}
	if err := o.deleteOrder(a.Db); err != nil {
		logError(w, http.StatusInternalServerError, err.Error())
		return
	}

	jsonResponse(w, http.StatusOK, map[string]string{"message": "deleted successfully"})
}

func (a *App) updateOrder(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		logError(w, http.StatusBadRequest, "Invalid id")
		return
	}

	var o order
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&o); err != nil {
		logError(w, http.StatusBadRequest, "Invalid resquest")
		return
	}
	defer r.Body.Close()
	o.Id = id

	if err := o.updateOrder(a.Db); err != nil {
		logError(w, http.StatusInternalServerError, err.Error())
		return
	}

	jsonResponse(w, http.StatusOK, o)
}

//helpers to log back the error and response
func logError(w http.ResponseWriter,code int,message string)  {
	jsonResponse(w, code, map[string]string{"error": message})
}

func jsonResponse(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}