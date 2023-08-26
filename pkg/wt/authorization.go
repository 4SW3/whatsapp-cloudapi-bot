package wt

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"

	"github.com/4d3v/gowtbot/pkg/helpers"
	"github.com/valyala/fasthttp"
)

var (
	headerNameXSign       = []byte("x-hub-signature-256")
	signaturePrefix       = []byte("sha256=")
	errNoXSignHeader      = errors.New("there is no x-sign header")
	errInvalidXSignHeader = errors.New("invalid x-sign header")
)

func auth(ctx *fasthttp.RequestCtx) error {
	sign := ctx.Request.Header.PeekBytes(headerNameXSign)
	if len(sign) == 0 {
		return errNoXSignHeader
	}

	// copy fasthttp body
	body := make([]byte, len(ctx.Request.Body()))
	copy(body, ctx.Request.Body())

	validSign, err := isValidSign(sign, body)
	if err != nil {
		return err
	}
	if !validSign {
		return errInvalidXSignHeader
	}

	return nil
}

func isValidSign(sign, body []byte) (bool, error) {
	// actualSign, err := hex.DecodeString(string(sign[len(signaturePrefix):]))
	actualSign, err := hex.DecodeString(helpers.B2S(sign[len(signaturePrefix):]))
	if err != nil {
		return false, err
	}
	return hmac.Equal(signBody(body), actualSign), nil
}

func signBody(body []byte) []byte {
	h := hmac.New(sha256.New, wtRepo.App.AppSecret)
	h.Reset()
	h.Write(body)
	return h.Sum(nil)
}
