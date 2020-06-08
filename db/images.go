package db

import (
	"github.com/masci/flickr"
	"sync"
)

var FLICKR_KEY string
var FLICKR_SECRET string

type FlickrClient struct {
	Mutex  sync.Mutex
	Client *flickr.FlickrClient
}

var FClient FlickrClient

func UploadImage() string {
	return ""
}
