package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"

	"github.com/jaredtao/Transer/services/baidu"
	"github.com/jaredtao/Transer/services/youdao"

	"github.com/jaredtao/Transer/services/transer"
)

const head = `<?xml version="1.0" encoding="utf-8"?>
<!DOCTYPE TS>
`

type location struct {
	Filename string `xml:"filename,attr"`
	Line     string `xml:"line,attr"`
}
type trans struct {
	Trans string `xml:",chardata"`
	Type  string `xml:"type,attr"`
}
type message struct {
	Locations []location `xml:"location"`
	Source    string     `xml:"source"`
	Trans     trans      `xml:"translation"`
}
type context struct {
	Names    string    `xml:"name"`
	Messages []message `xml:"message"`
}
type ts struct {
	XMLName  xml.Name  `xml:"TS"`
	Version  string    `xml:"version,attr"`
	Language string    `xml:"language,attr"`
	Contexts []context `xml:"context"`
}

//QtTransArgs args
type QtTransArgs struct {
	InputFile  string
	OutputFile string
	API        string
	ID         string
	Secret     string
	TargetLan  string
}

//Trans trans
func Trans(args QtTransArgs) {
	// filePath := "trans_zh.tr"
	data, err := ioutil.ReadFile(args.InputFile)
	if err != nil {
		fmt.Println(err)
		return
	}
	transData := &ts{}
	err = xml.Unmarshal(data, &transData)
	if err != nil {
		fmt.Println(err)
		return
	}
	//TODO  add trans flow
	input := &transer.TransInput{
		ID:     args.ID,
		Secret: args.Secret,
		To:     args.TargetLan,
	}
	alreadyTranNum := 0
	skipTranNum := 0
	baidu.ResetFailedCnt()
	transData.Language = args.TargetLan
	for i := 0; i < len(transData.Contexts); i++ {
		ctx := &transData.Contexts[i]
		for j := 0; j < len(ctx.Messages); j++ {

			// <message>
			//     <location filename="Qml/ContentData.js" line="14"/>
			//     <source>徽章</source>
			//     <translation type="unfinished"></translation>
			// </message>

			msg := &ctx.Messages[j]
			if msg.Trans.Type != "unfinished" {
				skipTranNum++
				fmt.Printf("translate  %d words skip.\r\n", skipTranNum)
				continue
			}
			if msg.Trans.Trans != "" {
				input.Query = msg.Trans.Trans
			} else {
				input.Query = msg.Source
			}
			if args.API == "baidu" {
				ans := baidu.Trans(input)
				if ans.Result != "" {
					msg.Trans.Trans = ans.Result
					msg.Trans.Type = ""
				}
			} else if args.API == "youdao" {
				ans := youdao.Trans(input)
				if ans.Result != "" {
					msg.Trans.Trans = ans.Result
					msg.Trans.Type = ""
				}
			}
			alreadyTranNum++
			fmt.Printf("translate %d words already.\r\n", alreadyTranNum)
		}
	}
	fmt.Printf("baidu translate %d words failed.\r\n", baidu.GetFailedCnt())

	// FIXME 把不对的都转换回来
	if transData.Language == "fra" {
		transData.Language = "fr"
	}
	if transData.Language == "cht" {
		transData.Language = "zh_TW"
	}
	if transData.Language == "ara" {
		transData.Language = "ar"
	}
	if transData.Language == "est" {
		transData.Language = "es"
	}
	if transData.Language == "jp" {
		transData.Language = "ja"
	}
	if transData.Language == "kor" {
		transData.Language = "ko"
	}

	res, err := xml.MarshalIndent(&transData, "", "    ")
	if err != nil {
		fmt.Println(err)
		return
	}
	// header := []byte(`<?xml version="1.0" encoding="utf-8"?>\r\n<!DOCTYPE TS>\r\n`)
	header := []byte(head)
	outFile := append(header, res...)
	err = ioutil.WriteFile(args.OutputFile, outFile, 0666)
	if err != nil {
		fmt.Println(err)
		return
	}
}
