/** @type {import('next').NextConfig} */
const nextConfig = {
  output: "standalone",
  reactStrictMode: true,
  env: {
    SUPABASE_URL: process.env.SUPABASE_URL,
    SUPABASE_KEY: process.env.SUPABASE_KEY,
    API_SERVER: process.env.API_SERVER,
    WEB_SERVER: process.env.WEB_SERVER,
  },
};

module.exports = nextConfig;
