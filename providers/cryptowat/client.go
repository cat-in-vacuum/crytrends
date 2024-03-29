package cryptowat

import (
	"encoding/json"
	"github.com/cat-in-vacuum/crytrade/providers"
	"github.com/cat-in-vacuum/crytrade/providers/cryptowat/types"
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
	"time"
)

const (
	// константы для логирования
	pkgName   = "cryptowat client"
	msgErrFmt = "error in %s; causer: %s; reqID: %s;"
	msgCommonFmt = "%s in " + pkgName + " => %s; reqID: %s;"

	// пути для апи
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
	addr        string
	exec        *http.Client
	Assets      *AssetsContainer
	RetryPolicy RetryPolicy
}

// todo реализация политики повторов
//  на случай если сервер поставщика временно отваливается
type RetryPolicy struct{}

type OHLCParams struct {
	Periods []string
	After   int
	Before  int
}

func (o OHLCParams) JoinPeriods() (out string) {
	for i := 0; i < len(o.Periods); i++ {
		out += o.Periods[i] + ","
	}
	return
}

func New(exec *http.Client, assets *AssetsContainer, policy RetryPolicy) *Client {
	return &Client{
		addr:        mainPath,
		exec:        exec,
		Assets:      assets,
		RetryPolicy: policy,
	}
}

// функция для загрузки данных по всему существующему в памяти списку ассетов
// todo вообще, было бы хорошо реализовать тонкую настройку параметров запроса для каждого ассета
//  но, для примера я этого не делал
func (c Client) GetAllOHLC(ctx providers.Context, params OHLCParams) (types.OHLCRespCommon, error) {
	// логируем весь список текущих ассетов перед началом загрузки,
	// что бы понимать начальные данные при ошибке
	id := ctx.GetID()
	log.Debug().Interface("current assets list", c.Assets.store).Msgf(msgCommonFmt, "start downloading OHLC from assets store", "client.GetAllOHLC()", id)

	// todo было бы хорошо реализовать pool под эти данные
	var out = make(types.OHLCRespCommon, len(c.Assets.store))
	if c.Assets.store == nil ||
		len(c.Assets.store) == 0 {
		return nil, errors.New("assets must be initialed and be not empty")
	}

	// todo тут реализовать кокурентный запрос для каждого ассета
	for market, pairs := range c.Assets.store {
		marketData := make(map[string]types.OHLCResp, len(pairs))
		for _, pair := range pairs {
			// тут может встпуать в бой запрос с политикой повтров
			// с.getOHLCRetryable()
			resp, err := c.getOHLC(ctx, market, pair, params)
			if err != nil {
				log.Error().Err(err).Strs("asset", []string{pair, market}).Msgf(msgErrFmt, pkgName, "client.getOHLC()", id)
				continue
			}
			resp.QueryTime = time.Now()
			marketData[pair] = *resp
			out[market] = marketData
		}
	}

	return out, nil
}

func (c Client) GetOHLCFromAsset(ctx providers.Context, market, pair string, params OHLCParams) (*types.OHLCResp, error) {
	return c.getOHLC(ctx, market, pair, params)
}

func (c Client) getOHLCRetryable() (*types.OHLCResp, error) {
	return nil, nil
}

// todo ф-ци самого запроса не долджна быть привязанна к конкретному роуту поставщика
//  стоит реализовтаь отдельный тип, который будет характеризовать все методы апи и потреблять тип
//  который будет описывать необходимый запрос
func (c Client) getOHLC(ctx providers.Context, market, pair string, params OHLCParams) (*types.OHLCResp, error) {
	var out types.OHLCResp
	var reqID string
	reqURL, err := url.Parse(mainPath)
	reqID = ctx.GetID()
	if err != nil {
		// если не удалось получить ИД прожолжаем выполнять запрос показывая, что есть проблема
		reqID = errors.Wrap(err, "error in  ctx.GetID(), getOHLC()").Error()
	}

	// todo убрать это отсюда и вынести в отдельную ф-цию
	reqURL.Path = path.Join(reqURL.Path, marketPath, market, pair, ohlcPath)
	q := reqURL.Query()
	q.Set("periods", params.JoinPeriods())
	q.Set("after", strconv.Itoa(params.After))
	q.Set("before", strconv.Itoa(params.Before))
	reqURL.RawQuery = q.Encode()

	log.Debug().Str("req url", reqURL.String()).Msgf(msgCommonFmt, "starting download data", "client.getOHLC()", reqID)
	resp, err := c.exec.Get(reqURL.String())
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("non 200 http.StatusCode")
	}

	log.Debug().Str("resp status", resp.Status).Msgf(msgCommonFmt, "starting download data", "client.getOHLC()", reqID)
	if err = json.Unmarshal(body, &out); err != nil {
		return nil, err
	}

	// переиспользуем соединение, что бы не выделять под него память повтороно при кажом запросе
	_, _ = io.Copy(ioutil.Discard, resp.Body)
	_ = resp.Body.Close()

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
