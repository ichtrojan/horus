new Vue({
    el: '#app',
    components: {
        'horus': httpVueLoader('public/js/components/Horus.vue'),
        'sort-bar': httpVueLoader('public/js/components/Sortbar.vue'),
        // 'content': httpVueLoader('/components/content.vue'),
        // 'footer': httpVueLoader('/components/footer.vue')
    },
    data: {
        message: 'Hello Vue'
    },
    methods: {},
})