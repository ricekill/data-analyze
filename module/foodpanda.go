package module

import (
	"crawler-project/common"
	"crawler-project/model"
	"database/sql"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

type Foodpanda struct {
	//分析数据目的域名
	Url string

	//坐标范围
	//用于统计香港数据
	Coordinates struct {
		North Coordinates
		East  Coordinates
		South Coordinates
		West  Coordinates
	}

	//定位增量
	Incremental float64
}

type Coordinates struct {
	Lat float64
	Lng float64
}

//foodpanda爬虫对象
var Fd *Foodpanda

//22.520000,114.110000 //北
//22.220000,114.110000 //南
//
//22.380000,113.900000 //西
//22.380000,114.286787 //东
func init() {
	Fd = new(Foodpanda)
	Fd.Url = `https://www.foodpanda.hk`

	Fd.Coordinates.East.Lng = 114.280000
	Fd.Coordinates.West.Lng = 113.900000
	Fd.Coordinates.South.Lat = 22.220000
	Fd.Coordinates.North.Lat = 22.520000

	Fd.Incremental = 0.01
}

//分析数据开始
func (f *Foodpanda) FoodpandaStart() {
	for lng := f.Coordinates.West.Lng; lng <= f.Coordinates.East.Lng; lng += f.Incremental {
		for lat := f.Coordinates.South.Lat; lat <= f.Coordinates.North.Lat; lat += f.Incremental {
			url := fmt.Sprintf(
				"%s/zh/restaurants/lat/%f/lng/%f/city/%s/address/test?verticalTab=restaurants",
				f.Url,
				lat,
				lng,
				"%E9%A6%99%E6%B8%AF")

			fmt.Println("/////////////////////////////////////////////////////////////////////")
			fmt.Printf("经纬度:%f,%f", lat, lng)
			common.Log.Infoln("url:", url)

			doc := f.getRequest(url)
			if doc == nil {
				continue
			}

			f.analyzeHtml(doc)
		}
	}

	return
}

func (f *Foodpanda) getRequest(url string) *goquery.Document {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		common.Log.Errorln("request err:", err)
		return nil
	}
	req.Header.Set("accept-language", "zh-CN,zh;q=0.9,en;q=0.8")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		common.Log.Errorln("http err:", err)
		return nil
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		common.Log.Errorln("goquery err:", err)
		return nil
	}

	return doc
}

func (f *Foodpanda) analyzeHtml(doc *goquery.Document) {
	gs := doc.Find(".restaurants__list .opened .opened a[class*=hreview-aggregate]")
	l := gs.Length()
	t := time.Duration(1) * time.Millisecond * 500
	if l == 0 {
		fmt.Println("该坐标没有餐厅")
		common.Log.Infoln("该坐标没有餐厅")
		time.Sleep(t)
		return
	}

	gs.Each(func(i int, s *goquery.Selection) {
		fmt.Println("--------------------:", i)

		pUrl, _ := s.Attr("href")
		pUrl = f.Url + pUrl
		fn := s.Find(".vendor-info .headline .name").Text()
		if f.checkNameInDB(fn) {
			ddoc := f.getRequest(pUrl)
			if ddoc != nil {
				f.detailHtml(ddoc)
			}
		}

		//休眠500毫秒
		time.Sleep(t)
	})
}

func (f *Foodpanda) detailHtml(ddoc *goquery.Document) {
	var foodpanda model.Foodpanda
	foodpanda.CreatedAt = time.Now().Unix()
	foodpanda.FoodType = ","
	info := ddoc.Find(".modal .modal-body .infos")
	fn := info.Find(".info-headline .vendor-name").Text()
	fmt.Println("商家名字:", fn)
	foodpanda.Name = fn
	fs := info.Find(".info-headline .ratings-component .rating strong").Text()
	fmt.Println("商家評分:", fs)
	fsFloat, _ := strconv.ParseFloat(fs, 64)
	foodpanda.Score = fsFloat
	fe := info.Find(".info-headline .ratings-component .count").Text()
	r, _ := regexp.Compile(`\d+`)
	feInt, _ := strconv.Atoi(r.FindString(fe))
	fmt.Println("商家評價:", feInt)
	foodpanda.Evaluation = feInt

	info.Find(".vendor-cuisines li").Each(func(k int, ts *goquery.Selection) {
		if k > 0 {
			ft := ts.Text()
			fmt.Println("菜品類型:", ft)
			foodpanda.FoodType += ft + ","
		}
	})

	//地址
	fa := ddoc.Find(".modal .modal-body .content .vendor-location").Text()
	fmt.Println("商家地址:", fa)
	foodpanda.Address = fa

	//坐标
	imgUrl, _ := ddoc.Find(".modal .modal-body .static-map-container img").Attr("data-img-url")
	rd, _ := regexp.Compile(`center=(\d+\.\d+),(\d+\.\d+)&`)
	fz := rd.FindStringSubmatch(imgUrl)
	fmt.Println("坐标:", fz[1], ",", fz[2])
	foodpanda.Latitude = fz[1]
	foodpanda.Longitude = fz[2]

	//添加数据
	f.insert(foodpanda)
}

func (f *Foodpanda) checkNameInDB(name string) bool {
	row := common.DB.QueryRow("select id from `analyze`.`foodpanda` where name = ?", name)
	var id int
	err := row.Scan(&id)
	if err == sql.ErrNoRows {
		return true
	} else if err != nil {
		common.Log.Errorln("mysql err:", err)
		return false
	}

	fmt.Println("该商家已存在id:", id)
	return false
}

func (f *Foodpanda) insert(foodpanda model.Foodpanda) {
	if foodpanda.Name != "" {
		sql := "INSERT INTO `analyze`.`foodpanda` " +
			"(`name`, `score`, `evaluation`, `food_type`, `address`, `latitude`, `longitude`, `banner`, `created_at`) " +
			" VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)" +
			" ON DUPLICATE KEY UPDATE `score`=?, `evaluation`=?, `created_at`=?"

		_, err := common.DB.Exec(sql, foodpanda.Name,
			foodpanda.Score,
			foodpanda.Evaluation,
			foodpanda.FoodType,
			foodpanda.Address,
			foodpanda.Latitude,
			foodpanda.Longitude,
			foodpanda.Banner,
			foodpanda.CreatedAt,
			foodpanda.Score,
			foodpanda.Evaluation,
			foodpanda.CreatedAt)

		if err != nil {
			common.Log.Errorln("mysql err:", err)
			return
		}
	}
}
