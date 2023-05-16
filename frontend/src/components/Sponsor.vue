<template>
<v-card class="sponsor background" width="170" flat :theme="layoutTheme" v-if="sponsor">
    <v-card-text class="text-body-2 font-weight-bold pa-2 text-center" v-if="sponsor === 'nirmata'">
        Hosted by <a href="https://nirmata.com" target="_blank"><img :src="`nirmata_${layoutTheme}.png`" height="14" style="margin-bottom: -3px;" /></a>
    </v-card-text>
    <v-card-text class="text-body-2 font-weight-bold pa-2 text-center" v-else>
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

<style scoped>
.background {
    background-color: rgb(var(--v-theme-background))
}

.sponsor {
  position: fixed;
  left: 50%;
  margin-left: -100px;
  bottom: 0px;
}
</style>