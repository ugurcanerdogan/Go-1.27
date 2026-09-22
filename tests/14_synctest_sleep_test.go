package tests

import (
	"fmt"
	"testing"
	"testing/synctest"
	"time"
)

// Go 1.27: testing/synctest.Sleep (Sentetik Zaman ve Deterministik Test)
//
// ─── NE DEĞİŞTİ? ───
// Go 1.25'te gelen `testing/synctest` paketi, eşzamanlı testleri sanal bir
// "zaman balonu" (bubble) içinde koşturur. Gerçek hayatta 0 milisaniye geçer!
//
// Go 1.27 ile `synctest.Sleep(d)` eklendi:
//     synctest.Sleep(d) == time.Sleep(d) + synctest.Wait()
//
// ─── ADIM ADIM TİMELİNE (ZAMAN ÇİZELGESİ) ───
//  T = 0s  : Worker başlatılır, sanal olarak 3 saniye uykuya geçer (`time.Sleep(3s)`).
//            Worker henüz işini bitirmedi, blokeli bekliyor.
//  T = 0s  : Test goroutine'i `synctest.Sleep(5s)` çağırır.
//  T = 3s  : Sanal saat 3s olduğunda worker uyanır, kanalını kapatır ve sonlanır.
//  T = 5s  : Sanal saat 5s'ye ulaşır.
//            BURADA ÇOK ÖNEMLİ KAVRAM: "DURULMA" (QUIESCENCE / SETTLE)
//            `synctest.Wait()` devreye girer. Bubble içindeki TÜM goroutine'lerin
//            ya tamamen sonlanmasını ya da blokeli bir duruma geçmesini bekler.
//            Eğer arkada hala aktif çalışan/hesaplama yapan goroutine kalsaydı test durulamazdı!
//  Sonuç   : Worker tamamen bitti, yarış durumu (race condition / flaky test) kalmadı.
//            Ve tüm bunlar gerçek dünyada sadece ~0ms sürdü!
//
// Çalıştırma:
//   go test -v -race ./tests/ -run TestSynctestSleep

func TestSynctestSleep(t *testing.T) {
	// synctest.Test: İzole bir sanal zaman balonu başlatır
	synctest.Test(t, func(t *testing.T) {
		virtualStart := time.Now()
		workerCompleted := make(chan struct{})

		// Worker 3 saniye sanal uykuya dalıyor:
		go func() {
			time.Sleep(3 * time.Second) // Sanal saatte 3 saniye bekler

			// Worker uyandığında sanal süreyi kontrol ediyoruz
			if elapsed := time.Since(virtualStart); elapsed != 3*time.Second {
				t.Errorf("Worker erken/geç uyandı: %s, beklenen: 3s", elapsed)
			}
			close(workerCompleted)
		}()

		// Test goroutine'i saati 5 saniye ileri sarıyor:
		// 3. saniyede worker uyanıp işini bitirecek, 5. saniyede test tam durulacak.
		synctest.Sleep(5 * time.Second)

		// Sanal saat tam 5 saniye ilerlemiş olmalı (gerçekte 0ms geçti):
		if totalElapsed := time.Since(virtualStart); totalElapsed != 5*time.Second {
			t.Errorf("Sanal saat ilerlemedi: %s, beklenen: 5s", totalElapsed)
		}

		// synctest.Sleep sayesinde worker'ın kesinlikle tamamlandığından %100 eminiz:
		select {
		case <-workerCompleted:
			// Başarılı: Worker tamamlandı ve duruldu
		default:
			t.Fatal("Hata: synctest.Sleep sonrasında worker hala durulmadı/çalışıyor!")
		}

		fmt.Println("Worker 3s'de bitti, ana test 5s sanal saatte duruldu (gerçek süre: ~0ms).")
	})
}
