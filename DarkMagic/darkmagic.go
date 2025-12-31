package main

import (
	"bufio"
	"context"
	"encoding/json"
	f "fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	ch "github.com/chromedp/chromedp"
)

func initLogger() *log.Logger {
	file, err := os.OpenFile(
		"darkmagic.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0644,
	)
	if err != nil {
		f.Print("[X]	LOG DOSYASI ACILAMADI!")
		log.Fatal("LOG DOSYASI ACILAMADI:", err)
	}

	return log.New(file, "[APP] ", log.Ldate|log.Ltime)
}
func main() {
	var target string
	var Ss string
	var screenshot []byte
	type TorResp struct {
		IsTor bool `json:"IsTor"`
	}
	var tr TorResp

	logger := initLogger()
	proxyURL, err := url.Parse("socks5://127.0.0.1:9050")
	httpTransport := &http.Transport{Proxy: http.ProxyURL(proxyURL)}
	httpClient := &http.Client{Transport: httpTransport, Timeout: 30 * time.Second}
	allocOpts := append(
		ch.DefaultExecAllocatorOptions[:],
		ch.ProxyServer("socks5://127.0.0.1:9050"),
		ch.Flag("headless", true),
		ch.NoSandbox,
	)

	allocCtx, allocCancel := ch.NewExecAllocator(context.Background(), allocOpts...)
	defer allocCancel()
	logger.Print("----------------------------------------PROGRAM STARTING----------------------------------------\n")
	ctx, cancel := ch.NewContext(allocCtx)
	defer cancel()
	f.Print("[-]	GIZLILIK KONTROLU YAPILIYOR LUTFEN BEKLEYINIZ...\n")

	tor, err := httpClient.Get("https://check.torproject.org/api/ip")
	if err != nil {
		f.Print("[!]	TOR CHECK API ERISILEMEDI\n")
		f.Print("[!]	KONTROLLERI SAGLAYIP TEKRAR DENEYINIZ!\n")
		logger.Printf("TOR CHECK API ERISILEMEDI: %v", err)
		return
	}

	defer tor.Body.Close()
	torCheck, err := io.ReadAll(tor.Body)

	if err != nil {
		f.Print("[!]	TOR BODY OKUNAMADI\n")
		logger.Printf("TOR BODY OKUNAMADI: %v", err)
		return
	}

	if err := json.Unmarshal(torCheck, &tr); err != nil {
		f.Print("[!]	TOR JSON PARSE HATASI\n")
		logger.Printf("TOR JSON PARSE HATASI: %v", err)
		return
	}

	if tr.IsTor {
		f.Printf("[✓]	GIZLILIK KONTROLU BASARILI | ISTOR: %v\n", tr.IsTor)
	} else {
		f.Printf("[X]	GIZLILIK KONTROLU BASARISIZ | ISTOR: %v\n", tr.IsTor)
		return
	}

	f.Print("\n[-]	HEDEF LISTENIZI GIRINIZ => ")
	f.Scan(&target)
	f.Print("\n")
	file, err := os.OpenFile(target, os.O_RDONLY, 0755)
	if err != nil {
		f.Print("\n[!]	HEDEF DOSYA ACILAMADI!\n\n")
		logger.Printf("HEDEF DOSYA ACILAMADI! %s", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		satir := scanner.Text()

		u, err := url.Parse(satir)
		if err != nil || u.Hostname() == "" {
			logger.Printf("URL PARSE HATASI: %s", satir)
			continue
		}
		safe := strings.ReplaceAll(u.Hostname(), ".", "_")

		req, err := http.NewRequest("GET", satir, nil)
		if err != nil {
			continue
		}
		req.Header.Set("User-Agent", "Mozilla/5.0 (Linux; Android 10; K) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/143.0.0.0 Mobile Safari/537.36")
		response, httpError := httpClient.Do(req)

		if httpError != nil {
			logger.Printf(
				"GEÇERSİZ DOMAIN / HTTP HATASI | URL: %s",
				satir,
			)
			f.Printf("[X]	BAGLANTI BASARISIZ => %s\n\n", satir)
			continue
		} else {

			logger.Printf(
				"DOMAIN: %s => %d", satir,
				response.StatusCode,
			)
		}

		response.Body.Close()

		f.Print("[✓]	DOMAIN: ", satir, " => ", response.StatusCode)
		if response.StatusCode == 502 {
			logger.Printf("BADE GATEWAY %v", response.StatusCode)
			f.Printf("\n[!]	BADE GATEWAY | EKRAN GORUNTUSU ALINAMAZ (!)\n\n")
		} else {
			if err := ch.Run(ctx, ch.Navigate(satir),
				ch.Sleep(45*time.Second), ch.EmulateViewport(1920, 1080),
				ch.WaitReady("body", ch.ByQuery),
				ch.FullScreenshot(&screenshot, 90)); err != nil {
				logger.Printf("EKRAN GORUNTUSU ALINAMADI : %v", err)
				f.Printf("\n[!]	EKRAN GORUNTUSU ALINAMADI : \n\n")
				continue
			} else {
				Ss = safe + ".png"
				os.WriteFile(Ss, screenshot, 0644)
				f.Printf("\n[✓]	EKRAN GORUNTUSU ALINDI. | DOSYA ADI: %s\n\n", Ss)
				logger.Print("EKRAN GORUNTUSU ALINDI. | DOSYA ADI: ", Ss)

			}

			if err := scanner.Err(); err != nil {
				logger.Printf("DOSYA OKUMA HATASI: %v", err)
				f.Print("[!]	DOSYA OKUMA HATASI!")
			}
		}
	}

	f.Print("[✓]	LISTE TAMAMLANDI PROGRAM KAPATILIYOR...")
	logger.Printf("----------------------------------------PROGRAM ENDING----------------------------------------\n\n")
}
