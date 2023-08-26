package wt

import (
	"bytes"
	"context"

	"github.com/4d3v/gowtbot/pkg/data"
	"github.com/4d3v/gowtbot/pkg/helpers"
	"github.com/4d3v/gowtbot/pkg/models"
	"github.com/valyala/fasthttp"
)

func WebhookHandler(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("text/plain; charset=utf8")

	// ctx.Response.Header.Set("Access-Control-Allow-Origin", "http://localhost:3000")
	// if bytes.Equal(ctx.Method(), []byte("OPTIONS")) {
	// 	ctx.Response.Header.Set("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
	// 	ctx.Response.Header.Set("Access-Control-Allow-Headers", "Accept, Content-Type, X-CSRF-Token, Authorization")
	// }

	if bytes.Equal(ctx.Path(), []byte("/ping")) && bytes.Equal(ctx.Method(), []byte("GET")) {
		pong(ctx)
		return
	}

	if bytes.Equal(ctx.Path(), []byte("/webhook")) && bytes.Equal(ctx.Method(), []byte("GET")) {
		hub(ctx)
		return
	}

	if bytes.Equal(ctx.Path(), []byte("/webhook")) && bytes.Equal(ctx.Method(), []byte("POST")) {
		wtbot(ctx)
		return
	}

	if bytes.Equal(ctx.Path(), []byte("/send")) && bytes.Equal(ctx.Method(), []byte("POST")) {
		wtsender(ctx)
		return
	}

	if bytes.Equal(ctx.Path(), []byte("/deallocate")) && bytes.Equal(ctx.Method(), []byte("PUT")) {
		wtdeallocator(ctx)
		return
	}

	if bytes.Equal(ctx.Path(), []byte("/block")) && bytes.Equal(ctx.Method(), []byte("PUT")) {
		wtblockunblock(ctx)
		return
	}

	ctx.Error("not found", fasthttp.StatusNotFound)
}

func pong(ctx *fasthttp.RequestCtx) {
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.Write([]byte("Pong!"))
}

func hub(ctx *fasthttp.RequestCtx) {
	hubChallenge := ctx.QueryArgs().Peek("hub.challenge")
	hubVerifyTkn := ctx.QueryArgs().Peek("hub.verify_token")

	if !bytes.Equal(hubVerifyTkn, wtRepo.App.HubToken) {
		ctx.Error("bad request", fasthttp.StatusBadRequest)
		return
	}

	ctx.Write(hubChallenge)
}

func wtbot(ctx *fasthttp.RequestCtx) {
	if err := auth(ctx); err != nil {
		ctx.Error(err.Error(), fasthttp.StatusUnauthorized)
		return
	}

	data, err := data.WebhookUnmarshalJSON(ctx.PostBody())
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}

	// whatsapp api data
	body := data.Entry[0].Changes[0]

	if !isFromWebhookMessages(body.Field) {
		ctx.Error("bad request", fasthttp.StatusBadRequest)
		return
	}

	if !hasContent(body.Value) {
		// ctx.Error("bad request", fasthttp.StatusBadRequest)
		ctx.SetStatusCode(fasthttp.StatusNoContent)
		return
	}

	if len(body.Value.Messages) > 0 {
		err := respond(context.TODO(), markSeen, body.Value.Messages[0])
		if err != nil {
			ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
			return
		}
	}

	phoneBlocked, err := checkPhoneBlocked(body.Value.Contacts[0].WaID)
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		return
	}
	if phoneBlocked {
		ctx.Error("phone blocked", fasthttp.StatusUnauthorized)
		return
	}

	for _, msg := range body.Value.Messages {
		dispatcher(msg, body.Value.Contacts[0])
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.Write([]byte("OK"))
}

