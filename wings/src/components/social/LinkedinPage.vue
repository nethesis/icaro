<template>
    <div class="ui segment">
        <p>
            <b>Code</b>: {{code}}</p>
    </div>
</template>

<script>
    import AuthMixin from './../../mixins/auth';
    export default {
        name: 'LinkedinPage',
        mixins: [AuthMixin],
        data() {
            // extract code and state from url
            var params = this.extractParams()

            if (params.code && params.state) {
                // extract wings preferences
                this.getPreferences(this.parseState(params.state), success => {
                    this.$parent.hotspot.name = success.body.hotspot_name
                }, error => {
                    console.error(error)
                })

                // make request to wax
                var url = this.createWaxURL(params.code, params.state, 'social/linkedin')

                // get user id
                this.$http.get(url).then(responseAuth => {
                    // exec dedalo login
                    this.doDedaloLogin({
                        id: responseAuth.body.user_id,
                        password: responseAuth.password || ''
                    }, success => {
                        console.log(success)
                    }, error => {
                        console.error(error)
                    })
                }, response => {
                    console.error(response)
                });
            } else {
                // get social login url
                params.li_client_id = '78vuct5wpb9uki'
                var url = this.getSocialLoginURL(params, 'linkedin')

                // open social login url
                window.location.replace(url)
            }
            return {
                code: params.code,
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