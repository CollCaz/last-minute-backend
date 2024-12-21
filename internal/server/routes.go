package server

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"

	"github.com/gin-gonic/gin"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := gin.Default()

	r.GET("/", s.HelloWorldHandler)
	r.POST("/", s.SubmitSolution)

	return r
}

func (s *Server) HelloWorldHandler(c *gin.Context) {
	resp := make(map[string]string)
	resp["message"] = "Hello World"

	c.JSON(http.StatusOK, resp)
}

func (s *Server) SubmitSolution(c *gin.Context) {
	// _, err := io.ReadAll(c.Request.Body)
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, err.Error())
	// 	return
	// }
	body := c.Request.Body
	codeBytes, err := io.ReadAll(body)
	if err != nil {
		log.Printf(err.Error())
	}
	f, _ := os.CreateTemp("", "*_code.cpp")
	if err != nil {
		log.Printf(err.Error())
		return
	}
	_, err = f.Write(codeBytes)
	if err != nil {
		log.Printf(err.Error())
		return
	}

	outFile, _ := os.CreateTemp("", "*_code.out")
	compileCmd := exec.Command("/bin/sh", "-c", fmt.Sprintf("g++ -o %s %s", outFile.Name(), f.Name()))
	runCmd := exec.Command(outFile.Name())
	_, _ = compileCmd.CombinedOutput()
	f.Close()
	outFile.Close()
	o, _ := runCmd.CombinedOutput()

	c.JSON(http.StatusOK, string(o))
}
