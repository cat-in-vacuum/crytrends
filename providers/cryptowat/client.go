package cryptowat

import (
	"encoding/json"
	"github.com/cat-in-vacuum/crytrade/providers/cryptowat/models"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	"sync"
)

const (
	pkgName = "cryptowat client"
	msgErrFmt = "error in %s; causer: %s"
	msgAllFmt = "%s in " + pkgName + " => %s"

	mainPath = "https://api.cryptowat.ch/"

	marketPath = "markets"
	ohlcPath   = "ohlc"
)

// вообще стоило использовать клиент https://github.com/cryptowatch/cw-sdk-go/blob/master/client/rest/rest.go
// т.к. все равно реализоавть лучше за короткий срок не получится. Но, раз уж я взялся, доделал хоть минимальный ф-ционал

type AssetsContainer struct {
	store map[string][]string
	// еслм мы заходим менять список ассетов динамически,
	// нужно защитить данные внтури мапы
	sync.RWMutex
}

type Client struct {
	addr   string
	exec   *http.Client
	assets *AssetsContainer
}

type OHLCParams struct {
	Periods []string
	After   int
	Before  int
}

func(o OHLCParams) JoinPeriods() (out string) {
	for i := 0; i < len(o.Periods); i ++  {
		out += o.Periods[i] + ","
	}
	return
}

func New(exec *http.Client, assets *AssetsContainer) *Client {
	return &Client{
		addr:   mainPath,
		exec:   exec,
		assets: assets,
	}
}

// функция для загрузки данных по всему существующему в памяти списку ассетов
// todo вообще, было бы хорошо реализовать тонкую настройку параметров запроса для каждого ассета
//  но, для примера я этого не делал
func (c Client) GetAllOHLC(params OHLCParams) (models.OHLCRespCommon, error) {
	// логируем весь список текущих ассетов перед началом загрузки,
	// что бы понимать начальные данные при ошибке
	log.Debug().Interface("current assets list", c.assets.store).Msgf(msgAllFmt, "start downloading OHLC from assets store", "client.GetAllOHLC()")

	// todo было бы хорошо реализовать pool под эти данные
	var out = make(models.OHLCRespCommon, len(c.assets.store))
	if c.assets.store == nil ||
		len(c.assets.store) == 0 {
		return nil, errors.New("assets must be initialed and be not empty")
	}

	// todo тут реализовать кокурентный запрос для каждого ассета
	for market, pairs := range c.assets.store {
		for _, pair := range pairs {
			resp, err:= c.getOHLC(market, pair, params)
			if err != nil {
				log.Error().Err(err).Msgf(msgErrFmt, pkgName, "client.getOHLC()")
				continue
			}
			out[market][pair] = *resp
		}
	}

	return out, nil
}

func (c Client) getOHLC(market, pair string, params OHLCParams) (*models.OHLCResp, error) {
	var out models.OHLCResp
	reqURL, err := url.Parse(mainPath)
	if err != nil {
		return nil, errors.Wrap(err, "error in build URL path, getOHLC()")
	}

	// todo убрать это отсюда и вынести в отдельную ф-цию
	reqURL.Path = path.Join(reqURL.Path, marketPath, market, pair, ohlcPath)
	q := reqURL.Query()
	q.Set("periods", params.JoinPeriods())
	q.Set("after", strconv.Itoa(params.After))
	q.Set("before", strconv.Itoa(params.Before))
	reqURL.RawQuery = q.Encode()

	log.Debug().Str("req url", reqURL.String()).Msgf(msgAllFmt, "starting download data", "client.getOHLC()")
	resp, err := c.exec.Get(reqURL.String())
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	log.Debug().Str("resp status", resp.Status).Msgf(msgAllFmt, "starting download data", "client.getOHLC()")
	if err = json.Unmarshal(body, &out); err != nil {
		return nil, err
	}

	// переиспользуем соединение, что бы не выделять под него память повтороно при кажом запросе
	io.Copy(ioutil.Discard, resp.Body)
	resp.Body.Close()

	return &out, nil
}

func (a AssetsContainer) GetCurrentAssets() map[string][]string {
	return a.store
}

func (a *AssetsContainer) SetAsset(market string, pairs ...string) {
	a.Lock()

	if a.store == nil {
		a.store = make(map[string][]string)
	}

	if existing, ok := a.store[market]; !ok {
		a.store[market] = pairs
	} else {
		a.store[market] = append(existing, pairs...)
	}

	a.Unlock()
}

// todo реализация возможности загрузки списка ассетов из файла
//  toml, json, etc.
func (a AssetsContainer) LoadAssetFromFile(marketList os.File) {}

// todo реализация возможности загрузки списка ассетов собственно из самого API Cryptowat
func (a AssetsContainer) LoadAssetFromAPI() {}
