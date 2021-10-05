import CryptoJS from 'crypto-js'

var protocol = "https://";
var host = window.location.host;

var AuthMixin = {
    methods: {
        getPreferences: function (params, success, error) {
            this.$http.get(protocol + host + '/wax/preferences' +
                '?digest=' + params.digest +
                '&uuid=' + params.uuid +
                '&sessionid=' + params.sessionid
            ).then(success, error);
        },
        deleteMarketingInfo: function (userId, params, success, error) {
            this.$http.delete(protocol + host + '/wax/marketings/' + userId +
                '?digest=' + params.digest +
                '&uuid=' + params.uuid +
                '&sessionid=' + params.sessionid
            ).then(success, error);
        },
        deleteSurveyInfo: function (userId, params, success, error) {
            this.$http.delete(protocol + host + '/wax/surveys/' + userId +
                '?digest=' + params.digest +
                '&uuid=' + params.uuid +
                '&sessionid=' + params.sessionid
            ).then(success, error);
        },
        addAdditionalInfo: function (userId, params, body, success, error) {
            this.$http.put(protocol + host + '/wax/additional/' + userId +
                '?digest=' + params.digest +
                '&uuid=' + params.uuid +
                '&sessionid=' + params.sessionid, body
            ).then(success, error);
        },
        extractParams: function () {
            var code = this.$route.query.code || null
            var state = this.$route.query.state || null

            var digest = this.$root.$options.hotspot.digest || null
            var uuid = this.$root.$options.hotspot.uuid || null
            var sessionid = this.$root.$options.hotspot.sessionid || null
            var nasid = this.$root.$options.hotspot.nasid || null
            var uamip = this.$root.$options.hotspot.uamip || null
            var uamport = this.$root.$options.hotspot.uamport || null
            var voucher = this.$root.$options.session && this.$root.$options.session['voucherCode'] ? this.$root.$options.session['voucherCode'] : null

            return {
                code: code,
                state: state,
                digest: digest,
                uuid: uuid,
                sessionid: sessionid,
                uamip: uamip,
                uamport: uamport,
                voucher: voucher,
                nasid: nasid
            }
        },
        parseState: function (state, base64) {
            if (base64) {
              state = atob(state)
              state = decodeURIComponent(state)
            }

            var digest = state.split('&')[0]
            var uuid = state.split('&')[1]
            var sessionid = state.split('&')[2]
            var uamip = state.split('&')[3]
            var uamport = state.split('&')[4]
            var voucher = state.split('&')[5]
            var nasid = state.split('&')[6]
            return {
                digest: digest,
                uuid: uuid,
                sessionid: sessionid,
                uamip: uamip,
                uamport: uamport,
                voucher: voucher,
                nasid: nasid
            }
        },
        createWaxURL: function (code, params, endpoint, reset, user) {
            var url = protocol + host + '/wax/register/' + endpoint + '/' + encodeURIComponent(code) +
                '?digest=' + params.digest +
                '&uuid=' + params.uuid +
                '&sessionid=' + params.sessionid +
                '&reset=' + reset +
                '&uamip=' + params.uamip +
                '&uamport=' + params.uamport +
                '&voucher_code=' + params.voucher +
                '&nasid=' + params.nasid +
                '&user=' + user
            return url
        },
        getSocialLoginURL: function (params, social) {
            var url = ''
            switch (social) {
                case 'facebook':
                    url = 'https://www.facebook.com/v2.11/dialog/oauth?' +
                        'client_id=' + params.fb_client_id +
                        '&state=' + encodeURIComponent(params.digest + "&" + params.uuid + "&" + params.sessionid + "&" + params.uamip + "&" + params.uamport + "&" + params.voucher + "&" + params.nasid) +
                        '&scope=email,public_profile,user_birthday,user_likes,user_location' +
                        '&redirect_uri=' + escape('https://' + window.location.host + '/wings/login/facebook')
                    break

                case 'linkedin':
                    url = 'https://www.linkedin.com/oauth/v2/authorization?' +
                        'client_id=' + params.li_client_id +
                        '&state=' + encodeURIComponent(params.digest + "&" + params.uuid + "&" + params.sessionid + "&" + params.uamip + "&" + params.uamport + "&" + params.voucher + "&" + params.nasid) +
                        '&scope=' + escape('r_liteprofile r_emailaddress w_member_social') +
                        '&redirect_uri=' + escape('http://' + window.location.host + '/wings/login/linkedin') +
                        '&response_type=code'
                    break

                case 'instagram':
                    url = 'https://api.instagram.com/oauth/authorize/?' +
                        'client_id=' + params.in_client_id +
                        '&state=' + btoa(encodeURIComponent(params.digest + "&" + params.uuid + "&" + params.sessionid + "&" + params.uamip + "&" + params.uamport + "&" + params.voucher + "&" + params.nasid)) +
                        '&scope=user_profile' +
                        '&redirect_uri=' + escape('https://' + window.location.host + '/wings/login/instagram') +
                        '&response_type=code'
                    break
            }
            return url
        },
        doDedaloLogin: function (user, callback, error, base64) {
            var params = this.extractParams()
            var ip = params.uamip || null
            var port = params.uamport || null
            var digest = params.digest || null
            var uuid = params.uuid || null
            var sessionid = params.sessionid || null

            if (params.state) {
                var state = this.parseState(params.state, base64)
                ip = state.uamip
                port = state.uamport
                digest = state.digest
                uuid = state.uuid
                sessionid = state.sessionid
            }
            var dedaloUrl = ip + ':' + port

            this.$http.get(protocol + host + '/wax/aaa/login' +
                '?digest=' + digest +
                '&uuid=' + uuid +
                '&sessionid=' + sessionid +
                '&username=' + encodeURIComponent(user.id)
            ).then(callback);

            // do dedalo login
            // this.$http.get('http://' + dedaloUrl + '/json/status').then(function (responseStatus) {
            //     // extract info to calculate response
            //     var chap_challenge = responseStatus.body.challenge;
            //     var string_to_hash = "00" + user.password + chap_challenge;
            //
            //     // calculate chap_password with challenge
            //     var response = CryptoJS.MD5(string_to_hash).toString();
            //
            //     // do dedalo login
            //     this.$http.get('http://' + dedaloUrl + '/json/logon?username=' + encodeURIComponent(user.id) +
            //         '&response=' + response).then(callback);
            // }, function (response) {
            //     callback(response)
            // });

            // callback({
            //   body: {
            //     clientState: 1
            //   }
            // })
        },
        doDedaloLogout: function (username, callback) {
            var params = this.extractParams()
            var ip = params.uamip || null
            var port = params.uamport || null
            var digest = params.digest || null
            var uuid = params.uuid || null
            var sessionid = params.sessionid || null

            if (params.state) {
                var state = this.parseState(params.state)
                ip = state.uamip
                port = state.uamport
                digest = state.digest
                uuid = state.uuid
                sessionid = state.sessionid
            }
            var dedaloUrl = ip + ':' + port

            // do dedalo logout
            this.$http.get(protocol + host + '/wax/aaa/logout' +
                '?digest=' + digest +
                '&uuid=' + uuid +
                '&sessionid=' + sessionid +
                '&username=' + username
            ).then(callback);
            //this.$http.get('http://' + dedaloUrl + '/json/logout').then(callback);
        },
        doTempSession: function (email, callback) {
            var params = this.extractParams()
            var ip = params.uamip || null
            var port = params.uamport || null
            var digest = params.digest || null
            var uuid = params.uuid || null
            var sessionid = params.sessionid || null

            if (params.state) {
                var state = this.parseState(params.state)
                ip = state.uamip
                port = state.uamport
                digest = state.digest
                uuid = state.uuid
                sessionid = state.sessionid
            }
            var dedaloUrl = ip + ':' + port

            // do dedalo temp session
            this.$http.get(protocol + host + '/wax/aaa/temp' +
                '?digest=' + digest +
                '&uuid=' + uuid +
                '&sessionid=' + sessionid +
                '&username=' + email
            ).then(callback);
            //this.$http.get('http://' + dedaloUrl + '/www/temporary.chi?username=' + email).then(callback);
        }
    }
};
export default AuthMixin;
