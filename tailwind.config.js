/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./assets/**/*.{html,js}",
    "./views/**/*.{html,js}",
  ],
  darkMode: 'class',
  theme: {
    extend: {
      colors: {
        base: 'rgba(186, 230, 253, .8)',
        navy: {
          light: '#1363DF',
          dark: '#06283D',
        },
      }
    }
  },
  plugins: [
    require('@tailwindcss/forms'),
  ],
}
