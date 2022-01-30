// package main

// import "fmt"

// func match_course(data map[string][]string) map[string]string {
// 	t2c := make(map[string]string)
// 	c2t := make(map[string]string)
// 	for k := range data {
// 		if _, temp := t2c[k]; !temp {
// 			check := make(map[string]struct{})
// 			dfs(k, &data, &check, &t2c, &c2t)
// 		}
// 	}
// 	return t2c
// }

// func dfs(now string, data *map[string][]string, check *map[string]struct{}, t2c *map[string]string, c2t *map[string]string) bool {
// 	for _, i := range (*data)[now] {
// 		if _, temp := (*check)[i]; !temp {
// 			(*check)[i] = struct{}{}
// 			if nxt, temp := (*c2t)[i]; !temp || dfs(nxt, data, check, t2c, c2t) {
// 				(*c2t)[i] = now
// 				(*t2c)[now] = i
// 				return true
// 			}
// 		}
// 	}
// 	return false
// }
