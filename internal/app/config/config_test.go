package config

import "fmt"

func ExampleNewConfig() {
	cfg := NewConfig()
	fmt.Println(cfg.ServerAddress)
	fmt.Println(*cfg.BaseURL)
	fmt.Println(*cfg.EnableHTTPS)

	// Output:
	// localhost:8080
	// http://localhost:8080
	// false

}
func ExampleGetServerAddress() {
	fmt.Println(GetServerAddress())
	// Output:
	// localhost:8080
}
func ExampleGetEnableHTTPS() {
	fmt.Println(GetEnableHTTPS())
	// Output:
	// false
}
func ExampleGetCertFile() {
	fmt.Println(GetCertFile())
	// Output:
	// cmd/shortener/certificate/localhost.crt
}
func ExampleGetKeyFile() {
	fmt.Println(GetKeyFile())
	// Output:
	// cmd/shortener/certificate/localhost.key
}
