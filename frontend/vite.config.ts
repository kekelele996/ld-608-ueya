import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";

// 本地开发：页面与统一 /api 请求都走 20108，由 Vite 代理到本地后端 21108。
// Docker 部署时改由 nginx.conf 的 location /api/ 反代到 http://backend:3000/。
export default defineConfig({
  plugins: [react()],
  server: {
    port: 20108,
    host: "0.0.0.0",
    proxy: {
      "/api": {
        target: "http://127.0.0.1:21108",
        changeOrigin: true
      }
    }
  },
  preview: {
    port: 20108,
    host: "0.0.0.0"
  }
});
