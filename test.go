package main

// import "fmt"

// func match_course(data *map[string][]string) map[string]string {
// 	t2c := make(map[string]string)
// 	c2t := make(map[string]string)
// 	for k := range *data {
// 		if _, is_exist := t2c[k]; !is_exist {
// 			check := make(map[string]struct{})
// 			dfs(k, data, &check, &t2c, &c2t)
// 		}
// 	}
// 	return t2c
// }

// func dfs(now string, data *map[string][]string, check *map[string]struct{}, t2c *map[string]string, c2t *map[string]string) bool {
// 	for _, i := range (*data)[now] {
// 		if _, is_exist := (*check)[i]; !is_exist {
// 			(*check)[i] = struct{}{}
// 			if nxt, is_exist := (*c2t)[i]; !is_exist || dfs(nxt, data, check, t2c, c2t) {
// 				(*c2t)[i] = now
// 				(*t2c)[now] = i
// 				return true
// 			}
// 		}
// 	}
// 	return false
// }

// func main() {
// 	m := make(map[string][]string)
// 	m["1"] = []string{"7"}
// 	m["7"] = []string{"123", "2", "3", "7"}
// 	m["2"] = []string{"2", "3", "7"}
// 	fmt.Println(match_course(&m))
// }
