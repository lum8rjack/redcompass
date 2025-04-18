import { onMounted, onUnmounted } from 'vue'

export function onClickOutside(element, callback) {
  function handleClick(event) {
    if (element.value && !element.value.contains(event.target)) {
      callback()
    }
  }

  onMounted(() => {
    document.addEventListener('click', handleClick)
  })

  onUnmounted(() => {
    document.removeEventListener('click', handleClick)
  })
} 