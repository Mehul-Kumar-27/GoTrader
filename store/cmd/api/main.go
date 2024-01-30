package main

import "gotrader/store/cmd/api/consumer"

//// This consumer has mainly two functions provide the data to the store and also add the data to a postgres database

// / First make a consumer to listen the specific topic and then add the data to the postgres database and
// / then add the data to the store
func main() {
	//Lets start the listining to the kafka topis
	consumer.ConsumeMessages()
}
