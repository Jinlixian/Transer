package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/jaredtao/Transer/services/baidu"
	"github.com/jaredtao/Transer/services/transer"
	"github.com/jaredtao/Transer/services/youdao"
)

const baiduID = "20190502000293463"
const baiduSecret = "0d2RvCho9XZNEO5GCGNs"

const youdaoID = "1bd659586c52ea1d"
const youdaoSecret = "5ZktXhHfLCpI0KnAdcxx4cPyGJwcVXaV"

var api = flag.String("api", "baidu", "baidu | youdao")
var userID = flag.String("userID", "20190502000293463", "your id")
var secret = flag.String("secret", "0d2RvCho9XZNEO5GCGNs", "your secret")
var originText = flag.String("text", "Hello World", "the text need translate")
var targetLan = flag.String("targetLang", "zh", "zh | en | ja | ko | fr | es | pt | it | ru | vi | de | ar | id")

var inputFile = flag.String("inputFile", "", "the input file need translate")
var outputFile = flag.String("outputFile", "", "need the output file")

// type QtTransArgs struct {
// 	InputFile  string
// 	OutputFile string
// 	API        string
// 	ID         string
// 	Secret     string
// 	TargetLan  string
// }

type FileTransArgs struct {
	InputFile  string
	OutputFile string
	API        string
	ID         string
	Secret     string
	TargetLan  string
}

func main() {
	flag.Parse()

	args := &FileTransArgs{}
	args.API = *api
	args.ID = *userID
	args.Secret = *secret
	args.InputFile = *inputFile
	args.OutputFile = *outputFile
	args.TargetLan = *targetLan

	if *inputFile == "" || *outputFile == "" {
		flag.PrintDefaults()
		return
	}
	data, err := ioutil.ReadFile(args.InputFile)
	// fmt.Println(string(data))
	inputList := strings.Split(string(data), "\n")

	// fmt.Println(inputList)
	fmt.Println("File:", args.InputFile, "lines:", len(inputList))
	if err != nil {
		fmt.Println(err)
		return
	}
	if string(data) == "" {
		fmt.Println("文档是空的")
		return
	}

	outString := ""

	for i := 0; i < len(inputList); i++ {
		input := &transer.TransInput{
			ID:     *userID,
			Secret: *secret,
			Query:  *originText,
			To:     *targetLan,
		}
		input.Query = string(inputList[i])
		outText := ""
		if input.Query == "" {
			outText = ""
		} else {
			var result *transer.TransOutput
			if *api == "baidu" {
				ok, to := baidu.LanConvertFromYouDao(*targetLan)
				if ok {
					input.To = to
				}
				result = baidu.Trans(input)
				if result.Result == "" {
					fmt.Println("retrying")
					ok, to := baidu.LanConvertFromYouDao(*targetLan)
					if ok {
						input.To = to
					}
					result = baidu.Trans(input)
				}
			} else if *api == "youdao" {
				result = youdao.Trans(input)
			} else {
				flag.PrintDefaults()
			}
			outText = result.Result
			fmt.Println(outText)
		}
		outString = outString + "\r\n" + outText
	}

	err = ioutil.WriteFile(args.OutputFile, []byte(outString), 0666)
	if err != nil {
		fmt.Println(err)
		return
	}
}
