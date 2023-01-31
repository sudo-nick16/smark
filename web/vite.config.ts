import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react-swc'

// https://vitejs.dev/config/
export default defineConfig({
    plugins: [react()],
    build: {
        outDir: '../extension/options',
        emptyOutDir: true,
        watch: {
            buildDelay: 0,
            clearScreen: true,
            skipWrite: false,
            include: ['src/**', 'src/*.ts', 'src/**/*.ts', 'src/**/*.tsx'], 
        }
    },
    base: ""
})
