package web

import (
	"../stats"
	"../utils"
	"fmt"
	"github.com/eXtern-OS/TokenMaster"
	"strconv"
	"sync"
)

// Company stats prepared
type CSPrepared struct {
	TD  string
	TDC string
	TDW string
	AD  string

	TR  string
	TRC string
	TRW string
	AR  string

	TRT  string
	TRTC string
	TRTW string
	ART  string

	TC  string
	TCC string
	TCW string
	AC  string

	DDJA int
	DDFE int
	DDMA int
	DDAP int
	DDMY int
	DDJN int
	DDJL int
	DDAU int
	DDSE int
	DDOC int
	DDNO int
	DDDE int

	DCJA int
	DCFE int
	DCMA int
	DCAP int
	DCMY int
	DCJN int
	DCJL int
	DCAU int
	DCSE int
	DCOC int
	DCNO int
	DCDE int

	C1 string
	C2 string
	C3 string
	C4 string

	C1D int
	C2D int
	C3D int
	C4D int

	C1C string
	C2C string
	C3C string
	C4C string
}

func (c *CSPrepared) Load(cid string, wg *sync.WaitGroup) {
	var cs stats.CompanyStats
	cs.Load(cid)

	if cs.TotalApps == 0 {
		c.AD, c.AR, c.ART, c.AC = "up", "up", "up", "up"
		c.C1, c.C2, c.C3, c.C4 = "No data", "No data", "No data", "No data"
		c.C1C, c.C2C, c.C3C, c.C4C = "No data", "No data", "No data", "No data"
		wg.Done()
		return
	}

	c.TD = strconv.Itoa(cs.TotalDownloads)
	tdc := (cs.TotalDownloads - cs.DownloadsDT[utils.KeyOffset()]) * 100 / (cs.DownloadsDT[utils.KeyOffset()] + 1) //Avoid dividing by zero
	c.TDC = fmt.Sprintf("%f", tdc)
	if tdc > 0 {
		c.TDW = c.TDC
		c.AD = "up"
	} else {
		c.TDW = "0"
		c.AD = "down"
	}

	c.TR = fmt.Sprintf("%f", cs.TotalRevenue)
	var trc = ((cs.TotalRevenue - cs.RevenueDT[utils.KeyOffset()]) * 100) / (cs.RevenueDT[utils.KeyOffset()] + 1) //Avoid dividing by zero
	c.TRC = fmt.Sprintf("%f", trc)
	if trc > 0 {
		c.TRW = strconv.Itoa(int(trc))
		c.AD = "up"
	} else {
		c.TRW = "0"
		c.AD = "down"
	}

	c.TRT = fmt.Sprintf("%f", cs.TotalRatings)
	TRTc := (cs.TotalRatings - cs.RatingsDT[utils.KeyOffset()]) * 100 / (cs.RatingsDT[utils.KeyOffset()] + 1) //Avoid dividing by zero
	c.TRTC = fmt.Sprintf("%f", TRTc)
	if TRTc > 0 {
		c.TRTW = c.TRTC
		c.AD = "up"
	} else {
		c.TRTW = "0"
		c.AD = "down"
	}

	c.TC = strconv.Itoa(cs.TotalComments)
	TCc := (cs.TotalComments - cs.CommentsDT[utils.KeyOffset()]) * 100 / (cs.CommentsDT[utils.KeyOffset()] + 1) //Avoid dividing by zero
	c.TCC = fmt.Sprintf("%f", TCc)
	if TCc > 0 {
		c.TCW = c.TCC
		c.AD = "up"
	} else {
		c.TCW = "0"
		c.AD = "down"
	}

	c.C1, c.C2, c.C3, c.C4, c.C1D, c.C2D, c.C3D, c.C4D, c.C1C, c.C2C, c.C3C, c.C4C = cs.GetCountriesDataSorted()

	c.DDJA, c.DDFE, c.DDMA, c.DDAP, c.DDMY, c.DDJN, c.DDJL, c.DDAU, c.DDSE, c.DDOC, c.DDNO, c.DDDE = cs.GetDDDataSorted()

	c.DCJA, c.DCFE, c.DCMA, c.DCAP, c.DCMY, c.DCJN, c.DCJL, c.DCAU, c.DCSE, c.DCOC, c.DCNO, c.DCDE = cs.GetCDDataSorted()

	wg.Done()
	return
}

type Account struct {
	Name   string
	Email  string
	PicURL string
	UID    string
}

func (a *Account) Load(uid string, wg *sync.WaitGroup) {
	if t, acc := TokenMaster.GetUserByUID(uid); t {
		a.Name = acc.Name
		a.Email = acc.Email
		a.PicURL = acc.AvatarURL
	}
	wg.Done()
	return
}
