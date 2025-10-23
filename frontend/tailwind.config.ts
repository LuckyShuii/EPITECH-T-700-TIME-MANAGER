import type { Config } from 'tailwindcss'
import daisyui from 'daisyui'

export default {
  content: [
    './index.html',
    './src/**/*.{vue,js,ts,jsx,tsx}'
  ],
  plugins: [
    require('daisyui')
  ],
  daisyui: {
    themes: [
      {
        light: {
          "primary": "#667eea",
          "secondary": "#764ba2",
          "accent": "#37cdbe",
          "neutral": "#3d4451",
          "base-100": "#ffffff",
          "info": "#3b82f6",
          "success": "#10b981",
          "warning": "#f59e0b",
          "error": "#ef4444",
        },
        dark: {
          "primary": "#667eea",
          "secondary": "#764ba2",
          "accent": "#37cdbe",
          "neutral": "#191d24",
          "base-100": "#2a303c",
          "info": "#3b82f6",
          "success": "#10b981",
          "warning": "#f59e0b",
          "error": "#ef4444",
        },
      },
    ],
  },
} satisfies Config