package wt

import (
	"context"
	"fmt"
	"time"

	"github.com/4d3v/gowtbot/pkg/helpers"
	"github.com/4d3v/gowtbot/pkg/models"
	"github.com/valyala/fasthttp"
)

const (
	graphURI              = "https://graph.facebook.com/v17.0"
	defaultRequestTimeout = 10 * time.Second
	setChoosenLang        = "set_choosen_lang"
	replyChooseLang       = "reply_choose_lang"
	markSeen              = "mark_seen"
)

var (
	client = fasthttp.Client{}
)

func respond(ctx context.Context, msgType string, msg models.Message) error {
	if msgType == markSeen {
		buf := []byte(`{"messaging_product":"whatsapp","status":"read","message_id":"`)
		buf = append(buf, msg.ID...)
		buf = append(buf, []byte(`"}`)...)

		return callApi(
			ctx,
			buf,
		)
	}

	buf := []byte(`{"messaging_product":"whatsapp","recipient_type":"individual","to":"`)
	buf = append(buf, msg.From...)

	if msgType == replyChooseLang {
		buf = append(buf, []byte(`","type":"interactive","interactive":{"type":"button","body":{"text":"*Bem vindo ao Rp Chauffeurs*\nObrigado por escolher nosso serviço. Por favor, Informe-nos o seu idioma preferido:\n\n*Welcome to Rp Chauffeurs!*\nThank you for choosing our personal chauffeur car service. Please let us know your preferred language:\n\n👇"},"action":{"buttons":[{"type":"reply","reply":{"id":"0","title":"Portuguese"}},{"type":"reply","reply":{"id":"1","title":"English"}}]}}}`)...)

		return callApi(
			ctx,
			buf,
		)
	}

	if msgType == setChoosenLang {
		ptLang := msg.Interactive.ButtonReply.ID == "0"

		if ptLang {
			buf = append(buf, []byte(`","type":"text","text":{"body":"*Perfeito! Em breve iremos te atender.*\nPara agilizar seu contato, escreva sobre qual serviço deseja agendar. Se possível, também inclua informações como seu nome, endereço, data, horário, número de passageiros/bagagens."}}`)...)
		} else {
			buf = append(buf, []byte(`","type":"text","text":{"body":"*Perfect! Soon we will assist you.*\nTo speed up your service, please write about which service you want to schedule. If possible, also include more info such as your name, address, date, time, number of passengers/luggages."}}`)...)
		}

		return callApi(
			ctx,
			buf,
		)
	}

	return nil
}

func callApi(ctx context.Context, reqBody []byte) error {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	req.SetRequestURI(fmt.Sprintf("%s/%s/messages", graphURI, wtRepo.App.WTPhoneID))
	req.Header.SetMethod(fasthttp.MethodPost)
	req.Header.Set("Accept-Encoding", "gzip")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", helpers.B2S(wtRepo.App.AccessTkn)))
	req.Header.Add("Content-Type", "application/json")
	req.SetBody(reqBody)

	res := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(res)

	dl, ok := ctx.Deadline()
	if !ok {
		dl = time.Now().Add(defaultRequestTimeout)
	}

	err := client.DoDeadline(req, res, dl)
	if err != nil {
		return fmt.Errorf("error doDeadline: %w", err)
	}

	if res.StatusCode() != fasthttp.StatusOK {
		return fmt.Errorf("unexpected response status %d", res.StatusCode())
	}

	return nil
}
