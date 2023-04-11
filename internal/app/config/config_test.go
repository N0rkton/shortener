package config

import "fmt"

func ExampleNewConfig() {
	cfg := NewConfig()
	fmt.Println(cfg.ServerAddress)
	fmt.Println(*cfg.BaseURL)

	// Output:
	// localhost:8080
	// http://localhost:8080

}
func ExampleGetServerAddress() {
	fmt.Println(GetServerAddress())
	// Output:
	// localhost:8080
}
