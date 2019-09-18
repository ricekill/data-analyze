package module

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"regexp"
)

func FoodpandaStart()  {
	fmt.Println("//////////////////////")

	url := "https://www.foodpanda.hk/zh/restaurants/lat/22.281184/lng/114.162595/city/%E9%A6%99%E6%B8%AF/address/sdfsdf?verticalTab=restaurants"

	req,err := http.NewRequest("GET",url,nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("accept-language", "zh-CN,zh;q=0.9,en;q=0.8")
	resp ,err :=  http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		panic(err)
	}

	analyzeHtml(doc)
}

func analyzeHtml(doc *goquery.Document)  {
	gs := doc.Find(".restaurants__list .opened .opened a[class*=hreview-aggregate]")
	l := gs.Length()
	if l == 0 {
		panic("error:")
	}

	gs.Each(func(i int, s *goquery.Selection) {
		fmt.Println("--------------------:", i)
		fn := s.Find(".vendor-info .headline .name").Text()
		fmt.Println("商家名字:", fn)
		fs := s.Find(".vendor-info .rating strong").Text()
		fmt.Println("商家評分:", fs)
		fe := s.Find(".vendor-info .count").Text()
		r, _ := regexp.Compile(`\d+`)
		fsint := r.FindString(fe)
		fmt.Println("商家評價:", fsint)

		s.Find(".vendor-characteristic span").Each(func(k int, ts *goquery.Selection) {
			ft := ts.Text()
			fmt.Println("菜品類型:", ft)
		})
	})
}

func detailStart()  {
	url := "https://www.foodpanda.hk/zh/restaurant/v6dr/nosh-central#restaurant-info"

	req,err := http.NewRequest("GET",url,nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("accept-language", "zh-CN,zh;q=0.9,en;q=0.8")
	resp ,err :=  http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		panic(err)
	}

	detailHtml(doc)
}

func detailHtml(ddoc *goquery.Document)  {
	info := ddoc.Find(".modal .modal-body .infos")
	fn := info.Find(".info-headline .vendor-name").Text()
	fmt.Println("商家名字:", fn)
	fs := info.Find(".info-headline .ratings-component .rating strong").Text()
	fmt.Println("商家評分:", fs)
	fe := info.Find(".info-headline .ratings-component .count").Text()
	r, _ := regexp.Compile(`\d+`)
	fsint := r.FindString(fe)
	fmt.Println("商家評價:", fsint)

	info.Find(".vendor-cuisines li").Each(func(k int, ts *goquery.Selection) {
		if k > 0 {
			ft := ts.Text()
			fmt.Println("菜品類型:", ft)
		}
	})

	//地址
	fa := ddoc.Find(".modal .modal-body .content .vendor-location").Text()
	fmt.Println("商家地址:", fa)

	//坐标
	imgUrl, _ := ddoc.Find(".modal .modal-body .static-map-container img").Attr("data-img-url")
	rd, _ := regexp.Compile(`center=(\d+\.\d+),(\d+\.\d+)&`)
	fz := rd.FindStringSubmatch(imgUrl)
	fmt.Println("坐标:", fz[1], ",", fz[2])
}
