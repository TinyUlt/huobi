package main

import (
	"github.com/go-gomail/gomail"
	"time"
	"services"
	"fmt"
	"models"
	//"math"
	"io/ioutil"
	"encoding/json"

	"os"
	//"container/list"
	//"container/list"
	//"reflect"
	//"github.com/ethereum/go-ethereum/swarm/network"
	"io"
)
type myMape map[float32]int
func (m myMape) String() string{

	copy := make(myMape)
	copy.Copy(m)

	str:=""
	for{
		lowest:=float32(100.0)
		for i := range copy{
			if i < lowest{
				lowest = i
			}
		}

		str += fmt.Sprintf("(%.2f¥%d家)", lowest, copy[lowest])

		delete(copy,lowest)

		if len(copy)==0{
			break
		}
	}
	return str
}
func (m myMape) Copy(src myMape){
	for k := range m {
		delete(m, k)
	}
	for k , v:= range src {
		m[k] = v
	}
}
var strarry =[] string{"d","e","f","g","h"}
var index = int(0)
var mailtoPassword ="Ww18767104183"
func getSendPerson() string{

	if index >= len(strarry){
		index = 0
	}
	emailfrom := fmt.Sprintf("%s%s",strarry[index],"@ultbtc.com")
	index++
	return emailfrom
}
// fileName:文件名字(带全路径)
// content: 写入的内容
func appendToFile(fileName string, content string) error {
	// 以只写的模式，打开文件
	f, err := os.OpenFile(fileName, os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("cacheFileList.yml file create failed. err: " + err.Error())
	} else {
		// 查找文件末尾的偏移量
		n, _ := f.Seek(0, io.SeekEnd)
		// 从末尾的偏移量开始写入内容
		_, err = f.WriteAt([]byte(content), n)
	}
	defer f.Close()
	return err
}


func sendMail(subject string, emails []string, data string){

	timeStr:=time.Now().Format("2006-01-02 15:04:05")


	for _,email := range emails{

		mailfrom := getSendPerson()
		m := gomail.NewMessage()
		m.SetHeader("From", mailfrom)
		m.SetHeader("To", email, email)
		m.SetAddressHeader("Cc", email, "Dan")
		m.SetHeader("Subject", subject)
		m.SetBody("text/html", data)
		//m.Attach("/home/Alex/lolcat.jpg")

		d := gomail.NewDialer("smtp.exmail.qq.com", 465, mailfrom, mailtoPassword)

		// Send the email to Bob, Cora and Dan.

		go func(m *gomail.Message){
			if err := d.DialAndSend(m); err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println(fmt.Sprintf("%s %s to %s:%s",timeStr,mailfrom, email, subject) )
		}(m)

	}

}
func getTradeName(tradeType int) string{
	if tradeType ==1{
		return "Sell"
	}else{
		return "Buy"
	}
}

type TradeFollow struct{
	emails [][]string
	stepEmails []string
	stepPrice float32
//	buyEmails []string
//	sellPrice float32
//	buyPrice float32
	step float64
	oldMaps map[int]myMape
}
func(m *TradeFollow)Update( data map[int][]models.USDTData){

	maps := make(map[int]myMape)
	for i,v := range data{
		maps[i]= make(myMape)
		for _, d := range v{
			_, ok := maps[i][d.Price]
			if !ok {
				maps[i][d.Price] = 0

			}
			maps[i][d.Price]++

		}
		isDiffrent:=false
		for j, f := range maps[i]{
			if m.oldMaps[i][j]!=f{
				isDiffrent = true
				m.oldMaps[i].Copy(maps[i])
				break
			}
		}
		if isDiffrent{
			mapsStr:=fmt.Sprint(maps[i])
			subjectStr:=fmt.Sprintf("%s%.2f btc%s:%s usdt%s",getTradeName(i),data[i][0].Price, btcPriceStr,mapsStr,usdtPriceStr )

			//sell
			if i == 1{

				if m.stepPrice !=data[i][0].Price{
					sendMail(subjectStr, m.stepEmails, fmt.Sprintln(data[i][0]))
					m.stepPrice = data[i][0].Price
				}
				timeStr:=time.Now().Format("2006-01-02 15:04:05")
				appendToFile(writeDataName, timeStr+" "+subjectStr + "\n")
			}

			sendMail(subjectStr, m.emails[i], fmt.Sprintln(data[i][0]))

		}

	}
	//m.oldMaps = maps
	//m.oldMaps = make(map[int]map[float32]int)
	//m.oldMaps.
	//copy( m.oldMaps[0] , maps[0])



	//if data.TradeType == "Sell"{

	//hightst := data[0][0].Price
	//
	//if math.Abs(float64(hightst - m.sellPrice)) >= m.step{
	//	m.sellPrice = hightst
	//
	//
	//	sendMail(fmt.Sprintf("%s:%f %s","Sell", hightst,data[][0].UserName), m.sellEmails, fmt.Sprintln(*data))
	//}
	//
	////}else {
	//hightst =data[0][0].Price
	//if math.Abs(float64(hightst - m.buyPrice)) >= m.step{
	//	m.buyPrice = hightst
	//
	//
	//	sendMail(fmt.Sprintf("%s:%f %s","Buy", hightst, data.Data[0].UserName), m.buyEmails, fmt.Sprintln(*data))
	//}

	//}

}
type PersonFollow struct{
	Name string
	sellEmails []string
	buyEmails []string

	sellPrice float32
	buyPrice float32
}

