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
	acc.Load(uid, &wg)
	wg.Wait()

	return gin.H{
		"name":        acc.Name,
		"email":       acc.Email,
		"profile_url": acc.PicURL,

		"total_downloads":        cs.TD,
		"total_downloads_change": cs.TDC,
		"total_downloads_width":  cs.TDW,
		"arrow_downloads":        cs.AD,

		"total_revenue":        cs.TR,
		"total_revenue_change": cs.TRC,
		"total_revenue_width":  cs.TRW,
		"arrow_revenue":        cs.AR,

		"total_ratings":        cs.TRT,
		"total_ratings_change": cs.TRTC,
		"total_ratings_width":  cs.TRW,
		"arrow_ratings":        cs.AR,

		"total_comments":        cs.TC,
		"total_comments_change": cs.TCC,
		"total_comments_width":  cs.TCW,
		"arrow_comments":        cs.AC,

		"data_main_jan": cs.DDJA,
		"data_main_feb": cs.DDFE,
		"data_main_mar": cs.DDMA,
		"data_main_apr": cs.DDAP,
		"data_main_may": cs.DDMY,
		"data_main_jun": cs.DDJN,
		"data_main_jul": cs.DDJL,
		"data_main_aug": cs.DDAU,
		"data_main_sep": cs.DDSE,
		"data_main_oct": cs.DDOC,
		"data_main_nov": cs.DDNO,
		"data_main_dec": cs.DDDE,

		"data_comments_jan": cs.DCJA,
		"data_comments_feb": cs.DCFE,
		"data_comments_mar": cs.DCMA,
		"data_comments_apr": cs.DCAP,
		"data_comments_may": cs.DCMY,
		"data_comments_jun": cs.DCJN,
		"data_comments_jul": cs.DCJL,
		"data_comments_aug": cs.DCAU,
		"data_comments_sep": cs.DCSE,
		"data_comments_oct": cs.DCOC,
		"data_comments_nov": cs.DCNO,
		"data_comments_dec": cs.DCDE,

		"country_1": cs.C1,
		"country_2": cs.C2,
		"country_3": cs.C3,
		"country_4": cs.C4,

		"country_1_downloads": cs.C1D,
		"country_2_downloads": cs.C2D,
		"country_3_downloads": cs.C3D,
		"country_4_downloads": cs.C4D,

		"country_1_change": cs.C1C,
		"country_2_change": cs.C2C,
		"country_3_change": cs.C3C,
		"country_4_change": cs.C4C,
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
