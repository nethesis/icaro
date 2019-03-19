<template>
  <div :id="'editor-preview-'+id" class="editor-preview">
    <div class="title">
      <h1>{{hotspot.title}}</h1>
      <img :src="hotspot.logo">
    </div>
    <div class="ui container" v-html="html"></div>
    <h5>{{ hotspot.business_name + " • " + hotspot.business_address + " • " + hotspot.business_email }}</h5>
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
  color: black !important;
  border: 1px solid #bbbbbb;
  padding: 15px;
}

h1 {
  color: white;
  border: none;
  padding: 0 0;
  font-weight: 700;
  text-transform: none;
}

h3 {
  color: #5a5a5a !important;
  font-size: 16px !important;
}

h5 {
  color: white;
  border: none;
  padding: 0 0;
  font-weight: 700;
  text-transform: none;
}

p {
  color: #5a5a5a !important;
  margin: 0 0 1em;
  line-height: 1.4285em;
  font-size: 12px;
}

img {
  margin-top: 20px !important;
  margin-bottom: 20px !important;
  width: 150px;
}

.ui.button {
  font-size: 15px;
}

.ui.small.image {
  width: 450px !important;
}
.ui.container {
  background: white;
  margin-left: 10px !important;
  margin-right: 10px !important;
  border-radius: 0.28571429rem;
  padding: 10px;
}
</style>
<style scoped src="semantic-ui-offline/semantic.min.css"></style>
