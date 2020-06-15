package leagues

import (
	"context"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/leighmacdonald/steamid"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"log"
	"net/url"
	"strings"
)

type ESEARank int

const (
	ESEARankNone   ESEARank = 0
	ESEARankOpen   ESEARank = 1
	ESEARankIM     ESEARank = 2
	ESEARankInvite ESEARank = 3
)

const (
	ESEAURL = "https://play.esea.net/index.php?s=search&source=users&query=%s"
)

func getESEA(ctx context.Context, steam steamid.SID64) (LeagueHistory, error) {
	lHist := LeagueHistory{
		League: "ESEA",
	}
	query := url.Values{}
	query.Add("url", fmt.Sprintf(ESEAURL, steamid.SID64ToSID3(steam)))
	u, _ := url.Parse("http://172.16.1.20:8050/render.html")
	u.RawQuery = query.Encode()
	resp, err := get(ctx, u.String(), nil)
	if err != nil {
		return lHist, errors.Wrapf(err, "Failed")
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return lHist, errors.Wrapf(err, "Failed to read response body: %v", err)
	}
	bodyStr := string(body)
	eseaID := parseESEASearch(bodyStr)
	if eseaID > 0 {
		parseHistory(eseaID)
	}
	lHist.Exists = !strings.Contains(strings.ToLower(bodyStr),
		"no ugc tf2 league history")
	return lHist, nil
}

func parseESEASearch(body string) int {
	dom, err := goquery.NewDocumentFromReader(strings.NewReader(body))
	if err != nil {
		logrus.WithFields(logrus.Fields{}).WithError(errors.WithStack(err)).Error("Could not parse esea search")
		return 0
	}
	dom.Find(".result-container").Each(func(i int, selection *goquery.Selection) {
		a := selection.Next().Text()
		log.Panicln(a)
	})
	return 0
}

func parseHistory(eseaID int) (ESEARank, int) {
	var rank ESEARank
	var count int
	_, err := get(context.Background(),
		fmt.Sprintf("https://play.esea.net/users/%d?tab=history", eseaID), nil)
	if err != nil {
		logrus.WithFields(logrus.Fields{}).WithError(errors.WithStack(err)).Error("Could not fetch esea history")
		return rank, count
	}
	return rank, count
}