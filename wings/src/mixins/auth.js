import CryptoJS from 'crypto-js'

var AuthMixin = {
    methods: {
        getPreferences(params, callback) {
            var host = 'hstest.neth.eu:8181' //window.location.host
            this.$http.get("http://" + host + '/wax/preferences' +
                '?digest=' + params.digest +
                '&uuid=' + params.uuid +
                '&sessionid=' + params.sessionid
            ).then(callback);
        },
        extractParams(searchURL) {
            var code = this.$route.query.code || null
            var state = this.$route.query.state || null

            // check if google login
            if (this.$route.path == '/login/google') {
                var params = this.$route.hash.substr(1).split('&') || null
                if (params.length > 1) {
                    var code = params[1].split('=')[1]
                    var state = unescape(params[0].split('=')[1])
                }
            } else {

            }

            var userId = '-'
            var digest = this.$root.$options.hotspot.digest || null
            var uuid = this.$root.$options.hotspot.uuid || null
            var sessionid = this.$root.$options.hotspot.sessionid || null

            return {
                code: code,
                state: state,
                digest: digest,
                uuid: uuid,
                sessionid: sessionid
            }
        },
        parseState(state) {
            var digest = state.split('&')[0]
            var uuid = state.split('&')[1]
            var sessionid = state.split('&')[2]
            return {
                digest: digest,
                uuid: uuid,
                sessionid: sessionid
            }
        },
        createWaxURL(code, state, endpoint) {
            var params = this.parseState(state)
            var url = 'http://' + window.location.host + ':8181/wax/register/' + endpoint + '/' + code +
                '?digest=' + params.digest +
                '&uuid=' + params.uuid +
                '&sessionid=' + params.sessionid
            return url
        },
        getSocialLoginURL(params, social) {
            var url = ''
            switch (social) {
                case 'facebook':
                    url = 'https://www.facebook.com/v2.11/dialog/oauth?' +
                        'client_id=' + params.fb_client_id +
                        '&state=' + encodeURIComponent(params.digest + "&" + params.uuid + "&" + params.sessionid) +
                        '&scope=email,public_profile,user_birthday,user_likes,user_location' +
                        '&redirect_uri=' + escape('http://' + window.location.host + '/login/facebook')
                    break

                case 'google':
                    url = 'https://accounts.google.com/o/oauth2/v2/auth?' +
                        'client_id=' + params.gl_client_id +
                        '&state=' + encodeURIComponent(params.digest + "&" + params.uuid + "&" + params.sessionid) +
                        '&scope=' + escape('profile email https://www.googleapis.com/auth/plus.login https://www.googleapis.com/auth/plus.me https://www.googleapis.com/auth/userinfo.email https://www.googleapis.com/auth/userinfo.profile') +
                        '&redirect_uri=' + escape('http://' + window.location.host + '/login/google') +
                        '&include_granted_scopes=true' +
                        '&response_type=token'
                    break

                case 'linkedin':
                    url = 'https://www.linkedin.com/oauth/v2/authorization?' +
                        'client_id=' + params.li_client_id +
                        '&state=' + encodeURIComponent(params.digest + "&" + params.uuid + "&" + params.sessionid) +
                        '&scope=' + escape('r_basicprofile r_emailaddress w_share') +
                        '&redirect_uri=' + escape('http://' + window.location.host + '/login/linkedin') +
                        '&response_type=code'
                    break
            }
            return url
        },
        doDedaloLogin(user, callback) {
            this.$http.get('http://10.1.0.1:3990/json/status').then(responseStatus => {
                // extract info to calculate response
                var chap_challenge = responseStatus.body.challenge;
                var string_to_hash = "00" + user.password + chap_challenge;

                // calculate chap_password with challenge
                var response = CryptoJS.MD5(string_to_hash).toString();

                // do dedalo login
                this.$http.get('http://10.1.0.1:3990/json/logon?username=' + user.id +
                    '&response=' + response).then(callback);
            }, response => {
                console.error(response)
            });
        }
    }
};
export default AuthMixin;