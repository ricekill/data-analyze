package module

import (
	"crawler-project/model"
	"fmt"
	"github.com/PuerkitoBio/goquery"
)

func analyze11(doc *goquery.Document) {
	var foods []model.FoodData
	doc.Find("#__next ul[class*=HomeFeedGrid-] li[class*=HomeFeedGrid-]").Each(func(i int, s *goquery.Selection) {
		if i >= 2 {
			fmt.Println("----------------------", i)
			var food model.FoodData
			food.FoodType = ","
			foodAnalyze(s, &food)
			fmt.Println("数据:", food)
			foods = append(foods, food)
		}
	})
	fmt.Println("总数据:", foods)
}

/*
商家分析
 */
func foodAnalyze (s *goquery.Selection, food *model.FoodData) {
	l := s.Find("li[class*=HomeFeedUICard-]").Length()
	if l == 0 {
		d, _ := s.Html()
		fmt.Println("空数据:", d)
	}
	s.Find("li[class*=HomeFeedUICard-]").Each(func(k int, ks *goquery.Selection) {
		if k == 0 {
			f, _ := ks.Find("span p").Html()
			food.Name = f
			fmt.Println("商家名字:", f)
		} else if k == 1 {
			ks.Find("span[class*=HomeFeedUICard-]").Each(func(n int, fs *goquery.Selection) {
				p, _ := fs.Find("span").Html()
				if n == 2 {
					//评分
					food.Score = p
					fmt.Println("评分:", p)
				} else if n == 4 {
					//评价
					food.Evaluation = p
					fmt.Println("评价:", p)
				} else if n > 7 && n%4 == 0 {
					//菜品类型
					food.FoodType += p + ","
					fmt.Println("菜品类型:", p)
				}
			})
		}
	})
}

func activityAnalyze(s *goquery.Selection, foods *[]model.FoodData)  {
	s.Find("ul[class*=Carousel-] li[class*=Slide-]").Each(func(j int, cs *goquery.Selection) {
		fmt.Println("----------------------", j)
		var food model.FoodData
		food.FoodType = ","

		foodAnalyze(cs, &food)

		fmt.Println("数据:", food)
		//foods = append(foods, food)
	})
}