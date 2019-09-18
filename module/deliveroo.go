package module

import (
	"crawler-project/common"
	"crawler-project/model"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

//商家json解析结构体
type jsonAnalyze struct {
	Props struct{
		InitialState struct{
			Home struct{
				Feed struct{
					Results struct{
						Data []struct{
							Header string `default:""`
							Blocks []struct{
								UiContent struct{
									Default struct{
										UiLines []struct{
											Text string `default:""`
											UiSpans []UiSpans
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}
}

//菜品类型分析结构体
type UiSpans struct {
	Text string
}

//json配置解析结构体
type configS struct {
	Area []struct{
		Name string `json:"name"`
		Area string `json:"Area"`
		Api string `json:"api"`
	} `json:"area"`
}

//分析开始
//获取需要爬虫的区域
//解析数据
func DeliverooStart()  {
	cs := jsonDecode()
	t := time.Duration(1)*time.Millisecond*500
	for _, confs := range cs.Area {
		//if confs.Api != "kowloon-city" {continue}
		fmt.Println("//////////////////////")
		fmt.Println("区域:", confs.Area)
		fmt.Println("地名:", confs.Name)
		fmt.Println("简写:", confs.Api)

		url := "https://deliveroo.hk/zh/restaurants/hong-kong/" + confs.Api

		resp, err := http.Get(url)
		if err != nil {
			panic(err)
		}

		if resp.StatusCode != 200 {
			panic("error")
		}

		doc, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			panic(err)
		}

		d := doc.Find("#__NEXT_DATA__").Text()
		//fmt.Println(d)
		var dd jsonAnalyze
		err = json.Unmarshal([]byte(d), &dd)
		if err != nil {
			panic(err)
		}
		//fmt.Println(dd)
		analyze(dd, confs.Area, confs.Name)

		//休眠500毫秒
		time.Sleep(t)
	}
}

//读取json配置
func jsonDecode() configS {
	f, err := os.Open("conf/area.json")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	d, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}
	var cs configS
	err = json.Unmarshal(d, &cs)
	if err != nil {
		panic(err)
	}
	return cs
}

//数据分析处理
func analyze(dd jsonAnalyze, area, place string) {
	var foods []model.Deliveroo
	var now = time.Now().Unix()
	l := len(dd.Props.InitialState.Home.Feed.Results.Data)
	if l <= 1 {
		fmt.Println("改地方没有商家入驻")
	} else {
		for i, di := range dd.Props.InitialState.Home.Feed.Results.Data {
			if i > 0 {
				header := di.Header
				if strings.Index(header, "間餐廳") > 0 {
					header = ""
				}
				for k, dv := range di.Blocks {
					UiLine := dv.UiContent.Default.UiLines
					if len(UiLine) != 3 {
						continue
					}
					var food model.Deliveroo
					food.FoodType = ","
					food.Area = area
					food.Place = place
					food.CreatedAt = now
					food.Banner = header

					fmt.Println("-------------", k)
					fmt.Println("参加活动:", header)
					f := UiLine[0].Text
					fmt.Println("商家名字:", f)
					food.Name = f

					analyzeFoodType(&food, UiLine[1].UiSpans)

					foods = append(foods, food)
				}
			}
		}

		//入库
		insertAndUpdate(foods)
	}
}

//分析菜品类型
func analyzeFoodType(food *model.Deliveroo, UiSpans []UiSpans)  {
	if len(UiSpans) == 0 {
		food.FoodType = ""
		return
	}
	if UiSpans[0].Text == "" {
		//有评分
		for n, dft := range UiSpans {
			p := dft.Text
			if n == 2 {
				//评分
				fmt.Println("评分:", p)
				vkk := strings.Split(p, " ")
				food.Score = vkk[0]
			} else if n == 4 {
				//评价
				fmt.Println("评价:", p)
				food.Evaluation = p
			} else if n > 7 && n%4 == 0 {
				//菜品类型
				fmt.Println("菜品类型:", p)
				food.FoodType += p + ","
			}
		}
	} else {
		//没有评分
		fmt.Println("评分:")
		food.Score = ""
		fmt.Println("评价:")
		food.Evaluation = ""
		for n, dft := range UiSpans {
			p := dft.Text
			if n%4 == 0 {
				//菜品类型
				fmt.Println("菜品类型:", p)
				food.FoodType += p + ","
			}
		}
	}
}

//添加数据,若存在更新
func insertAndUpdate(foods []model.Deliveroo)  {
	for _, food := range foods {
		if food.Name != "" {
			sql := "INSERT INTO `analyze`.`deliveroo` " +
				"(`name`, `score`, `evaluation`, `food_type`, `area`, `place`, `banner`, `created_at`) " +
				" VALUES (?, ?, ?, ?, ?, ?, ?, ?)" +
				" ON DUPLICATE KEY UPDATE `score`=?, `evaluation`=?, `banner`=?, `created_at`=?"

			_, err := common.DB.Exec(sql, food.Name,
				food.Score,
				food.Evaluation,
				food.FoodType,
				food.Area,
				food.Place,
				food.Banner,
				food.CreatedAt,
				food.Score,
				food.Evaluation,
				food.Banner,
				food.CreatedAt)

			if err != nil {
				fmt.Println("数据库操作失败:", err)
			}
		}
	}
}
