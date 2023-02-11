/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./content/**/*.{html,js}", "./layouts/**/*.{html,js}"],
  theme: {
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
