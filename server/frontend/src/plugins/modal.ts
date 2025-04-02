import Modal from '@/components/modal/modal.ts'
import type { App } from 'vue'

export default {
  install: (app: App) => {
    app.component('Modal', Modal)
  }
}
