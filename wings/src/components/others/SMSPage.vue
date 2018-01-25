<template>
    <div class="ui segment form">
        <div v-if="!dedaloRequested">
            <div v-if="!codeRequested" class="inline field" v-bind:class="{ error: errors.badInput }">
                <label>{{ $t("sms.number") }}</label>
                <div class="ui big left icon input">
                    <input v-model="authSMS" type="tel" :placeholder="$t('sms.insert_number')">
                    <i class="talk icon"></i>
                </div>
            </div>
            <button v-if="!codeRequested" v-on:click="getCode()" class="ui big button">{{ $t("sms.get_code") }}</button>
            <div v-if="errors.badNumber" class="ui tiny icon negative message">
                <i class="remove icon"></i>
                <div class="content">
                    <div class="header">
                        {{ $t("sms.error_code") }}
                    </div>
                    <p>{{ $t("sms.error_code_sub") }}</p>
                </div>
            </div>
            <div v-if="codeRequested">
                <div class="inline field">
                    <label>Code</label>
                    <div class="ui big left icon input">
                        <input v-model="authCode" type="number" :placeholder="$t('sms.insert_your_code')">
                        <i class="braille icon"></i>
                    </div>
                </div>
            </div>
            <button v-on:click="getCode(true)" v-if="authReset && resetDone != 'true'" class="ui red button auth-code-cont">
                {{ $t("sms.reset_code") }}
            </button>
            <div class="ui divider"></div>
            <button v-on:click="execLogin()" :disabled="isDisabled()" class="big ui green button">
                {{ $t("sms.start_navigate") }}
            </button>
        </div>
        <div v-if="dedaloRequested">
            <div v-if="!authorized && !errors.dedaloError" class="ui active centered inline text loader">{{ $t("sms.auth_progress") }}...</div>
            <div v-if="authorized" class="ui icon positive message">
                <i class="check icon"></i>
                <div class="content">
                    <div class="header">
                        {{ $t("sms.auth_success") }}
                    </div>
                    <p>{{ $t("sms.auth_success_sub") }}...</p>
                </div>
            </div>
            <div v-if="errors.dedaloError" class="ui icon negative message">
                <i class="remove icon"></i>
                <div class="content">
                    <div class="header">
                        {{ $t("sms.auth_error") }}
                    </div>
                    <p>{{ $t("sms.auth_error_sub") }}</p>
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
            var badInput = false
            var authReset = false
            var resetDone = false

            return {
                authorized: authorized,
                codeRequested: codeRequested,
                dedaloRequested: dedaloRequested,
                authSMS: '',
                authCode: '',
                authReset: authReset,
                resetDone: resetDone,
                errors: {
                    badNumber: badNumber,
                    badCode: badCode,
                    dedaloError: dedaloError,
                    badInput: badInput
                }
            }
        },
        methods: {
            isDisabled() {
                return this.authSMS.length == 0 || this.authCode.length == 0
            },
            getCode(reset) {
                this.errors.badNumber = false
                if (!this.authSMS.startsWith('+')) {
                    this.errors.badInput = true
                    return
                }
                var params = this.extractParams()

                // make request to wax
                var url = this.createWaxURL(this.authSMS, params, 'sms', reset)

                // get user id
                this.$http.get(url).then(responseAuth => {
                    this.codeRequested = true
                    this.authReset = responseAuth.body.exists
                    this.resetDone = responseAuth.body.reset
                }, error => {
                    this.codeRequested = false
                    this.errors.badNumber = true
                    console.error(error)
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