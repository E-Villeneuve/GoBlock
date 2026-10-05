package main

func main() {
	//bc := NewBlockchain()
	//defer bc.db.Close() // Waits until the end of the main function to close the database connection

	cli := CLI{}
	cli.Run()
}
