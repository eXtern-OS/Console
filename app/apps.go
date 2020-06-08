package app

import (
	"../db"
	"../payment"
	"../utils"
	"mime/multipart"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"time"
)

func (a *Application) CreateUID() {
	a.AppId = utils.Makehash(a.Name + a.Description + strconv.Itoa(int(time.Now().UnixNano())))
}
func (a *Application) InitVersions(version, uid, rnotes, paurl string) {
	var vr = VersionRecord{
		AppId:        a.AppId,
		Version:      version,
		MaintainerID: uid,
		ReleaseNotes: rnotes,
		ReleaseIndex: 0,
		PackageURL:   paurl,
	}
	var x []VersionRecord
	x = append(x, vr)
	a.Version.AppId = a.AppId
	a.Version.CurrentVersion = vr
	a.Version.History = VersionHistory{Versions: x}
	return
}

func (a *Application) MakeSlug() bool {
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		a.Slug = strconv.Itoa(int(time.Now().UnixNano()))
		return os.MkdirAll(filepath.Join("/packages", a.Slug), os.ModePerm) == nil
	}
	if _, err := os.Stat("/packages/" + reg.ReplaceAllString(a.Name, "")); err == nil {
		a.Slug = reg.ReplaceAllString(a.Name, "")
		return os.MkdirAll(filepath.Join("/packages", a.Slug), os.ModePerm) == nil
	} else {
		a.Slug = strconv.Itoa(int(time.Now().UnixNano()))
		return os.MkdirAll(filepath.Join("/packages", a.Slug), os.ModePerm) == nil
	}
}

func CreateFreeApp(name, description, package_url, screenshots, app_version, version_description, uid string, appIcon, appCover *multipart.FileHeader) {
	var a = Application{
		Name:        name,
		Description: description,
		Rating:      0,
		PaymentType: payment.PaymentType{
			Price:           0,
			Monthly:         false,
			Yearly:          false,
			Once:            false,
			Free:            true,
			SubscriptionUID: "",
		},
		Downloads: 0,
		Status:    "review",
		Slug:      "",
	}
	a.CreateUID()
	a.InitVersions(app_version, uid, version_description, package_url)
	a.IconURL = db.UploadImage()
}
