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
    generateUUID() {
      var d = new Date().getTime();
      var uuid = 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, function (c) {
        var r = (d + Math.random() * 16) % 16 | 0;
        d = Math.floor(d / 16);
        return (c == 'x' ? r : (r & 0x3 | 0x8)).toString(16);
      });
      return uuid;
    },
    getLoginIcon(userType) {
        var icon = 'fa fa-user'
        switch (userType) {
          case 'admin':
            icon = 'fa fa-graduation-cap'
            break;
          case 'reseller':
            icon = 'fa fa-user'
            break;
          case 'customer':
            icon = 'fa fa-briefcase'
            break;
          case 'desk':
            icon = 'fa fa-coffee'
            break;
        }
        return icon
      },

  }
};
export default UtilService;
