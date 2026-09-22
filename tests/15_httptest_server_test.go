package tests

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Go 1.27: httptest.NewTestServer (Bellek İçi / In-Memory Test Sunucusu)
//
// ─── ESKİ httptest.NewServer vs YENİ httptest.NewTestServer FARKI ───
// Soru: Tek fark TCP portu açmaması mı? Kod kullanımı aynı gibi?
// Cevap: Evet, çağırma API'si alışık olduğumuz `srv.Client()` stilini korur,
// ancak altyapıda 3 kritik iyileştirme vardır:
//
// 1. GERÇEK TCP PORTU AÇILMAZ (In-Memory Pipe):
//    - Eskisi OS üzerinde rastgele bir loopback TCP portu (örn: 127.0.0.1:54321) açardı.
//    - Yüzlerce paralel test koşan CI sunucularında port tükenmesi (port exhaustion)
//      ve işletim sistemi güvenlik duvarı (firewall) uyarıları çıkardı.
//    - Yeni `NewTestServer` bellek içi sanal ağ kullanır, sıfır port tüketir ve çok daha hızlıdır!
//
// 2. OTOMATİK TEMİZLEME (t.Cleanup):
//    - Eskiden `defer srv.Close()` yazmayı unutursanız port asılı kalır ve kaynak sızardı.
//    - Yeni fonksiyon ilk parametre olarak `t` (*testing.T) alır ve test bitince
//      sunucuyu kendiliğinden kapatır (`defer` yazmanıza gerek kalmaz).
//
// 3. SYNCTEST UYUMLULUĞU:
//    - Bellek içi çalıştığı için `testing/synctest` sanal zaman balonunun içine
//      sokulabilir; zaman aşımı içeren HTTP testlerini 0 milisaniyede test edebilirsiniz!
//
// Çalıştırma:
//   go test -v -race ./tests/ -run TestNewTestServer

func TestNewTestServer(t *testing.T) {
	// Basit bir HTTP handler
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(w, `{"status":"ok","source":"in-memory-network"}`)
	})

	// ─── 1. Bellek İçi Sunucuyu Başlatma ───
	// İlk parametre testing.TB (t) alır, böylece otomatik cleanup kaydeder.
	// defer srv.Close() yazmamıza gerek YOKTUR!
	srv := httptest.NewTestServer(t, handler)

	// ─── 2. İstemci ile İstek Atma ───
	// srv.Client() bellek içi ağa bağlı özel bir http.Client döndürür.
	// İstek adresi olarak srv.URL veya doğrudan "http://example.com/" kullanılabilir;
	// gerçek DNS araması yapılmaz, doğrudan bellek içi handler'a bağlanır.
	resp, err := srv.Client().Get(srv.URL)
	if err != nil {
		t.Fatalf("İstek başarısız: %v", err)
	}
	defer resp.Body.Close()

	// ─── 3. Yanıtı Doğrulama ───
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Yanıt gövdesi okunamadı: %v", err)
	}

	fmt.Printf("Sunucu Adresi (Sanal): %s\n", srv.URL)
	fmt.Printf("Gelen Yanıt          : %s\n", string(bodyBytes))

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Beklenmeyen status: %d", resp.StatusCode)
	}
}
