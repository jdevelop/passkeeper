<template>
  <b-container class="bg-white h-75" fluid>
    <b-card title="Manage" class="h-100" title-tag="h2">
      <b-row>
        <b-col sz="10" lg="4" offset-lg="4">
          <b-btn size="lg" variant="success" @click="backup()">Backup</b-btn>&nbsp;
          <b-btn size="lg" variant="danger" @click="popRestore()">Restore</b-btn>
          <input type="file" hidden @change="restore()" ref="restoreF" />
        </b-col>
      </b-row>
      <b-row class="mt-5">
        <b-col class="mt-5">
          <b-btn size="lg" variant="danger" v-b-toggle.collapse-card-access>Change Card Password</b-btn>
        </b-col>
      </b-row>
      <b-collapse id="collapse-card-access">
        <b-row class="mt-5">
          <b-col sz="10" lg="4" offset-lg="4">
            <b-jumbotron
              bg-variant="danger"
              header="Dangerous Area"
              lead="Change RFID card password"
              class="pwdform"
            >
              <b-form-group id="input-group-1" description="Password" label-for="password">
                <b-form-input
                  id="password"
                  type="password"
                  v-model="passwordUpdate.password"
                  placeholder="Password"
                ></b-form-input>
              </b-form-group>
              <b-form-group id="input-group-2" description="Confirmation" label-for="confirmation">
                <b-form-input
                  id="confirmation"
                  type="password"
                  v-model="passwordUpdate.confirm"
                  placeholder="Confirmation"
                ></b-form-input>
              </b-form-group>
              <b-btn variant="success" @click="changeCardPassword()">CHANGE</b-btn>
            </b-jumbotron>
          </b-col>
        </b-row>
      </b-collapse>
    </b-card>
  </b-container>
</template>

<script>
import { REST } from "@/js/restapi.js";
import { FileUploadService } from "v-file-upload";
export default {
  data() {
    return {
      passwordUpdate: {
        password: "",
        confirm: ""
      }
    };
  },
  methods: {
    changeCardPassword() {
      if (
        this.passwordUpdate.password == "" ||
        this.passwordUpdate.password != this.passwordUpdate.confirm
      ) {
        this.$bvModal.msgBoxOk("Password empty or mismatch", {
          title: "Error processing passwords",
          size: "sm",
          buttonSize: "lg",
          okVariant: "danger",
          okTitle: "Close",
          headerClass: "p-2 border-bottom-0",
          footerClass: "p-2 border-top-0",
          centered: true
        });
        return;
      }
      const t = this;
      t.$bvModal
        .msgBoxConfirm("Are you sure?", {
          title: "Updating passwords on the RFID card could lead to data loss",
          size: "sm",
          buttonSize: "lg",
          okVariant: "danger",
          okTitle: "Yes, proceed",
          cancelTitle: "Skip",
          headerClass: "p-2 border-bottom-0",
          footerClass: "p-2 border-top-0",
          centered: true
        })
        .then(confirm => {
          if (confirm) {
            REST.UpdateCardPassword(
              t.passwordUpdate,
              () => {
                t.$bvModal.msgBoxOk("Successfully updated!", {
                  title: "Password update done",
                  size: "sm",
                  buttonSize: "lg",
                  okVariant: "success",
                  okTitle: "Close",
                  headerClass: "p-2 border-bottom-0",
                  footerClass: "p-2 border-top-0",
                  centered: true
                });
              },
              fail => {
                t.$bvModal.msgBoxOk("Update failed", {
                  title: fail,
                  size: "sm",
                  buttonSize: "lg",
                  okVariant: "danger",
                  okTitle: "Close",
                  headerClass: "p-2 border-bottom-0",
                  footerClass: "p-2 border-top-0",
                  centered: true
                });
              }
            );
          }
        });
    },
    backup() {
      REST.Backup();
    },
    popRestore() {
      this.$refs.restoreF.click();
    },
    restore() {
      const fileUpload = new FileUploadService(
        `${process.env.VUE_APP_API_URL}/restore`
      );
      fileUpload.upload(this.$refs.restoreF.files[0]).then(() => {
        this.$router.push({ path: "/" });
      });
    }
  }
};
</script>

<style>
.pwdform {
  color: white !important;
  font-weight: bolder;
}
</style>
