<template>
<v-card class="sponsor" width="150" flat :theme="layoutTheme" style="background-color: rgb(var(--v-theme-background))" v-if="sponsor">
    <v-card-text class="text-body-2 font-weight-bold pa-2 text-center">
        <span v-html="sponsor"></span>
    </v-card-text>
</v-card>
</template>

<script setup lang="ts">
import { layoutTheme } from '@/config';
import { resolveAPI } from '@/utils';
import { ref } from 'vue';

const api = resolveAPI()
const sponsor = ref('')

fetch(`${api}/sponsor`, {
    method: "POST",
    mode: "cors",
    cache: "no-cache",
    headers: {
      "Content-Type": "text/plain",
    },
  }).then((resp) => {
    if (resp.status > 200) {
        return '';
    }

    return resp.text()
  }).then((text) => {
    sponsor.value = text
  }).catch(console.error)
</script>

<style>
.sponsor {
  position: fixed;
  left: 50%;
  margin-left: -100px;
  bottom: 0px;
}
</style>