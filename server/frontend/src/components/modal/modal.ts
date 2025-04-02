import { defineComponent } from 'vue'

export default defineComponent({
  name: 'Modal',
  props: {
    title: {
      type: String,
      default: ''
    },
    show: {
      type: Boolean,
      default: false
    }
  },
  emits: ['update:show'],
  setup(_props, { emit }) {
    const closeModal = () => {
      emit('update:show', false)
    }

    return {
      closeModal
    }
  }
})
