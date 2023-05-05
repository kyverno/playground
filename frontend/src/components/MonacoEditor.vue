<template>
  <div class="monaco-editor-vue3" :style="style"></div>
</template>

<script>
import { defineComponent, computed, toRefs } from "vue";
import * as monaco from "monaco-editor";

export default defineComponent({
  name: "MonacoEditor",
  props: {
    uri: { type: Object },
    width: { type: [String, Number], default: "100%" },
    height: { type: [String, Number], default: "100%" },
    value: String,
    language: { type: String, default: "javascript" },
    theme: { type: String, default: "vs" },
    options: {
      type: Object,
      default() {
        return {};
      },
    },
  },
  emits: ["editorWillMount", "editorDidMount", "change", "update:value"],
  setup(props) {
    const { width, height } = toRefs(props);
    const style = computed(() => {
      const fixedWidth = width.value.toString().includes("%")
        ? width.value
        : `${width.value}px`;
      const fixedHeight = height.value.toString().includes("%")
        ? height.value
        : `${height.value}px`;
      return {
        width: fixedWidth,
        height: fixedHeight,
        "text-align": "left",
      };
    });
    return {
      style,
    };
  },
  mounted() {
    this.initMonaco();
  },
  beforeUnmount() {
    this.editor && this.editor.dispose();
  },
  methods: {
    initMonaco() {
      this.$emit("editorWillMount", monaco);
      const { value, language, theme, options } = this;
      let model = null;
      if (this.uri) {
        model = monaco.editor.getModel(this.uri);
      }
      if (!model) {
        model = monaco.editor.createModel(value, language, this.uri);
      }

      this.editor = monaco.editor.create(this.$el, {
        value: value,
        language: language,
        theme: theme,
        ...options,
        model,
      });

      // @event `change`
      const editor = this._getEditor();
      editor &&
        editor.onDidChangeModelContent((event) => {
          const value = editor.getValue();
          if (this.value !== value) {
            this.$emit("change", value, event);
            this.$emit("update:value", value);
          }
        });

      this.$emit("editorDidMount", this.editor);
    },
    _setModel(value, original) {
      const { language } = this;
      const originalModel = monaco.editor.createModel(original, language, this.uri);
      const modifiedModel = monaco.editor.createModel(value, language, this.uri);
      this.editor.setModel({
        original: originalModel,
        modified: modifiedModel,
      });
    },
    _setValue(value) {
      let editor = this._getEditor();
      if (editor) return editor.setValue(value);
    },
    _getValue() {
      let editor = this._getEditor();
      if (!editor) return "";
      return editor.getValue();
    },
    _getEditor() {
      return this.editor;
    },
    _setOriginal() {
      const { original } = this.editor.getModel();
      original.setValue(this.original);
    },
  },
  watch: {
    options: {
      deep: true,
      handler(options) {
        this.editor.updateOptions(options);
      },
    },
    value() {
      this.value !== this._getValue() && this._setValue(this.value);
    },
    original() {
      this._setOriginal();
    },
    language() {
      if (!this.editor) return;

      monaco.editor.setModelLanguage(this.editor.getModel(), this.language);
    },
    theme() {
      monaco.editor.setTheme(this.theme);
    },
  },
});
</script>
