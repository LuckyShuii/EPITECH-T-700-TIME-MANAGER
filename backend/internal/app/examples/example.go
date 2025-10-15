package example

func Greet(name string) string {
	if name == "" {
		return "Hello, world!"
	}
	return "Hello, " + name + "!"
}
