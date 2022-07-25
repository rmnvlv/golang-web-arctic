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
      },
      backgroundImage: {
        'arctic-sea-1': "url('/a/images/arctic-sea-1.jpg')",
        'arctic-sea-2': "url('/a/images/arctic-sea-2.jpg')",
        'arctic-sea-3': "url('/a/images/arctic-sea-3.jpg')",
        'arctic-sea-4': "url('/a/images/arctic-sea-4.jpg')",
        'arctic-sea-5': "url('/a/images/arctic-sea-5.jpg')",
      }
    }
  },
  plugins: [
    require('@tailwindcss/forms'),
  ],
}
