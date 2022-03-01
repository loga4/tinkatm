package atm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"tinkoffbot/pkg/config"

	"github.com/go-redis/redis"
	tgapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	jsoniter "github.com/json-iterator/go"
	"go.uber.org/zap"
)

type TinkoffClient struct {
	tgbot  *tgapi.BotAPI
	logger *zap.Logger
	rds    *redis.Client
	cfg    *config.Config
	client *http.Client

	atms map[string]map[string]ATM
}

func NewTinkoffClient(cfg *config.Config, client *http.Client, logger *zap.Logger, tgbot *tgapi.BotAPI, rds *redis.Client, atms map[string]map[string]ATM) *TinkoffClient {
	return &TinkoffClient{cfg: cfg, client: client, logger: logger, tgbot: tgbot, rds: rds, atms: atms}
}

type TRedisClient redis.Client

func (self *TinkoffClient) SendRequest(currency string) error {
	redisKey := fmt.Sprintf("atms:%s:%d", currency, self.cfg.Telegram.Chat)
	get(self.rds, redisKey, &self.atms)

	req := &config.TRequest{
		Bounds: config.Bounds{
			BottomLeft: config.Point{
				Lat: self.cfg.Bounds.BottomLeft.Lat,
				Lng: self.cfg.Bounds.BottomLeft.Lng,
			},
			TopRight: config.Point{
				Lat: self.cfg.Bounds.TopRight.Lat,
				Lng: self.cfg.Bounds.TopRight.Lng,
			},
		},
		Filters: config.Filters{
			Banks: []string{
				"tcs",
			},
			Currencies: []string{
				currency,
			},
		},
		Zoom: self.cfg.Zoom,
	}

	data, err := json.Marshal(req)
	if err != nil {
		return err
	}

	resp, err := self.client.Post("https://api.tinkoff.ru/geo/withdraw/clusters", "application/json", bytes.NewReader(data))
	if err != nil {
		return err
	}

	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var r TinkResponse

	if err := jsoniter.Unmarshal(result, &r); err != nil {
		return err
	}

	if len(r.Payload.Clusters) == 0 {
		return nil
	}

	checkedAtms := make(map[string]string)

	var index = 0
	for _, cluster := range r.Payload.Clusters {

		if len(cluster.Points) == 0 {
			return nil
		}

		//all atms
		for _, point := range cluster.Points {
			_, foundCur := self.atms[currency]
			if foundCur == false {
				self.atms[currency] = make(map[string]ATM)
			}

			_, found := self.atms[currency][point.ID]
			if found == false {
				self.atms[currency][point.ID] = ATM{
					ID:      point.ID,
					Address: point.Address,
					Limit:   0,
					Index:   index,
					Location: Point{
						Lat: float32(point.Location.Lat),
						Lng: float32(point.Location.Lng),
					},
				}
				index++
			}

			checkedAtms[point.ID] = point.ID

			if atm, ok := self.atms[currency][point.ID]; ok {
				var messages []string

				for _, limit := range point.Limits {

					if limit.Currency != currency {
						continue
					}

					link := fmt.Sprintf("[link](https://www.tinkoff.ru/maps/atm/?latitude=%f&longitude=%f&currency=%s&partner=tcs)\n", atm.Location.Lat, atm.Location.Lng, limit.Currency)

					if atm.Limit > limit.Amount {
						fmt.Printf("Decreased limit limit\n%#s\n\n", atm)

						msg := fmt.Sprintf("*Лимит уменьшился*, %d %s\n\n", limit.Amount, limit.Currency)
						msg += fmt.Sprintf("Address: %s\n", atm.Address)
						msg += link

						messages = append(messages, msg)

					} else if atm.Limit < limit.Amount {
						fmt.Printf("Increased limit limit\n%#s\n\n", atm)

						msg := fmt.Sprintf("*Завезли бабло!*, %d %s\n\n", limit.Amount, limit.Currency)
						msg += fmt.Sprintf("Address: %s\n", atm.Address)
						msg += link

						messages = append(messages, msg)
					}

					atm.Limit = limit.Amount
				}

				self.atms[currency][point.ID] = atm

				if len(messages) != 0 {
					tgmsg := tgapi.NewMessage(self.cfg.Telegram.Chat, strings.Join(messages, "\n"))
					tgmsg.ParseMode = tgapi.ModeMarkdown
					self.tgbot.Send(tgmsg)
				}
			}
		}
	}

	//checked
	for _, atm := range self.atms[currency] {

		_, ok := checkedAtms[atm.ID]
		if !ok {
			fmt.Printf("Removed unused id")

			delete(self.atms[currency], atm.ID)

			msg := fmt.Sprintf("*Бабло закончилось*, %s\n\n", currency)
			msg += fmt.Sprintf("Address: %s\n", atm.Address)
			msg += fmt.Sprintf("[link](https://www.tinkoff.ru/maps/atm/?latitude=%f&longitude=%f&currency=%s&partner=tcs)", atm.Location.Lat, atm.Location.Lng, currency)

			tgmsg := tgapi.NewMessage(self.cfg.Telegram.Chat, msg)
			tgmsg.ParseMode = tgapi.ModeMarkdown
			self.tgbot.Send(tgmsg)
		}
	}

	checkedAtms = nil

	set(self.rds, redisKey, self.atms)

	return nil
}

func set(rds *redis.Client, key string, value interface{}) error {
	p, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return rds.Set(key, p, -1).Err()
}

func get(rds *redis.Client, key string, dest interface{}) error {
	p, err := rds.Get(key).Result()
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(p), dest)
}
