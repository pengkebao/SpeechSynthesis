package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	//本地合成
	invoice := "创客宝多码付到账100.00元"
	newvoice := ""
	for _, v := range invoice {
		fmt.Println(string(v), v)
		if v >= 46 && v <= 59 {
			newvoice += string(v)
		}
	}
	fmt.Println("newvoice", newvoice)
	filename := synthesis(newvoice)
	fmt.Println(filename)
}

var unit [5]string = [...]string{"个", "十", "百", "千", "万"}
var conversion map[string]string = map[string]string{
	"十": "tts_ten.mp3", "百": "tts_hundred.mp3", "千": "tts_thousand.mp3", "万": "tts_ten_thousand.mp3",
	"点": "tts_dot.mp3", "元": "tts_yuan.mp3", "0": "tts_0.mp3", "1": "tts_1.mp3", "2": "tts_2.mp3", "3": "tts_3.mp3",
	"4": "tts_4.mp3", "5": "tts_5.mp3", "6": "tts_6.mp3", "7": "tts_7.mp3", "8": "tts_8.mp3", "9": "tts_9.mp3", "base": "tts_pre.mp3"}

func synthesis(invoice string) string {
	voice := strings.Split(invoice, ".")
	voice1 := voice[0]
	voice2 := ""
	if len(voice) > 1 {
		voice2 = voice[1]
	}

	voice1len := len(voice1)
	outvoice := ""
	if voice1len > 0 && voice1len <= 8 {
		for _, v := range voice1 {
			voice1len--
			if voice1len > 4 {
				outvoice += string(v)
				if string(v) != "0" {
					outvoice += unit[voice1len-4]
				}
			} else {
				outvoice += string(v)
				if voice1len > 0 && string(v) != "0" {
					outvoice += unit[voice1len]
				}
			}
		}
		//去掉重复的0
		outvoice = dedupZero(outvoice)
		//去掉末尾的0
		if len(outvoice) > 1 {
			outvoice = strings.TrimRight(outvoice, "0")
		}
		voice2len := len(voice2)
		if voice2len > 0 && voice2 != "00" {
			outvoice += "点"
			for _, v := range voice2 {
				voice2len--
				outvoice += string(v)
			}
		}
		outvoice += "元"
		fmt.Println(outvoice)

		tmpfile, _ := os.Create("tmp.mp3")
		defer tmpfile.Close()
		x, err := os.Open("mp3/" + conversion["base"])
		if err != nil {
			fmt.Println(err)
		}
		io.Copy(tmpfile, x)
		x.Close()
		for _, v := range outvoice {
			f, err := os.Open("mp3/" + conversion[string(v)])
			if err != nil {
				fmt.Println(err)
			}
			io.Copy(tmpfile, f)
			f.Close()
		}
		tmpfile.Close()
		return "tmp.mp3"
	}
	return "tts_default.mp3"
}

func dedupZero(str string) string {
	if !strings.Contains(str, "00") {
		return str
	}
	str = strings.Replace(str, "00", "0", -1)
	return dedupZero(str)
}
