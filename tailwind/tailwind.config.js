/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "../templates/**/*.{gohtml,html}"
  ],
  theme: {
    extend: {
      fontFamily: {
        poppins: ['Poppins']
      }
    },
  },
  plugins: [],
}

