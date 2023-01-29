package argostranslate

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"gtt/internal/lock"
	"io/ioutil"
	"net/http"
)

const (
	textURL = "https://translate.argosopentech.com/translate"
)

type ArgosTranslate struct {
	srcLang    string
	dstLang    string
	EngineName string
	SoundLock  *lock.Lock
}

func (t *ArgosTranslate) GetEngineName() string {
	return t.EngineName
}

func (t *ArgosTranslate) GetAllLang() []string {
	return lang
}

func (t *ArgosTranslate) GetSrcLang() string {
	return t.srcLang
}

func (t *ArgosTranslate) GetDstLang() string {
	return t.dstLang
}

func (t *ArgosTranslate) SetSrcLang(srcLang string) {
	t.srcLang = srcLang
}

func (t *ArgosTranslate) SetDstLang(dstLang string) {
	t.dstLang = dstLang
}

func (t *ArgosTranslate) SwapLang() {
	t.srcLang, t.dstLang = t.dstLang, t.srcLang
}

func (t *ArgosTranslate) Translate(message string) (
	translation string,
	definition string,
	partOfSpeech string,
	err error) {
	var data interface{}

	res, err := http.Post(textURL,
		"application/json",
		bytes.NewBuffer([]byte(fmt.Sprintf(`{
			"q": "%s",
			"source": "%s",
			"target": "%s" }`, message, langCode[t.srcLang], langCode[t.dstLang]))))
	if err != nil {
		return "", "", "", err
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", "", "", err
	}
	if err = json.Unmarshal(body, &data); err != nil {
		return "", "", "", err
	}

	if len(data.(map[string]interface{})) > 0 {
		translation += fmt.Sprintf("%v",
			data.(map[string]interface{})["translatedText"])

		return translation, definition, partOfSpeech, nil
	}
	return "", "", "", errors.New("Translation not found")
}