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
    getLoginIcon(accountType) {
      var icon = 'fa fa-user'
      switch (accountType) {
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
    getUserTypeIcon(userType) {
      var icon = 'fa fa-user'
      switch (userType) {
        case 'facebook':
          icon = 'fa fa-facebook'
          break;
        case 'google':
          icon = 'fa fa-google'
          break;
        case 'linkedin':
          icon = 'fa fa-linkedin'
          break;
        case 'email':
          icon = 'fa fa-envelope'
          break;
        case 'mobile-alt':
          icon = 'fa fa-envelope'
          break;
      }
      return icon
    },
    generatePassword() {
      var length = 8,
        charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789",
        retVal = "";
      for (var i = 0, n = charset.length; i < length; ++i) {
        retVal += charset.charAt(Math.floor(Math.random() * n));
      }
      return retVal;
    },
    getInputType(value) {
      if (value === "true" || value === "false") {
        return 'checkbox'
      } else {
        return 'text'
      }
    }
  }
};
export default UtilService;
