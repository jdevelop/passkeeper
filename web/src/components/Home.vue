<template>
  <b-container class="bg-light h-75" fluid>
    <b-card title="Passwords" class="h-100" title-tag="h2">
      <b-row>
        <b-col class="d-flex justify-content-center">
          <b-table-lite
            :fields="fields"
            :items="items"
            small
            responsive="sm"
            head-variant="dark"
            class="datatable"
          >
            <template slot="actions" slot-scope="v">
              <b-btn size="sm" variant="light" @click="changePassword(v.item)">
                <font-awesome-icon icon="edit" />
              </b-btn>|
              <b-btn size="sm" variant="light" @click="removePassword(v.item)">
                <font-awesome-icon icon="trash" />
              </b-btn>
            </template>
            <template slot="table-caption">
              <div class="d-flex justify-content-center">
                <b-btn variant="dark" v-b-modal.add-new>Add new</b-btn>
              </div>
            </template>
          </b-table-lite>
        </b-col>
      </b-row>
    </b-card>
    <b-modal
      id="add-new"
      ref="passDialog"
      header-bg-variant="dark"
      header-text-variant="light"
      title="Add new credentials"
      centered
      @ok="saveCredentials"
    >
      <b-form>
        <b-form-group id="input-group-1" label-for="input-1" description="Password service">
          <b-alert variant="danger" :show="modalErrors.serviceField">{{ modalErrors.serviceField }}</b-alert>
          <b-form-input
            v-bind:class="{'border border-danger': modalErrors.serviceField }"
            id="input-1"
            type="text"
            v-model="credentials.service"
            placeholder="Service"
          ></b-form-input>
        </b-form-group>

        <b-form-group id="input-group-2" description="Password" label-for="input-2">
          <b-alert
            variant="danger"
            :show="modalErrors.passwordField"
          >{{ modalErrors.passwordField }}</b-alert>
          <b-form-input
            v-bind:class="{'border border-danger': modalErrors.passwordField }"
            id="input-2"
            type="password"
            v-model="credentials.secret"
            required
            placeholder="Password"
          ></b-form-input>
        </b-form-group>

        <b-form-group id="input-group-3" description="Confirm password" label-for="input-3">
          <b-alert variant="danger" :show="modalErrors.confirmField">{{ modalErrors.confirmField }}</b-alert>
          <b-form-input
            v-bind:class="{'border border-danger': modalErrors.confirmField }"
            id="input-3"
            type="password"
            v-model="credentials.confirm"
            required
            placeholder="Confirm"
          ></b-form-input>
        </b-form-group>

        <b-form-group id="input-group-4" description="Comment" label-for="input-4">
          <b-form-input
            id="input-4"
            type="text"
            v-model="credentials.comment"
            placeholder="Comment"
          ></b-form-input>
        </b-form-group>
      </b-form>
    </b-modal>
  </b-container>
</template>

<script>
import { REST } from "@/js/restapi.js";
import { Model } from "@/js/models.js";
const columnNames = [
  {
    key: "service",
    label: "Service",
    sortable: false,
    thClass: "col-4",
    class: "v-middle"
  },
  {
    key: "comment",
    label: "Comment",
    sortable: false,
    class: "v-middle"
  },
  {
    key: "actions",
    label: "",
    sortable: false,
    thClass: "col-1",
    class: "v-middle"
  }
];
export default {
  name: "Home",
  mounted() {
    const t = this;
    REST.ListCredentials(function(data) {
      t.items = data;
    });
  },
  data() {
    return {
      modalErrors: {},
      credentials: Model.Credentials(),
      fields: columnNames,
      items: []
    };
  },
  methods: {
    saveCredentials(bvModalEvt) {
      const t = this;
      t.modalErrors = {};
      if (t.credentials.service == "") {
        t.modalErrors = { serviceField: "Please provide service name" };
        bvModalEvt.preventDefault();
        return;
      }
      if (t.credentials.secret == "") {
        t.modalErrors = { passwordField: "Please provide password" };
        bvModalEvt.preventDefault();
        return;
      }
      if (t.credentials.secret != t.credentials.confirm) {
        t.modalErrors = { confirmField: "Password mismatch" };
        bvModalEvt.preventDefault();
        return;
      }
      REST.SaveCredentials(t.credentials, function() {
        t.credentials = Model.Credentials();
        REST.ListCredentials(function(data) {
          t.items = data;
        });
      });
    },
    changePassword(pwd) {
      this.credentials = pwd;
      this.$refs.passDialog.show();
    },
    removePassword(pwd) {
      const t = this;
      this.$bvModal
        .msgBoxConfirm(`Please confirm removal of ${pwd.service}`, {
          title: "Password removal",
          size: "sm",
          buttonSize: "sm",
          okVariant: "danger",
          okTitle: "YES",
          cancelTitle: "NO",
          footerClass: "p-2",
          hideHeaderClose: false,
          centered: true
        })
        .then(value => {
          if (value) {
            REST.RemoveCredentials(pwd.id, () =>
              REST.ListCredentials(d => (t.items = d))
            );
          }
        })
        .catch(err => {
          console.log("Error on list", err);
        });
    }
  }
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style>
@media only screen and (min-width: 1000px) {
  .datatable {
    width: 800px;
  }
  .v-middle {
    vertical-align: middle !important;
  }
}
</style>