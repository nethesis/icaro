<template>
  <div>
    <h2>{{ msg }}</h2>
    <a role="menuitem" tabindex="-1" href="" data-toggle="modal" data-target="#HSmodifyModal">{{ $t('modify') }}</a>
  </div>
</template>

<script>
  export default {
    name: 'HotspotDetails',
    data() {
      return {
        msg: 'Hotspot ' + this.$route.params.id + ' details'
      }
    }
  }
</script>

<style scoped>

</style>