package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/Sirupsen/logrus"

	"github.com/urfave/cli"
)

// https://www.migrationsverket.se/ansokanbokning/valjtyp?sprak=en&bokningstyp=2&enhet=Z209&sokande=1

type applicationPurpose int

const (
	fingerAndPhoto applicationPurpose = 2
	alienPassport  applicationPurpose = 4

	migrationsverketHost string = "www.migrationsverket.se"
	basePath             string = "/ansokanbokning/valjtyp?sprak=en"
)

var locationMap = map[string]string{
	"sundbyberg": "Z209",
	"norrköping": "Z083",
}

// Check is the ommand to check if there is empty time slot
func Check() cli.Command {
	flags := []cli.Flag{
		cli.StringFlag{
			Name:  "location, l",
			Usage: "The location of the migration board unit",
			Value: "sundbyberg",
		},
		cli.IntFlag{
			Name:  "number, n",
			Usage: "The number of people",
			Value: 1,
		},
	}
	return cli.Command{
		Name:   "check",
		Usage:  "Check if there is any empty time slot",
		Action: action,
		Flags:  flags,
	}
}

func action(ctx *cli.Context) error {
	err := checkTimeSlot(ctx)

	return err
}

func checkTimeSlot(ctx *cli.Context) error {
	queryPath := composePath(ctx.String("location"), ctx.Int("number"))
	client := &http.Client{}
	reqURL := &url.URL{
		Host:   migrationsverketHost,
		Path:   queryPath,
		Scheme: "https",
	}
	req, err := http.NewRequest("GET", reqURL.String(), bytes.NewBuffer([]byte("")))
	if err != nil {
		logrus.Error(err)
		return err
	}
	logrus.Infof("Querying URL https://%s%s", migrationsverketHost, queryPath)
	resp, err := client.Do(req)
	if err != nil {
		logrus.Error(err)
		return err
	}
	defer resp.Body.Close()
	logrus.Debugf("Return code is %d", resp.StatusCode)
	// Get the response body as a string
	dataInBytes, err := ioutil.ReadAll(resp.Body)
	pageContent := string(dataInBytes)
	//logrus.Debug(pageContent)
	if strings.Contains(pageContent, "No time slots available") {
		logrus.Info("Sorry, no time slots available")
	} else {
		logrus.Infof("There are time slots available, please click following link")
		logrus.Infof("https://www.migrationsverket.se/English/Contact-us/Book-an-appointment-before-you-visit-us.html")
	}
	return nil
}

func composePath(loc string, noOfPeople int) string {
	purposeStr := "&bokningstyp=2"
	finalURL := fmt.Sprintf("%s%s&enhet=%s&sokande=%d", basePath, purposeStr, locationMap[loc], noOfPeople)
	return finalURL
}
