<template>
  <v-dialog v-model="dialog" width="600px">
    <template v-slot:activator="{ props }">
      <v-btn v-bind="props" prepend-icon="mdi-web">from URL</v-btn>
    </template>

    <v-card theme="light">
      <v-card-text>
        <v-text-field label="URL" v-model="url" />
      </v-card-text>
      <v-card-actions>
        <v-btn @click="dialog = false">Close</v-btn>
        <v-spacer />
        <v-btn @click="onLoad" color="primary" :loading="loading">Load Content</v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>
<script setup lang="ts">
import { ref } from "vue";

const emit = defineEmits(["click"]);

const dialog = ref<boolean>(false);
const loading = ref<boolean>(false);
const url = ref<string>("");

const onLoad = () => {
    if (!url.value) return;

    loading.value = true

    fetch(url.value).then(resp => resp.text()).then(body => {
        emit('click', body)
        url.value = ''
        dialog.value = false
    }).catch(err => console.error(err)).finally(() => {
        loading.value = false
    })
}
</script>
