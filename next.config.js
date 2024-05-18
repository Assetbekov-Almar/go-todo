/** @type {import('next').NextConfig} */
const nextConfig = {
    output: 'export',
    basePath: '/go-todo',
    images: {
        loader: 'custom',
        loaderFile: './custom-image-loader.js',
    },
};

module.exports = nextConfig;
