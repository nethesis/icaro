var LoginService = {
    methods: {
        execLogin(body, success, error) {
            this.$http.post('https://' + window.location.host + '/api/login', body).then(success, error);
        },
        execLogout(token, success, error) {
            this.$http.post('https://' + window.location.host + '/api/logout', {}, {
                headers: {
                    'Token': token
                }
            }).then(success, error);
        },
        execGetInfo(id, token, success, error) {
            this.$http.get('https://' + window.location.host + '/api/accounts/' + id, {
                headers: {
                    'Token': token
                }
            }).then(success, error);
        },
        execChangePassword(password, id, token, success, error) {
            this.$http.put('https://' + window.location.host + '/api/accounts/' + id, {
                password: password
            }, {
                headers: {
                    'Token': token
                }
            }).then(success, error);
        }
    }
};
export default LoginService;