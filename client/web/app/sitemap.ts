import { MetadataRoute } from "next";

export default function sitemap(): MetadataRoute.Sitemap {
  return [
    {
      url: "https://uptimatic.aeria.my.id/",
      lastModified: new Date(),
    },
    {
      url: "https://uptimatic.aeria.my.id/auth/login",
      lastModified: new Date(),
    },
    {
      url: "https://uptimatic.aeria.my.id/auth/register",
      lastModified: new Date(),
    },
    {
      url: "https://uptimatic.aeria.my.id/uptime",
      lastModified: new Date(),
    },
  ];
}
