<template>
  <div>
    <h2 class="ui header center aligned">{{$t('reviews')}}</h2>

    <!-- STARS -->
    <h4 v-show="!reviewLeave" class="ui header center aligned">{{$t('review_text_before')}}</h4>
    <div v-show="!reviewLeave" class="ui form">
      <div class="ui one column stackable center aligned page grid field">
        <div class="ui column twelve wide">
          <i
            @click="setStars(1)"
            :class="['star', stars >= 1 ? 'yellow' : 'outline', 'icon big pointer']"
          ></i>
          <i
            @click="setStars(2)"
            :class="['star', stars >= 2 ? 'yellow' : 'outline', 'icon big pointer']"
          ></i>
          <i
            @click="setStars(3)"
            :class="['star', stars >= 3 ? 'yellow' : 'outline', 'icon big pointer']"
          ></i>
          <i
            @click="setStars(4)"
            :class="['star', stars >= 4 ? 'yellow' : 'outline', 'icon big pointer']"
          ></i>
          <i
            @click="setStars(5)"
            :class="['star', stars >= 5 ? 'yellow' : 'outline', 'icon big pointer']"
          ></i>
        </div>
      </div>
    </div>
    <!-- END STARS -->
    <!-- TEXTAREA -->
    <h4
      v-show="!reviewLeave && stars > 0"
      class="ui header center aligned"
    >{{$t('review_text_after_1')}}</h4>
    <h5
      v-show="!reviewLeave && stars > 0"
      class="ui header center aligned no-mg-top"
    >{{$t('review_text_after_2')}}</h5>
    <div v-show="!reviewLeave && stars > 0" class="ui form">
      <div class="ui field">
        <textarea :placeholder="$t('leave_blank_optional')" v-model="message"></textarea>
      </div>
    </div>
    <!-- END TEXTAREA -->
    <!-- SUBMIT -->
    <div class="ui form">
      <div class="ui one column stackable center aligned page grid field"></div>
      <div
        v-show="!reviewLeave && stars > 0"
        class="ui one column stackable center aligned page grid"
      >
        <div class="column twelve wide">
          <button @click="setReview()" class="ui button large green">{{$t('submit')}}</button>
        </div>
      </div>
    </div>
    <!-- END SUBMIT -->
    <!-- THANKS -->
    <div v-show="reviewLeave" class="ui message success">
      <div class="content text-center">
        <div class="header">{{$t('thank_you_review')}}</div>
        <p>{{$t('thank_you_review_text')}}.</p>
      </div>
    </div>
    <!-- END THANKS -->
    <!-- EXTERNAL -->
    <div v-show="reviewLeave && urls.length > 0 && stars > threshold" class="ui form">
      <h4 class="ui header center aligned">{{$t('review_good_external')}}</h4>
      <div class="ui center aligned field grid">
        <a
          target="_blank"
          :href="u"
          v-for="u in urls"
          :key="u"
          :class="['ui', parseUrl(u, true), 'button']"
        >
          <i :class="[parseUrl(u, false), 'icon']"></i>
          {{parseUrl(u, false, true) | capitalize}}
        </a>
      </div>
    </div>
    <!-- END EXTERNAL -->
  </div>
</template>

<script>
var psl = require("psl");

export default {
  name: "Reviews",
  data() {
    return {
      loading: true,
      stars: 0,
      message: null,
      reviewLeave: false,
      urls: [],
      threshold: 3
    };
  },
  mounted() {
    this.getInfo();
  },
  methods: {
    getInfo() {
      this.loading = true;
      this.getTokenInfo(
        "reviews",
        this.$route.params.token,
        function(success) {
          this.loading = false;
          this.hotspot_name = success.hotspot_name;
          this.hotspot_logo = success.hotspot_logo;
          this.urls = success.body.urls || [];
          this.threshold = success.body.threshold || 3;
          $("body").css("background-color", success.bg_color);
        },
        function(error) {
          this.loading = false;
          this.$router.push("/not-found");
        }
      );
    },
    setStars(value) {
      this.stars = value;
    },
    parseUrl(u, color, name) {
      var url = new URL(u);
      var parsed = psl.parse(url.hostname);
      var sld = parsed.sld;

      if (name) {
        if (sld == "tripadvisor") {
          sld = "TripAdvisor";
        }
        return sld;
      } else {
        switch (sld) {
          case "facebook":
            return "facebook";
            break;
          case "google":
            return color ? "google red" : "google";
            break;
          case "tripadvisor":
            return color ? "tripadvisor teal" : "tripadvisor";
            break;
          default:
            return color ? "pencil alternate orange" : "pencil alternate";
            break;
        }
      }
    },
    setReview() {
      this.setTokenReview(
        {
          stars: this.stars,
          message: this.message
        },
        this.$route.params.token,
        function(success) {
          this.reviewLeave = true;
        },
        function(error) {
          this.reviewLeave = true;
        }
      );
    }
  }
};
</script>

<style scoped>
.star:hover {
  color: #ffdc77 !important;
}
</style>
