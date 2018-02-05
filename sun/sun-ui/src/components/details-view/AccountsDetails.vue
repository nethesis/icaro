<template>
  <div>
    <h2>Manager:
      <strong class="soft">{{ info.data.username }}</strong>
    </h2>

    <div class="row row-cards-pf">
      <div class="col-xs-12 col-sm-12 col-md-6">
        <div class="card-pf card-pf-accented">
          <div class="card-pf-heading">
            <h2 class="card-pf-title">
              {{info.data.name}}
              <div v-if="!info.isLoading" :class="[getLoginIcon(info.data.type), 'right']"  data-toggle="tooltip" data-placement="left" :title="$t(info.data.type)"></div>
              <div v-if="info.isLoading" class="spinner spinner-sm right"></div>
            </h2>
          </div>
          <div v-if="!info.isLoading" class="card-pf-body">
            <p>
              <dt>{{ $t("account.hotspot") }}</dt>
              <dd><a :href="'#/hotspots/' + info.data.hotspot_id">{{info.data.hotspot_id}}</a></dd>
            </p>
            <p>
              <dt>{{ $t("account.email") }}</dt>
              <dd>{{info.data.email}}</dd>
            </p>
            <p>
              <dt>{{ $t("account.uuid") }}</dt>
              <dd>{{info.data.uuid}}</dd>
            </p>
            <p>
              <dt>{{ $t("account.created") }}</dt>
              <dd>{{info.data.created}}</dd>
            </p>
          </div>
          <div v-if="!info.isLoading" class="card-pf-footer">
            <div class="dropdown card-pf-time-frame-filter">
              <account-action details="false" :obj="info.data" :update="getInfo"></account-action>
            </div>
            <p>
              <a href="#" class="card-pf-link-with-icon">
              </a>
            </p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
  import AccountService from '../../services/account';
  import StorageService from '../../services/storage';
  import UtilService from '../../services/util';
  import AccountAction from '../../directives/AccountAction.vue';

  export default {
    name: 'AccountsDetails',
    mixins: [AccountService, StorageService, UtilService],
    components: {
       accountAction: AccountAction
    },
    data() {
      // get account info
      this.getInfo()

      return {
        info: {
          isLoading: true,
          data: {}
        },
      }
    },
    // enable tooltips after rendering
    updated: function () {
        $('[data-toggle="tooltip"]').tooltip()
    },
    methods: {
      getInfo() {
        this.accountGet(this.$route.params.id, success => {
          this.info.data = success.body
          this.info.isLoading = false
        }, error => {
          console.log(error.body)
        })
      },
    }
  }

</script>

