/** @type {import('next').NextConfig} */
const nextConfig = {
  reactStrictMode: true,
  env: {
    API_SERVER_HOST: process.env.API_SERVER_HOST,
    API_SERVER_PORT: process.env.API_SERVER_PORT,
  }
}

module.exports = nextConfig
