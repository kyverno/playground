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
import { MarkerSeverity, editor } from "monaco-editor";
import { EngineResponse } from "@/types";

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

const handleEditorErrors = () => {
  const markers = editor
    .getModelMarkers({ owner: "yaml" })
    .filter((m) => m.severity == MarkerSeverity.Error)
    .map(
      (m) => `${m.resource.path} L${m.startColumn}:${m.startLineNumber}: ${m.message}`
    );

  return markers;
};

const submit = (): void => {
  const errors = handleEditorErrors();
  if (errors.length) {
    // emit("on-error", new Error(`<b>YAML validation failed, check the errors below.</b><br />${errors.join('<br />')}`));
    emit("on-error", new Error(`YAML validation failed, please check your manifests.`));
    return;
  }

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
      context: props.context,
    }),
    method: "POST",
    mode: "cors",
    cache: "no-cache",
    headers: {
      "Content-Type": "application/json",
    },
  })
    .then((resp) => {
      if (resp.status > 300) {
        resp.text().then((err) => emit("on-error", new Error(`ServerError: ${err}`)));
        return;
      }

      return resp
        .json()
        .catch(() => ({}))
        .then((content: EngineResponse) => {
          emit("on-response", content);
        });
    })
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
