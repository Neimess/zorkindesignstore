const { createProxyMiddleware } = require('http-proxy-middleware');
console.log('✅ setupProxy.js загружен');

module.exports = function(app) {
  app.use(
    '/api',
    createProxyMiddleware({
      target: 'https://dev.api.inspireforge.ru',
      changeOrigin: true,
      secure: false,
    })
  );
};
