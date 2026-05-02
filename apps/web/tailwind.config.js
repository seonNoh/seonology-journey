/** @type {import('tailwindcss').Config} */
export default {
  content: ['./index.html', './src/**/*.{ts,tsx}'],
  theme: {
    extend: {
      colors: {
        sakura: {
          50: '#fff5f7',
          100: '#ffe3ea',
          200: '#ffc7d4',
          300: '#ffa3b8',
          400: '#ff7896',
          500: '#f25c7a',
          600: '#d94560',
          700: '#b3354c',
          800: '#8a283a',
          900: '#601c29',
        },
      },
      fontFamily: {
        sans: ['Inter', 'Pretendard', 'Noto Sans JP', 'sans-serif'],
      },
    },
  },
  plugins: [],
}
