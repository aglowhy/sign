package selenium

import (
	"fmt"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/firefox"
	"os"
	"strings"
	"time"
)

// Start a Selenium WebDriver server instance (if one is not already
// running).
const (
	// These paths will be different on your system.
	seleniumPath    = "/apps/selenium-server-standalone-3.141.59.jar"
	geckoDriverPath = "/apps/geckodriver"
	//geckoDriverPath = "/apps/geckodriver.exe"
	port = 8080
)

func Test() {
	opts := []selenium.ServiceOption{
		//selenium.StartFrameBuffer(),           // Start an X frame buffer for the browser to run in.
		selenium.GeckoDriver(geckoDriverPath), // Specify the path to GeckoDriver in order to use Firefox.
		selenium.Output(os.Stderr),            // Output debug information to STDERR.
	}
	selenium.SetDebug(true)
	service, err := selenium.NewSeleniumService(seleniumPath, port, opts...)
	if err != nil {
		panic(err) // panic is used only as an example and is not otherwise recommended.
	}
	defer service.Stop()

	// Connect to the WebDriver instance running locally.
	capabilities := firefox.Capabilities{
		Binary: "/apps/firefox/firefox",
		Args: []string{
			"--headless",
		},
		//Prefs: map[string]interface{}{
		//	"general.useragent.override": "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:76.0) Gecko/20100101 Firefox/76.0",
		//},
	}

	caps := selenium.Capabilities{
		"browserName": "firefox",
		//"firefoxOptions": map[string]interface{}{
		//	"mobileEmulation": map[string]interface{}{
		//		"deviceName": "iPhone X",
		//	},
		//},
	}
	caps.AddFirefox(capabilities)
	wd, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", port))
	if err != nil {
		panic(err)
	}
	defer wd.Quit()

	//fmt.Println(capabilities.Profile)
	//os.Exit(1)

	// Navigate to the simple playground interface.
	if err := wd.Get("http://bbs.guilinlife.com/member.php?mod=logging&action=login&mobile=2"); err != nil {
		panic(err)
	}

	cookies, _ := wd.GetCookies()
	fmt.Println(cookies)

	// Get a reference to the text box containing code.
	elem, err := wd.FindElement(selenium.ByID, "loginform")
	if err != nil {
		panic(err)
	}

	// Remove the boilerplate code already in the text box.
	if err := elem.Clear(); err != nil {
		panic(err)
	}

	// Enter some new code in text box.
	err = elem.SendKeys(`
		package main
		import "fmt"

		func main() {
			fmt.Println("Hello WebDriver!\n")
		}
	`)
	if err != nil {
		panic(err)
	}

	// Click the run button.
	btn, err := wd.FindElement(selenium.ByCSSSelector, "#run")
	if err != nil {
		panic(err)
	}
	if err := btn.Click(); err != nil {
		panic(err)
	}

	// Wait for the program to finish running and get the output.
	outputDiv, err := wd.FindElement(selenium.ByCSSSelector, "#output")
	if err != nil {
		panic(err)
	}

	var output string
	for {
		output, err = outputDiv.Text()
		if err != nil {
			panic(err)
		}
		if output != "Waiting for remote server..." {
			break
		}
		time.Sleep(time.Millisecond * 100)
	}

	fmt.Printf("%s", strings.Replace(output, "\n\n", "\n", -1))

	// Example Output:
	// Hello WebDriver!
	//
	// Program exited.
}
