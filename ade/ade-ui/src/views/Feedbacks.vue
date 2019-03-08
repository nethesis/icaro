<template>
  <div>
    <h2 class="ui header center aligned">{{$t('feedbacks')}}</h2>

    <div v-show="feedbackLeave" class="ui message success">
      <div class="content text-center">
        <div class="header">{{$t('thank_you_feedback')}}</div>
        <p>{{$t('thank_you_feedback_text')}}.</p>
      </div>
    </div>

    <h4 v-show="!feedbackLeave" class="ui header center aligned">{{$t('feedback_text')}}</h4>
    <div v-show="!feedbackLeave" class="ui form">
      <div class="field">
        <textarea :placeholder="$t('leave_blank_mandatory')" v-model="message"></textarea>
        <span class="counter">{{message.length}}/400</span>
      </div>
      <div class="ui one column stackable center aligned page grid">
        <div class="column twelve wide">
          <button
            :disabled="message.length < 5 || message.length > 400"
            @click="setFeedback()"
            class="ui button large green"
          >{{$t('submit')}}</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: "Feedback",
  data() {
    return {
      loading: false,
      message: "",
      feedbackLeave: false
    };
  },
  mounted() {
    this.getInfo();
  },
  methods: {
    getInfo() {
      this.loading = true;
      this.getTokenInfo(
        "feedbacks",
        this.$route.params.token,
        function(success) {
          this.loading = false;
          this.hotspot_name = success.hotspot_name;
          this.hotspot_logo = success.hotspot_logo;
          $("body").css("background-color", success.bg_color);
        },
        function(error) {
          this.loading = false;
          this.$router.push("/not-found");
        }
      );
    },
    setFeedback() {
      this.setTokenFeedback(
        {
          message: this.message
        },
        this.$route.params.token,
        function(success) {
          this.feedbackLeave = true;
        },
        function(error) {
          this.feedbackLeave = true;
        }
      );
    }
  }
};
</script>

<style scoped>
.counter {
  text-align: right;
  display: block;
  color: #777777;
}
</style>
