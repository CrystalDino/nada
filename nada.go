package main

import "nada/controller"

func main() {
	// r := gin.New()

	// // Global middleware
	// r.Use(gin.Logger())
	// r.Use(gin.Recovery())

	// r.GET("/test_token", func(c *gin.Context) {
	// 	//token-------------------------------------
	// 	u := &models.User{
	// 		Id:          10101,
	// 		Name:        "dino",
	// 		Cell:        "13912345678",
	// 		Email:       "crystal.dino@hotmail.com",
	// 		CTime:       time.Now().Unix(),
	// 		MTime:       time.Now().Unix(),
	// 		LastLoginIp: "127.0.0.1",
	// 	}
	// 	token, err := core.TokenMake(u, "dtime", "password", "transcode", "stat")
	// 	log.Println(token, err)
	// 	du, err := core.TokenValidate(token)
	// 	log.Println(du, err)
	// 	c.JSON(http.StatusOK, gin.H{
	// 		"message": "pong",
	// 	})
	// })

	// v1 := r.Group("/v1")
	// v1.GET("/1", func(c *gin.Context) {
	// 	c.String(http.StatusOK, "%s", "there is v1 test")
	// })
	// v2 := r.Group("/v2")
	// v2.GET("/2", func(c *gin.Context) {
	// 	c.String(http.StatusOK, "%s", "there is v2 test")
	// })

	controller.Run()
}
