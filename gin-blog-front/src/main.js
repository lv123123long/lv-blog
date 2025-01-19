// custom style
import './styles/index.css'
import './styles/common.css'
import './styles/animate.css'

// unocss
// unocss是一个高性能的原子css引擎
// Tailwind css 是一个实用工具优先的css引擎
import 'uno.css'
import '@unocss/reset/tailwind.css'

// vue
import { createApp } from 'vue'

// 导入路由和状态管理
import { router } from './router'
import { pinia } from './store'
// 导入根组件
import App from './App.vue'

// 创建并配置Vue应用
const app = createApp(App)
// 使用路由配置
app.use(router)
// 使用状态管理库
app.use(pinia)
// 把Vue应用挂载到HTML文档中的 #app 元素上
app.mount('#app')
