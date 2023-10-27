/** @type {import('next').NextConfig} */
const nextConfig = {
  output: "standalone",
  reactStrictMode: true,
  env: {
    API_SERVER: process.env.API_SERVER,
    WEB_SERVER: process.env.WEB_SERVER,
  },
};

module.exports = nextConfig;
