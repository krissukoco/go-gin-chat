// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
    css: ['~/assets/css/main.css'],
    appConfig: {
        apiBaseUrl: process.env.API_BASE_URL || 'http://localhost:8000',
    }
})
