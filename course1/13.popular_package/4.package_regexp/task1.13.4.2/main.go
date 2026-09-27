package main

import (
	"fmt"
	"regexp"
)

type Ad struct {
	Title       string
	Description string
}

func main() {
	ads := []Ad{
		{
			Title:       "Куплю велосипед MeRiDa",
			Description: "Куплю велосипед meriDA в хорошем состоянии.",
		},
		{
			Title:       "Продам ВаЗ 2101",
			Description: "Продам ваз 2101 в хорошем состоянии.",
		},
		{
			Title:       "Продам БМВ",
			Description: "Продам бМв в хорошем состоянии.",
		},
		{
			Title:       "Продам macBook pro",
			Description: "Продам macBook PRO в хорошем состоянии.",
		},
	}

	ads = censorAds(ads, map[string]string{
		"велосипед merida": "телефон Apple",
		"ваз":              "ВАЗ",
		"бмв":              "BMW",
		"macbook pro":      "Macbook Pro",
	})

	for _, ad := range ads {
		fmt.Println(ad.Title)
		fmt.Println(ad.Description)
		fmt.Println()
	}
}

func censorAds(ads []Ad, censor map[string]string) []Ad {
	for i := range ads{

		for old, new := range censor{
			re := regexp.MustCompile("(?i)"+old)
			ads[i].Title = re.ReplaceAllString(ads[i].Title, new)
		}

		for old, new := range censor{
			re := regexp.MustCompile("(?i)"+old)
			ads[i].Description = re.ReplaceAllString(ads[i].Description, new)
		}
	}
	return ads
}
