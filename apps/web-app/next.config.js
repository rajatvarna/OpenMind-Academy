/** @type {import('next').NextConfig} */
const nextConfig = {
  reactStrictMode: true,
  // This enables the standalone output mode, which is optimized for Docker.
  // It creates a '.next/standalone' folder with a minimal server and dependencies.
  output: 'standalone',
};

module.exports = nextConfig;
