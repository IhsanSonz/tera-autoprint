package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/gin-gonic/gin"
)

type RequestBody struct {
	PrinterLocation string
	PDF             []byte
}

var serviceStatus = "Hello from github.com/IhsanSonz :8088 :)"

func isReadable(filename string) bool {
	f, err := os.Open(filename)
	if err != nil {
		return false
	}
	defer f.Close()
	return true
}

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	// Define an API endpoint to identify that the service is running
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, serviceStatus)
	})

	r.POST("/autoprint-pdf", func(c *gin.Context) {
		_, header, err := c.Request.FormFile("pdf")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": true,
				"msg":   "failed create pdf",
				"hint":  "params error: invalid file",
			})
			return
		}

		if header.Header.Get("Content-Type") != "application/pdf" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": true,
				"msg":   "failed create pdf",
				"hint":  "params error: invalid type for file pdf (application/pdf needed)",
			})
			return
		}

		pdfDir := "./pdf"
		if _, err := os.Stat(pdfDir); os.IsNotExist(err) {
			err := os.MkdirAll(pdfDir, 0777)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": true,
					"msg":   "failed create pdf",
					"hint":  "failed to create pdf directory",
				})
				return
			}
		}

		filename := time.Now().Format("20060102150405") + fmt.Sprintf("%05d", rand.Intn(90000)+1000) + ".pdf"
		filename = filepath.Join(pdfDir, filename)
		err = c.SaveUploadedFile(header, filename)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": true,
				"msg":   "failed create pdf",
				"hint":  "failed to save uploaded file",
			})
			return
		}

		if _, err := os.Stat(filename); os.IsNotExist(err) || !isReadable(filename) {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": true,
				"msg":   "failed create pdf",
				"hint":  "failed to find the uploaded file",
			})
			return
		}

		printTo := c.Request.FormValue("print_to")
		if printTo == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": true,
				"msg":   "failed create pdf",
				"hint":  "params error: invalid print_to",
			})
			return
		}

		printSettings := c.Request.FormValue("print_settings")
		if printSettings == "" {
			printSettings = "1x"
		}

		command := fmt.Sprintf(".\\SumatraPDF-3.3.3-64.exe -print-settings \"%s\" -print-to \"%s\" .\\%s", printSettings, printTo, filename)

		// Process the valid input
		log.Println("pdf:", filename)
		log.Println("printTo: ", printTo)
		log.Println("printSettings: ", printSettings)
		log.Println("command:", command)

		err = exec.Command("powershell", "-Command", command).Run()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": true,
				"msg":   "failed to print pdf",
				"hint":  "failed to execute print command",
			})
			return
		}

		err = os.Remove(filename)
		if err != nil {
			log.Println(err)
		}

		c.JSON(http.StatusOK, gin.H{
			"error": false,
			"msg":   "AutoPrint PDF activated",
		})
	})

	// Start the API server in a separate Goroutine
	go func() {
		err := r.Run(":8088") // Change the port as needed
		if err != nil {
			fmt.Println("Error starting the API server:", err)
		} else {
			fmt.Println("Service running in port 8088")
		}
	}()

	// Create a Fyne GUI
	myApp := app.New()

	myWindow := myApp.NewWindow("Service Status")
	myLabel := widget.NewLabel(serviceStatus)

	// Load the icon from a file
	icon, _ := fyne.LoadResourceFromPath("./github-logo.png")
	myWindow.SetIcon(icon)

	// Create a container for the GUI elements
	content := container.NewVBox(
		myLabel,
	)

	myWindow.SetContent(content)
	myWindow.SetFixedSize(true)
	myWindow.Hide()
	// myWindow.SetAlwaysOnTop(true)

	myWindow.ShowAndRun()
}
