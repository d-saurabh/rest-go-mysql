# rest-go-mysql
An end to end Go REST api with MySql.

## Purpose and functionality
This article will be a good start to understand how to implement a REST api's using Go with Mysql.
This is a very basic introduction though meaning ful interms of concepts. This REST api exposes endpoints for an Order service, which includes basic CRUD operations as follows:

1. [GET] /order/{id}
2. [GET] /orders
3. [POST] /order
4. [PUT] /order/{id}
5. [DELETE] /order/{id}


##Database script

We will create a simple MySql table called ```orders``` with following fields
1. id - int type auto incremented column also serve as a Primary Key
2. title - title of the order
3. status - status of the order

The Db script to create a table: 
```mysql
CREATE TABLE orders ( 
id int not null auto_increment, 
title varchar(50) not null, 
status bool,
constraint pk_example primary key (id) 
);
```

##Dependencies

1. Gorilla Mux Routes **mux**
2. MySql driver **mysql**

```
go get github.com/gorilla/mux github.com/go-sql-driver/mysql
```
**Note**: Install mysql to $GOPATH in case it gives any problem 


##Scaffolding

This is how your project structure should look like 

```flow js
┌── app.go
├── main.go
└── dbmodel.go
``` 
The app.go file will hold reference to our main libraries 

```go
app.go 

type App struct {
	Router *mux.Router
	Db *sql.DB
}
```

Add two functions called Init and Start as follows:

```go
app.go 

func (a *App) Init(user,password,db string){
}

func (a *App) Start(route string) {
}
```
Init will initialize a Db connection and Start will begin execution of our application

The final app.go looks like this 
```go
package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

//template structure
type App struct {
	Router *mux.Router
	Db *sql.DB
}

//functions of App struct
func (a *App) Init(user,password,db string){
}

func (a *App) Start(route string) {
}
```
We will add three environment variables called DB_USER_NAME,DB_PASSWORD and DB_NAME.
Now we will create a main.go file which is the entry point for our application:

```go
main.go

package main

import "os"

func main() {
	a := App()
	a.Init(os.Getenv("DB_USER_NAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"))

	a.Start(":8000")

}
```

Lets now add a db model representing our table order. We will also add the basic CRUD methods for order entity

```go
dbmodel.go 

package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

type order struct{
	Id int `json:"id"`
	Title string `json:"title"`
	Status bool `json:"status"`
}

func (o *order) createOrder(repository *sql.DB) error{
}

func (o *order) getOrderById(repository *sql.DB) (order,error){
}
func (o *order) updateOrder(repository *sql.DB) error{
}

func (o *order) getOrder(repository *sql.DB) ([]order,error){
}

func (o *order) deleteOrder(repository *sql.DB) error{	
}
```

Lets add some real code for the database access layer

```go
dbmodel.go

package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)


type order struct{
	Id int `json:"id"`
	Title string `json:"title"`
	Status bool `json:"status"`
}

func (o *order) createOrder(repository *sql.DB) error{
	insert,err := repository.Query("INSERT INTO orders(title,status) VALUES(?,?)",o.Title,o.Status)

	if err!=nil{
		return err
	}

	defer insert.Close()

	return nil
}

func (o *order) updateOrder(repository *sql.DB) error{
	_, err := repository.Exec("UPDATE orders SET title=?, status=? WHERE id=?",o.Title,o.Status,o.Id)

	if err!=nil{
		return err
	}

	return  nil
}

func (o *order) getOrderById(repository *sql.DB) (order,error){
	var obj order
	err := repository.QueryRow("SELECT * FROM orders WHERE id=?",o.Id).Scan(&obj.Id,&obj.Title,&obj.Status)

	if err!=nil{
		return obj,err
	}

	return obj,nil

}
func getOrders(repository *sql.DB) ([]order,error){
	orders := []order{}
	rows,err := repository.Query("SELECT * from orders")

	if err!=nil{
		return orders,err
	}

	defer rows.Close()

	for rows.Next()  {
		var temp order

		err = rows.Scan(&temp.Id,&temp.Title,&temp.Status)
		checkErr(err)
		orders = append(orders, temp)
	}

	return orders, nil
}

func (o *order) deleteOrder(repository *sql.DB) error{
	_, err := repository.Exec("DELETE FROM orders WHERE id=?",o.Id)
	if err!=nil{
		return err
	}
	return nil
}
```

Now lets add endpoints and handlers for them
api execute the following command in
#####1. Create Order
```go
app.go

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
```
#####2. Get Order
```go
app.go

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
```
#####3. Get Order(s)
```go
app.go


func (a *App) getOrders(w http.ResponseWriter,r *http.Request) {

	orders, err := getOrders(a.Db)
	if err!=nil{
		logError(w,http.StatusBadRequest,"Some error")
	}

	jsonResponse(w,http.StatusOK,orders)
}
```
#####4. Update Order
```go
app.go

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
```
#####5. Delete Order
```go
app.go

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
```

#####Some helper methods

These helper methods will be used to log errors and response back to the caller
```go
app.go

func logError(w http.ResponseWriter,code int,message string)  {
	jsonResponse(w, code, map[string]string{"error": message})
}

func jsonResponse(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
```
With this simply run the api and test it with tools like Postman/Insomnia

                                    Thank You!           







