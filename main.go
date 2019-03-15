package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type Audio struct {
	body *bytes.Buffer
}

func NewAudio() *Audio {
	a := new(Audio)
	a.body = bytes.NewBuffer(make([]byte, 0))
	return a
}

func (a *Audio) WriteTo(w io.Writer) (n int64, err error) {
	n, err = a.body.WriteTo(w)
	if err != nil {
		return
	}
	return
}

func main() {
	http.HandleFunc("/", OutAudio)
	http.ListenAndServe("127.0.0.1:9000", nil)
}

var (
	unit       = [...]string{"个", "十", "百", "千", "万"}
	conversion = map[string]string{"default": "tts_default.mp3",
		"十": "tts_ten.mp3", "百": "tts_hundred.mp3", "千": "tts_thousand.mp3", "万": "tts_ten_thousand.mp3",
		"点": "tts_dot.mp3", "元": "tts_yuan.mp3", "0": "tts_0.mp3", "1": "tts_1.mp3", "2": "tts_2.mp3", "3": "tts_3.mp3",
		"4": "tts_4.mp3", "5": "tts_5.mp3", "6": "tts_6.mp3", "7": "tts_7.mp3", "8": "tts_8.mp3", "9": "tts_9.mp3", "base": "tts_pre.mp3"}
	sounds = map[string][]byte{}
)

func OutAudio(w http.ResponseWriter, r *http.Request) {
	invoice := r.FormValue("s")
	if len(invoice) < 1 {
		return
	}
	//invoice := "创客宝多码付到账18元"
	newvoice := ""
	for _, v := range invoice {
		fmt.Println(string(v), v)
		if v >= 46 && v <= 59 {
			newvoice += string(v)
		}
	}
	if len(newvoice) < 1 {
		w.Write([]byte("数据错误"))
		return
	}
	audio := synthesis(newvoice)
	w.Header().Add("Content-Disposition", "audio/mp3;filename=aa.mp3")
	w.Header().Add("Content-type", "audio/mp3")
	n, err := audio.WriteTo(w)
	fmt.Println(n, err)
}

func init() {
	for k, v := range conversion {
		var err error
		f, err := os.Open("mp3/" + v)
		if err != nil {
			fmt.Println(err)
		}
		sounds[k], err = ioutil.ReadAll(f)
		if err != nil {
			panic(err)
		}
		f.Close()
	}
}

func synthesis(invoice string) *Audio {
	voice := strings.Split(invoice, ".")
	voice1 := voice[0]
	voice2 := ""
	if len(voice) > 1 {
		voice2 = voice[1]
	}
	voice1len := len(voice1)
	outvoice := ""
	audio := NewAudio()
	if voice1len > 0 && voice1len <= 8 {
		for _, v := range voice1 {
			voice1len--
			if voice1len > 4 {
				outvoice += string(v)
				if string(v) != "0" {
					outvoice += unit[voice1len-4]
				}
			} else {
				if voice1len == 4 && len(outvoice) > 0 && string(v) == "0" {
					if len(outvoice) > 1 {
						outvoice = strings.TrimRight(outvoice, "0")
					}
					outvoice += unit[voice1len]
				}
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
		if strings.HasPrefix(outvoice, "1十") {
			outvoice = strings.TrimPrefix(outvoice, "1")
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

		audio.body.Write(sounds["base"])
		for _, v := range outvoice {
			audio.body.Write(sounds[string(v)])
		}
		return audio
	}
	audio.body.Write(sounds["default"])
	return audio
}

func dedupZero(str string) string {
	if !strings.Contains(str, "00") {
		return str
	}
	str = strings.Replace(str, "00", "0", -1)
	return dedupZero(str)
}
