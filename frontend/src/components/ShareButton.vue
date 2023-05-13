<template>
  <v-dialog v-model="dialog" width="600px">
    <template v-slot:activator>
      <v-btn
        @click="share"
        prepend-icon="mdi-share"
        :variant="variant"
        :block="block"
        title="create a shareable URL with the current input"
        >Share</v-btn
      >
    </template>

    <v-card title="Share Policy">
      <v-divider />
      <v-card-text>
        <v-text-field label="URL" v-model="url" />
      </v-card-text>
      <v-card-actions>
        <v-btn @click="dialog = false">Close</v-btn>
        <v-spacer />
        <v-tooltip
          :model-value="copied"
          location="top"
          text="Copied"
          :open-on-hover="false"
        >
          <template v-slot:activator="{ props }">
            <v-btn
              variant="tonal"
              :color="btnColor"
              @click="copy(url)"
              :disabled="!isSupported"
              v-bind="props"
              >Copy URL to Clipboard</v-btn
            >
          </template>
        </v-tooltip>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>
<script setup lang="ts">
import { ref } from "vue";
import * as lzstring from "lz-string";
import { useRouter } from "vue-router";
import { useClipboard } from "@vueuse/core";
import { PropType } from "vue";
import { btnColor } from '@/config'

const props = defineProps({
  policy: { type: String, default: "" },
  resource: { type: String, default: "" },
  context: { type: String, default: "" },
  variant: { type: String as PropType<"outlined" | "text"> },
  block: { type: Boolean },
});

const dialog = ref<boolean>(false);
const loading = ref<boolean>(false);
const url = ref<string>("");

const router = useRouter();

const { copy, copied, isSupported } = useClipboard({ source: url });

const share = () => {
  loading.value = true;
  const content = JSON.stringify({
    policy: props.policy,
    resource: props.resource,
    context: props.context,
  });

  const compressed = lzstring.compressToBase64(content);

  url.value =
    window.location.origin +
    "/" +
    router.resolve({ name: "home", query: { content: compressed } }).href;
  dialog.value = true;
  loading.value = false;
};
</script>
