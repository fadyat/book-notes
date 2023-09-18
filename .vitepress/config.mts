import {defineConfig} from 'vitepress'

// https://vitepress.dev/reference/site-config
export default defineConfig({
    title: "Book notes",
    description: "Some book notes for better learning",
    themeConfig: {
        // https://vitepress.dev/reference/default-theme-config
        nav: [
            {text: 'Home', link: '/'},
        ],

        socialLinks: [
            {icon: 'github', link: 'https://github.com/fadyat/book-notes'},
        ]
    },
    ignoreDeadLinks: true,
    base: '/book-notes/',
})
