<template>
  <div>
    <h2 :style="textStyle" class="ui header center aligned">{{$t('feedbacks')}}</h2>

    <div v-show="feedbackLeave" class="ui message success">
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
            :style="buttonStyle"
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
      textColor: '#4A4A4A',
      textFont: 'Roboto'
    };
  },
  mounted() {
    this.getInfo();
  },
  computed: {
    textStyle: function () {
      return {
        color: this.textColor,
        'font-family': this.textFont
      }
    },
    buttonStyle: function () {
      return {
        'font-family': this.textFont
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
          this.hotspot_name = success.hotspot_name;
          this.hotspot_logo = success.hotspot_logo;
          $("body").css("background-color", success.bg_color);
          this.textColor = success.text_color || '#4A4A4A';
          this.textFont = success.text_style || 'Roboto';
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
