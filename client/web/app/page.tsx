import { CTASection } from "@/components/landing/CTASection";
import { FeaturesSection } from "@/components/landing/FeatureSection";
import { HeroSection } from "@/components/landing/HeroSection";

export default function Home() {
  return (
    <main className="relative">
      <HeroSection />
      <FeaturesSection />
      <CTASection />
    </main>
  );
}
