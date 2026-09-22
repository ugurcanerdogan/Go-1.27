package main

import (
	"fmt"
	"os"
	"runtime"
	"runtime/pprof"
)

// Go 1.27: goroutineleak Profili Artık Genel Kullanımda (Kararlı!)
//
// ─── pprof.Lookup VE WriteTo ARKA PLANDA NASIL ÇALIŞIR? ───
// 1. `pprof.Lookup("goroutineleak")`:
//    Go çalışma zamanında (runtime) kayıtlı hazır profil toplayıcıları sorgular
//    (örneğin "heap", "goroutine", "threadcreate" gibi). Go 1.27 ile "goroutineleak"
//    resmi ve standart bir profil toplayıcı olarak runtime'a eklendi.
//
// 2. GC Taraması ve Sızıntının Profile Düşmesi:
//    Go'nun Garbage Collector'ı (GC) belleği tararken nesnelerin birbirine olan
//    erişim bağlarını (reachability graph) inceler:
//    - Eğer bir goroutine bir kanal (channel) veya mutex üzerinde kilitlenip UYUYORSA,
//    - VE şu anda aktif çalışan hiçbir goroutine bu kanal/mutex adresine ULAŞAMIYORSA,
//    GC bu goroutine'in bir daha asla uyanamayacağını kesin olarak anlar!
//    Bu goroutine'ler runtime tarafından doğrudan "goroutineleak" profiline eklenir.
//
// 3. `prof.WriteTo(os.Stdout, 1)`:
//    Toplanan bu sızıntı profilini insan tarafından okunabilir metin formatında
//    verilen io.Writer'a (burada os.Stdout) yazar.
//
// Web servislerinde izlemek için: http://localhost:6060/debug/pprof/goroutineleak

// Sızıntı üreten fonksiyon:
// Kanal yerel değişkendir, dışarıya referansı verilmez.
// İçine değer yazılmaya çalışılır fakat alıcısı olmadığı için goroutine sonsuza dek kilitlenir.
func leakyGoroutine() {
	ch := make(chan int) // Unbuffered kanal
	ch <- 100            // Alıcı yok, goroutine burada sonsuza dek kilitlenir!
}

func main() {
	// 1. 'goroutineleak' profili runtime'dan alınır
	prof := pprof.Lookup("goroutineleak")
	if prof == nil {
		fmt.Println("Hata: goroutineleak profili bulunamadı!")
		return
	}

	fmt.Println("=== 1. Sızıntı Oluşturuluyor ===")
	go leakyGoroutine() // 1 adet sızan goroutine başlatıyoruz

	// Scheduler'ın goroutine'i çalıştırmasına izin veriyoruz
	runtime.Gosched()

	// GC çalıştırılır: GC erişilemeyen primitifleri ve kilitli goroutine'leri tarayıp profile kaydeder
	runtime.GC()

	// ─── 2. Profil Çıktısını Ekrana Yazma ───
	fmt.Println("\n=== 2. Goroutine Leak Profil Çıktısı ===")
	// Çıktıda 'total 1' ve leakyGoroutine satırını göreceğiz
	if err := prof.WriteTo(os.Stdout, 1); err != nil {
		fmt.Println("Yazma hatası:", err)
	}

	fmt.Println("\n-> Görüldüğü gibi runtime, takılmış ve uyanamayacak goroutine'i tespit etti.")
}
