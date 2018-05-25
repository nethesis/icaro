<template>
    <div class="ui segment">
        <div v-if="!authorized && !dedaloError" class="ui active centered inline text loader">{{ $t("social.auth_progress") }}...</div>
        <div v-if="authorized" class="ui icon positive message">
            <i class="check icon"></i>
            <div class="content">
                <div class="header">
                    {{ $t("social.auth_success") }}
                </div>
                <p>{{ $t("social.auth_success_sub") }}...</p>
            </div>
        </div>
        <div v-if="dedaloError" class="ui icon negative message">
            <i class="remove icon"></i>
            <div class="content">
                <div class="header">
                    {{ $t("social.auth_error") }}
                </div>
                <p v-html="$t('social.auth_error_sub')"></p>
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
</template>

<script>
    import AuthMixin from './../../mixins/auth';
    import {
        setTimeout
    } from 'timers';
    export default {
        name: 'InstagramPage',
        mixins: [AuthMixin],
        data() {
            var authorized = false
            var dedaloError = false

            // extract code and state
            var params = this.extractParams()

            if (params.code && params.state) {
                // extract wings preferences
                this.getPreferences(this.parseState(params.state), success => {
                    this.$parent.hotspot.name = success.body.hotspot_name
                    this.$parent.hotspot.disclaimers = success.body.disclaimers
                    this.$parent.hotspot.preferences = success.body.preferences
                    this.$root.$options.hotspot.disclaimers = success.body.disclaimers
                    this.$root.$options.hotspot.preferences = success.body.preferences
                    this.hotspot.disclaimers = success.body.disclaimers
                    $("body").css("background-color", success.body.preferences.captive_7_background ||
                        '#2a87be');
                }, error => {
                    this.authorized = false
                    console.error(error)
                })

                // make request to wax
                var url = this.createWaxURL(params.code, this.parseState(params.state), 'social/instagram')

                // get user id
                this.$http.get(url).then(responseAuth => {
                    this.userId = responseAuth.body.user_db_id
                    // exec dedalo login
                    this.doDedaloLogin({
                        id: responseAuth.body.user_id,
                        password: responseAuth.password || ''
                    }, responseDedalo => {
                        if (responseDedalo.body.clientState == 1) {
                            this.authorized = true
                            this.dedaloError = false
                        } else {
                            this.authorized = false
                            this.dedaloError = true
                        }
                    }, error => {
                        this.authorized = false
                        this.dedaloError = true
                        console.error(error)
                    })
                }, error => {
                    this.authorized = false
                    this.dedaloError = true
                    console.error(error)
                });
            } else {
                // get social login url
                params.in_client_id = this.$root.$options.hotspot.socials.instagram_client_id
                var url = this.getSocialLoginURL(params, 'instagram')

                // open social login url
                window.location.replace(url)
            }
            return {
                authorized: authorized,
                dedaloError: dedaloError,
                hotspot: {
                    disclaimers: this.$root.$options.hotspot.disclaimers
                },
                userId: 0
            }
        },
        methods: {
            deleteInfo() {
                // extract code and state
                var params = this.extractParams()
                this.deleteMarketingInfo(this.userId, this.parseState(params.state), function (success) {
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

    .text-center {
        text-align: center;
    }

    textarea {
        min-height: 150px !important;
    }
</style>