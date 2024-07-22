import * as monaco from 'monaco-editor'

function fixAddCommand(editor: monaco.editor.IStandaloneCodeEditor, context: string): { dispose(): void } {
  const addCommand = editor.addCommand

  editor.addCommand = function addCommand_hijacked(keybinding, handler) {
    return addCommand.call(this, keybinding, handler, context)
  }

  return {
    dispose: () => {
      editor.addCommand = addCommand
    }
  }
}

export const useEditorFix = (editor: monaco.editor.IStandaloneCodeEditor, id: string) => {
  const editorFocusedContextKeyName = `__isEditorFocused-${id}`
  const isEditorFocused = editor.createContextKey<boolean>(editorFocusedContextKeyName, false)
  const onBlurDisposable = editor.onDidBlurEditorWidget(() => isEditorFocused.set(false))
  const onFocusDisposable = editor.onDidFocusEditorText(() => isEditorFocused.set(true))
  const disposeAddCommandFix = fixAddCommand(editor, editorFocusedContextKeyName)

  return () => {
    onBlurDisposable.dispose()
    onFocusDisposable.dispose()
    disposeAddCommandFix.dispose()
    editor.dispose()
  }
}
