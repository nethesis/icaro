<template>
    <div class="ui segment form">
        <div v-if="!dedaloRequested">
            <div v-if="choosedMode && !authSMS" class="inline field" v-bind:class="{ error: errors.badInput }">
                <label>{{ $t("sms.prefix") }}</label>
                <div class="ui fluid pointing search selection dropdown select-state">
                    <input type="hidden" name="country">
                    <i class="dropdown icon"></i>
                    <div class="default text">{{$t('sms.select_state')}}</div>
                    <div class="menu">
                        <div @click="setPrefix(c.dial_code)" v-for="c in countries" v-bind:key="c.code" class="item">
                            <i :class="[c.code.toLowerCase(), 'flag']"></i>{{c.name}} ({{c.dial_code}})
                        </div>
                    </div>
                </div>
            </div>
            <div v-if="choosedMode" class="inline field">
                <label>{{ $t("sms.number") }}</label>
                <div class="ui big left icon input">
                    <input v-model="authSMS" type="tel" :placeholder="$t('sms.insert_number')">
                    <i class="talk icon"></i>
                </div>
            </div>
            <button v-if="!codeRequested && !choosedMode" v-on:click="chooseMode()" class="ui big button request-code">{{ $t("sms.not_have_code") }}</button>
            <button v-if="!codeRequested && !choosedMode" v-on:click="chooseMode(true)" class="ui big button request-code">{{ $t("sms.have_code") }}</button>
            <button v-if="!codeRequested && choosedMode" v-on:click="getCode(true)" class="ui big button request-code">{{ $t("sms.get_code") }}</button>
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
                    <label>{{ $t("sms.code") }}</label>
                    <div class="ui big left icon input">
                        <input v-model="authCode" type="number" :placeholder="$t('sms.insert_your_code')">
                        <i class="braille icon"></i>
                    </div>
                </div>
                <div class="ui compact message info no-margin-top">
                    <div class="content">
                        <div class="header">
                            {{$t('sms.wait')}}
                        </div>
                        <p>{{$t('sms.we_are_sending_sms_code')}}</p>
                    </div>
                </div>
            </div>
            <div class="ui divider"></div>
            <button v-on:click="execLogin()" :disabled="isDisabled()" class="big ui green button">
                {{ $t("sms.start_navigate") }}
            </button>
            <div v-if="authReset && resetDone != 'true'" class="ui divider"></div>
            <button v-on:click="getCode(true)" v-if="authReset && resetDone != 'true'" class="ui red button">
                {{ $t("sms.reset_code") }}
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
            <div v-if="authorized">
                <h3>{{ $t("login.disclaimer_marketing") }}</h3>
                <div class="inline field">
                    <textarea readonly class="text-center" v-model="hotspot.disclaimers.marketing_use"></textarea>
                </div>
                <button v-on:click="deleteInfo()" class="ui big button red">{{ $t("login.decline") }}</button>
                <button v-on:click="accept()" class="ui big button green">{{ $t("login.accept") }}</button>
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
            var params = this.extractParams()

            this.getPreferences(params, success => {
                this.$parent.hotspot.disclaimers = success.body.disclaimers
                this.$root.$options.hotspot.disclaimers = success.body.disclaimers
                this.hotspot.disclaimers = success.body.disclaimers
            }, error => {
                console.error(error)
            })

            return {
                authorized: false,
                choosedMode: (this.$route.query.num && this.$route.query.code) ? true : false,
                codeRequested: this.$route.query.code || false,
                dedaloRequested: false,
                authPrefix: '',
                authSMS: this.$route.query.num || '',
                authCode: this.$route.query.code || '',
                authReset: this.$route.query.code || false,
                userId: this.$route.query.user || 0,
                resetDone: false,
                errors: {
                    badNumber: false,
                    badCode: false,
                    dedaloError: false,
                    badInput: false
                },
                countries: require('./../../i18n/countries.json'),
                hotspot: {
                    disclaimers: this.$root.$options.hotspot.disclaimers
                },
            }
        },
        methods: {
            isDisabled() {
                return this.authSMS.length == 0 || this.authCode.length == 0
            },
            setPrefix(prefix) {
                this.authPrefix = prefix
            },
            chooseMode(haveCode) {
                this.choosedMode = true
                if (haveCode) {
                    this.codeRequested = true
                }
                setTimeout(function () {
                    $('.ui.dropdown')
                        .dropdown();
                }, 100)
            },
            getCode(reset) {
                this.errors.badNumber = false
                if (!(this.authPrefix + this.authSMS).startsWith('+')) {
                    this.errors.badInput = true
                    return
                }
                var params = this.extractParams()

                // make request to wax
                var url = this.createWaxURL(this.authPrefix + this.authSMS, params, 'sms', reset)

                // get user id
                this.$http.get(url).then(responseAuth => {
                    this.codeRequested = true
                    this.authReset = responseAuth.body.exists
                    this.resetDone = responseAuth.body.reset
                    this.userId = responseAuth.body.user_db_id
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
                    id: this.authPrefix + this.authSMS,
                    password: this.authCode || ''
                }, responseDedalo => {
                    if (responseDedalo.body.clientState == 1) {
                        this.authorized = true
                        this.errors.dedaloError = false
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
            },
            deleteInfo() {
                // extract code and state
                var params = this.extractParams()
                this.deleteMarketingInfo(this.userId, params, function (success) {
                    this.accept()
                }, function (error) {
                    console.error(error)
                    if (error.status == 404) {
                        this.accept()
                    }
                })
            },
            accept() {
                // open redir url
                window.location.replace(this.$root.$options.hotspot.preferences
                    .captive_1_redir)
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

    .select-state {
        max-width: 270px;
        margin: 0 auto !important;
        margin-bottom: 1em !important;
        margin-top: 0.25em !important;
    }

    .request-code {
        margin-bottom: 10px !important;
    }

    .no-margin-top {
        margin-top: 0px;
    }
</style>