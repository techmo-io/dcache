package main

func main() {
	d := NewHTTPDriver(NewStore(Clock{}))
	d.Run()
}
