/** @type {import('next').NextConfig} */
const nextConfig = {
  output: 'export',
  reactStrictMode: true,
  env: {
    API_SERVER_HOST: process.env.API_SERVER_HOST,
    API_SERVER_PORT: process.env.API_SERVER_PORT,
    WEB_SERVER_HOST: process.env.WEB_SERVER_HOST,
    WEB_SERVER_PORT: process.env.WEB_SERVER_PORT,
  },
};

module.exports = nextConfig;
