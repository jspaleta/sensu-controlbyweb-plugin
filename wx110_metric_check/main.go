package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"github.com/spf13/cobra"
	"golang.org/x/net/html/charset"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

var (
	metricName string
	stateURL   string
	sensor1    bool
	sensor2    bool
	sensor3    bool
	allSensors bool
	verbose    bool
)

type stateXML struct {
	XMLName xml.Name `xml:"datavalues"`
	Units   string   `xml:"units"`
	Sensor1 string   `xml:"sensor1"`
	Sensor2 string   `xml:"sensor2"`
	Sensor3 string   `xml:"sensor3"`
	Time    string   `xml:"time"`
}

func main() {

	//// Open our xmlFile
	//xmlFile, err := os.Open("state.xml")
	//// if we os.Open returns an error then handle it
	//if err != nil {
	//	fmt.Println(err)
	//}

	//fmt.Println("Successfully Opened state.xml")
	//// defer the closing of our xmlFile so that we can parse it later on
	//defer xmlFile.Close()

	//xml_data, _ := ioutil.ReadAll(xmlFile)
	//var state stateXML

	//reader := bytes.NewReader(xml_data)
	//decoder := xml.NewDecoder(reader)
	//decoder.CharsetReader = charset.NewReaderLabel
	//err = decoder.Decode(&state)
	//if err != nil {
	//	fmt.Println("decoder error:", err)
	//	os.Exit(1)
	//}
	//fmt.Printf("state: %v\n", state)
	//os.Exit(1)

	rootCmd := configureRootCommand()
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	} else {
	}
}

func configureRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "wx110_metric_check",
		Short: "Retrieve WX-110 as graphite metric",
		Long:  `Retrieve Particle string variable and output in graphite plaintext format`,
		RunE:  run,
	}

	cmd.Flags().StringVarP(&stateURL,
		"url",
		"u",
		"",
		"WX-110 state url: (ex: http://x.x.x.x:80/state.xml)")

	_ = cmd.MarkFlagRequired("url")

	cmd.Flags().StringVarP(&metricName,
		"metric",
		"m",
		"",
		"Optional Metric Name, if not set will be determined from hostname.variable")

	cmd.Flags().BoolVar(&verbose,
		"verbose",
		false,
		"Enable verbose output")

	cmd.Flags().BoolVar(&sensor1,
		"sensor1",
		false,
		"Enable sensor1 output")

	cmd.Flags().BoolVar(&sensor2,
		"sensor2",
		false,
		"Enable sensor2 output")

	cmd.Flags().BoolVar(&sensor3,
		"sensor3",
		false,
		"Enable sensor3 output")

	cmd.Flags().BoolVar(&allSensors,
		"all",
		false,
		"Enable all sensor output")

	return cmd
}

func run(cmd *cobra.Command, args []string) error {
	if len(args) != 0 {
		_ = cmd.Help()
		return fmt.Errorf("invalid argument(s) received")
	}

	var output stateXML
	err := xw110DeviceState(&output)
	if err != nil {
		return err
	}
	if metricName == "" {
		metricName, err = os.Hostname()
	}
	if sensor1 || allSensors {
		if s, e := strconv.ParseFloat(output.Sensor1, 32); e == nil {
			fmt.Printf("%s %.2f %s\n", metricName+".sensor1", s, output.Time)
		}
	}
	if sensor2 || allSensors {
		if s, e := strconv.ParseFloat(output.Sensor2, 32); e == nil {
			fmt.Printf("%s %.2f %s\n", metricName+".sensor2", s, output.Time)
		}
	}
	if sensor3 || allSensors {
		if s, e := strconv.ParseFloat(output.Sensor3, 32); e == nil {
			fmt.Printf("%s %.2f %s\n", metricName+".sensor3", s, output.Time)

		}
	}
	return err
}

func xw110DeviceState(output *stateXML) error {
	var (
		body []byte
		err  error
	)
	body, err = makeRequest(stateURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(2)
	}
	reader := bytes.NewReader(body)
	decoder := xml.NewDecoder(reader)
	decoder.CharsetReader = charset.NewReaderLabel
	err = decoder.Decode(&output)

	if verbose {
		fmt.Printf("Response: %s\n", body)
		fmt.Printf("Var:%v\n", output)
	}
	return err
}

func makeRequest(urlStr string) ([]byte, error) {
	resp, err := http.Get(urlStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	if resp.StatusCode != 200 {
		err = fmt.Errorf("Failed Request %s StatusCode: %v", urlStr, resp.StatusCode)
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	return body, err
}
