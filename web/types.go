package web

import (
	"../stats"
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

	DMJA int
	DMFE int
	DMMA int
	DMAP int
	DMMY int
	DMJN int
	DMJL int
	DMAU int
	DMSE int
	DMOC int
	DMNO int
	DMDE int

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

	C1 int
	C2 int
	C3 int
	C4 int

	C1C int
	C2C int
	C3C int
	C4C int
}

func (c *CSPrepared) Load(cid string, wg *sync.WaitGroup) {
	var cs stats.CompanyStats
	cs.Load(cid)
}

type Account struct {
	Name   string
	Email  string
	PicURL string
}

func (a *Account) Load(wg *sync.WaitGroup) {

}
