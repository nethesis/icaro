<template>
    <div class="ui segment">
        <div v-if="!dedaloRequested">
            <div class="ui big left icon input">
                <input v-model="authSMS" type="tel" placeholder="Insert your number">
                <i class="talk icon"></i>
            </div>
            <button v-on:click="getCode()" class="ui big button auth-code-cont">Get Code</button>
            <div v-if="errors.badNumber" class="ui tiny icon negative message">
                <i class="remove icon"></i>
                <div class="content">
                    <div class="header">
                        Error sending verification code
                    </div>
                    <p>Your cellphone number is invalid (use the international prefix)</p>
                </div>
            </div>
            <div v-if="codeRequested" class="auth-code-cont">
                <div class="ui big left icon input">
                    <input v-model="authCode" type="number" placeholder="Insert your code">
                    <i class="braille icon"></i>
                </div>
            </div>
            <div class="ui divider"></div>
            <button v-on:click="execLogin()" :disabled="isDisabled()" class="big ui green button">
                Start Navigate
            </button>
        </div>
        <div v-if="dedaloRequested">
            <div v-if="!authorized && !errors.dedaloError" class="ui active centered inline text loader">Authorization in progress...</div>
            <div v-if="authorized" class="ui icon positive message">
                <i class="check icon"></i>
                <div class="content">
                    <div class="header">
                        You are successfully authenticated
                    </div>
                    <p>In a few seconds you will be redirected...</p>
                </div>
            </div>
            <div v-if="errors.dedaloError" class="ui icon negative message">
                <i class="check icon"></i>
                <div class="content">
                    <div class="header">
                        Error on authentication
                    </div>
                    <p>Something went wrong on authentication process </p>
                </div>
            </div>
        </div>
    </div>
</template>

<script>
    import AuthMixin from './../../mixins/auth';
    import {
        setTimeout
    } from 'timers';
    export default {
        name: 'SMSPage',
        mixins: [AuthMixin],
        data() {
            var authorized = false
            var codeRequested = false
            var dedaloRequested = false
            var dedaloError = false
            var badNumber = false
            var badCode = false

            return {
                authorized: authorized,
                codeRequested: codeRequested,
                dedaloRequested: dedaloRequested,
                authSMS: '',
                authCode: '',
                errors: {
                    badNumber: badNumber,
                    badCode: badCode,
                    dedaloError: dedaloError
                }
            }
        },
        methods: {
            isDisabled() {
                return this.authSMS.length == 0 || this.authCode.length == 0
            },
            getCode() {
                this.errors.badNumber = false
                if (!this.authSMS.startsWith('+')) {
                    this.errors.badNumber = true
                }
                var params = this.extractParams()

                // make request to wax
                var url = this.createWaxURL(this.authSMS, params, 'sms')

                // get user id
                this.$http.get(url).then(responseAuth => {
                    this.codeRequested = true
                    this.authCode = responseAuth.body.password || ''
                }, response => {
                    this.codeRequested = false
                    this.errors.badNumber = true
                    console.error(response)
                });
            },
            execLogin() {
                this.dedaloRequested = true
                this.authorized = false
                this.errors.dedaloError = false
                this.errors.badCode = false

                // exec dedalo login
                this.doDedaloLogin({
                    id: this.authSMS,
                    password: this.authCode || ''
                }, responseDedalo => {
                    if (responseDedalo.body.clientState == 1) {
                        this.authorized = true
                        this.errors.dedaloError = false
                        setTimeout(function () {
                            // open redir url
                            window.location.replace(this.$root.$options.hotspot.preferences
                                .captive_redir)
                        }.bind(this), 2500)
                    } else {
                        this.authorized = false
                        this.errors.dedaloError = true
                        this.errors.badCode = true
                    }
                }, error => {
                    this.authorized = false
                    this.errors.dedaloError = true
                    console.error(error)
                })
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

    .auth-code-cont {
        margin-top: 15px !important;
        margin: 0;
    }
</style>