func(p *PersonFollow)Find(data *models.USDTData, tradeType int){

	if tradeType == 1{
		if p.sellPrice != data.Price{
			p.sellPrice = data.Price

			sendMail(fmt.Sprintf("%s %s:%f",p.Name, getTradeName(tradeType), data.Price), p.sellEmails, fmt.Sprintln(*data))
		}
	}else{
		if p.buyPrice != data.Price{
			p.buyPrice = data.Price

			sendMail(fmt.Sprintf("%s %s:%f",p.Name, getTradeName(tradeType), data.Price), p.buyEmails, fmt.Sprintln(*data))
		}
	}


}
func(p *PersonFollow)UnFind(tradeType int){


	oldPrice:=float32(0.0)
	if tradeType == 1{
		if p.sellPrice != -1{
			oldPrice = p.sellPrice
			p.sellPrice = -1
			sendMail(fmt.Sprintf("%s %s:%s",p.Name,getTradeName(tradeType),"Cancel"), p.sellEmails, fmt.Sprintf("%f Cancel", oldPrice))
		}
	}else{
		if p.buyPrice != -1{
			oldPrice = p.buyPrice
			p.buyPrice = -1
			sendMail(fmt.Sprintf("%s %s:%s",p.Name,getTradeName(tradeType),"Cancel"), p.buyEmails, fmt.Sprintf("%f Cancel", oldPrice))
		}
	}


}
func findName(datas []models.USDTData, name string) *models.USDTData{
	for _, data := range datas {
		if data.UserName==name{
			return &data
		}
	}
	return nil
}
type AveFollow struct{
	sellEmails []string
	buyEmails []string
	pageCount int
	aveSellPriceList []float32
	aveBuyPriceList []float32
}
func(a *AveFollow)SetPrice(tradeType int, avePrice, histPrice float32){

	if tradeType == 1{

		if len(a.sellEmails) != 0 && avePrice != a.aveSellPriceList[0] {


			var temp  = make([]float32, 30)
			temp[0] = avePrice
			copy(temp[1:], a.aveSellPriceList[:len(a.aveSellPriceList)-1])
			a.aveSellPriceList = temp[:]
			sendMail(fmt.Sprintf("%s:a %f,l %f, h %f",getTradeName(tradeType),avePrice, a.aveSellPriceList[1], histPrice), a.sellEmails, fmt.Sprintln(a.aveSellPriceList))

		}
	}else{

		if len(a.buyEmails) != 0 && avePrice != a.aveBuyPriceList[0] {


			var temp  = make([]float32, 30)
			temp[0] = avePrice
			copy(temp[1:], a.aveBuyPriceList[:len(a.aveBuyPriceList)-1])
			a.aveBuyPriceList = temp[:]
			sendMail(fmt.Sprintf("%s:a %f,l %f, h %f",getTradeName(tradeType),avePrice, a.aveBuyPriceList[1], histPrice), a.buyEmails, fmt.Sprintln(a.aveBuyPriceList))

		}
	}

}
var btcPriceStr string
var usdtPriceStr string
var writeDataName string
func main(){



	//fmt.Print( services.GetTicker("btcusdt"))
	//return
	//temp := make( map[int]int)
	//temp[0] = 1
	//temp[1] = 1
	//temp[2] = 1
	//temp[3] = 1
	//temp[4] = 1
	//temp[5] = 1
	//
	//for k := range temp{
	//	delete(temp,k)
	//	fmt.Println(temp)
	//}
	//fmt.Println(temp)
	//return
	//var temp2 []float32
	//temp2[0]=-1
	//temp2[1]=-1
	////temp2 = append(temp2, temp)
	//copy(temp2[2:], temp[:])
	//fmt.Print(temp2)
	//return

	fileName:="config.json"
	writeDataName ="usdtData"
	if len(os.Args)>2{

		fileName = os.Args[1]
		writeDataName = os.Args[2]
	}

	bytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Println("ReadFile: ", err.Error())
		return

	}
	configData := models.ConfigData{}
	if err := json.Unmarshal(bytes, &configData); err != nil {
		fmt.Println("Unmarshal: ", err.Error())
		return
	}

	CoinId := configData.CoinId
	TradeFollowSellEmial :=configData.TradeFollowSellEmial
	TradeFollowBuyEmial := configData.TradeFollowBuyEmial
	TradeFollowStepEmial := configData.TradeFollowStepEmail
	//PersonFollowName := configData.PersonFollowName
	//PersonFollowSellEmail := configData.PersonFollowSellEmail
	//PersonFollowBuyEmail := configData.PersonFollowBuyEmail
	Step := configData.Step
	AveFollowSellEmail:=configData.AveFollowSellEmail
	AveFollowBuyEmail:=configData.AveFollowBuyEmail
	PageCount := configData.PageCount

	tradeFollow := TradeFollow{

		emails:make([][]string,2),// TradeFollowSellEmial,
		oldMaps:make(map[int]myMape),
		//buyEmails:TradeFollowBuyEmial,
		//sellPrice:-1,
		//buyPrice:-1,
		step:Step,
	}
	tradeFollow.oldMaps[0] = make(myMape)
	tradeFollow.oldMaps[1] = make(myMape)
	tradeFollow.emails[1] = TradeFollowSellEmial
	tradeFollow.emails[0] = TradeFollowBuyEmial
	tradeFollow.stepEmails = TradeFollowStepEmial
	//personFollow :=PersonFollow{
	//	Name:PersonFollowName,
	//	sellEmails:PersonFollowSellEmail,
	//	buyEmails:PersonFollowBuyEmail,
	//	sellPrice:-1,
	//	buyPrice:-1,
	//}

	aveFollow:=AveFollow{
		sellEmails:AveFollowSellEmail,
		buyEmails:AveFollowBuyEmail,
		pageCount:PageCount,
		aveSellPriceList:make([]float32, 30),
		aveBuyPriceList:make([]float32, 30),
	}

	//btcOldPrice:=float64(0)

	//var BuyData []models.USDTData;
	for {
		_, _,  usdtPriceStr= services.GetMarketPrice()

		btcPrice:=services.GetTicker("btcusdt").Tick.Close

		//attachStr:=""
		//if btcPrice > btcOldPrice{
		//	attachStr = "↑"
		//}else if btcPrice < btcOldPrice{
		//	attachStr = "↓"
		//
		//}
		//btcOldPrice = btcPrice
		btcPriceStr =fmt.Sprintf("%.2f$",btcPrice)
		traderData :=make(map[int][]models.USDTData)
		traderData[0] = make([]models.USDTData,0)
		traderData[1] = make([]models.USDTData,0)
		for tradeType:=0; tradeType<2; tradeType++{

			//allPrice:=int32(0)
			//allCount:=int(0)
			//highstPrice:=float32(0)
			for page:=1;page <= aveFollow.pageCount;page++{

				data:=services.GetUstdData(CoinId,tradeType,page)

				if len(data.Data)==0{
					continue
				}

				tempdata := make([]models.USDTData, len(data.Data)+len(traderData[tradeType]))
				copy(tempdata[:], traderData[tradeType])
				copy(tempdata[len(traderData[tradeType]):], data.Data)
				traderData[tradeType] = tempdata

				//if page == 0{
				//	tradeFollow.Update(&data)
				//	highstPrice = data.Data[0].Price
				//}

				//for _, d := range data.Data {
				//	allPrice += int32(d.Price*100)
				//	allCount++
				//}
				//if page == aveFollow.pageCount-1{
				//	aveFollow.SetPrice(tradeType,float32(allPrice)/float32(allCount)/100,highstPrice)
				//}
				//
				//personData :=findName(data.Data, personFollow.Name)
				//if personData!=nil{
				//	personFollow.Find(personData ,tradeType)
				//	//break
				//}else if page == aveFollow.pageCount-1{
				//	personFollow.UnFind(tradeType)
				//}
			}
		}
		tradeFollow.Update(traderData)
		time.Sleep(1 * time.Second)
	}
}