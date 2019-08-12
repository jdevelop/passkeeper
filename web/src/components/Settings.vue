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
    </b-card>
  </b-container>
</template>

<script>
import { REST } from "@/js/restapi.js";
import { FileUploadService } from "v-file-upload";
export default {
  data() {
    return {};
  },
  methods: {
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
</style>
