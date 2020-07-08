var UtilService = {
  methods: {
    tableLangs() {
      return {
        nextText: this.$i18n.t("next"),
        prevText: this.$i18n.t("prev"),
        ofText: this.$i18n.t("of"),
        rowsPerPageText: this.$i18n.t("rows_per_page"),
        globalSearchPlaceholder: this.$i18n.t("search_table"),
        allText: this.$i18n.t("all")
      };
    },
    resizeImage(img, maxWidth, maxHeight) {
      var srcWidth = img.width;
      var srcHeight = img.height;

      var resizeWidth = srcWidth;
      var resizeHeight = srcHeight;

      var aspect = resizeWidth / resizeHeight;

      if (resizeWidth > maxWidth) {
        resizeWidth = maxWidth;
        resizeHeight = resizeWidth / aspect;
      }
      if (resizeHeight > maxHeight) {
        aspect = resizeWidth / resizeHeight;
        resizeHeight = maxHeight;
        resizeWidth = resizeHeight * aspect;
      }

      return {
        width: resizeWidth,
        height: resizeHeight
      };
    },
    uploadImageLangs() {
      return {
        drag: "<h1>" + this.$i18n.t("upload_drag") + "</h1>",
        change: this.$i18n.t("upload_change"),
        remove: this.$i18n.t("upload_remove"),
        fileType: this.$i18n.t("upload_file_not_supported"),
        fileSize: this.$i18n.t("upload_file_exceed")
      };
    },
    generateUUID() {
      var d = new Date().getTime();
      var uuid = "xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx".replace(
        /[xy]/g,
        function(c) {
          var r = (d + Math.random() * 16) % 16 | 0;
          d = Math.floor(d / 16);
          return (c == "x" ? r : (r & 0x3) | 0x8).toString(16);
        }
      );
      return uuid;
    },
    getLoginIcon(accountType) {
      var icon = "fa fa-user";
      switch (accountType) {
        case "admin":
          icon = "fa fa-graduation-cap";
          break;
        case "reseller":
          icon = "fa fa-user";
          break;
        case "customer":
          icon = "fa fa-briefcase";
          break;
        case "desk":
          icon = "fa fa-coffee";
          break;
      }
      return icon;
    },
    getUserTypeIcon(userType) {
      var icon = "fa fa-user";
      switch (userType) {
        case "facebook":
          icon = "fa fa-facebook";
          break;
        case "instagram":
          icon = "fa fa-instagram";
          break;
        case "linkedin":
          icon = "fa fa-linkedin";
          break;
        case "email":
          icon = "fa fa-envelope";
          break;
        case "sms":
          icon = "fa fa-commenting";
          break;
      }
      return icon;
    },
    getPrefTypeIcon(prefType) {
      var icon = "";
      switch (prefType) {
        case "facebook_login":
          icon = "fa fa-facebook-square login-pref-option";
          break;
        case "facebook_login_page":
          icon = "fa fa-thumbs-up login-pref-option";
          break;
        case "instagram_login":
          icon = "fa fa-instagram login-pref-option";
          break;
        case "linkedin_login":
          icon = "fa fa-linkedin-square login-pref-option";
          break;
        case "email_login":
          icon = "fa fa-envelope-square login-pref-option";
          break;
        case "sms_login":
          icon = "fa fa-commenting login-pref-option";
          break;
        case "voucher_login":
          icon = "fa fa-wpforms login-pref-option";
          break;
        case "temp_code_login":
          icon = "fa fa-barcode login-pref-option";
          break;
        case "CoovaChilli-Bandwidth-Max-Down":
          icon = "fa fa-arrow-circle-down login-pref-option";
          break;
        case "CoovaChilli-Bandwidth-Max-Up":
          icon = "fa fa-arrow-circle-up login-pref-option";
          break;
        case "CoovaChilli-Max-Navigation-Time":
          icon = "fa fa-clock-o login-pref-option";
          break;
        case "CoovaChilli-Max-Total-Octets":
          icon = "fa fa-database login-pref-option";
          break;
      }
      return icon;
    },
    generatePassword() {
      var length = 8,
        charset =
          "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789",
        retVal = "";
      for (var i = 0, n = charset.length; i < length; ++i) {
        retVal += charset.charAt(Math.floor(Math.random() * n));
      }
      return retVal;
    },
    getInputType(key, value) {
      var type = "text";
      switch (key) {
        case "facebook_login":
        case "instagram_login":
        case "linkedin_login":
        case "email_login":
        case "sms_login":
        case "voucher_login":
        case "temp_code_login":
        case "auto_login":
        case "auth_renew":
          type = "checkbox";
          break;
        case "sms_login_max":
        case "sms_login_threshold":
        case "temp_session_duration":
        case "user_expiration_days":
        case "voucher_expiration_days":
        case "CoovaChilli-Bandwidth-Max-Down":
        case "CoovaChilli-Bandwidth-Max-Up":
        case "CoovaChilli-Max-Navigation-Time":
        case "CoovaChilli-Max-Total-Octets":
          type = "number";
          break;
        case "captive_1_redir":
          type = "url";
          break;
      }
      return type;
    },
    urltoFile(dataURI, name) {
      var byteString;
      if (dataURI.split(",")[0].indexOf("base64") >= 0)
        byteString = atob(dataURI.split(",")[1]);
      else byteString = unescape(dataURI.split(",")[1]);

      // separate out the mime component
      var mimeString =
        dataURI &&
        dataURI
          .split(",")[0]
          .split(":")[1]
          .split(";")[0];

      // write the bytes of the string to a typed array
      var ia = new Uint8Array(byteString.length);
      for (var i = 0; i < byteString.length; i++) {
        ia[i] = byteString.charCodeAt(i);
      }

      var blob = new Blob([ia], {
        type: mimeString
      });
      return new File([blob], name);
    },
    createCSV(columns, rows) {
      var cs = [];
      for (var c in columns) {
        if (columns[c].field !== "") cs.push(columns[c].label);
      }

      var rs = [];
      for (var r in rows) {
        var row = [];
        for (var c in columns) {
          if (columns[c].field !== "") row.push(rows[r][columns[c].field]);
        }
        rs.push(row);
      }

      return {
        cols: cs,
        rows: rs
      };
    },
    downloadCSV(columns, rows, name) {
      var header = '"' + columns.join('","') + '"' + "\r\n";
      var body = "";
      for (var r in rows) {
        body += '"' + rows[r].join('","') + '"' + "\r\n";
      }
      var finalCSV = "data:text/csv;charset=utf-8," + header + body;
      var encodedUri = encodeURI(finalCSV);
      var link = document.getElementById("download-csv");
      link.setAttribute("href", encodedUri);
      var now = new Date();
      var date =
        now.getFullYear() + "_" + (now.getMonth() + 1) + "_" + now.getDate();
      link.setAttribute("download", name + date + ".csv");
      link.click();
    },
    accountTypeColor(account) {
      switch (account) {
        case "whatsapp":
          return "#21ba45";

        case "facebook":
          return "#3b5998";

        case "instagram":
          return "#49769c";

        case "linkedin":
          return "#1f88be";

        case "email":
          return "#db2828";

        case "sms":
          return "#fbbd08";

        case "mac":
          return "#767676";

        case "voucher":
          return "#00b5ad";
      }
    }
  }
};
export default UtilService;
