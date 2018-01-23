<template>
    <div class="ui segment">
        <div v-if="voucherAvailable && !voucherValidated">
            <h3>Insert Voucher code</h3>
            <div class="ui big left icon input">
                <input v-model="authCode" type="email" placeholder="Insert voucher code">
                <i class="braille icon"></i>
            </div>
            <button v-on:click="validateCode()" class="ui big button auth-code-cont">Validate Code</button>
            <div v-if="badCode" class="ui tiny icon negative message">
                <i class="remove icon"></i>
                <div class="content">
                    <div class="header">
                        Error validating voucher code
                    </div>
                    <p>It seems that your voucher code is not correct</p>
                </div>
            </div>
        </div>
        <div v-if="!voucherAvailable || (voucherAvailable && voucherValidated)">
            <h3>Choose your Login</h3>
            <div class="ui relaxed list">
                <div v-if="hotspot.preferences.facebook_login == 'true'" class="item">
                    <router-link to="/login/facebook" class="ui facebook button big fluid">
                        <i class="facebook icon"></i>
                        Facebook
                    </router-link>
                </div>
                <div v-if="hotspot.preferences.google_login == 'true'" class="item">
                    <router-link to="/login/google" class="ui google plus button big fluid">
                        <i class="google plus icon"></i>
                        Google Plus
                    </router-link>
                </div>
                <div v-if="hotspot.preferences.linkedin_login == 'true'" class="item">
                    <router-link to="/login/linkedin" class="ui linkedin button big fluid">
                        <i class="linkedin icon"></i>
                        LinkedIn
                    </router-link>
                </div>
            </div>
            <div class="ui divider"></div>
            <div class="ui relaxed list">
                <div v-if="hotspot.preferences.sms_login == 'true'" class="item">
                    <router-link to="/login/sms" class="ui button green big fluid">
                        <i class="talk icon"></i>
                        SMS
                    </router-link>
                </div>
                <div v-if="hotspot.preferences.email_login == 'true'" class="item">
                    <router-link to="/login/email" class="ui button yellow big fluid">
                        <i class="mail icon"></i>
                        Email
                    </router-link>
                </div>
            </div>
        </div>
    </div>
</template>

<script>
    import AuthMixin from './../mixins/auth';
    export default {
        name: 'LoginPage',
        mixins: [AuthMixin],
        data() {
            var voucherAvailable = false
            var voucherValidated = false
            var badCode = false
            var authCode = ''
            if (this.$root.$options.hotspot.preferences.voucher_login == "true") {
                voucherAvailable = true
            }
            return {
                hotspot: {
                    preferences: this.$root.$options.hotspot.preferences
                },
                voucherAvailable: voucherAvailable,
                voucherValidated: voucherValidated,
                authCode: authCode,
                badCode: badCode
            }
        },
        methods: {
            validateCode() {
                this.badCode = false
                if (!this.authCode.length > 1) {
                    this.badCode = true
                }
                var params = this.extractParams()

                // make request to wax
                var url = this.createWaxURL(this.authCode, params, 'voucher')

                // get user id
                this.$http.get(url).then(responseAuth => {
                    this.voucherValidated = true
                }, response => {
                    this.voucherValidated = false
                    this.badCode = true
                    console.error(response)
                });
            }
        }
    }
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
    h1,
    h2 {
        font-weight: normal;
    }

    ul {
        list-style-type: none;
        padding: 0;
    }

    li {
        display: inline-block;
        margin: 0 10px;
    }

    a {
        color: #42b983;
    }

    .divider {
        margin: 20px;
    }

    .item {
        margin: 10px;
    }

    .auth-code-cont {
        margin-top: 15px !important;
        margin: 0;
    }
</style>