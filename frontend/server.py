from http.server import SimpleHTTPRequestHandler, HTTPServer
import socketserver
import logging

def end_headers(self):
    # Code pour désactiver le cache
    self.send_header("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
    self.send_header("Pragma", "no-cache")
    self.send_header("Expires", "0")
    super().end_headers()


class CustomHandler(SimpleHTTPRequestHandler):
    def do_GET(self):
        try:
            super().do_GET()
        except BrokenPipeError:
            logging.warning("Client disconnected before the response was sent.")

if __name__ == "__main__":
    PORT = 9000
    with HTTPServer(("localhost", PORT), CustomHandler) as httpd:
        print(f"Serving on http://localhost:{PORT}")
        httpd.serve_forever()
