package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"

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

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	// Define an API endpoint to identify that the service is running
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, serviceStatus)
	})

	r.POST("/autoprint-pdf", func(c *gin.Context) {
		var body RequestBody

		// Get the printer_location from the form data
		if printerLocation := c.PostForm("printer_location"); printerLocation != "" {
			body.PrinterLocation = printerLocation
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "printer_location is required"})
			return
		}

		// Get the pdf file from the request body
		if pdf, err := io.ReadAll(c.Request.Body); err == nil {
			body.PDF = pdf
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "failed to read pdf file"})
			return
		}

		// Process the valid input
		fmt.Println("Printer Location:", body.PrinterLocation)
		fmt.Println("PDF:", body.PDF)

		// Check if the PDF is a valid file
		if len(body.PDF) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "PDF file is required"})
			return
		}

		// Replace "path/to/your/exefile.exe" with the actual path to your executable.
		exePath := "./SumatraPDF-3.3.3-64.exe"

		// Create a new command to run the executable.
		cmd := exec.Command(exePath)

		// Set the standard output and error to os.Stdout and os.Stderr, respectively.
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		// Run the command.
		if err := cmd.Run(); err != nil {
			// fmt.Printf("Error running the executable: %v\n", err)
			c.String(http.StatusNotFound, "AutoPrint PDF activation failed")
			return
		}

		// fmt.Println("Executable ran successfully.")
		c.String(http.StatusOK, "AutoPrint PDF activated")
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
