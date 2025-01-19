<script>
// 导入reactive 函数，创建响应式对象
import { reactive } from 'vue'

const props = defineProps({
    show: { type: Boolean, default: false },   // 是否显示通知，默认为false
    position: {            // 通知的位置，可以是top，bottom 默认是top
        type: String,
        default: 'top',
        validator: value => ['top', 'bottom'].includes(value),
    },
    align: {               // 通知的对齐方式
        type: String,
        default: 'center',
        validator: value => ['left', 'center', 'right'].includes(value)
    },
    timeout: {             // 通知显示时间，默认2500毫秒
        type: Number,
        default: 2500,
    },
    queue: {               // 是否启用通知队列
        type: Boolean,
        default: true,
    },
    zIndex: { type: Number, default: 100 },     // 通知的 z-index，默认为 100。
    closeable: { type: Boolean, default: false },  // 是否允许关闭通知
})

// 定义一个响应式对象flux
const flux = reactive({
    events: [],      // 存储通知事件的数组
    // 分别添加不同类型的通知
    success: content => flux.add('success', content),
    info: content => flux.add('info', content),
    warning: content => flux.add('warning', content),
    error: content => flux.add('error', content),

    show: (type, content) => flux.add(type, content),  // 添加通知，可以是任意类型
    add: (type, content) => {                  // 添加通知到 events 数组，并在timeout后移除
        if (!props.queue)
            flux.events = []
        setTimeout(() => {
            const event = { id: Date.now(), content, type }
            flux.events.push(event)
            setTimeout(() => flux.remove(event), props.timeout)
        }, 100)
    },
    remove: (event) => {            // 从events数组中移除指定通知
        flux.events = flux.events.filter(e => e.id !== event.id)
    },
})

// 暴露flux对象的方法，使得父组件可以调用这些方法
defineExpose({
    show: flux.show,
    success: flux.success,
    info: flux.info,
    warning: flux.warning,
    error: flux.error,
})
</script>

<template>
    <Teleport to="body">
        <div
            v-show="flux.events.length"
            class="pointer-events-none fixed w-4/5 sm:w-[400px]"
            :class="{
                'left-1/2 -translate-x-1/2': align === 'center',
                'left-16': align === 'left',
                'right-16': align === 'right',
                'top-4': position === 'bottom',
                'bottom-6': position === 'bottom',
            }"
            :style="{ zIndex }"
        >
        <TransitionGroup
            tag="ul"
            enter-active-class="transition ease-out duration-200"
            leave-active-class="transition ease-in duration-200 absolute w-full"
            :enter-class="position === 'bottom'
          ? 'transform translate-y-3 opacity-0'
          : 'transform -translate-y-3 opacity-0'"
        enter-to-class="transform translate-y-0 opacity-100"
        leave-class="transform translate-y-0 opacity-100"
        :leave-to-class="position === 'bottom'
          ? 'transform translate-y-1/4 opacity-0'
          : 'transform -translate-y-1/4 opacity-0'"
        move-class="ease-in-out duration-200"
        class="inline-block w-full"
      >
        <li
          v-for="event in flux.events" :key="event.id"
          :class="{
            'pb-2': position === 'bottom',
            'pt-2': position === 'top',
          }"
        >
          <slot :type="event.type" :content="event.content">
            <div class="pointer-events-auto w-full overflow-hidden rounded-lg bg-white ring-1 ring-black ring-opacity-5">
              <div class="flex justify-between px-4 py-3">
                <div class="flex items-center">
                  <div
                    class="mr-6 h-6 w-6"
                    :class="{
                      'i-mdi:check-circle text-green': event.type === 'success',
                      'i-mdi:information-outline text-blue': event.type === 'info',
                      'i-mdi:alert-outline text-yellow': event.type === 'warning',
                      'i-mdi:alert-circle-outline text-red': event.type === 'error',
                    }"
                  />
                  <div class="ml-1">
                    <slot>
                      <div> {{ event.content }} </div>
                    </slot>
                  </div>
                </div>
                <button
                  v-if="closeable"
                  class="i-mdi:close h-5 w-5 flex items-center justify-center rounded-full rounded-full p-1 font-bold text-gray-400 hover:text-gray-600"
                  @click="flux.remove(event)"
                />
              </div>
            </div>
          </slot>
        </li>
      </TransitionGroup>
    </div>
  </Teleport>
</template>
