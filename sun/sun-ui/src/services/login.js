var LoginService = {
    methods: {
        execLogin(body, success, error) {
            this.$http.post('https://hstest.neth.eu/api/login', body).then(success, error);
        },
        execLogout(token, success, error) {
            this.$http.post('https://hstest.neth.eu/api/logout', {}, {
                headers: {
                    'Token': token
                }
            }).then(success, error);
        },
        execGetInfo(id, token, success, error) {
            this.$http.get('https://hstest.neth.eu/api/accounts/' + id, {
                headers: {
                    'Token': token
                }
            }).then(success, error);
        }
    }
};
export default LoginService;