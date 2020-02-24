<template>
  <div id="app" class="ui segment container" :style="containerStyle">
    <div v-show="loading" class="ui active dimmer">
      <div class="ui loader"></div>
    </div>
    <h1 v-show="!loading" class="ui header center aligned white" :style="titleStyle">{{hotspot_name}}</h1>
    <img v-show="!loading" :src="hotspot_logo" class="ui image centered small">
    <div v-show="!loading" class="ui">
      <router-view class="ui"></router-view>
    </div>
  </div>
</template>

<script>
export default {
  name: "Feedback",
  data() {
    return {
      loading: false,
      hotspot_name: "",
      hotspot_logo: "",
      titleColor: '#4A4A4A',
      textFont: 'Roboto',
      containerBgColor: '#ffffffdb'
    };
  },
  created() {
    document.title = CONFIG.APP_NAME;
  },
  mounted() {
    $("body").css("background-color", "rgba(0,0,0,.85)");
    this.getInfo();
  },
  computed: {
    titleStyle: function () {
      return {
        color: this.titleColor,
        'font-family': this.textFont
      }
    },
    containerStyle: function () {
      return {
        'background-color': this.containerBgColor
      }
    },
  },
  methods: {
    getInfo() {
      this.loading = true;
      this.getTokenInfo(
        this.$route.name,
        this.$route.params.token,
        function(success) {
          this.loading = false;
          this.hotspot_name = success.body.hotspot_name;
          this.hotspot_logo = success.body.hotspot_logo;
          $("body").css("background-color", success.body.bg_color || '#2a87be');

          // background image
          if (success.body.bg_image) {
            $("body").css("height", "100vh");
            $("body").css("background-size", "cover");
            $("body").css("background-position", "center");
            $("body").css("background-image", 'url("' + success.body.bg_image + '")');
          }
          this.titleColor = success.body.title_color || '#4A4A4A'
          this.textFont = success.body.text_style || 'Roboto';
          this.containerBgColor = success.body.container_bg_color || '#ffffffdb';
        },
        function(error) {
          this.loading = false;
          this.$router.push("/not-found");
        }
      );
    }
  }
};
</script>

<style>
.white {
  color: white !important;
}
.pointer {
  cursor: pointer !important;
}
.text-center {
  text-align: center !important;
}
.no-mg-top {
  margin-top: 0px !important;
}
</style>
