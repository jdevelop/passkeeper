<template>
  <b-container class="bg-white h-75" fluid>
    <b-card title="Manage" class="h-100" title-tag="h2">
      <b-row>
        <b-col sz="10" lg="4" offset-lg="4">
          <b-btn size="lg" variant="success" @click="backup()">Backup</b-btn>&nbsp;
          <b-btn size="lg" variant="danger" @click="restore()">Restore</b-btn>
        </b-col>
      </b-row>
    </b-card>
  </b-container>
</template>

<script>
import { REST } from "@/js/restapi.js";
export default {
  data() {
    return {};
  },
  methods: {
    backup() {
      REST.Backup();
    },
    restore() {
      console.log("Restore");
    }
  }
};
</script>

<style>
</style>
