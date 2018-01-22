<template>
    <div class="ui segment">
        <div v-if="!authorized" class="ui active centered inline text loader">Authorization in progress...</div>
        <div v-if="authorized" class="ui icon positive message">
            <i class="check icon"></i>
            <div class="content">
                <div class="header">
                    You are successfully authenticated
                </div>
                <p>In a few seconds you will be redirected...</p>
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
        name: 'LinkedInPage',
        mixins: [AuthMixin],
        data() {
            var authorized = false

            // extract code and state
            var params = this.extractParams()

            if (params.code && params.state) {
                // extract wings preferences
                this.getPreferences(this.parseState(params.state), success => {
                    this.$parent.hotspot.name = success.body.hotspot_name
                    this.$parent.hotspot.preferences = success.body.preferences
                    this.$root.$options.hotspot.preferences = success.body.preferences
                }, error => {
                    this.authorized = false
                    console.error(error)
                })

                // make request to wax
                var url = this.createWaxURL(params.code, this.parseState(params.state), 'social/linkedin')

                // get user id
                this.$http.get(url).then(responseAuth => {
                    // exec dedalo login
                    this.doDedaloLogin({
                        id: responseAuth.body.user_id,
                        password: responseAuth.password || ''
                    }, responseDedalo => {
                        if (responseDedalo.body.clientState == 1) {
                            this.authorized = true
                            setTimeout(function () {
                                // open redir url
                                window.location.replace(this.$root.$options.hotspot.preferences
                                    .captive_redir)
                            }.bind(this), 2500)
                        } else {
                            this.authorized = false
                        }
                    }, error => {
                        this.authorized = false
                        console.error(error)
                    })
                }, response => {
                    this.authorized = false
                    console.error(response)
                });
            } else {
                // get social login url
                params.fb_client_id = this.$root.$options.hotspot.socials.linkedin_client_id
                var url = this.getSocialLoginURL(params, 'linkedin')

                // open social login url
                window.location.replace(url)
            }
            return {
                authorized: authorized,
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
</style>