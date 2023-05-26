module.exports = {
  root: true,
  env: {
    node: true
  },
  extends: ['plugin:prettier/recommended', 'plugin:vue/vue3-essential', 'eslint:recommended', '@vue/eslint-config-typescript'],
  rules: {
    'vue/multi-word-component-names': 'off',
    'vue/no-dupe-keys': 'off',
    'prettier/prettier': ['error', { semi: false, endOfLine: 'auto' }]
  },
  plugins: ['prettier']
}
