/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    './src/**/*.{html,js,svelte,ts}',
    './index.html',
  ],
  theme: {
    extend: {
      colors: {
        'k8s-running': '#10b981',
        'k8s-pending': '#f59e0b',
        'k8s-failed': '#ef4444',
        'k8s-success': '#22c5e',
      },
    },
  },
  plugins: [],
}
