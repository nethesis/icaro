<template>
    <div class="ui segment form">
        <div v-if="voucherAvailable && !voucherValidated && !authorized && !dedaloError">
            <h3>{{ $t("login.voucher_title") }}</h3>
            <div class="inline field" v-bind:class="{ error: badInput }">
                <label>Voucher</label>
                <div class="ui big left icon input">
                    <input v-model="authCode" type="email" :placeholder="$t('login.insert_voucher')">
                    <i class="braille icon"></i>
                </div>
            </div>
            <button v-on:click="validateCode()" class="ui big button">{{ $t("login.validate_code") }}</button>
            <div v-if="badCode" class="ui tiny icon negative message">
                <i class="remove icon"></i>
                <div class="content">
                    <div class="header">
                        {{ $t("login.error_voucher_code") }}
                    </div>
                    <p>{{ $t("login.error_voucher_code_sub") }}</p>
                </div>
            </div>
        </div>
        <div v-if="!voucherAvailable || (voucherAvailable && voucherValidated) && !authorized && !dedaloError">
            <h3>{{ $t("login.choose_login") }}</h3>
            <div class="ui relaxed list">
                <div v-if="hotspot.preferences.facebook_login == 'true'" class="item">
                    <div @click="changeRoute('/login/facebook', false)" class="ui facebook button big fluid">
                        <i class="facebook icon"></i>
                        Facebook
                    </div>
                </div>
                <div v-if="hotspot.preferences.instagram_login == 'true'" class="item">
                    <div @click="changeRoute('/login/instagram', false)" class="ui instagram button big fluid">
                        <i class="instagram icon"></i>
                        Instagram
                    </div>
                </div>
                <div v-if="hotspot.preferences.linkedin_login == 'true'" class="item">
                    <div @click="changeRoute('/login/linkedin', false)" class="ui linkedin button big fluid">
                        <i class="linkedin icon"></i>
                        LinkedIn
                    </div>
                </div>
            </div>
            <div class="ui divider"></div>
            <div class="ui relaxed list">
                <div v-if="hotspot.preferences.sms_login == 'true'" class="item">
                    <div @click="changeRoute('/login/sms', false)" class="ui button green big fluid">
                        <i class="talk icon"></i>
                        SMS
                    </div>
                </div>
                <div v-if="hotspot.preferences.email_login == 'true'" class="item">
                    <div @click="changeRoute('/login/email', false)" class="ui button yellow big fluid">
                        <i class="mail icon"></i>
                        Email
                    </div>
                </div>
            </div>
            <div v-if="hotspot.preferences.voucher_code_login == 'true'" class="ui divider"></div>
            <div class="ui relaxed list">
                <div v-if="hotspot.preferences.voucher_code_login == 'true'" class="item">
                    <div @click="changeRoute('/login', true)" class="ui button teal big fluid">
                        <i class="barcode icon"></i>
                        {{ $t("login.code") }}
                    </div>
                </div>
            </div>
        </div>
        <div v-if="dedaloError" class="ui icon negative message">
            <div class="content">
                <div class="header">
                    {{ $t("social.auth_error") }}
                </div>
                <p>{{ $t("social.auth_error_sub") }}</p>
            </div>
        </div>
        <button v-if="dedaloError" v-on:click="back()" class="ui big button green">{{ $t("login.back") }}</button>

        <div v-if="authorized" class="ui icon positive message">
            <i class="check icon"></i>
            <div class="content">
                <div class="header">
                    {{ $t("social.auth_success") }}
                </div>
                <p>{{ $t("social.auth_success_sub") }}...</p>
            </div>
        </div>
        <div v-if="authorized">
            <h3>{{ $t("login.disclaimer_marketing") }}</h3>
            <div class="inline field">
                <textarea readonly class="text-center" v-model="hotspot.disclaimers.marketing_use"></textarea>
            </div>
            <div class="ui checkbox">
                <input v-model="conditions" type="checkbox" name="example">
                <label>{{ $t("login.disclaimer_privacy_accept") }}</label>
            </div>
            <br />
            <button v-on:click="navigate()" class="ui big button green adjust-top">{{ $t("login.navigate") }}</button>
        </div>
    </div>
</template>

<script>
    import AuthMixin from "./../mixins/auth";
    export default {
        name: "LoginPage",
        mixins: [AuthMixin],
        data: function () {
            var voucherAvailable = false;
            var voucherValidated = false;
            var badCode = false;
            var badInput = false;
            var authorized = false;
            var dedaloError = false;

            var authCode = "";
            if (this.$root.$options.hotspot.preferences.voucher_login == "true") {
                voucherAvailable = true;
            }
            return {
                hotspot: {
                    preferences: this.$root.$options.hotspot.preferences,
                    disclaimers: this.$root.$options.hotspot.disclaimers
                },
                voucherAvailable: voucherAvailable,
                voucherValidated: voucherValidated,
                authCode: authCode,
                badCode: badCode,
                badInput: badInput,
                authorized: authorized,
                dedaloError: dedaloError,
                userId: 0,
                conditions: false
            };
        },
        methods: {
            back() {
                this.voucherAvailable = this.$root.$options.hotspot.preferences.voucher_login
                this.voucherValidated = false
                this.authCode = ""
                this.badCode = false
                this.badInput = false
                this.authorized = false
                this.dedaloError = false
            },
            changeRoute: function (path, withCode) {
                this.$root.$options.session["loginDest"] = path;
                this.$root.$options.hotspot.preferences.voucher_login = withCode.toString()
                this.$router.push({
                    path: "login/disclaimer"
                });
            },
            validateCode: function () {
                this.badCode = false;
                if (this.authCode.length == 0) {
                    this.badInput = true;
                    return;
                }
                var params = this.extractParams();

                // make request to wax
                var url = this.createWaxURL(this.authCode, params, "voucher");

                // get user id
                this.$http.get(url).then(
                    function (responseAuth) {
                        this.$root.$options.session["voucherCode"] = responseAuth.body.code;
                        this.userId = responseAuth.body.user_db_id

                        if (responseAuth.body.type == "auth") {
                            // do login immediately
                            var context = this
                            this.doDedaloLogin({
                                    id: responseAuth.body.code,
                                    password: responseAuth.body.code || ""
                                },
                                function (responseDedalo) {
                                    if (responseDedalo && responseDedalo.body && responseDedalo.body.clientState ==
                                        1) {
                                        context.authorized = true;
                                        context.dedaloError = false;
                                    } else {
                                        context.authorized = false;
                                        context.dedaloError = true;
                                    }
                                },
                                function (error) {
                                    this.authorized = false;
                                    this.dedaloError = true;
                                    console.error(error);
                                }
                            );
                        } else {
                            this.voucherValidated = true;
                        }
                    },
                    function (error) {
                        this.voucherValidated = false;
                        this.badCode = true;
                        console.error(error);
                    }
                );
            },
            navigate() {
                if (this.conditions) {
                    this.accept()
                } else {
                    this.deleteInfo()
                }
            },
            deleteInfo: function () {
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
            accept: function () {
                // open redir url
                window.location.replace(this.$root.$options.hotspot.preferences
                    .captive_1_redir)
            }
        }
    };
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

    .adjust-top {
        margin-top: 10px;
    }
</style>