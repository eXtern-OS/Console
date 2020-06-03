package web

import (
	"github.com/gin-gonic/gin"
	"sync"
)

func RenderIndex(uid, tid string) gin.H {
	var wg sync.WaitGroup

	var cs CSPrepared
	var acc Account

	wg.Add(2)
	cs.Load(tid, &wg)
	acc.Load(&wg)
	wg.Wait()

	return gin.H{
		"name":        "foxclore",
		"email":       "foxclore@zoho.com",
		"profile_url": "https://pbs.twimg.com/profile_images/1257996915167260673/ps_y8zs5_400x400.jpg",

		"total_downloads":       "-4.2",
		"total_downloads_count": "9270",
		"total_downloads_width": "55",
		"arrow_downloads":       "down",

		"total_revenue":        "8444",
		"total_revenue_change": "+1.6",
		"total_revenue_width":  "23",
		"arrow_revenue":        "up",

		"total_ratings":        "5.0",
		"total_ratings_change": "+1.3",
		"total_ratings_width":  "14",
		"arrow_ratings":        "up",

		"total_comments":        "1325",
		"total_comments_change": "+7.9",
		"total_comments_width":  "38",
		"arrow_comments":        "up",

		"data_main_jan": 4,
		"data_main_feb": 8,
		"data_main_mar": 2,
		"data_main_apr": 9,
		"data_main_may": 11,
		"data_main_jun": 6,
		"data_main_jul": 7,
		"data_main_aug": 8,
		"data_main_sep": 9,
		"data_main_oct": 15,
		"data_main_nov": 6,
		"data_main_dec": 8,

		"data_comments_jan": 14,
		"data_comments_feb": 20,
		"data_comments_mar": 2,
		"data_comments_apr": 8,
		"data_comments_may": 1,
		"data_comments_jun": 7,
		"data_comments_jul": 3,
		"data_comments_aug": 6,
		"data_comments_sep": 11,
		"data_comments_oct": 5,
		"data_comments_nov": 9,
		"data_comments_dec": 4,

		"country_1": "USA",
		"country_2": "China",
		"country_3": "Australia",
		"country_4": "Canada",

		"country_1_downloads": 666,
		"country_2_downloads": 5421,
		"country_3_downloads": 234,
		"country_4_downloads": 100,

		"country_1_change": "-1.9",
		"country_2_change": "+4.6",
		"country_3_change": "-2.1",
		"country_4_change": "+3.7",
	}
}

func RenderApplicationPage() gin.H {
	return gin.H{
		"app_name":            "Photos",
		"app_description":     "Photos app for eXternOS",
		"app_rating":          4.9,
		"app_rating_width":    95,
		"total_app_downloads": 1400,
		"total_app_revenue":   "$10900",
		"total_app_comments":  600,
		"top_app_country":     "USA (68%)",
		"update_type":         "url",
	}
}
