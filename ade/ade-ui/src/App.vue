<template>
  <div id="app" class="ui container">
    <div v-show="loading" class="ui active dimmer">
      <div class="ui loader"></div>
    </div>
    <h1 v-show="!loading" class="ui header center aligned white">{{hotspot_name}}</h1>
    <img v-show="!loading" :src="hotspot_logo" class="ui image centered small">
    <div v-show="!loading" class="ui segments">
      <router-view class="ui segment"></router-view>
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
      hotspot_logo: ""
    };
  },
  mounted() {
    $("body").css("background-color", "rgba(0,0,0,.85)");
    this.getInfo();
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
          $("body").css("background-color", success.body.bg_color);
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
