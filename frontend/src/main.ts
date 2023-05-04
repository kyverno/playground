import App from './App.vue'
import editorWorker from 'monaco-editor/esm/vs/editor/editor.worker?worker';
import jsonWorker from 'monaco-editor/esm/vs/language/json/json.worker?worker';
import yamlWorker from './worker/yaml.worker.js?worker?worker';
import { setDiagnosticsOptions } from 'monaco-yaml';

import { createApp } from 'vue'
import { registerPlugins } from '@/plugins'

import clusterpolicy from './schemas/clusterpolicy.json'

setDiagnosticsOptions({
    enableSchemaRequest: true,
    hover: true,
    completion: true,
    validate: true,
    format: true,
    schemas: [
        { ...clusterpolicy, uri: '/schemas/clusterpolicy.json', fileMatch: ['policy.yaml'] }
    ]
});

// @ts-ignore
self.MonacoEnvironment = {
    getWorker(_: any, label: string) {
        if (label === 'yaml') {
            return new yamlWorker();
        }
        if (label === 'json') {
            return new jsonWorker();
        }
        return new editorWorker();
    },
};

const app = createApp(App)

registerPlugins(app)

app.mount('#app')
