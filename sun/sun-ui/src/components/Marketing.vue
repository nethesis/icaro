<template>
  <div>
    <h2>{{ msg }}</h2>
    <div v-if="isLoading" class="spinner spinner-lg"></div>
    <form class="form-horizontal">
      <div
        v-if="(user.account_type == 'admin') || (user.account_type == 'reseller') && !isLoading"
        class="form-group"
      >
        <label v-if="!isLoading" class="col-sm-3 control-label" for="textInput-markup">Hotspot</label>
        <div v-if="!isLoading" class="col-sm-4">
          <select v-on:change="getMarketingInfo()" v-model="hotspotSearchId" class="form-control">
            <option
              v-for="hotspot in hotspots"
              v-bind:key="hotspot.id"
              v-bind:value="hotspot.id"
            >{{ hotspot.name }} - {{ hotspot.description}}</option>
          </select>
        </div>
      </div>
    </form>

    <form class="form-horizontal" v-on:submit.prevent="saveMarketingInfo()">
      <h3 class="subtitle">{{ $t('marketing.preferences') }}</h3>

      <div v-if="!isLoading && isLoadingMarketing" class="form-group">
        <label class="col-sm-2 control-label" for="textInput-markup">
          <div class="spinner spinner-sm"></div>
        </label>
      </div>

      <div v-if="!isLoadingMarketing" class="form-group">
        <label
          v-if="!isLoadingMarketing"
          class="col-sm-3 control-label"
          for="textInput-markup"
        >{{$t('marketing.ask_for_reason_country')}}</label>
        <div v-if="!isLoadingMarketing" class="col-sm-4">
          <input
            v-model="marketingPrefs.marketing_0_reason_country"
            class="form-control"
            type="checkbox"
          >
        </div>
      </div>

      <div v-if="!isLoadingMarketing" class="form-group">
        <label
          v-if="!isLoadingMarketing"
          class="col-sm-3 control-label"
          for="textInput-markup"
        >{{$t('marketing.enable_service')}}</label>
        <div v-if="!isLoadingMarketing" class="col-sm-4">
          <input v-model="marketingPrefs.marketing_1_enabled" class="form-control" type="checkbox">
        </div>
      </div>

      <div v-if="!isLoadingMarketing && marketingPrefs.marketing_1_enabled" class="form-group">
        <label v-if="!isLoadingMarketing" class="col-sm-3 control-label" for="textInput-markup">
          {{$t('marketing.email_feedback')}}
          <span class="fa fa-envelope-square login-pref-option"></span>
        </label>
        <div v-if="!isLoadingMarketing" class="col-sm-4">
          <input
            :required="marketingPrefs.marketing_1_enabled"
            v-model="marketingPrefs.marketing_2_feedback_email"
            class="form-control"
            type="text"
          >
        </div>
      </div>

      <div v-if="!isLoadingMarketing && marketingPrefs.marketing_1_enabled" class="form-group">
        <label v-if="!isLoadingMarketing" class="col-sm-3 control-label" for="textInput-markup">
          {{$t('marketing.email_first')}}
          <span class="fa fa-send login-pref-option"></span>
        </label>
        <div v-if="!isLoadingMarketing" class="col-sm-4">
          <input
            v-model="marketingPrefs.marketing_3_first_email_enabled"
            @change="handleFeedbackEditor()"
            class="form-control"
            type="checkbox"
          >
        </div>
      </div>

      <div
        v-if="!isLoadingMarketing && marketingPrefs.marketing_3_first_email_enabled && marketingPrefs.marketing_1_enabled"
        class="form-group"
      >
        <label
          v-if="!isLoadingMarketing"
          class="col-sm-3 control-label"
          for="textInput-markup"
        >{{$t('marketing.email_first_after')}}</label>
        <div v-if="!isLoadingMarketing" class="col-sm-4">
          <input
            v-model="marketingPrefs.marketing_6_first_after"
            class="form-control"
            type="number"
          >
        </div>
      </div>

      <br
        v-show="!isLoadingMarketing && marketingPrefs.marketing_3_first_email_enabled && marketingPrefs.marketing_1_enabled"
      >

      <div
        v-show="!isLoadingMarketing && marketingPrefs.marketing_3_first_email_enabled && marketingPrefs.marketing_1_enabled"
        class="form-group"
      >
        <label
          v-show="!isLoadingMarketing"
          class="col-sm-3 control-label"
          for="textInput-markup"
        >{{$t('marketing.email_first_template')}}</label>
        <div v-show="!isLoadingMarketing" class="col-sm-6">
          <editor-menu-bar :editor="feedbackEditor">
            <div class="menubar" slot-scope="{ commands, isActive }">
              <button
                type="button"
                class="btn btn-default"
                :class="{ 'is-active': isActive.bold() }"
                @click="commands.bold"
              >
                <span class="fa fa-bold"/>
              </button>

              <button
                type="button"
                class="btn btn-default"
                :class="{ 'is-active': isActive.italic() }"
                @click="commands.italic"
              >
                <span class="fa fa-italic"/>
              </button>

              <button
                type="button"
                class="btn btn-default space-btn"
                :class="{ 'is-active': isActive.paragraph() }"
                @click="commands.paragraph"
              >
                <span class="fa fa-paragraph"/>
              </button>

              <button
                type="button"
                class="btn btn-default space-btn"
                :class="{ 'is-active': isActive.heading({ level: 1 }) }"
                @click="commands.heading({ level: 1 })"
              >H1</button>

              <button
                type="button"
                class="btn btn-default"
                :class="{ 'is-active': isActive.heading({ level: 2 }) }"
                @click="commands.heading({ level: 2 })"
              >H2</button>

              <button
                type="button"
                class="btn btn-default"
                :class="{ 'is-active': isActive.heading({ level: 3 }) }"
                @click="commands.heading({ level: 3 })"
              >H3</button>
            </div>
          </editor-menu-bar>
          <editor-content id="feedbackEditorContent" :editor="feedbackEditor"/>
        </div>
      </div>

      <div
        v-show="!isLoadingMarketing && marketingPrefs.marketing_3_first_email_enabled && marketingPrefs.marketing_1_enabled"
        class="form-group"
      >
        <label
          v-if="!isLoadingMarketing"
          class="col-sm-3 control-label"
          for="textInput-markup"
        >{{$t('marketing.email_first_template_preview')}}</label>
        <div v-if="!isLoadingMarketing" class="col-sm-6">
          <editor-preview :id="'feedback'" :hotspot="hotspot" :html="feedbackEditorHTML"></editor-preview>
          <button
            @click="testMail('feedback')"
            class="btn btn-primary test-mail"
          >{{$t('marketing.test_mail')}}</button>
          <p
            class="test-mail"
          >{{$t('marketing.test_mail_to')}}: {{marketingPrefs.marketing_2_feedback_email}}</p>
        </div>
      </div>

      <br v-if="marketingPrefs.marketing_3_first_email_enabled">
      <br v-if="marketingPrefs.marketing_3_first_email_enabled">

      <div v-if="!isLoadingMarketing && marketingPrefs.marketing_1_enabled" class="form-group">
        <label v-if="!isLoadingMarketing" class="col-sm-3 control-label" for="textInput-markup">
          {{$t('marketing.email_second')}}
          <span class="fa fa-send login-pref-option"></span>
        </label>
        <div v-if="!isLoadingMarketing" class="col-sm-4">
          <input
            v-model="marketingPrefs.marketing_4_second_email_enabled"
            @change="handleReviewEditor()"
            class="form-control"
            type="checkbox"
          >
        </div>
      </div>

      <div
        v-if="!isLoadingMarketing && marketingPrefs.marketing_4_second_email_enabled && marketingPrefs.marketing_1_enabled"
        class="form-group"
      >
        <label
          class="col-sm-3 control-label"
          for="textInput-modal-markup"
        >{{$t('marketing.email_second_after')}}</label>
        <div class="col-sm-7">
          <span class="span-radio">
            <input
              v-model="marketingPrefs.marketing_7_second_after"
              required
              class="form-check-input"
              type="radio"
              name="timeExpiration"
              id="timeExpiration1"
              value="days"
            >
            <label class="form-check-label" for="timeExpiration1">{{$t('marketing.after_days')}}</label>
          </span>
          <span class="span-radio">
            <input
              v-model="marketingPrefs.marketing_7_second_after"
              required
              class="form-check-input"
              type="radio"
              name="timeExpiration"
              id="timeExpiration2"
              value="expiration"
            >
            <label
              class="form-check-label"
              for="timeExpiration2"
            >{{$t('marketing.after_expiration')}}</label>
          </span>
        </div>
      </div>

      <div
        v-if="!isLoadingMarketing && marketingPrefs.marketing_4_second_email_enabled && marketingPrefs.marketing_7_second_after == 'days' && marketingPrefs.marketing_1_enabled"
        class="form-group"
      >
        <label class="col-sm-3 control-label" for="textInput-modal-markup"></label>
        <div class="col-sm-4">
          <input
            v-model="marketingPrefs.marketing_8_second_after_days"
            type="number"
            id="textInput-modal-markup"
            class="form-control"
          >
        </div>
      </div>

      <br
        v-show="!isLoadingMarketing && marketingPrefs.marketing_4_second_email_enabled && marketingPrefs.marketing_1_enabled"
      >

      <div
        v-show="!isLoadingMarketing && marketingPrefs.marketing_4_second_email_enabled && marketingPrefs.marketing_1_enabled"
        class="form-group"
      >
        <label
          v-show="!isLoadingMarketing"
          class="col-sm-3 control-label"
          for="textInput-markup"
        >{{$t('marketing.email_second_template')}}</label>
        <div v-show="!isLoadingMarketing" class="col-sm-6">
          <editor-menu-bar :editor="reviewEditor">
            <div class="menubar" slot-scope="{ commands, isActive }">
              <button
                type="button"
                class="btn btn-default"
                :class="{ 'is-active': isActive.bold() }"
                @click="commands.bold"
              >
                <span class="fa fa-bold"/>
              </button>

              <button
                type="button"
                class="btn btn-default"
                :class="{ 'is-active': isActive.italic() }"
                @click="commands.italic"
              >
                <span class="fa fa-italic"/>
              </button>

              <button
                type="button"
                class="btn btn-default space-btn"
                :class="{ 'is-active': isActive.paragraph() }"
                @click="commands.paragraph"
              >
                <span class="fa fa-paragraph"/>
              </button>

              <button
                type="button"
                class="btn btn-default space-btn"
                :class="{ 'is-active': isActive.heading({ level: 1 }) }"
                @click="commands.heading({ level: 1 })"
              >H1</button>

              <button
                type="button"
                class="btn btn-default"
                :class="{ 'is-active': isActive.heading({ level: 2 }) }"
                @click="commands.heading({ level: 2 })"
              >H2</button>

              <button
                type="button"
                class="btn btn-default"
                :class="{ 'is-active': isActive.heading({ level: 3 }) }"
                @click="commands.heading({ level: 3 })"
              >H3</button>
            </div>
          </editor-menu-bar>
          <editor-content id="reviewEditorContent" :editor="reviewEditor"/>
        </div>
      </div>

      <div
        v-show="!isLoadingMarketing && marketingPrefs.marketing_4_second_email_enabled && marketingPrefs.marketing_1_enabled"
        class="form-group"
      >
        <label
          v-if="!isLoadingMarketing"
          class="col-sm-3 control-label"
          for="textInput-markup"
        >{{$t('marketing.email_second_template_preview')}}</label>
        <div v-if="!isLoadingMarketing" class="col-sm-6">
          <editor-preview :id="'review'" :hotspot="hotspot" :html="reviewEditorHTML"></editor-preview>
          <button
            @click="testMail('review')"
            class="btn btn-primary test-mail"
          >{{$t('marketing.test_mail')}}</button>
          <p
            class="test-mail"
          >{{$t('marketing.test_mail_to')}}: {{marketingPrefs.marketing_2_feedback_email}}</p>
        </div>
      </div>

      <br v-if="marketingPrefs.marketing_4_second_email_enabled">
      <br v-if="marketingPrefs.marketing_4_second_email_enabled">

      <div v-if="!isLoadingMarketing && marketingPrefs.marketing_1_enabled" class="form-group">
        <label v-if="!isLoadingMarketing" class="col-sm-3 control-label" for="textInput-markup">
          {{$t('marketing.sms')}}
          <span class="fa fa-commenting login-pref-option"></span>
        </label>
        <div v-if="!isLoadingMarketing" class="col-sm-4">
          <input
            v-model="marketingPrefs.marketing_5_sms_enabled"
            class="form-control"
            type="checkbox"
          >
        </div>
      </div>

      <div v-if="!isLoadingMarketing && marketingPrefs.marketing_1_enabled" class="form-group">
        <label
          v-if="!isLoadingMarketing"
          class="col-sm-3 control-label"
          for="textInput-markup"
        >{{$t('marketing.threshold')}}</label>
        <div v-if="!isLoadingMarketing" class="col-sm-4">
          <span
            @click="setStars(1)"
            :class="['fa', marketingPrefs.marketing_9_threshold >= 1 ? 'fa-star star-filled' :  'fa-star-o star']"
          ></span>
          <span
            @click="setStars(2)"
            :class="['fa', marketingPrefs.marketing_9_threshold >= 2 ? 'fa-star star-filled' :  'fa-star-o star']"
          ></span>
          <span
            @click="setStars(3)"
            :class="['fa', marketingPrefs.marketing_9_threshold >= 3 ? 'fa-star star-filled' :  'fa-star-o star']"
          ></span>
          <span
            @click="setStars(4)"
            :class="['fa', marketingPrefs.marketing_9_threshold >= 4 ? 'fa-star star-filled' :  'fa-star-o star']"
          ></span>
          <span
            @click="setStars(5)"
            :class="['fa', marketingPrefs.marketing_9_threshold >= 5 ? 'fa-star star-filled' :  'fa-star-o star']"
          ></span>
        </div>
      </div>

      <div v-if="!isLoadingMarketing && marketingPrefs.marketing_1_enabled" class="form-group">
        <label
          v-if="!isLoadingMarketing"
          class="col-sm-3 control-label"
          for="textInput-markup"
        >{{$t('marketing.first_external_url')}}</label>
        <div v-if="!isLoadingMarketing" class="col-sm-4">
          <input
            placeholder="https://www.tripadvisor.it/MYRESTAURANT"
            v-model="marketingPrefs.marketing_10_first_url"
            class="form-control"
            type="url"
          >
        </div>
      </div>
      <div v-if="!isLoadingMarketing && marketingPrefs.marketing_1_enabled" class="form-group">
        <label
          v-if="!isLoadingMarketing"
          class="col-sm-3 control-label"
          for="textInput-markup"
        >{{$t('marketing.second_external_url')}}</label>
        <div v-if="!isLoadingMarketing" class="col-sm-4">
          <input
            placeholder="https://www.booking.com/hotel/MYHOTEL"
            v-model="marketingPrefs.marketing_11_second_url"
            class="form-control"
            type="url"
          >
        </div>
      </div>
      <div v-if="!isLoadingMarketing && marketingPrefs.marketing_1_enabled" class="form-group">
        <label
          v-if="!isLoadingMarketing"
          class="col-sm-3 control-label"
          for="textInput-markup"
        >{{$t('marketing.third_external_url')}}</label>
        <div v-if="!isLoadingMarketing" class="col-sm-4">
          <input v-model="marketingPrefs.marketing_12_third_url" class="form-control" type="url">
        </div>
      </div>

      <div v-if="!isLoadingMarketing" class="form-group">
        <label v-if="!isLoadingMarketing" class="col-sm-3 control-label" for="textInput-markup"></label>
        <div class="col-sm-3">
          <span
            v-show="marketingPrefs.isLoading"
            class="spinner spinner-sm spinner-inline modal-spinner"
          ></span>
          <button
            :disabled="marketingPrefs.isLoading"
            class="btn btn-primary"
            type="submit"
          >{{$t('marketing.save')}}</button>
        </div>
      </div>
    </form>
  </div>
