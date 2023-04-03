package config

import "fmt"

func ExampleNewConfig() {
	cfg := NewConfig()
	fmt.Println(cfg)

}
func ExampleGetServerAddress() {
	GetServerAddress()
	//Output:
	//localhost:8080
}
