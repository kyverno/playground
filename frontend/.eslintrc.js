/* eslint-env node */
require('@rushstack/eslint-patch/modern-module-resolution')

module.exports = {
  root: true,
  extends: [
    'plugin:vue/vue3-essential',
    'eslint:recommended',
    '@vue/eslint-config-typescript',
    '@vue/eslint-config-prettier/skip-formatting'
  ],
  rules: {
    'vue/multi-word-component-names': 'off',
    'vue/no-dupe-keys': 'off',
    'prettier/prettier': ['error', { semi: false, endOfLine: 'auto' }]
  },
  parserOptions: {
    ecmaVersion: 'latest'
  }
}
