<template>
    <div class="ui segment">
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
</template>

<script>
    export default {
        name: 'LoginPage',
        data() {
            return {
                hotspot: {
                    preferences: this.$root.$options.hotspot.preferences
                }
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
</style>