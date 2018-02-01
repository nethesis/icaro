var UtilService = {
  methods: {
    tableLangs() {
      return {
        nextText: this.$i18n.t('next'),
        prevText: this.$i18n.t('prev'),
        ofText: this.$i18n.t('of'),
        rowsPerPageText: this.$i18n.t('rows_per_page'),
        globalSearchPlaceholder: this.$i18n.t('search'),
      }
    },
  }
};
export default UtilService;
