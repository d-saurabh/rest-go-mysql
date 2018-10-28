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
		if err !=nil{
			return orders,err
		}
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