import type { Config } from 'tailwindcss';

export default {
  content: ['./src/**/*.{js,jsx,ts,tsx}'],
  theme: {
    extend: {
      colors: {
        primary: '#911841', // primary color
        secondary: '#e68f24', // secondary color
        'primary-dark': '#750e30', // darker shade of primary color
        'secondary-dark': '#cc6e00', // darker shade of secondary color
        'primary-600': '#87123a', // slightly darker shade of primary color
        'primary-500': '#9c1c48', // slightly lighter shade of primary color        
      },
      fontFamily: {
        lora: ['Lora', 'serif'],
        jost: ['Jost', 'sans-serif'],
      },
    },
  },
  plugins: [],
} satisfies Config;
