import { defineConfig, loadEnv } from 'vite'
import vue from '@vitejs/plugin-vue' // 导入 Vite 的 Vue 插件

import path from 'node:path'  // 导入node.js的path模块，用于处理文件路径
import unocss from 'unocss/vite'
import viteCompression from 'vite-plugin-compression'  // 导入 Vite 的压缩插件
import { visualizer } from 'rollup-plugin-visualizer'  // 导入 Rollup 的可视化插件

// https://vite.dev/config/
// export default defineConfig({
//   plugins: [vue()],
// })
// 定义vite配置
export default defineConfig((configEnv) => {
  const env = loadEnv(configEnv.mode, process.cwd())  // 加载环境变量

  return {
    base: env.VITE_PUBLIC_PATH || '/',   // 设置项目的基础路径，默认为 '/'
    resolve: {  // 配置路径别名
      alias: {
        '@': path.resolve(path.resolve(process.cwd()), 'src'),   // 将 '@' 别名指向 src 目录
        '~': path.resolve(process.cwd()), // 将 '~' 别名指向项目根目录
      },
    },
    plugins: [   // 配置插件
      vue(),     // 使用 Vite 的 Vue 插件
      unocss(),  // 使用 UnoCSS 的 Vite 插件
      viteCompression({ algorithm: 'gzip'}),   // 使用 Vite 的压缩插件，使用 gzip 算法
      visualizer({ open: false, gzipSize: true, brotliSize: true}),  // 使用 Rollup 的可视化插件，不自动打开，显示 gzip 和 brotli 压缩大小
    ],
    server: {     // 配置开发服务器
      host: '0.0.0.0',
      port: 3333,
      open: false,  // 启动服务器时是否自动打开浏览器
      proxy: {      // 配置代理
        '/api': {   // 代理路径为 /api 的请求
          target: env.VITE_BACKEND_URL,   // 代理目标地址，从环境变量中读取
          changeOrigin: true,             // 是否改变源地址
        },
      },
    },
    build: {  // 配置构建选项
      chunkSizeWarningLimit: 1024, // chunk大小警告的限制 单位 kb
    },
    esbuild: {  // 配置 ESBuild 选项
      drop: ['debugger'], // 移除代码中的 debugger 语句
    },
  }
})


// vite 构建工具的配置文件，配置了项目的基础路径，别名，插件，服务器设置，构建选项和ESBuild选项
// 环境变量文件中的变量名必须以 VITE_ 开头，这是 Vite 的约定。