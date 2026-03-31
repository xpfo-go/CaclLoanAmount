import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

const REPO_BASE = '/CaclLoanAmount/'

export default defineConfig(({ command }) => ({
  base: command === 'build' ? REPO_BASE : '/',
  plugins: [react()],
}))
