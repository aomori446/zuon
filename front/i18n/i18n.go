package i18n

import (
	"embed"
	"encoding/json"
	"fmt"
	
	"fyne.io/fyne/v2"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

//go:embed locales/*.json
var localeFS embed.FS

var bundle *i18n.Bundle
var localizer *i18n.Localizer

func Init() {
	bundle = i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	
	_, err := bundle.LoadMessageFileFS(localeFS, "locales/active.en.json")
	if err != nil {
		fmt.Println("Error loading en locale:", err)
	}
	_, err = bundle.LoadMessageFileFS(localeFS, "locales/active.zh.json")
	if err != nil {
		fmt.Println("Error loading zh locale:", err)
	}
	_, err = bundle.LoadMessageFileFS(localeFS, "locales/active.ja.json")
	if err != nil {
		fmt.Println("Error loading ja locale:", err)
	}
	_, err = bundle.LoadMessageFileFS(localeFS, "locales/active.my.json")
	if err != nil {
		fmt.Println("Error loading my locale:", err)
	}
	
	currentLang := fyne.CurrentApp().Preferences().StringWithFallback("language", "zh")
	SetLang(currentLang)
}

func SetLang(lang string) {
	localizer = i18n.NewLocalizer(bundle, lang)
	fyne.CurrentApp().Preferences().SetString("language", lang)
}

func GetLang() string {
	return fyne.CurrentApp().Preferences().StringWithFallback("language", "zh")
}

func T(messageID string) string {
	return localizer.MustLocalize(&i18n.LocalizeConfig{
		MessageID: messageID,
	})
}

func Tf(messageID string, args map[string]interface{}) string {
	return localizer.MustLocalize(&i18n.LocalizeConfig{
		MessageID:    messageID,
		TemplateData: args,
	})
}
