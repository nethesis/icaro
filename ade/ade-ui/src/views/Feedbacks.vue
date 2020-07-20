<template>
  <div>
    <h2 :style="textStyle" class="ui header center aligned">{{$t('feedbacks')}}</h2>

    <div v-show="feedbackLeave" class="thank-you">
      <div class="content text-center">
        <div class="header" :style="textStyle">{{$t('thank_you_feedback')}}</div>
        <p :style="textStyle">{{$t('thank_you_feedback_text')}}</p>
      </div>
    </div>

    <h4 v-show="!feedbackLeave" class="ui header center aligned" :style="textStyle">{{$t('feedback_text')}}</h4>
    <div v-show="!feedbackLeave" class="ui form">
      <div class="field">
        <textarea :placeholder="$t('leave_blank_mandatory')" v-model="message"></textarea>
        <span class="counter" :style="textStyle">{{message.length}}/400</span>
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
      feedbackLeave: false,
      textColor: '#4A4A4A'
    };
  },
  mounted() {
    this.getInfo();
  },
  computed: {
    textStyle: function () {
      return {
        color: this.textColor
      }
    }
  },
  methods: {
    getInfo() {
      this.loading = true;
      this.getTokenInfo(
        "feedbacks",
        this.$route.params.token,
        function(success) {
          this.loading = false;
          this.hotspot_name = success.body.hotspot_name;
          this.hotspot_logo = success.body.hotspot_logo;
          $("body").css("background-color", success.body.bg_color);
          this.textColor = success.body.text_color || '#4A4A4A'
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
.thank-you {
  margin-bottom: 25px;
}
</style>
