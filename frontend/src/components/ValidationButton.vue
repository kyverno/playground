<template>
  <v-btn
    size="large"
    prepend-icon="mdi-play"
    color="primary"
    :loading="loading"
    class="play"
    rounded
    @click="submit"
    >Start Validation</v-btn
  >
</template>

<script setup lang="ts">
import { ref } from "vue";

const props = defineProps({
  policy: { type: String, default: "" },
  resource: { type: String, default: "" },
  context: { type: String, default: "" },
});

const emit = defineEmits(["on-response", "on-error"]);

const loading = ref<boolean>(false);

const api: string = import.meta.env.VITE_API_HOST || "";

const submit = (): Promise<any> => {
  loading.value = true;

  return fetch(`${api}/engine`, {
    body: JSON.stringify({
      policy: props.policy,
      resources: props.resource,
      context: props.context,
    }),
    method: "POST",
    mode: "cors",
    cache: "no-cache",
    headers: {
      "Content-Type": "application/json",
    },
  })
    .then((resp) => resp.json().catch(() => ({})))
    .then((content) => emit("on-response", content))
    .catch((err) => emit("on-error", err))
    .finally(() => (loading.value = false));
};
</script>

<style scoped>
.play {
  position: fixed;
  bottom: 45px;
  right: 40px;
}
</style>
