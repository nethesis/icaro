<template>
  <div id="captive-preview">
    <div class="ui segment container">
      <div class="ui center aligned">
        <h2>{{extractValue('captive_2_title')}}</h2>
        <img :src="extractValue('captive_3_logo')" class="ui centered image tiny">
      </div>
      <div class="ui route-container">
        <div class="ui">
          <div>
            <div>
              <h3>{{extractValue('captive_4_subtitle')}}</h3>
              <img class="ui centered small image" :src="extractValue('captive_5_banner')">
              <p v-html="extractValue('captive_6_description')"></p>
              <div class="ui divider"></div>
              <div class="ui checkbox">
                <input type="checkbox" id="agree-terms" />
                <label for="agree-terms" class="inline">
                  <span class="terms-space">{{$t('hotspot.i_agree')}}</span>
                </label><span class="terms-space">(<a class="terms">{{$t('hotspot.tos_and_privacy')}}</a>)</span>
              </div>
              <div class="ui divider"></div>
              <router-link to="/login" class="big ui green button">{{ $t("hotspot.start_navigate") }}</router-link>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
<script>
export default {
  name: "CaptivePortal",
  props: ["obj"],
  data() {
    return {};
  },
  methods: {
    extractValue(key) {
      for (var p in this.obj) {
        var pref = this.obj[p];
        if (pref.key == key) {
          return pref.value;
        }
      }
    }
  }
};
</script>

<style scoped>
#captive-preview {
  -webkit-font-smoothing: antialiased !important;
  -moz-osx-font-smoothing: grayscale !important;
  text-align: center !important;
  color: #eff7fc !important;
  margin-bottom: 30px !important;
  padding: 25px 65px;
}

h3 {
  color: #4A4A4A;
  font-size: 16px !important;
}

p {
  color: #4A4A4A;
  margin: 0 0 1em;
  line-height: 1.4285em;
  font-size: 12px;
}

.terms-space {
  margin-right: 3px;
  cursor: default;
}


a.terms {
  text-decoration: underline !important;
  cursor: pointer;
}

img {
  margin-top: 20px !important;
  margin-bottom: 20px !important;
}

.ui.button {
  font-size: 15px;
}

.ui.small.image {
  width: 450px !important;
}

.ui.segment {
  border: none !important;
  box-shadow: none !important;
  -webkit-box-shadow: none !important;
}

label.inline {
  display: inline !important;
}
</style>
<style scoped src="semantic-ui-offline/semantic.min.css"></style>
