<template>
  <v-btn @click="download">{{ label }}</v-btn>
</template>

<script setup lang="ts">
const props = defineProps({
  label: { type: String, required: true },
  content: { type: String, required: true },
  filename: { type: String, required: true },
  contentType: { type: String, default: 'text/yaml' }
})

const download = () => {
  var a = window.document.createElement('a')
  a.href = window.URL.createObjectURL(new Blob([props.content], { type: 'application/yaml' }))
  a.download = props.filename

  document.body.appendChild(a)
  a.click()

  document.body.removeChild(a)
}
</script>
