package ui

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func ShowLoginWindow(a fyne.App, onLoginSuccess func(token string)) {
	w := a.NewWindow("Login - zuon")
	
	statusLabel := widget.NewLabel("Please login to continue")
	loginBtn := widget.NewButton("Login with GitHub", nil)
	progressBar := widget.NewProgressBarInfinite()
	progressBar.Hide()

	loginBtn.OnTapped = func() {
		loginBtn.Disable()
		statusLabel.SetText("Opening browser...")
		progressBar.Show()

		go func() {
			// 1. Get Login URL
			resp, err := http.Get("http://localhost:8080/auth/login")
			if err != nil {
				handleLoginError(w, loginBtn, statusLabel, progressBar, "Server error")
				return
			}
			defer resp.Body.Close()

			var loginData struct {
				RequestID string `json:"request_id"`
				LoginURL  string `json:"login_url"`
			}
			json.NewDecoder(resp.Body).Decode(&loginData)

			// 2. Open browser
			u, _ := url.Parse(loginData.LoginURL)
			fmt.Println("Opening:", loginData.LoginURL)
			a.OpenURL(u)

			// 3. Polling
			ticker := time.NewTicker(2 * time.Second)
			defer ticker.Stop()

			timeout := time.After(5 * time.Minute)

			for {
				select {
				case <-ticker.C:
					pollResp, err := http.Get(fmt.Sprintf("http://localhost:8080/auth/poll?req_id=%s", loginData.RequestID))
					if err != nil {
						continue
					}
					var pollResult struct {
						Status string `json:"status"`
						Token  string `json:"token"`
					}
					json.NewDecoder(pollResp.Body).Decode(&pollResult)
					pollResp.Body.Close()

					if pollResult.Status == "success" {
						fyne.Do(func() {
							onLoginSuccess(pollResult.Token)
							w.Close()
						})
						return
					}
				case <-timeout:
					handleLoginError(w, loginBtn, statusLabel, progressBar, "Login timeout")
					return
				}
			}
		}()
	}

	w.SetContent(container.NewCenter(
		container.NewVBox(
			widget.NewLabelWithStyle("Welcome to ZUON", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
			statusLabel,
			loginBtn,
			progressBar,
		),
	))

	w.Resize(fyne.NewSize(300, 200))
	w.CenterOnScreen()
	w.Show()
}

func handleLoginError(w fyne.Window, btn *widget.Button, label *widget.Label, bar *widget.ProgressBarInfinite, msg string) {
	fyne.Do(func() {
		btn.Enable()
		label.SetText(msg)
		bar.Hide()
	})
}

