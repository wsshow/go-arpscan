package utils

import (
	"log"

	"golang.org/x/text/encoding/simplifiedchinese"
)

type Buffer struct {
	Data  []byte
	Start int
}

func (b *Buffer) PrependBytes(n int) []byte {
	length := cap(b.Data) + n
	newData := make([]byte, length)
	copy(newData, b.Data)
	b.Start = cap(b.Data)
	b.Data = newData
	return b.Data[b.Start:]
}

func NewBuffer() *Buffer {
	return &Buffer{}
}

func ReverseString(s string) (result string) {
	for _, v := range s {
		result = string(v) + result
	}
	return
}

func ConvertGBK2StrFromStr(gbkStr string) string {
	b, err := simplifiedchinese.GBK.NewDecoder().String(gbkStr)
	if err != nil {
		log.Println(err)
	}
	return b
}

func IsUtf8(data []byte) bool {
	fPreNum := func(data byte) int {
		var mask byte = 0x80
		var num int = 0
		//8bit中首个0bit前有多少个1bits
		for i := 0; i < 8; i++ {
			if (data & mask) == mask {
				num++
				mask = mask >> 1
			} else {
				break
			}
		}
		return num
	}
	i := 0
	for i < len(data) {
		if (data[i] & 0x80) == 0x00 {
			// 0XXX_XXXX
			i++
			continue
		} else if num := fPreNum(data[i]); num > 2 {
			// 110X_XXXX 10XX_XXXX
			// 1110_XXXX 10XX_XXXX 10XX_XXXX
			// 1111_0XXX 10XX_XXXX 10XX_XXXX 10XX_XXXX
			// 1111_10XX 10XX_XXXX 10XX_XXXX 10XX_XXXX 10XX_XXXX
			// 1111_110X 10XX_XXXX 10XX_XXXX 10XX_XXXX 10XX_XXXX 10XX_XXXX
			// preNUm() 返回首个字节的8个bits中首个0bit前面1bit的个数，该数量也是该字符所使用的字节数
			i++
			for j := 0; j < num-1; j++ {
				//判断后面的 num - 1 个字节是不是都是10开头
				if (data[i] & 0xc0) != 0x80 {
					return false
				}
				i++
			}
		} else {
			//其他情况说明不是utf-8
			return false
		}
	}
	return true
}
