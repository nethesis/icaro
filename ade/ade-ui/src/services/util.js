var UtilService = {
    methods: {
        getTokenInfo(type, token, success, error) {
            this.$http
                .get(this.$root.$options.api_scheme +
                    this.$root.$options.api_host + '/ade/' + type + '/' + token)
                .then(success, error);
        },
        setTokenFeedback(body, token, success, error) {
            this.$http
                .post(
                    this.$root.$options.api_scheme +
                    this.$root.$options.api_host +
                    "/ade/feedbacks/" + token, body
                )
                .then(success, error);
        },
        setTokenReview(body, token, success, error) {
            this.$http
                .post(
                    this.$root.$options.api_scheme +
                    this.$root.$options.api_host +
                    "/ade/reviews/" + token, body
                )
                .then(success, error);
        },
    }
};
export default UtilService;