// this function is using a map to cache phone numbers
// so it gains performance by not having to query firebase store every time
func dispatcher(msg models.Message, contact models.Contact) error {
	// if wtRepo.Phones does not contain client phone number then check if its in firebase store in repo.DBRepo variable
	botControl, ok := checkPhoneHash(wtRepo.Phones, contact.WaID)
	if ok {
		err := wtRepo.DBRepo.UpdateChat(contact.WaID, msg)
		if err != nil {
			return err
		}

		if !botControl {
			return nil
		}

		// if msg.Type is interactive means the client has chosen a language by pressing the button
		if msg.Type == "interactive" {
			return updateChatPlusBotHelper(contact.WaID, msg)
		}

		// if msg.Type is anything else and botControl is true means the client has not chosen a language yet
		// so we send him again an interactive message asking him to choose a language
		return respond(context.TODO(), replyChooseLang, msg)
	}

	// user is not in wtRepo.Phones cache, so try to check if it exists in firebase store
	found, botControl, err := getUserHelper(contact.WaID)
	if err != nil {
		return err
	}
	if !found {
		return createUserHelper(contact, msg)
	}

	wtRepo.Phones[contact.WaID] = botControl

	err = wtRepo.DBRepo.UpdateChat(contact.WaID, msg)
	if err != nil {
		return err
	}

	if !botControl {
		return nil
	}

	if msg.Type == "interactive" {
		return respond(context.TODO(), setChoosenLang, msg)
	}

	return respond(context.TODO(), replyChooseLang, msg)
}

func isFromWebhookMessages(field string) bool {
	return field == "messages"
}

func hasContent(value models.Value) bool {
	// facebook sends messages with no content
	// TODO check the reason for this
	return (value.Messages != nil && len(value.Messages) > 0) &&
		(value.Contacts != nil && len(value.Contacts) > 0)
}

func checkPhoneBlocked(phone string) (bool, error) {
	// check if wtRepo.BlockedPhones is loaded
	// otherwise call wtRepo.DBRepo.GetBlockedPhones() (map[string]bool, error)
	if !wtRepo.BlockedPhones["LOADED"] {
		blockedPhones, err := wtRepo.DBRepo.GetBlockedPhones()
		if err != nil {
			return false, err
		}

		wtRepo.BlockedPhones = blockedPhones
		wtRepo.BlockedPhones["LOADED"] = true
	}

	if wtRepo.BlockedPhones[phone] {
		return true, nil
	}

	return false, nil
}

func checkPhoneHash(m map[string]bool, k string) (bool, bool) {
	val, ok := m[k]
	return val, ok
}

func getUserHelper(phone string) (bool, bool, error) {
	usrPhone, botControl, err := wtRepo.DBRepo.GetUser(phone)
	if err != nil {
		return true, botControl, err
	}
	if usrPhone == "" {
		return false, botControl, nil
	}
	return true, botControl, nil
}

func createUserHelper(contact models.Contact, msg models.Message) error {
	err := wtRepo.DBRepo.CreateUser(contact, msg)
	if err != nil {
		return err
	}
	wtRepo.Phones[contact.WaID] = true
	return respond(context.TODO(), replyChooseLang, msg)
}

func updateChatPlusBotHelper(phone string, message models.Message) error {
	err := wtRepo.DBRepo.UpdateChatPlusBotControl(phone, message, false)
	if err != nil {
		return err
	}
	wtRepo.Phones[phone] = false
	return respond(context.TODO(), setChoosenLang, message)
}

// used to send messages to clients
func wtsender(ctx *fasthttp.RequestCtx) {
	body := ctx.PostBody()

	data, err := data.AdmMessageUnmarshalJSON(body)
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}

	err = wtRepo.DBRepo.UpdateChat(data.To, data)
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		return
	}

	err = callApi(ctx, body)
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		return
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.Write([]byte("OK"))
}

// used to clear the chat of a client (delete chat from a specific phone number)
func wtdeallocator(ctx *fasthttp.RequestCtx) {
	body := ctx.PostBody()
	body = body[14 : len(body)-3] // e.g { "phone": "5521999990000" } body[14:len(body)-3] = 5521999990000

	// removing a phone number from firestore if it exists
	err := wtRepo.DBRepo.RemoveUser(helpers.B2S(body))
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		return
	}

	// remove phone from cache if exists
	delete(wtRepo.Phones, helpers.B2S(body))

	ctx.SetStatusCode(fasthttp.StatusNoContent)
	ctx.Write(nil)
}

// used to block or unblock a phone number
func wtblockunblock(ctx *fasthttp.RequestCtx) {
	body := ctx.PostBody()

	// e.g { "phone": "5521999990000", "block": true }
	// parse body to get phone number and block value
	data, err := data.AdmBlockUnblockUnmarshalJSON(body)
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}

	if err := wtRepo.DBRepo.BlockOrUnblockUser(data); err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		return
	}

	if data.Block {
		wtRepo.BlockedPhones[data.Phone] = true
	} else {
		delete(wtRepo.BlockedPhones, data.Phone)
	}

	ctx.SetStatusCode(fasthttp.StatusNoContent)
	ctx.Write(nil)
}
