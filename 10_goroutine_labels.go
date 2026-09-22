package main

import (
	"context"
	"fmt"
	"runtime"
	"runtime/pprof"
	"strings"
)

// Go 1.27: Traceback Başlıklarında Goroutine Labels (Etiketler)
//
// ─── pprof.Do NEDİR VE NE İŞE YARAR? (USE CASE) ───
// Soru: `pprof.Do` yeni bir goroutine mi başlatır?
// Cevap: HAYIR! `pprof.Do` yeni goroutine AÇMAZ.
//
// Amacı: Mevcut goroutine'e ve context'e anahtar-değer etiketleri (labels) iliştirir
// ve verilen fonksiyonu BU ETİKETLERLE çalıştırır.
//
// Gerçek Hayat Senaryosu (Use Case):
// Web sunucunuzda (HTTP middleware) veya Kafka consumer'ınızda gelen her isteğe ait
// `request_id`, `user_id` veya `tenant_id` bilgisini `pprof.Labels` ile tanımlarsınız.
//
// ─── GO 1.27 İLE NE DEĞİŞTİ? ───
// Go 1.26'ya kadar bu etiketler SADECE pprof CPU profil dosyalarında (.pb.gz) görünürdü.
// Sistemde beklenmedik bir `panic` patladığında veya loglara `runtime.Stack` basıldığında,
// o stack trace'in HANGİ İSTEĞE (hangi request ID'ye) ait olduğunu görmek imkansızdı!
//
// Go 1.27 ile birlikte (`go.mod` içinde `go 1.27` varsa):
// Traceback başlık satırına bu etiketler OTOMATİK eklenir:
//   `goroutine 1 [running] {request: 33}:`
// Böylece loglarda crash olan goroutine'in hangi istek yüzünden çöktüğü anında anlaşılır!
//
// ─── GÜVENLİK / OPT-OUT ───
// Etiketlerde hassas bilgi (token vb.) loglanmasın istenirse:
// `GODEBUG=tracebacklabels=0` ile kapatılabilir.

func main() {
	ctx := context.Background()

	// 1. İstek etiketlerini tanımlıyoruz (örn: HTTP middleware içinde)
	labels := pprof.Labels("request", "33")

	// 2. pprof.Do: Bu etiketleri mevcut goroutine'e bağlayıp fonksiyonu çalıştırır
	pprof.Do(ctx, labels, func(context.Context) {
		// Mevcut goroutine'in stack dökümünü alıyoruz
		buf := make([]byte, 1024)
		n := runtime.Stack(buf, false)
		stackTrace := string(buf[:n])

		fmt.Println("=== Traceback Çıktısı ===")
		fmt.Print(stackTrace)

		// Başlıkta etiketin bulunup bulunmadığını doğruluyoruz
		hasLabel := strings.Contains(stackTrace, "{request: 33}")
		fmt.Println("\nBaşlıkta {request: 33} etiketi var mı? ->", hasLabel)
	})
}
