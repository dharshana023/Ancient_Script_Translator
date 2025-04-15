const http = require('http');
const httpProxy = require('http-proxy');

// Create a proxy server
const proxy = httpProxy.createProxyServer({});

// Error handling
proxy.on('error', function(err, req, res) {
  console.error('Proxy error:', err);
  res.writeHead(500, {
    'Content-Type': 'text/plain'
  });
  res.end('Proxy error: ' + err.message);
});

// Create the server that will proxy requests
const server = http.createServer(function(req, res) {
  console.log('Received request:', req.method, req.url);
  
  // Proxy to the Go server running on localhost:5000
  proxy.web(req, res, { target: 'http://localhost:5000' });
});

// Listen on port 5000
console.log('Proxy server listening on port 5000');
server.listen(5000);