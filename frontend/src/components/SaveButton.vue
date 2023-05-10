<template>
  <v-tooltip
    :model-value="persisted"
    location="bottom"
    text="Persisted current state locally"
    :open-on-hover="false"
  >
    <template v-slot:activator="{ props }">
      <v-btn
        prepend-icon="mdi-content-save"
        @click="persist"
        v-bind="props"
        >Save</v-btn
      >
    </template>
  </v-tooltip>
</template>
<script setup lang="ts">
import { ref } from "vue";
import { useConfig } from "@/config"

const props = defineProps({
  policy: { type: String, default: "" },
  resource: { type: String, default: "" },
  context: { type: String, default: "" },
});

const persisted = ref<boolean>(false);

const config = useConfig()

const persist = () => {
  config.policy.value = props.policy
  config.resource.value = props.resource
  config.context.value = props.context

  persisted.value = true
  setTimeout(() => {
    persisted.value = false
  }, 2000)
}
</script>
