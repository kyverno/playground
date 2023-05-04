<template>
<div style="height: 100%; position: relative;">
  <MonacoEditor
    language="yaml"
    theme="vs-dark"
    :value="props.modelValue"
    @update:value="(event: string) => emit('update:modelValue', event)"
    :options="options"
    ref="monaco"
  />
  <v-card class="config" theme="dark" color="black" v-if="false">
    <v-card-text class="my-0 py-1">
      <v-switch v-model="autocompleteOnEnter" label="autocomplete" hide-details density="compact" />
    </v-card-text>
  </v-card>
</div>
</template>

<script setup lang="ts">
import MonacoEditor from "monaco-editor-vue3";
import { ref, watch } from "vue";
import { editor, Uri, KeyCode } from "monaco-editor";

const props = defineProps({
  modelValue: { type: String, default: "" },
});

const emit = defineEmits(["update:modelValue"]);

const monaco = ref(null);
const uri = Uri.parse("policy.yaml");
let model = editor.getModel(uri);
if (!model) {
  model = editor.createModel(props.modelValue, "yaml", uri);
}

watch(props, (props) => {
  if (!monaco.value) return;

  // @ts-ignore
  if (props.modelValue !== monaco.value._getValue()) {
    // @ts-ignore
    monaco.value._setValue(props.modelValue);
  }
});


const autocompleteOnEnter = ref(true);
const eventRegistered = ref(true);

watch(monaco, (e: any) => {
  if (!e) return;

  const edit = e._getEditor();

  edit.setModel(model);

  if (eventRegistered.value) return;

  edit.onKeyUp((e: any) => {
    if (!autocompleteOnEnter.value) return;

    const position = edit.getPosition();
    const text = edit.getModel().getLineContent(position.lineNumber).trim();
    if (e.keyCode === KeyCode.Enter && !text) {
      edit.trigger("", "editor.action.triggerSuggest", "");
    }
  });

  eventRegistered.value = true
});

const options = {
  colorDecorators: true,
  lineHeight: 24,
  tabSize: 2,
};
</script>

<style scoped>
.config {
  position: absolute;
  bottom: 60px;
  width: 80%;
  margin-left: 10%;
}
</style>