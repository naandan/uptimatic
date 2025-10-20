export default function PrivacyPolicy() {
  return (
    <div className="py-20 bg-white">
      <div className="max-w-4xl mx-auto px-6 text-slate-700">
        <h1 className="text-3xl font-semibold text-slate-800 mb-6 text-center">
          Kebijakan Privasi Uptimatic
        </h1>
        <p className="text-center text-slate-600 mb-12">
          Terakhir diperbarui: 17 Oktober 2025
        </p>

        <div className="space-y-10 text-base leading-relaxed">
          <section>
            <h2 className="text-xl font-semibold text-slate-800 mb-2">
              1. Pendahuluan
            </h2>
            <p>
              Uptimatic menghargai privasi Anda. Kebijakan privasi ini
              menjelaskan bagaimana kami mengumpulkan, menggunakan, dan
              melindungi data pribadi Anda saat menggunakan layanan pemantauan
              uptime kami.
            </p>
          </section>

          <section>
            <h2 className="text-xl font-semibold text-slate-800 mb-2">
              2. Data yang Kami Kumpulkan
            </h2>
            <p>Kami dapat mengumpulkan data berikut:</p>
            <ul className="list-disc list-inside mt-2 space-y-1">
              <li>Alamat email untuk tujuan notifikasi dan autentikasi.</li>
              <li>Data situs atau endpoint API yang Anda pantau.</li>
              <li>Log performa dan uptime untuk analisis layanan.</li>
            </ul>
          </section>

          <section>
            <h2 className="text-xl font-semibold text-slate-800 mb-2">
              3. Cara Kami Menggunakan Data
            </h2>
            <p>Data yang dikumpulkan digunakan untuk:</p>
            <ul className="list-disc list-inside mt-2 space-y-1">
              <li>Memberikan laporan performa dan uptime secara real-time.</li>
              <li>
                Mengirim notifikasi downtime melalui email atau integrasi pihak
                ketiga.
              </li>
              <li>
                Meningkatkan kualitas layanan Uptimatic melalui analisis teknis.
              </li>
            </ul>
          </section>

          <section>
            <h2 className="text-xl font-semibold text-slate-800 mb-2">
              4. Penyimpanan dan Keamanan Data
            </h2>
            <p>
              Kami menggunakan infrastruktur global yang aman dengan enkripsi
              data dan sistem backup berkala untuk memastikan integritas
              informasi Anda. Data Anda tidak akan dijual atau dibagikan kepada
              pihak ketiga tanpa izin.
            </p>
          </section>

          <section>
            <h2 className="text-xl font-semibold text-slate-800 mb-2">
              5. Integrasi Pihak Ketiga
            </h2>
            <p>
              Uptimatic dapat terhubung dengan layanan pihak ketiga seperti
              Slack atau sistem notifikasi lain. Kami hanya membagikan data yang
              diperlukan untuk fungsi tersebut, dan pengguna memiliki kendali
              penuh atas integrasinya.
            </p>
          </section>

          <section>
            <h2 className="text-xl font-semibold text-slate-800 mb-2">
              6. Hak Pengguna
            </h2>
            <p>
              Anda berhak untuk mengakses, memperbarui, atau menghapus data
              pribadi Anda kapan saja melalui dasbor Uptimatic atau dengan
              menghubungi tim kami.
            </p>
          </section>

          <section>
            <h2 className="text-xl font-semibold text-slate-800 mb-2">
              7. Perubahan Kebijakan
            </h2>
            <p>
              Kami dapat memperbarui kebijakan privasi ini dari waktu ke waktu.
              Pembaruan akan diumumkan melalui situs resmi Uptimatic dan mulai
              berlaku segera setelah dipublikasikan.
            </p>
          </section>

          <section>
            <h2 className="text-xl font-semibold text-slate-800 mb-2">
              8. Kontak Kami
            </h2>
            <p>
              Jika Anda memiliki pertanyaan atau permintaan terkait privasi,
              silakan hubungi kami di:{" "}
              <a
                href="mailto:nandanramdani608@gmail.com"
                className="text-primary hover:underline"
              >
                nandanramdani608@gmail.com
              </a>
              .
            </p>
          </section>
        </div>
      </div>
    </div>
  );
}
