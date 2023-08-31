/*
 */
package main

import "fmt"

func main() {
	fmt.Println("sqlite3://na:na@nowhere/C:\\Temp\\sqlite\\G2C.db")
	fmt.Println(`sqlite3://na:na@nowhere/C:\\Temp\\sqlite\\G2C.db`)
	fmt.Println(`sqlite3://na:na@nowhere/C:\Temp\sqlite\G2C.db`)
}
