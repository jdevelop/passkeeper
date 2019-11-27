<template>
  <b-container fluid id="app" class="vh-100">
    <b-navbar type="dark" variant="dark">
      <b-navbar-brand>PassKeeper 2.0</b-navbar-brand>
      <b-navbar-nav small>
        <b-nav-item>
          <router-link to="/" class="menulinks" exact active-class="current-link">Dashboard</router-link>
        </b-nav-item>
        <b-nav-item>
          <router-link to="/settings" class="menulinks" exact active-class="current-link">Manage</router-link>
        </b-nav-item>
      </b-navbar-nav>
    </b-navbar>
    <router-view />
  </b-container>
</template>

<script>
export default {
  name: "app"
};
</script>

<style>
#app {
  font-family: "Avenir", Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  margin-top: 5x;
}
.menulinks {
  font-size: 12pt;
  text-decoration: none;
  color: burlywood !important;
}
.current-link {
  font-weight: bolder;
  color: goldenrod;
}
</style>
