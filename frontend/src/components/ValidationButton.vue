<template>
  <v-btn
    size="large"
    :prepend-icon="icon"
    :color="color"
    :loading="loading"
    class="play"
    rounded
    @click="submit"
    >Start Validation</v-btn
  >
</template>

<script setup lang="ts">
import { ref, watch } from "vue";

const props = defineProps({
  policy: { type: String, default: "" },
  resource: { type: String, default: "" },
  context: { type: String, default: "" },
  errorState: { type: Boolean, default: false },
});

const emit = defineEmits(["on-response", "on-error"]);

const loading = ref<boolean>(false);
const color = ref<string | undefined>("primary");
const icon = ref<string | undefined>("mdi-play");

watch(props, ({ errorState }: { errorState: boolean }) => {
  if (errorState) {
    color.value = "error";
    icon.value = "mdi-alert-circle-outline";
    return;
  }

  color.value = "primary";
  icon.value = "mdi-play";
});

const api: string = import.meta.env.VITE_API_HOST || "";

const submit = (): void => {
  if (!props.policy.trim()) {
    emit("on-error", new Error("Policy is required"));
    return;
  }

  if (!props.resource.trim()) {
    emit("on-error", new Error("Resource is required"));
    return;
  }

  loading.value = true;

  fetch(`${api}/engine`, {
    body: JSON.stringify({
      policy: props.policy,
      resources: props.resource,
      context: JSON.parse(props.context),
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
