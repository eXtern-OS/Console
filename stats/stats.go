package stats

type CountryRec map[string]int

type CountryStats struct {
	Total   CountryRec
	TotalDT map[string]CountryRec
}

type CompanyStats struct {
	TotalDownloads int
	TotalRevenue   float64
	TotalComments  int
	TotalRatings   float64
	TotalApps      int

	DownloadsDT map[string]int
	RevenueDT   map[string]float64
	RatingsDT   map[string]float64
	CommentsDT  map[string]int

	Country CountryStats
}

func (c *CompanyStats) Load(uid string) {

}
