<template>
  <div>
    <h1 class="ui header center aligned">{{$t('not_found')}}</h1>
    <h4 class="ui header center aligned">{{$t('not_found_text')}}</h4>
  </div>
</template>

<script>
export default {
  name: "NotFound",
  mounted() {
    $("body").css("background-color", "rgba(0,0,0,.85)");
  }
};
</script>

<style scoped>
</style>
