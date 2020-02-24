<template>
  <div :id="'editor-preview-'+id" class="editor-preview">
    <div class="ui container">
      <div class="title">
        <h1>{{hotspot.title}}</h1>
      </div>
      <div v-html="html"></div>
      <h5>{{ hotspot.business_name + " • " + hotspot.business_address + " • " + hotspot.business_email }}</h5>
    </div>
  </div>
</template>
<script>
export default {
  name: "EditoPreview",
  props: ["html", "hotspot", "id"],
  data() {
    return {};
  },
  mounted() {
    $("#editor-preview-" + this.id).css("background-color", this.hotspot.color);
    $(".editor-preview").css("color", this.hotspot.textColor);
    $("h1").css("color", this.hotspot.titleColor);
    $(".ui.container").css("background-color", this.hotspot.containerBgColor);
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
.editor-preview {
  -webkit-font-smoothing: antialiased !important;
  -moz-osx-font-smoothing: grayscale !important;
  text-align: center !important;
  border: 1px solid #bbbbbb;
  padding: 15px;
}

h1 {
  border: none;
  padding: 0 0;
  font-weight: 700;
  text-transform: none;
}

h3 {
  font-size: 16px !important;
}

h5 {
  border: none;
  padding: 0 0;
  font-weight: 700;
  text-transform: none;
}

p {
  margin: 0 0 1em;
  line-height: 1.4285em;
  font-size: 12px;
}

.ui.button {
  font-size: 15px;
}

.ui.small.image {
  width: 450px !important;
}
.ui.container {
  margin-top: 20px !important;
  margin-bottom: 20px !important;
  margin-left: auto !important;
  margin-right: auto !important;
  border-radius: 0.28571429rem;
  padding: 10px;
  width: 70%;
}
</style>
<style scoped src="semantic-ui-offline/semantic.min.css"></style>
