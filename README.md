# DarkMagic

DarkMagic, Tor ağı üzerinden verilen URL listesini kontrol eden ve erişilebilen sitelerin tam sayfa ekran görüntüsünü alan bir Go (Golang) uygulamasıdır.

Program çalışmaya başlamadan önce Tor bağlantısının aktif olup olmadığını kontrol eder. Tor aktif değilse program çalışmaz.

## Özellikler

- Tor SOCKS5 proxy kullanımı (127.0.0.1:9050)
- Tor çıkış IP kontrolü
- Dosyadan URL okuma
- HTTP durum kodu kontrolü
- Chromedp ile full page screenshot
- Loglama (darkmagic.log)
- Headless Chrome desteği

## Gereksinimler

- Linux
- Go 1.20 veya üzeri
- Tor servisi
- Google Chrome / Chromium

## Kurulum

Tor servisini başlat:

sudo systemctl start tor
sudo systemctl enable tor

Programı çalıştır:

go run main.go

## Kullanım

Hedef URL listesini içeren bir dosya oluştur:

https://example.com
https://github.com
https://torproject.org

Program çalıştıktan sonra hedef dosya adını gir:

targets.yaml

## Çıktılar

- Ekran görüntüleri hostname.png formatında kaydedilir
- Tüm işlemler darkmagic.log dosyasına yazılır

## Notlar

- 502 Bad Gateway dönen sitelerde ekran görüntüsü alınmaz
- Tor ağı nedeniyle sayfa yüklenmesi yavaş olabilir
- Bazı siteler Tor çıkış IP’lerini engelleyebilir

## Uyarı

Bu araç yalnızca izinli ve yasal hedeflerde kullanılmalıdır. Yetkisiz kullanım yasal sorumluluk doğurabilir.
