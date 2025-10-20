import type { Metadata } from "next";
import { Poppins } from "next/font/google";
import "./globals.css";
import { ClientLayout } from "@/components/ClientLayout";
import { Toaster } from "sonner";
import { AuthProvider } from "@/context/AuthContext";
import Script from "next/script";

const poppins = Poppins({
  subsets: ["latin"],
  variable: "--font-poppins",
  weight: ["400", "500", "600", "700"],
});

export const metadata: Metadata = {
  title: "Uptimatic - Monitor Website & Server",
  description:
    "Uptimatic membantu kamu memonitor uptime website dan server secara realtime.",
  keywords: [
    "uptime",
    "monitoring",
    "website",
    "server",
    "alerts",
    "uptimatic",
  ],
  authors: [{ name: "Uptimatic Team", url: "https://uptimatic.aeria.my.id" }],
  creator: "Uptimatic Team",
  publisher: "Uptimatic",
  openGraph: {
    title: "Uptimatic - Monitor Website & Server",
    description:
      "Uptimatic membantu kamu memonitor uptime website dan server secara realtime.",
    url: "https://uptimatic.aeria.my.id",
    siteName: "Uptimatic",
    images: [
      {
        url: "https://uptimatic.aeria.my.id/og-image.png",
        width: 1200,
        height: 630,
        alt: "Uptimatic - Monitor Website & Server",
      },
    ],
    locale: "id_ID",
    type: "website",
  },
  twitter: {
    card: "summary_large_image",
    title: "Uptimatic - Monitor Website & Server",
    description:
      "Uptimatic membantu kamu memonitor uptime website dan server secara realtime.",
    images: ["https://uptimatic.aeria.my.id/og-image.png"],
    creator: "@uptimatic",
  },
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  const jsonLd = {
    "@context": "https://schema.org",
    "@type": "WebSite",
    name: "Uptimatic",
    url: "https://uptimatic.aeria.my.id",
    author: {
      "@type": "Organization",
      name: "Uptimatic Team",
    },
    description:
      "Uptimatic membantu kamu memonitor uptime website dan server secara realtime.",
  };

  return (
    <html lang="id" className="scroll-smooth" data-scroll-behavior="smooth">
      <body className={`${poppins.variable} antialiased`}>
        <AuthProvider>
          <ClientLayout>{children}</ClientLayout>
        </AuthProvider>
        <Toaster />
        {/* JSON-LD structured data */}
        <Script
          id="json-ld"
          type="application/ld+json"
          dangerouslySetInnerHTML={{ __html: JSON.stringify(jsonLd) }}
        />
      </body>
    </html>
  );
}
