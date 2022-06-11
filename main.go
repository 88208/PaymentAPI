package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
	"golang.org/x/net/context"
	"net/http"
)

type details struct {
	ID      int32  `json:"id" `
	Email   string `json:"email"`
	Payment int    `json:"payment"`
}

func main() {

	conn, err := connectDB()
	if err != nil {
		return
	}

	router := gin.Default()

	router.Use(dbMiddleware(*conn))
	router.POST("sum", SumOfPayments)
	router.GET("getcustomer", GetCustomer)
	router.GET("getpayment", GetPayment)

	gin.SetMode(gin.ReleaseMode)
	router.Run(":3000")
}

func connectDB() (c *pgx.Conn, err error) {
	conn, err := pgx.Connect(context.Background(), "postgresql://postgres:8822rj@localhost:5432/postgres")
	if err != nil || conn == nil {
		fmt.Println("Error connecting to DB")
		fmt.Println(err.Error())
	}
	_ = conn.Ping(context.Background())
	return conn, err
}

func dbMiddleware(conn pgx.Conn) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("db", conn)
		c.Next()
	}
}
func SumOfPayments(c *gin.Context) {
	db, _ := c.Get("db")
	conn := db.(pgx.Conn)

	var sum details
	email := c.PostForm("email")
	sum.Email = email
	//sum.Payment = payment

	rows, err := conn.Query(context.Background(), "SELECT sum(payment) FROM customer,payment WHERE email = $1", sum.Email)
	if err != nil {
		fmt.Printf("error in query")
		fmt.Println(err)
	}

	var Details []details
	for rows.Next() {
		data := details{}
		err = rows.Scan(&data.Payment)

		if err != nil {
			fmt.Printf("error in scan ")
			fmt.Println(err)
			continue
		}

		Details = append(Details, data)

	}

	if err != nil {
		fmt.Println(err)
	}
	c.JSON(http.StatusOK, gin.H{"salaries": Details})
}
func GetCustomer(c *gin.Context) {
	db, _ := c.Get("db")
	conn := db.(pgx.Conn)
	rows, err := conn.Query(context.Background(), "SELECT email FROM customer limit 10")
	if err != nil {
		fmt.Printf("error in query")
		fmt.Println(err)
	}

	var Details []details
	for rows.Next() {
		data := details{}
		err = rows.Scan(&data.Email)
		if err != nil {
			fmt.Printf("error in scan ")
			fmt.Println(err)
			continue
		}
		Details = append(Details, data)
	}
	if err != nil {
		fmt.Println(err)
	}
	c.JSON(http.StatusOK, gin.H{"items": Details})
}

func GetPayment(c *gin.Context) {
	db, _ := c.Get("db")
	conn := db.(pgx.Conn)
	rows, err := conn.Query(context.Background(), "SELECT  payment FROM payment")
	if err != nil {
		fmt.Printf("error in query")
		fmt.Println(err)
	}

	var Details []details
	for rows.Next() {
		data := details{}
		err = rows.Scan(&data.Payment)
		if err != nil {
			fmt.Printf("error in scan ")
			fmt.Println(err)
			continue
		}
		Details = append(Details, data)
	}
	if err != nil {
		fmt.Println(err)
	}
	c.JSON(http.StatusOK, gin.H{"items": Details})
}
