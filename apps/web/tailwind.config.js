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
        sky: {
          50: '#f0f9ff',
          100: '#e0f2fe',
          200: '#bae6fd',
          300: '#7dd3fc',
          400: '#38bdf8',
          500: '#0ea5e9',
          600: '#0284c7',
          700: '#0369a1',
          800: '#075985',
          900: '#0c4a6e',
        },
        warm: {
          50: '#fafaf9',
          100: '#f5f5f4',
          200: '#e7e5e4',
          300: '#d6d3d1',
          400: '#a8a29e',
          500: '#78716c',
          600: '#57534e',
          700: '#44403c',
          800: '#292524',
          900: '#1c1917',
        },
      },
      fontFamily: {
        sans: ['Inter', 'Pretendard', 'Noto Sans JP', 'sans-serif'],
        display: ['Gowun Dodum', 'Inter', 'sans-serif'],
      },
      borderRadius: {
        card: '12px',
        btn: '12px',
      },
    },
  },
  plugins: [],
}
