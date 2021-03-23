Vue.$cookies.config('7d')
new Vue({
    el: '#app',
    components: {
        'horus': httpVueLoader('public/js/components/Horus.vue'),
        'sort-bar': httpVueLoader('public/js/components/Sortbar.vue'),
        // 'content': httpVueLoader('/components/content.vue'),
        // 'footer': httpVueLoader('/components/footer.vue')
    },
    data: {
        loggedIn : false,
        key : "",
        cookie: "",
        error: false
    },
    methods: {
        login(key){
            let formData = new FormData();
            formData.append('key', this.key);
            const request = {
                method: "POST",
                body: formData
            };
            fetch("./login", request)
                .then((response) => response.json())
                .then((data) => {
                    if (data.status){
                        this.$cookies.set(this.$cookies.get('horus'))
                        this.loggedIn = true
                        this.cookie = "1"
                    }else{
                        this.error = true
                    }
                });
        },
        logout(){
            this.cookie = "0"
            this.$cookies.remove('horus')
            fetch("./logout",)
            .then((response) => response.json())
            .then((data) => {
                if (data.status){
                    console.log("Logged out")
                    this.loggedIn = false
                }else{
                    console.log("Error")
                }
            });
        }
    },
    created(){
        if(this.$cookies.get('horus') != null){
            this.loggedIn = true
        }
    }
})