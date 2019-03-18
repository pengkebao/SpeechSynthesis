package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
)

var (
	unitStr      = []string{"tts_ten.mp3", "tts_hundred.mp3", "tts_thousand.mp3", "tts_ten_thousand.mp3"}
	digitalStr   = []string{"tts_0.mp3", "tts_1.mp3", "tts_2.mp3", "tts_3.mp3", "tts_4.mp3", "tts_5.mp3", "tts_6.mp3", "tts_7.mp3", "tts_8.mp3", "tts_9.mp3"}
	otherStr     = []string{"tts_dot.mp3", "tts_pre.mp3", "tts_default.mp3", "tts_yuan.mp3"}
	unitVoice    [][]byte
	digitalVoice [][]byte
	otherVoice   [][]byte
)

func main() {
	http.HandleFunc("/", OutAudio)
	http.ListenAndServe("127.0.0.1:9000", nil)
}

func OutAudio(w http.ResponseWriter, r *http.Request) {
	invoice := r.FormValue("s")
	if len(invoice) < 1 {
		return
	}
	w.Header().Add("Content-Disposition", "audio/mp3;filename=aa.mp3")
	w.Header().Add("Content-type", "audio/mp3")
	audio := DigitalProcess(invoice)
	audio.WriteTo(w)
}

func init() {
	unitVoice = InitVoice(unitStr)
	digitalVoice = InitVoice(digitalStr)
	otherVoice = InitVoice(otherStr)
}

func InitVoice(srcStr []string) (dstVoice [][]byte) {
	dstVoice = make([][]byte, len(srcStr))
	for k, v := range srcStr {
		f, err := os.Open("mp3/" + v)
		if err != nil {
			fmt.Println(err)
		}
		tmpVoice, err := ioutil.ReadAll(f)
		if err != nil {
			panic(err)
		}
		dstVoice[k] = tmpVoice
		f.Close()
	}
	return
}

/**
处理数字音频
*/
func DigitalProcess(n string) (outVoice *bytes.Buffer) {
	src := strings.Split(n, ".")
	integer, err := strconv.Atoi(src[0])
	if err != nil {
		fmt.Println(err)
	}
	outVoice = bytes.NewBuffer(make([]byte, 0))
	if integer >= 100000000 {
		outVoice.Write(otherVoice[2])
		return outVoice
	}
	decimal := 0
	if len(src) > 1 {
		decimal, err = strconv.Atoi(src[1])
		if err != nil {
			panic(err)
		}
	}
	integerArray := createIntegerArray(integer)
	baseVoice := bytes.NewBuffer(make([]byte, 0))
	baseVoice.Write(otherVoice[1])
	//处理万位
	data0 := integerArray[:4]
	data0Len := len(data0)
	for k, v := range data0 {
		if k != 0 && outVoice.Len() > 0 && data0[k-1] == 0 && v > 0 {
			outVoice.Write(digitalVoice[data0[k-1]])
		}
		if v > 0 {
			if k == 2 {
				if data0[0] != 0 || data0[1] != 0 || v != 1 {
					outVoice.Write(digitalVoice[v])
				}
			} else {
				outVoice.Write(digitalVoice[v])
			}
			if k < data0Len-1 && v > 0 {
				outVoice.Write(unitVoice[data0Len-k-2])
			}
		}
	}
	if outVoice.Len() > 0 {
		outVoice.Write(unitVoice[3]) //万
	}
	//处理千位
	data1 := integerArray[4:]
	data1Len := len(data1)
	for k, v := range data1 {
		if k != 0 && outVoice.Len() > 0 && data1[k-1] == 0 && v > 0 {
			outVoice.Write(digitalVoice[data1[k-1]])
		}
		if v > 0 {
			if k == 2 {
				if data1[0] != 0 || data1[1] != 0 || v != 1 {
					outVoice.Write(digitalVoice[v])
				}
			} else {
				outVoice.Write(digitalVoice[v])
			}
			if k < data1Len-1 && v > 0 {
				outVoice.Write(unitVoice[data1Len-k-2])
			}
		}
	}
	if outVoice.Len() < 1 {
		outVoice.Write(digitalVoice[0])
	}
	/**
	处理小数部份 只报两位
	*/
	if decimal > 0 {
		outVoice.Write(otherVoice[0]) //点
		decimalArray := createDecimalArray(decimal, 2)
		for k, v := range decimalArray {
			outVoice.Write(digitalVoice[v])
			if k == 1 {
				break
			}
		}
	}
	if outVoice.Len() > 0 {
		outVoice.Write(otherVoice[3]) //元
	}
	baseVoice.Write(outVoice.Bytes())
	return baseVoice
}

/**
生成小数部份数组
@param n 小数整型
@param length 要保留小数的位数
*/
func createDecimalArray(n int, length int) []int {
	s := strconv.Itoa(n)
	data := make([]int, 0)
	for k, v := range s {
		tmp, _ := strconv.Atoi(string(v))
		data = append(data, tmp)
		if k == length-1 {
			break
		}
	}
	return data
}

/**
生成整数部份数组
*/
func createIntegerArray(n int) []int {
	data := make([]int, 0)
	data = append(data, n/10000000%10)
	data = append(data, n/1000000%10)
	data = append(data, n/100000%10)
	data = append(data, n/10000%10)
	data = append(data, n/1000%10)
	data = append(data, n/100%10)
	data = append(data, n/10%10)
	data = append(data, n/1%10)
	return data
}
