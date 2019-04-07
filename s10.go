//
// s10.go
// Copyright (C) 2019 Yongwen Zhuang <zeoman@163.com>
//
// Distributed under terms of the MIT license.
//

package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/tealeg/xlsx"
	"golang.org/x/text/encoding/traditionalchinese"
	"golang.org/x/text/transform"
	"io/ioutil"
	"math/rand"
	"os"
	"strconv"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

//convert UTF-8 to BIG5
func EncodeBig5(s []byte) []byte {
	I := bytes.NewReader(s)
	O := transform.NewReader(I, traditionalchinese.Big5.NewEncoder())
	d, _ := ioutil.ReadAll(O)
	return d
}

func StringToByte(str string, min int, max int) []byte {
	var res int
	if str == "" {
		res = (rand.Intn(max-min) + min)
	} else {
		ret, _ := strconv.ParseInt(str, 10, 16)
		res = int(ret)
	}
	if res > max || res < min {
		res = (rand.Intn(max-min) + min)
	}

	buf := make([]byte, 2)
	binary.LittleEndian.PutUint16(buf, uint16(res))
	return buf
}

func main() {

	buff := bytes.NewBuffer([]byte{})
	cnt := 0
	excelFileName := "s10.xlsx"
	xlFile, err := xlsx.OpenFile(excelFileName)
	check(err)
	if _, err := os.Stat("TPERSON.S10"); err == nil {
		f, err := os.Open("TPERSON.S10")
		check(err)
		defer f.Close()
		f.Seek(4, 0)
		buf := make([]byte, 260)
		for i := 0; i < 110; i++ {
			_, err := f.Read(buf)
			check(err)
			if buf[0] == 0 && buf[1] == 0 {
				break
			} else {
				cnt++
				buff.Write(buf)
			}
		}
	}
	// For more granular writes, open a file for writing.
	// 775 884
	f, err := os.Create("TPERSON.S10")
	check(err)
	defer f.Close()

	// You can `Write` byte slices as you'd expect.
	head := []byte{0x9e, 0x00, 0x00, 0x00}
	f.Write(head)
	sheet := xlFile.Sheets[0]
	for num, row := range sheet.Rows {
		if num > 1 && row.Cells[0].String() != "" {
			cnt++
			// 头像
			buff.Write(StringToByte(row.Cells[5].String(), 1, 600))
			// 父亲母亲配偶家族兄弟仇敌
			for i := 0; i < 6; i++ {
				buff.Write([]byte{255, 255})
			}
			// 诞生 登场 死亡年份
			for i := 7; i <= 9; i++ {
				buff.Write(StringToByte(row.Cells[i].String(), 120, 300))
			}

			str := row.Cells[6].String()
			var start int
			if str != "" {
				start, _ = strconv.Atoi(str)
			} else {
				start = 70
			}

			// 声望 统帅 武力 智力 政治 魅力
			buff.Write(StringToByte(row.Cells[10].String(), 1, 500))
			for i := 11; i <= 15; i++ {
				buff.Write(StringToByte(row.Cells[i].String(), start, 99))
			}
			buff.Write([]byte{0, 4})
			buff.Write(StringToByte(row.Cells[3].String(), 0, 1))
			for i := 16; i <= 26; i++ {
				res, _ := strconv.ParseInt(row.Cells[i].String(), 10, 8)
				buff.Write([]byte{byte(res)})
			}
			// teji_num := rand.Intn(6) + 3
			// 1,2,4,8,16,32,64,128
			msks := []byte{1, 2, 4, 8, 32, 64, 128}
			for i := 0; i < 5; i++ {
				var tmp byte
				num := rand.Intn(6)
				for j := 0; j < num; j++ {
					tmp = tmp | msks[rand.Intn(6)]
				}
				buff.Write([]byte{tmp})
			}
			tmp := rand.Intn(4)
			if tmp == 0 {
				buff.Write([]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 255, 255, 0, 0, 0, 0})
			} else if tmp == 1 {
				buff.Write([]byte{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 255, 255, 0, 0, 0, 0})
			} else if tmp == 2 {
				buff.Write([]byte{2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 255, 255, 0, 0, 0, 0})
			} else if tmp == 3 {
				buff.Write([]byte{3, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 255, 255, 0, 0, 0, 0})
			}
			// 信念、类型
			for i := 27; i <= 28; i++ {
				res, _ := strconv.ParseInt(row.Cells[i].String(), 10, 8)
				buff.Write([]byte{byte(res)})
			}
			buff.Write([]byte{0, 0, 0, 0, 0, 0, 0, 0, 0})
			xing := []byte(row.Cells[0].String())
			ming := []byte(row.Cells[1].String())
			zi := []byte(row.Cells[2].String())
			liezhuan := []byte(row.Cells[6].String())
			buff.Write(EncodeBig5(xing)[:4])
			buff.Write([]byte{0})
			buff.Write(EncodeBig5(ming)[:4])
			buff.Write([]byte{0})
			buff.Write(EncodeBig5(zi)[:4])
			buff.Write([]byte{0})
			// 读音
			for i := 0; i < 10; i++ {
				buff.Write([]byte{0})
			}
			buff.Write([]byte{0})
			if row.Cells[6].String() == "" {
				liezhuan = []byte(fmt.Sprintf("第%d號穿越者", cnt))
			}
			buff.Write(EncodeBig5(liezhuan)[:150])
			buff.Write([]byte{0})
			buff.Write([]byte{0})
			buff.Write([]byte{0})

		}
	}
	buff.WriteTo(f)
	for ; cnt < 110; cnt++ {
		for i := 0; i < 260; i++ {
			f.Write([]byte{0})
		}
	}

	fmt.Print(`欢迎使用三国志X自定义武将生成器（破产版）
此工具由YEO开发（耗时1天）
自定义武将支持excel导入
请修改s10.xlsx文件
`)

	// Issue a `Sync` to flush writes to stable storage.
	f.Sync()
	fmt.Print("导入成功！按回车结束！")
	reader := bufio.NewReader(os.Stdin)
	reader.ReadString('\n')
}
