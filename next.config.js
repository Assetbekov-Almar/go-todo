/** @type {import('next').NextConfig} */
const nextConfig = {
    output: 'export',
    images: {
        remotePatterns: [
            {
                protocol: 'https',
                hostname: 'todo-go.fly.dev/',
                port: '8080',
                pathname: '/static/**',
            },
        ],
    },
};

module.exports = nextConfig;