</template>
<script>
import MarketingService from "../services/marketing";
import UtilService from "../services/util";
import StorageService from "../services/storage";
import HotspotService from "../services/hotspot";
import PreferenceService from "../services/preference";

import EditorPreview from "../directives/EditorPreview.vue";

import { Editor, EditorContent, EditorMenuBar } from "tiptap";
import { Heading, Bold, Italic } from "tiptap-extensions";

export default {
  name: "Unit",
  mixins: [
    HotspotService,
    UtilService,
    MarketingService,
    StorageService,
    PreferenceService
  ],
  components: {
    EditorContent,
    EditorMenuBar,
    EditorPreview
  },
  beforeRouteLeave(to, from, next) {
    if (this.feedbackEditor) {
      this.feedbackEditor.destroy();
      this.feedbackEditor = null;
    }

    if (this.reviewEditor) {
      this.reviewEditor.destroy();
      this.reviewEditor = null;
    }
    next();
  },
  data() {
    return {
      msg: this.$i18n.t("menu.marketing"),
      isLoading: true,
      isLoadingMarketing: true,
      marketingPrefs: {
        marketing_1_enabled: false,
        marketing_2_feedback_email: "",
        marketing_3_first_email_enabled: false,
        marketing_4_second_email_enabled: false,
        marketing_6_first_after: false,
        marketing_7_second_after: "expiration",
        marketing_8_second_after_days: 4,
        marketing_5_sms_enabled: false,
        marketing_9_threshold: 3,
        marketing_10_first_url: "",
        marketing_11_second_url: "",
        marketing_12_third_url: "",
        isLoading: false
      },
      hotspotSearchId: 0,
      hotspots: [],
      hotspot: {},
      user: this.get("loggedUser") || null,
      feedbackEditor: null,
      reviewEditor: null,
      feedbackEditorHTML: "",
      reviewEditorHTML: ""
    };
  },
  mounted() {
    // get hotspots list
    var context = this;
    this.getAllHotspots(function() {
      context.getMarketingInfo();
    });
  },
  methods: {
    handleFeedbackEditor() {
      if (this.marketingPrefs.marketing_3_first_email_enabled) {
        this.feedbackEditor = new Editor({
          extensions: [
            new Heading({ levels: [1, 2, 3] }),
            new Bold(),
            new Italic()
          ],
          content: this.marketingPrefs.marketing_13_feedback_body_text,
          onUpdate: ({ getJSON, getHTML }) => {
            this.feedbackEditorHTML = getHTML();
          }
        });
        this.feedbackEditorHTML = this.marketingPrefs.marketing_13_feedback_body_text;
      } else {
        this.feedbackEditor.destroy();
        this.feedbackEditor = null;
      }
    },
    handleReviewEditor() {
      if (this.marketingPrefs.marketing_4_second_email_enabled) {
        this.reviewEditor = new Editor({
          extensions: [
            new Heading({ levels: [1, 2, 3] }),
            new Bold(),
            new Italic()
          ],
          content: this.marketingPrefs.marketing_14_review_body_text,
          onUpdate: ({ getJSON, getHTML }) => {
            this.reviewEditorHTML = getHTML();
          }
        });
        this.reviewEditorHTML = this.marketingPrefs.marketing_14_review_body_text;
      } else {
        this.reviewEditor.destroy();
        this.reviewEditor = null;
      }
    },
    setStars(value) {
      this.marketingPrefs.marketing_9_threshold = value;
    },
    getAllHotspots(callback) {
      this.hotspotGetAll(
        null,
        null,
        null,
        success => {
          this.hotspots = success.body.data;
          var hsId = this.get("marketing_hotspot_id") || this.hotspots[0].id;
          this.hotspotSearchId = hsId;
          this.isLoading = false;

          callback();
        },
        error => {
          console.error(error);
          callback();
        }
      );
    },
    getMarketingInfo() {
      //get hotspot
      var hsId = this.hotspotSearchId;
      this.hotspot = this.hotspots.filter(function(h) {
        return h.id == hsId;
      })[0];

      this.isLoadingMarketing = true;
      if (this.feedbackEditor) {
        this.feedbackEditor.destroy();
        this.feedbackEditor = null;
      }

      if (this.reviewEditor) {
        this.reviewEditor.destroy();
        this.reviewEditor = null;
      }

      this.hsPrefGet(
        this.hotspotSearchId,
        success => {
          for (var p in success.body) {
            var pref = success.body[p];

            if (pref.key == "captive_2_title") {
              this.hotspot.title = pref.value;
            }
            if (pref.key == "captive_3_logo") {
              this.hotspot.logo = pref.value;
            }
            if (pref.key == "captive_7_background") {
              this.hotspot.color = pref.value;
            }
          }

          this.marketingPrefGet(
            this.hotspotSearchId,
            success => {
              for (var i in success.body) {
                var pref = success.body[i];
                if (pref.value == "true") pref.value = true;

                if (pref.value == "false") pref.value = false;
                this.marketingPrefs[pref.key] = pref.value;
              }

              this.marketingPrefs.isLoading = false;
              this.isLoadingMarketing = false;
              this.set("marketing_hotspot_id", this.hotspotSearchId);

              if (this.marketingPrefs.marketing_3_first_email_enabled) {
                this.feedbackEditor = new Editor({
                  extensions: [
                    new Heading({ levels: [1, 2, 3] }),
                    new Bold(),
                    new Italic()
                  ],
                  content: this.marketingPrefs.marketing_13_feedback_body_text,
                  onUpdate: ({ getJSON, getHTML }) => {
                    this.feedbackEditorHTML = getHTML();
                  }
                });
                this.feedbackEditorHTML = this.marketingPrefs.marketing_13_feedback_body_text;
              }

              if (this.marketingPrefs.marketing_4_second_email_enabled) {
                this.reviewEditor = new Editor({
                  extensions: [
                    new Heading({ levels: [1, 2, 3] }),
                    new Bold(),
                    new Italic()
                  ],
                  content: this.marketingPrefs.marketing_14_review_body_text,
                  onUpdate: ({ getJSON, getHTML }) => {
                    this.reviewEditorHTML = getHTML();
                  }
                });
                this.reviewEditorHTML = this.marketingPrefs.marketing_14_review_body_text;
              }
            },
            error => {
              console.error(error.body);
              this.isLoadingMarketing = false;
            }
          );
        },
        error => {
          console.error(error.body);
        }
      );
    },
    saveMarketingInfo() {
      this.marketingPrefs.isLoading = true;

      var promises = [];
      for (var i in this.marketingPrefs) {
        if (
          i == "marketing_13_feedback_body_text" &&
          this.marketingPrefs["marketing_3_first_email_enabled"]
        ) {
          this.marketingPrefs[i] = this.feedbackEditorHTML;
        }
        if (
          i == "marketing_14_review_body_text" &&
          this.marketingPrefs["marketing_4_second_email_enabled"]
        ) {
          this.marketingPrefs[i] = this.reviewEditorHTML;
        }

        var pref = {
          key: i,
          value: this.marketingPrefs[i]
        };

        if (pref.key != "isLoading") {
          promises.push(
            new Promise((resolve, reject) => {
              pref.value = pref.value.toString();

              this.marketingPrefModify(
                this.hotspotSearchId,
                pref,
                success => {
                  resolve(success);
                },
                error => {
                  reject(error);
                }
              );
            })
          );
        }
      }

      // exec promises
      var context = this;
      Promise.all(promises).then(function(response) {
        context.marketingPrefs.isLoading = true;
        context.getMarketingInfo();
      });
    },
    testMail(type) {
      this.marketingPrefModify(
        type,
        {
          body: this[type + "EditorHTML"]
        },
        success => {
          console.log(success);
        },
        error => {
          console.error(error);
        }
      );
    }
  }
};
</script>
<style scoped>
.subtitle {
  padding-left: 15px !important;
}
.star {
  font-size: 20px;
}
.star:hover {
  color: #f5c12e;
  cursor: pointer;
}
.star-filled {
  font-size: 20px;
  color: #f0ab00;
}
.star-filled:hover {
  cursor: pointer;
  color: #f5c12e;
}

.menubar {
  margin-bottom: 8px;
}
.space-btn {
  margin-left: 5px;
}
#feedbackEditorContent {
  padding: 5px;
  border: 1px solid #bbbbbb;
  background: white;
}
#reviewEditorContent {
  padding: 5px;
  border: 1px solid #bbbbbb;
  background: white;
}
.test-mail {
  margin-top: 8px;
}
</style>