/* eslint-env node */
require('@rushstack/eslint-patch/modern-module-resolution')

module.exports = {
  root: true,
  'extends': [
    'eslint:recommended',
    '@vue/eslint-config-typescript',
  ],
  rules: {
    'vue/multi-word-component-names': ['error', {
      'ignores': 'all'
    }],
  },
  parserOptions: {
    ecmaVersion: 'latest'
  }
}
