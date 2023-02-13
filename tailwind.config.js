/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./content/**/*.{html,js}", "./layouts/**/*.{html,js}"],
  theme: {
    fontFamily: {
      sans: ["Roboto", "sans-serif"],
    },
    extend: {},
  },
  plugins: [],
  // TODO: remove this after dev
  safelist: [
    {
      pattern: /.*/,
    },
  ],
}
