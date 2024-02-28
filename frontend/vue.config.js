const { defineConfig } = require('@vue/cli-service');

const isDevMode = process.env.NODE_ENV === 'development';

const devServer = isDevMode
  ? {
      proxy: {
        '/api': {
          target: 'http://localhost:3000',
          changeOrigin: true,
          ws: true,
          // pathRewrite: { '^/api': '' },
        },
      },
    }
  : {};

module.exports = defineConfig({
  transpileDependencies: true,
  devServer,
});
