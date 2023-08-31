/*
 */
package main

import "fmt"

func main() {
	fmt.Print("sqlite3://na:na@nowhere/C:\\Temp\\sqlite\\G2C.db")
	fmt.Print(`sqlite3://na:na@nowhere/C:\\Temp\\sqlite\\G2C.db`)
	fmt.Print(`sqlite3://na:na@nowhere/C:\Temp\sqlite\G2C.db`)
}
