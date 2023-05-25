<template>
  <div class="monaco-editor-vue3" :style="style"></div>
</template>

<script>
import { defineComponent, computed, toRefs } from "vue";
import * as monaco from "monaco-editor";

export default defineComponent({
  props: {
    width: { type: [String, Number], default: "100%" },
    height: { type: [String, Number], default: "100%" },
    value: String,
    original: String,
    language: { type: String, default: "javascript" },
    theme: { type: String, default: "vs" },
    options: {
      type: Object,
      default() {
        return {};
      },
    },
  },
  emits: ["editorWillMount", "editorDidMount", "change", "update:value", "update:original"],
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
      const { value, language, theme, options, original } = this;

      this.editor = monaco.editor.createDiffEditor(this.$el, {
        value: value,
        language: language,
        automaticLayout: true,
        theme: theme,
        ...options,
      });
      this._setModel(value, original)

      // @event `change`
      const editor = this._getEditor();
      editor &&
        editor.onDidChangeModelContent(() => {
          const value = editor.getValue();
          if (this.value !== value) {
            this.$emit("update:value", value);
          }
        });

      // @event `change`
      const orignalEditor = this._getOriginalEditor();
      orignalEditor &&
        orignalEditor.onDidChangeModelContent(() => {
            const value = orignalEditor.getValue();
            if (this.original !== value) {
              this.$emit("update:original", value);
            }
          });


      this.$emit("editorDidMount", this.editor);
    },
    _setModel(value, original) {
      const { language } = this;
      const originalModel = monaco.editor.createModel(original, language);
      const modifiedModel = monaco.editor.createModel(value, language);
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
    _getOriginal() {
      let editor = this._getOriginalEditor();
      if (!editor) return "";
      return editor.getValue();
    },
    _getEditor() {
      if (!this.editor) return null
      return this.editor._modifiedEditor
    },
    _getOriginalEditor() {
      if (!this.editor) return null
      return this.editor._originalEditor
    },
    _setOriginal(value) {
      let editor = this._getOriginalEditor();
      if (editor) return editor.setValue(value);
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
      this.original !== this._getOriginal() && this._setOriginal(this.original);
    },
    language() {
      if (!this.editor) return;
      
      monaco.editor.setModelLanguage(original, this.language)
      monaco.editor.setModelLanguage(modified, this.language)
    },
    theme() {
      monaco.editor.setTheme(this.theme);
    },
  },
});
</script>
