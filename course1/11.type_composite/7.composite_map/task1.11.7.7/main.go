package main


type User struct {
	Nickname string
	Age      int
	Email    string
}

func main() {
	users := []User{
		{"Heller9409", 58, "chesleyhartmann@tromp.io"},
		{"Heller9409", 58, "chesleyhartmann@tromp.io"},
		{"Hermiston3154", 55, "desmondgerhold@schimmel.name"},
		{"Hermiston3154", 55, "desmondgerhold@schimmel.name"},
		{"Prohaska8581", 41, "trinitybrakus@boehm.com"},
		{"Prohaska8581", 41, "trinitybrakus@boehm.com"},
		{"Sporer9227", 25, "hanswaelchi@lakin.info"},
		{"Grimes6069", 35, "edwinajones@crona.org"},
		{"Sporer9227", 25, "hanswaelchi@lakin.info"},
		{"Grimes6069", 35, "edwinajones@crona.org"},
	}
	getUniqueUsers(users)
}

func getUniqueUsers(users []User) []User {
	uniqueuserSl := []User{}
	uniqueuserMap := make(map[string]int)
	for _, val := range users {
		if uniqueuserMap[val.Nickname] == 0 {
			uniqueuserSl = append(uniqueuserSl, val)
			uniqueuserMap[val.Nickname]++
		}
	}
	return uniqueuserSl
	
}